package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hiddengearz/jsubfinder/core"
	"github.com/spf13/cobra"
)

func init() {
	//rootCmd.AddCommand(cmdExec)
	searchExec.PersistentFlags().StringVarP(&inputFile, "inputFile", "f", "", "File containing domains")
	searchExec.PersistentFlags().StringVarP(&url, "url", "u", "", "url to check")

}

//Search the provided URL's
var searchExec = &cobra.Command{
	Use:   "search",
	Short: "Search javascript/URL's for domains",
	Long:  `Execute the command specified`,
	Run: func(cmd *cobra.Command, arguments []string) {

	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			defer core.TimeTrack(time.Now(), "JSubfinder")
		}
		safetyChecks()
		getURLs()
		fmt.Println(core.Urls)
	},
}

//Retrieve the URL's needed to be searched
func getURLs() {
	if url != "" { //if -u isn't empty
		core.Urls = append(core.Urls, url)
	} else if inputFile != "" { //else if -f isn't empty
		core.Urls = core.ReadFile(inputFile)
	} else { //finally try from input being piped
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			core.Urls = append(core.Urls, strings.ToLower(sc.Text()))
		}
	}
}
