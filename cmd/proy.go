package cmd

import (
	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/spf13/cobra"
)

//Start JSubFiner in proxy mode
var proxyExec = &cobra.Command{
	Use:   "proxy",
	Short: "Run JSubfinder as a proxy",
	Long:  `Execute the command specified`,
	Run: func(cmd *cobra.Command, arguments []string) {

	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := safetyChecks()
		if err != nil {
			l.Log.Fatal(err)
		}
	},
}
