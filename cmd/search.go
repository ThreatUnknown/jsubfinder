package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/hiddengearz/jsubfinder/core"
	C "github.com/hiddengearz/jsubfinder/core"
	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/spf13/cobra"
)

func init() {
	//rootCmd.AddCommand(cmdExec)
	searchExec.PersistentFlags().StringVarP(&C.InputFile, "inputFile", "f", "", "File containing domains")
	searchExec.PersistentFlags().StringSliceVarP(&C.InputURLs, "url", "u", []string{}, "Url to check")
	searchExec.PersistentFlags().BoolVarP(&C.Crawl, "crawl", "c", false, "Enable crawling")
	searchExec.PersistentFlags().IntVarP(&C.Threads, "threads", "t", 5, "Ammount of threads to be used")
	searchExec.PersistentFlags().BoolVarP(&C.Greedy, "greedy", "g", false, "Check all files for URL's not just Javascript")

}

//Search the provided URL's
var searchExec = &cobra.Command{
	Use:   "search",
	Short: "Search javascript/URL's for domains",
	Long:  `Execute the command specified`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		C.Command = "search"
		err := safetyChecks() //safety checks
		if err != nil {
			l.Log.Fatal(err)
		}

		err = getURLs() //Get the URL's from input file, -u flag or stdins
		if err != nil {
			l.Log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, arguments []string) {
		if C.Debug {
			defer core.TimeTrack(time.Now(), "searchExec")
		}
		C.ExecSearch() //Start the search
	},
}

//Retrieve the URL's needed to be searched
func getURLs() (err error) {
	if len(C.InputURLs) != 0 { //if -u isn't empty
		return
	} else if C.InputFile != "" { //else if -f isn't empty
		core.InputURLs, err = core.ReadFile(C.InputFile)
		return
	} else if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 { //finally try from input being piped (stdin)
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			core.InputURLs = append(core.InputURLs, strings.ToLower(sc.Text()))
		}
		return
	} else {
		return errors.New("no URL's provided")
	}

}
