package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listFamilySamplesCmd represents the listFamilySamples command
var listFamilySamplesCmd = &cobra.Command{
	Use:   "familySamples",
	Short: "Will list all the indexed samples for a family",
	Long: `Will list all the indexed samples for a family. 
	
Example usage:
- malpedia_cli familySamples flame
- malpedia_cli familySamples stuxnet --json
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("familySamples requires a family name")
		}

		family, err := util.GetFamilyName(args[0], apiKey)
		if err != nil {
			log.Fatalf("unable to find family: %s", args[0])
		}

		formattedEndpoint := fmt.Sprintf(types.EndpointListFamilySamples, family)

		resp, err := util.HttpGetQuery(types.Endpoint(formattedEndpoint), apiKey)
		if err != nil {
			log.Fatalf("unable to get family listing: %s", err)
		}

		if jsonFormat {
			util.PrettyPrintJson(resp)
		} else {
			samples := &types.FamilySamples{}
			err := json.Unmarshal(resp, samples)
			if err != nil {
				log.Fatal("failed to unmarshal response")
			}

			for _, sample := range *samples {
				log.WithFields(logrus.Fields{
					"Status":  sample.Status,
					"Version": sample.Version,
				}).Infof("Family Information for %s sample %s", args[0], sample.Sha256)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listFamilySamplesCmd)
}
