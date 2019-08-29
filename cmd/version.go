package cmd

// DONE

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "returns the current version of malpedia",
	Long: `returns the current version of malpedia.
Example usage:  

- malpedia_cli version
`,
	Run: func(cmd *cobra.Command, args []string) {
		buf, err := util.HttpGetQuery(types.EndpointVersion, apiKey)
		if err != nil {
			log.Fatal(err)
		}

		if jsonFormat {
			err = util.PrettyPrintJson(buf)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			version := &types.Version{}
			err = json.Unmarshal(buf, version)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Date: " + version.Date)
			fmt.Println("Version: " + strconv.Itoa(version.Version))
		}

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
