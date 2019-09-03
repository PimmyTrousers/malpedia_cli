package cmd

import (
	"fmt"
	"os"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

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
		}

		f, err := os.Open(args[0])
		if err != nil {
			log.Fatal("failed to acquire file handle on yara rule")
		}

		resp, err := util.HttpRawFileUpload(types.EndpointScanYara, apiKey, f, f.Name())
		if err != nil {
			fmt.Println(err)
			log.Fatal("failed to scan rule against samples")
		}

		if jsonFormat {
			util.PrettyPrintJson(resp)
		} else {
			// TODO: Implement structured output
		}
	},
}

func init() {
	rootCmd.AddCommand(scanYaraCmd)
}
