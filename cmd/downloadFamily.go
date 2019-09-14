package cmd

import (
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

var (
	familyRaw bool
)

// downloadFamilyCmd represents the downloadFamily command
var downloadFamilyCmd = &cobra.Command{
	Use:   "downloadFamily",
	Short: "Will take a known family and download all of the samples within malpedia",
	Long: `Will take a known family and download all of the samples within malpedia. 
	
Example Usage:
	- malpedia_cli downloadFamily ursnif
	- malpedia_cli downloadFamily emotet
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("downloadFamily requires 1 argument only")
		} else if !util.IsAPIKeyValid(apiKey) {
			log.Fatal("apikey is required")
		}

		name, err := util.GetFamilyName(args[0], apiKey)
		if err != nil {
			log.Fatal(err)
		}

		samples, err := getHashes(name)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range *samples {
			for _, sample := range v {
				log.Infof("downloading %s sample: %s", k, sample)

				states, err := util.DownloadSample(sample, apiKey)
				if err != nil {
					log.Warnf("unable to download %s", sample)
					continue
				}

				log.Infof("writing %s sample: %s", k, sample)

				if familyRaw {
					err = util.DumpRaw(states, sample)
				} else {
					err = util.DumpZip(states, sample, sample+".zip")
				}

				if err != nil {
					log.Warnf("unable to write %s", sample)
					continue
				}
			}
		}

	},
}

func init() {

	rootCmd.AddCommand(downloadFamilyCmd)
	downloadFamilyCmd.PersistentFlags().BoolVarP(&familyRaw, "raw", "r", false, "the returned family will have all of it's samples written to disk in raw")
}
