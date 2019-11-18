package main

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
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGetActors(apiKey)
	},
}

func runGetActors(key string) error {
	// TODO: make apiKey a param for testability
	resp, err := util.HttpGetQuery(types.EndpointGetActors, key)
	if err != nil {
		return err
	}

	if jsonFormat {
		fmt.Println(string(resp))
		err = util.PrettyPrintJson(resp)
		if err != nil {
			return err
		}
	} else {
		var actors types.Actors
		err = json.Unmarshal(resp, &actors)
		if err != nil {
			return err
		}
		for _, actor := range actors {
			fmt.Println(actor)
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(getActorsCmd)
}
