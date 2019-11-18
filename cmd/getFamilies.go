package main

// DONE

import (
	"encoding/json"
	"fmt"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"

	"github.com/spf13/cobra"
)

// getFamiliesCmd represents the getFamilies command
var getFamiliesCmd = &cobra.Command{
	Use:   "families",
	Short: "families will returns all malware families tracked in Malpedia",
	Long: `families will returns all malware families tracked in Malpedia

Usage examples: 
- malpedia_cli families --json
- malpedia_cli families`,
	Run: func(cmd *cobra.Command, args []string) {
		if !util.IsAPIKeyValid(apiKey) {
			log.Fatal("apikey is required")
		}

		resp, err := util.HttpGetQuery(types.EndpointGetFamilies, apiKey)
		if err != nil {
			log.Fatal(err)
		}

		if jsonFormat {
			err = util.PrettyPrintJson(resp)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			var families types.Families
			err = json.Unmarshal(resp, &families)
			if err != nil {
				log.Fatal(err)
			}
			for name := range families {
				fmt.Println(name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getFamiliesCmd)
}
