package cmd

import (
	"fmt"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

// getYaraCmd represents the getYara command
var getYaraCmd = &cobra.Command{
	Use:   "getYara",
	Short: "getYara will return all yara rules for a specific family",
	Long: `getYara will return all yara rules for a specific family. Only a
a single family can be passed at a time. The yara rules will be dropped in a zip
file by default but the user may specify having the rules dropped in their raw format. 
If no output name is passed the commands argument will be used as the 
name for the zip file. 

Example Usage:
- malpedia_cli getYara ursnif 
- malpedia_cli getYara emotet -r 
- malpedia_cli getYara njrat -o njrat.zip
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("getYara expects a single family as an argument")
		}

		familyName, err := util.GetFamilyName(args[0])
		if err != nil {
<<<<<<< HEAD
			log.Fatalf("family %s not found", args[0])
=======
			log.Fatalf("unable to get yara rules: %s", err)
		} else {
			log.Printf("successfully wrote rules to: %s", outputFileName)
>>>>>>> 09f1726... rough implementation of post request and scanBinary against yara rules
		}

		formattedEndpoint := fmt.Sprintf(types.EndpointGetYaraRulesForFamily, familyName)

		// TODO: currently dont have any internet, so need to parse this result
		_, err = util.HttpGetQuery(types.Endpoint(formattedEndpoint), apiKey)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
<<<<<<< HEAD
	rootCmd.AddCommand(getYaraCmd)
=======
	rootCmd.AddCommand(getYaraRulesCmd)

	getYaraRulesCmd.Flags().BoolVarP(&zipDownload, "zip", "z", false, "will download a zip file with all Yara rules")
	getYaraRulesCmd.Flags().StringVarP(&outputFileName, "output", "o", "yara_rules", "output location of the yara rules file")
>>>>>>> 09f1726... rough implementation of post request and scanBinary against yara rules
}
