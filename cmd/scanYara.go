package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

var testInput = `
{
    "COZY_FANCY_BEAR_Hunt": {
        "win.xtunnel": {
            "2016-04-25/4845761c9bed0563d0aa83613311191e075a9b58861e80392914d61a21bad976": {
                "matches": true,
                "matched_strings": [
                    "0x00150195:$s6|34352e33322e3132392e313835"
                ]
            },
            "2016-05-05/40ae43b7d6c413becc92b07076fa128b875c8dbb4da7c036639eccf5a9fc784f": {
                "matches": true,
                "matched_strings": [
                    "0x00150195:$s7|32332e3232372e3139362e323137"
                ]
            }
        }
    }
}
`

type YaraMatches struct {
	RuleName RuleMetaData
}

type RuleMetaData struct {
	Families map[string]Family
}

type Family struct {
	Matches        bool     `json:"matches"`
	MatchedStrings []string `json:"matched_strings"`
}

// scanYaraCmd represents the scanYara command
var scanYaraCmd = &cobra.Command{
	Use:   "scanYara",
	Short: "Will upload a yara rule to malpedia and queue it up to be scanned against all samples contained in malpedia",
	Long: `Will upload a yara rule to malpedia and queue it up to be scanned against all samples contained in malpedia.
Rules must be raw text and multiple rules can be contained in a single file. Imported Yara modules are most likely not supported.

Example usage:
- malpedia_cli scanYara myRule.yar
- malpedia_cli scanYara myRule.yar --json
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("scanYara requires a single file as an argument")
		} else if !util.IsAPIKeyValid(apiKey) {
			log.Fatal("apikey is required")
		}

		// f, err := os.Open(args[0])
		// if err != nil {
		// 	log.Fatal("failed to acquire file handle on yara rule")
		// }

		// resp, err := util.HttpRawFileUpload(types.EndpointScanYara, apiKey, f, f.Name())
		// if err != nil {
		// 	fmt.Println(err)
		// 	log.Fatal("failed to scan rule against samples")
		// }

		if jsonFormat {
			// util.PrettyPrintJson(resp)
		} else {
			result := &map[string]interface{}{}
			err := json.Unmarshal([]byte(testInput), result)
			if err != nil {
				fmt.Println(err)
				log.Fatal("failed to unmarshal data")
			}
			for k, v := range *result {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanYaraCmd)
}

// func customUnmarshal(buf map[string]interface{}) (YaraMatches, error) {
// 	return nil, nil
// }
