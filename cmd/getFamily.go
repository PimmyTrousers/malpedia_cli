package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getFamilyCmd represents the getFamily command
var getFamilyCmd = &cobra.Command{
	Use:   "getFamily",
	Short: "getFamily will return information about a malware family",
	Long: `getFamily will return information about a malware family.

Usage examples: 
- malpedia_cli getFamily ursnif
- malpedia_cli getFamily emotet --json`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("getFamily requires 1 argument only")
		}

		name, err := util.GetFamilyName(args[0])
		if err != nil {
			log.Fatal(err)
		}

		endpoint := fmt.Sprintf(types.EndpointGetFamily, name)
		resp, err := util.HttpGetQuery(types.Endpoint(endpoint), apiKey)
		if err != nil {
			log.Fatal(err)
		}
		if jsonFormat {
			err = util.PrettyPrintJson(resp)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			family := &types.Family{}
			err = json.Unmarshal(resp, family)
			if err != nil {
				log.Fatal(err)
			}
			log.WithFields(logrus.Fields{
				"Attribution": family.Attribution,
				"Description": family.Description,
				"Common Name": family.CommonName,
				"Sources":     family.Sources,
			}).Infof("Family Information for %s", name)
		}

	},
}

func init() {
	rootCmd.AddCommand(getFamilyCmd)
}
