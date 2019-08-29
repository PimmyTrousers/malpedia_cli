package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/pimmytrousers/malpedia_cli/types"
	"github.com/pimmytrousers/malpedia_cli/util"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getActorCmd represents the getActor command
var getActorCmd = &cobra.Command{
	Use:   "getActor",
	Short: "Will return metadata about a specific actor",
	Long: `This will make 2 requests, one to check that the actor ID is valid and another
to request the metadata if the actor ID is valid

Example Usage:
- malpedia_cli getActor apt28
- malpedia_cli getActor apt28 --json
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("getActor takes in only a single argument")
		}

		name, err := util.GetActorName(args[0])
		if err != nil {
			log.Fatal(err)
		}

		endpoint := types.EndpointGetActor
		endpoint = fmt.Sprintf(endpoint, name)
		res, err := util.HttpGetQuery(types.Endpoint(endpoint), apiKey)
		if err != nil {
			if err == types.ErrResourceNotFound {
				log.Fatal("No record of threat actor")
			}

		}

		actor := &types.Actor{}
		err = json.Unmarshal(res, actor)
		if err != nil {
			log.Fatal("failed to parse response data")
		}

		log.WithFields(
			logrus.Fields{
				"description": actor.Description,
				"country":     actor.Meta.Country,
			},
		).Info("actor informtation")

		for name, family := range actor.Families {
			log.WithFields(
				logrus.Fields{
					"name":        name,
					"common name": family.CommonName,
				},
			).Infof("family results for %s", args[0])
		}

	},
}

func init() {
	rootCmd.AddCommand(getActorCmd)
}
