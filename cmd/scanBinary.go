package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

// scanBinaryCmd represents the scanBinary command
var scanBinaryCmd = &cobra.Command{
	Use:   "scanBinary",
	Short: "scanBinary will upload a file to malpedia and scan it against all the yara rules there",
	Long:  `scanBinary will upload a file to malpedia and scan it against all the yara rules there`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("scanBinary requires a single file as an argument")
		}

		f, err := os.Open(args[0])
		if err != nil {
			log.Fatal("failed to acquire file handle")
		}

		resp, err := util.HttpPostQuery(types.EndpointScanBinary, apiKey, f, "fileToBeScanned.bin")
		fmt.Printf("%s\n", resp)
		if err != nil {
			fmt.Println(err)
			log.Fatal("failed to scan file against yara rules")
		}

		if jsonFormat {
			util.PrettyPrintJson(resp)
		} else {
			matches := &map[string]types.YaraMatchesValue{}
			err := json.Unmarshal(resp, matches)
			if err != nil {
				fmt.Println(err)
				log.Fatal("failed to unmarshal response")
			}

			for name, match := range *matches {
				if match.Match {
					log.Info("%s matched yara rule %s", args[0], name)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanBinaryCmd)
}
