package cmd

// DONE

import (
	"fmt"
	"os"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"

	"github.com/spf13/cobra"
)

var zipDownload bool
var outputFileName string

// getYaraCmd represents the getYara command
var getYaraRulesCmd = &cobra.Command{
	Use:   "getYaraRules",
	Short: "Will download all Yara rules contained within Malpedia",
	Long: `Will download all Yara rules contained within Malpedia.
Can pass a zip flag to request all the Yara rules as a zip file.
All requests require a TLP level of eithe white green or amber.

Example usage:
- malpedia_cli getYaraRules white
- malpedia_cli getYaraRules green -o yara_rules
- malpedia_cli getYaraRules amber -z -o yara_rules.zip`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("getYaraRules requires an argument of TLP level")
		}
		var tlpLevel string

		switch args[0] {
		case "white":
			tlpLevel = "tlp_white"
		case "green":
			tlpLevel = "tlp_green"
		case "amber":
			tlpLevel = "tlp_amber"
		default:
			log.Fatal("Invalid TLP level")
		}

		// Always get zip since its compressed and we can decompress it
		endpoint := fmt.Sprintf(types.EndpointGetYaraZip, tlpLevel)

		resp, err := util.HttpGetQuery(types.Endpoint(endpoint), apiKey)
		if err != nil {
			log.Fatal(err)
		}

		if zipDownload {
			err = writeZip(resp, outputFileName)
		} else {
			err = writeYaraRules(resp)
		}
		if err != nil {
			log.Fatalf("unable to write yara rules: %s", err)
		}

		log.Printf("successfully wrote yara rules to %s", outputFileName)
	},
}

func writeZip(buf []byte, location string) error {
	f, err := os.Create(location)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func writeYaraRules(buf []byte) error {
	zipFile := outputFileName + "temp"
	err := writeZip(buf, zipFile)
	if err != nil {
		return err
	}

	defer os.Remove(zipFile)

	err = util.Unzip(zipFile, outputFileName)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(getYaraRulesCmd)

	getYaraRulesCmd.Flags().BoolVarP(&zipDownload, "zip", "z", false, "will download a zip file with all Yara rules")
	getYaraRulesCmd.Flags().StringVarP(&outputFileName, "output", "o", "yara_rules", "ouput location of the yara rules file")
}
