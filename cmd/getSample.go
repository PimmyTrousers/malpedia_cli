package cmd

// DONE

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/yeka/zip"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

var (
	raw      bool
	password string
	outzip   string
)

// getSampleCmd represents the getSample command
var getSampleCmd = &cobra.Command{
	Use:   "getSample",
	Short: "returns a sample and potentially unpacked versions given a hash",
	Long: `getSample will default to returning samples in a zip with the password infected, but it 
can also return the samples in raw given the "-r" flag. On top of this a user can provide a password
that will be used to zip encrypt the samples as well as a output filename for the zip file.

Example usage:
	- malpedia_cli getSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -p infected123 -o samples1234.zip
	- malpedia_cli getSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -r 
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("getSample takes a single hash")
		}

		hashtype, err := util.GetHashType(args[0])
		if err != nil {
			log.Fatal("failed to parse hash")
		} else if hashtype != util.SHA256 && hashtype != util.MD5 {
			log.Fatal("getSample requires SHA256 or MD5")
		}

		formattedEndpoint := fmt.Sprintf(types.EndpointGetSampleRaw, args[0])

		res, err := util.HttpGetQuery(types.Endpoint(formattedEndpoint), apiKey)
		if err != nil {
			log.Fatal("failed to make a request to malpedia")
		}

		jsonBody := make(map[string]string)

		err = json.Unmarshal(res, &jsonBody)
		if err != nil {
			log.Fatal("failed to process json body")
		}

		if !raw {
			err := dumpZip(&jsonBody)
			if err != nil {
				log.Fatal("failed to dump samples to disk")
			}
			log.Printf("wrote sample to %s", outzip)
		} else {
			err := dumpRaw(&jsonBody)
			if err != nil {
				log.Fatal("failed to dump samples to disk")
			}
		}
	},
}

func dumpRaw(elems *map[string]string) error {
	for k, v := range *elems {
		buf, err := util.Base64DecodeContent(&v)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(k, buf.Bytes(), 0644)
		if err != nil {
			return err
		}

		log.Printf("wrote sample to %s", k)
	}

	return nil
}

func dumpZip(elems *map[string]string) error {
	fzip, err := os.Create(outzip)
	if err != nil {
		return err
	}

	zipw := zip.NewWriter(fzip)
	defer zipw.Close()

	for k, v := range *elems {
		buf, err := util.Base64DecodeContent(&v)
		if err != nil {
			return err
		}
		// TODO: apparently the unzip binary in UNIX does not support AES encryption
		w, err := zipw.Encrypt(k, password, zip.AES128Encryption)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(getSampleCmd)

	getSampleCmd.PersistentFlags().BoolVarP(&raw, "raw", "r", false, "the returning file will be written to disk in raw")
	getSampleCmd.PersistentFlags().StringVarP(&password, "samplePassword", "p", "infected", "sample password that will be used to encrypt the sample")
	getSampleCmd.PersistentFlags().StringVarP(&outzip, "output", "o", "samples.zip", "name of the resulting zip file")
}
