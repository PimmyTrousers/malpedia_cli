package main

import (
	"fmt"
	"os"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

// scanYaraAgainstFamilyCmd represents the scanYaraAgainstFamily command
var scanYaraAgainstFamilyCmd = &cobra.Command{
	Use:   "scanYaraAgainstFamily",
	Short: "Behaves similarly to scanYara except its targeted to a single family gvien by the user",
	Long: `Behaves similarly to scanYara except its targeted to a single family gvien by the user.
This can offer performance benefits since the subset of samples being scanned against will be much smaller.

Example usage:
- malpedia_cli scanYaraAgainstFamily ursnif myRule.yar
- malpedia_cli scanYaraAgainstFamily emotet myRule.yar --json
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("scanYaraAgainstFamily requires a family name and a path to a yara rule")
		}

		name, err := util.GetFamilyName(args[0], apiKey)
		if err != nil {
			log.Fatal(err)
		}

		formattedEndpoint := fmt.Sprintf(types.EndpointScanYaraAgainstFamily, name)

		f, err := os.Open(args[1])
		if err != nil {
			log.Fatal("failed to acquire file handle on yara rule")
		}

		resp, err := util.HttpRawFileUpload(types.Endpoint(formattedEndpoint), apiKey, f, f.Name())
		if err != nil {
			fmt.Println(err)
			log.Fatalf("failed to scan rule against %s samples", name)
		}

		if jsonFormat {
			util.PrettyPrintJson(resp)
		} else {
			// TODO: Implement structured output
		}
	},
}

func init() {
	rootCmd.AddCommand(scanYaraAgainstFamilyCmd)
}
