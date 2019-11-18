package main

// DONE

import (
	"github.com/pimmytrousers/malpedia_cli/util"
	"github.com/spf13/cobra"
)

var (
	sampleRaw    bool
	password     string
	sampleOutzip string
)

// getSampleCmd represents the getSample command
var getSampleCmd = &cobra.Command{
	Use:   "downloadSample",
	Short: "returns a sample and potentially unpacked versions given a hash",
	Long: `downloadSample will default to returning samples in a zip named samples.zip, but it 
can also return the samples in raw given the "-r" flag. 

Example usage:
	- malpedia_cli downloadSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -o samples1234.zip
	- malpedia_cli downloadSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -r 
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("downloadSample takes a single hash")
		} else if !util.IsAPIKeyValid(apiKey) {
			log.Fatal("apikey is required")
		}

		states, err := util.DownloadSample(args[0], apiKey)
		if err != nil {
			log.Fatal("unable to download sample")
		}

		var dumpErr error

		if !sampleRaw {
			dumpErr = util.DumpZip(states, args[0], sampleOutzip)
		} else {
			dumpErr = util.DumpRaw(states, args[0])
		}

		if dumpErr != nil {
			log.Fatal("failed to dump samples to disk")
		}

		log.Infof("successfully wrote sample states for %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(getSampleCmd)
	getSampleCmd.PersistentFlags().BoolVarP(&sampleRaw, "raw", "r", false, "the returning file will be written to disk in raw")
	getSampleCmd.PersistentFlags().StringVarP(&sampleOutzip, "output", "o", "samples.zip", "name of the resulting zip file")
}
