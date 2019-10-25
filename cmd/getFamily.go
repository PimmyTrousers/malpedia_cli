package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"

	"github.com/spf13/cobra"
)

var printHashes bool

// getFamilyCmd represents the getFamily command
var getFamilyCmd = &cobra.Command{
	Use:   "family",
	Short: "family will return information about a malware family",
	Long: `family will return information about a malware family.

Usage examples: 
- malpedia_cli family ursnif
- malpedia_cli family njrat --samples
- malpedia_cli family emotet --json`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("family requires 1 argument only")
		} else if !util.IsAPIKeyValid(apiKey) {
			log.Fatal("family requres an API key")
		}

		name, err := util.GetFamilyName(args[0], apiKey)
		if err != nil {
			log.Fatal(err)
		}

		endpoint := fmt.Sprintf(types.EndpointGetFamily, name)
		resp, err := util.HttpGetQuery(types.Endpoint(endpoint), apiKey)
		if err != nil {
			log.Fatal(err)
		}
		if jsonFormat {
			err = util.PrettyPrintJson(resp)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			family := &types.Family{}
			err = json.Unmarshal(resp, family)
			if err != nil {
				log.Fatal(err)
			}

			t := table.NewWriter()
			t.SetAllowedColumnLengths([]int{20, 128})

			t.Style().Options.SeparateColumns = true
			t.Style().Options.DrawBorder = true
			t.SetStyle(table.StyleRounded)
			t.Style().Options.SeparateColumns = true
			t.Style().Options.DrawBorder = true
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Field", "Value"})

			t.AppendRow(table.Row{"Common Name", family.CommonName + "\n"})

			t.AppendRow(table.Row{"Last Updated", family.Updated + "\n"})

			if len(family.AltNames) > 0 {
				for i, alias := range family.AltNames {
					key := fmt.Sprintf("Alias %d", i+1)
					val := alias
					if i == len(family.AltNames)-1 {
						val += "\n"
					}
					t.AppendRow(table.Row{key, val})
				}
			}

			if len(family.Attribution) > 0 {
				for i, attr := range family.Attribution {
					key := fmt.Sprintf("Attribution %d", i+1)
					val := attr
					if i == len(family.Attribution)-1 {
						val += "\n"
					}
					t.AppendRow(table.Row{key, val})
				}
			}

			family.Description = strings.ReplaceAll(family.Description, "\t", " ")
			family.Description = strings.ReplaceAll(family.Description, "\n", " ")
			family.Description = strings.ReplaceAll(family.Description, "\r", "")
			if family.Description != "" {
				t.AppendRow(table.Row{"Description", family.Description + "\n"})
			}

			if len(family.Urls) > 0 {
				for i, url := range family.Urls {
					key := fmt.Sprintf("Reference %d", i+1)
					val := url
					if i == len(family.Urls)-1 {
						val += "\n"
					}
					t.AppendRow(table.Row{key, val})
				}
			}

			// TODO: This is probably over the top
			if printHashes {
				hashes, err := getHashes(name)
				if err == nil {
					for status, list := range *hashes {
						for i, hash := range list {
							key := fmt.Sprintf("%s %d", status, i+1)
							val := hash
							if i == len(list)-1 {
								val += "\n"
							}
							t.AppendRow(table.Row{key, val})
						}
					}
				}
			}

			t.Render()
		}

	},
}

func init() {
	getFamilyCmd.PersistentFlags().BoolVarP(&printHashes, "samples", "s", false, "Will print the hashes of the samples in malpedia for thee famiyl requested")
	rootCmd.AddCommand(getFamilyCmd)
}

func getHashes(family string) (*map[string][]string, error) {
	formattedEndpoint := fmt.Sprintf(types.EndpointListFamilySamples, family)

	resp, err := util.HttpGetQuery(types.Endpoint(formattedEndpoint), apiKey)
	if err != nil {
		return nil, err
	}

	samples := &[]types.ListFamilySamples{}
	err = json.Unmarshal(resp, samples)
	if err != nil {
		return nil, err
	}

	var hashes = map[string][]string{}

	for _, sample := range *samples {
		hashes[sample.Status] = append(hashes[sample.Status], sample.Sha256)
	}

	return &hashes, nil
}
