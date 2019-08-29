package cmd

// DONE

import (
	"errors"
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
	Use:   "getYara",
	Short: "Will download all Yara rules specified by the conditions passed",
	Long: `Will download all Yara rules specified by the conditions passed.
Can pass a zip flag to request all the Yara rules as a zip file.
All requests require a TLP level of eithe white green or amber.

Example usage:
- malpedia_cli getYara tlp white
- malpedia_cli getYara tlp green -o yara_rules
- malpedia_cli getYara tlp amber -z -o yara_rules.zip
- malpedia_cli getYara family ursnif 
- malpedia_cli getYara family emotet
- malpedia_cli getYara family njrat -z -o njrat.zip`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("getYara requires two arguments")
		}

		var err error
		if args[0] == "tlp" {
			err = getYaraByTLP(args[1])
		} else if args[0] == "family" {
			err = getYaraByFamily(args[1])
		} else {
			log.Fatalf("invalid type specifier: %s", args[0])
		}

		if err != nil {
			log.Fatalf("unable to get yara rules: %s", err)
		} else {
			log.Printf("successfully wrote rules to: %s", outputFileName)
		}
	},
}

func getYaraByFamily(family string) error {
	familyName, err := util.GetFamilyName(family, apiKey)
	if err != nil {
		log.Fatalf("family %s not found", family)
	}

	formattedEndpoint := fmt.Sprintf(types.EndpointGetYaraRulesForFamily, familyName)

	resp, err := util.HttpGetQuery(types.Endpoint(formattedEndpoint), apiKey)
	if err != nil {
		return err
	}

	if zipDownload {
		err = writeZip(resp, outputFileName)
	} else {
		err = writeYaraRules(resp)
	}
	if err != nil {
		return err
	}

	return nil
}

func getYaraByTLP(tlpLevel string) error {

	switch tlpLevel {
	case "white":
		tlpLevel = "tlp_white"
	case "green":
		tlpLevel = "tlp_green"
	case "amber":
		tlpLevel = "tlp_amber"
	default:
		return errors.New("invalid TLP level")
	}

	// Always get zip since its compressed and we can decompress it
	endpoint := fmt.Sprintf(types.EndpointGetYaraZip, tlpLevel)

	resp, err := util.HttpGetQuery(types.Endpoint(endpoint), apiKey)
	if err != nil {
		return err
	}

	if zipDownload {
		err = writeZip(resp, outputFileName)
	} else {
		err = writeYaraRules(resp)
	}

	if err != nil {
		return err
	}

	return nil
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
	getYaraRulesCmd.Flags().StringVarP(&outputFileName, "output", "o", "yara_rules", "output location of the yara rules file")
}
