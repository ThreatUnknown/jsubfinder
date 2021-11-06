package cmd

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/hiddengearz/jsubfinder/core"
	C "github.com/hiddengearz/jsubfinder/core"
	"github.com/spf13/cobra"
)

func init() {
	//rootCmd.AddCommand(cmdExec)
	searchExec.PersistentFlags().StringVarP(&C.InputFile, "inputFile", "f", "", "File containing domains")
	searchExec.PersistentFlags().StringVarP(&C.Url, "url", "u", "", "Url to check")

}

//Search the provided URL's
var searchExec = &cobra.Command{
	Use:   "search",
	Short: "Search javascript/URL's for domains",
	Long:  `Execute the command specified`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		safetyChecks()
		getURLs()
	},
	Run: func(cmd *cobra.Command, arguments []string) {
		if C.Debug {
			defer core.TimeTrack(time.Now(), "searchExec")
		}
		C.ExecSearch(C.Threads, C.OutputFile)
	},
}

//Retrieve the URL's needed to be searched
func getURLs() {
	if C.Url != "" { //if -u isn't empty
		core.Urls = append(C.Urls, C.Url)
	} else if C.InputFile != "" { //else if -f isn't empty
		core.Urls = core.ReadFile(C.InputFile)
	} else { //finally try from input being piped
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			core.Urls = append(core.Urls, strings.ToLower(sc.Text()))
		}
	}
}
