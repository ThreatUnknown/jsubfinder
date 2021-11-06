package cmd

import (
	"github.com/spf13/cobra"
)

var cmdExec = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command",
	Long:  `Execute the command specified`,
	Run: func(cmd *cobra.Command, arguments []string) {

	},
}

func init() {
	rootCmd.AddCommand(cmdExec)
	cmdExec.PersistentFlags().IntVarP(&threads, "threads", "t", 5, "Ammount of threads to be used")
	cmdExec.PersistentFlags().StringVarP(&outputFile, "outputFile", "o", "results.txt", "name/location to store the file")
	cmdExec.PersistentFlags().BoolVarP(&greedy, "greedy", "g", false, "Check all files for URL's not just Javascript")
	cmdExec.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	cmdExec.PersistentFlags().BoolVarP(&crawl, "crawl", "c", false, "Enable crawling")
	cmdExec.PersistentFlags().BoolVarP(&findSecrets, "secrets", "s", false, "Check results for secrets e.g api keys")
	cmdExec.PersistentFlags().StringVarP(&sig, "sig", "S", "~/.jsf_signatures.yaml", "Location of signatures for finding secrets")

}
