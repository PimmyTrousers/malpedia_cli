package cmd

// DONE

import (
	"encoding/json"
	"fmt"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"

	"github.com/spf13/cobra"
)

// getActorsCmd represents the getActors command
var getActorsCmd = &cobra.Command{
	Use:   "actors",
	Short: "will return a list of all actors tracked in malpedia",
	Long: `actors will return a list of all actors tracked in malpedia

Example Usage: 
- malpedia_cli actors
- malpedia_cli actors --json`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := util.HttpGetQuery(types.EndpointGetActors, apiKey)
		if err != nil {
			log.Fatal(err)
		}

		if jsonFormat {
			fmt.Println(string(resp))
			err = util.PrettyPrintJson(resp)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			var actors types.Actors
			err = json.Unmarshal(resp, &actors)
			if err != nil {
				log.Fatal(err)
			}
			for _, actor := range actors {
				fmt.Println(actor)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(getActorsCmd)
}
