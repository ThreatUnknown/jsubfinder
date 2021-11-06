package cmd

import (
	"os"

	"github.com/hiddengearz/jsubfinder/core"
	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Laelap",
		Short: "Agent for the [redacted] project",
		Long:  `[redacted]`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			f, err := os.OpenFile("log.info", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				l.Log.Error(err)
			}

			if debug && silent {
				l.Log.Fatal("Please choose Debug mode or silent mode. Not both.")
			} else if debug {

				l.InitDetailedLogger(f)
				l.Log.SetLevel(logrus.DebugLevel)

				l.Log.Debug("Debug mode enabled")
			} else if silent {
				l.Log.SetLevel(logrus.PanicLevel)
			}

		},
	}
	threads     int
	inputFile   string
	url         string
	outputFile  string
	greedy      bool
	debug       bool
	crawl       bool
	findSecrets bool
	sig         string
	silent      bool
)

func init() {
	l.Log = logrus.New()
	rootCmd.AddCommand(searchExec)
	rootCmd.AddCommand(proxyExec)
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 5, "Ammount of threads to be used")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "outputFile", "o", "", "name/location to store the file")
	rootCmd.PersistentFlags().BoolVarP(&greedy, "greedy", "g", false, "Check all files for URL's not just Javascript")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().BoolVarP(&crawl, "crawl", "c", false, "Enable crawling")
	rootCmd.PersistentFlags().BoolVarP(&findSecrets, "secrets", "s", false, "Check results for secrets e.g api keys")
	rootCmd.PersistentFlags().BoolVarP(&silent, "silent", "S", false, "Enable Silent mode")
	rootCmd.PersistentFlags().StringVar(&sig, "sig", "~/.jsf_signatures.yaml", "Location of signatures for finding secrets")

}

func Execute() error {
	return rootCmd.Execute()
}

//Things to check before running any code.
func safetyChecks() {

	//Check if we can write to the outputFile
	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_WRONLY, 0666)
		if err != nil {
			if os.IsPermission(err) {
				l.Log.Fatal(err)
			}

		}
		file.Close()
	}

	//Ensure you don't provide both url and input file
	if inputFile != "" && url != "" {
		l.Log.Fatal("Provide either -f or -u, you can't provide both")
	}

	//ensure signature file exists
	if findSecrets {
		core.ConfigSigs.ParseConfig(sig) //https://github.com/eth0izzle/shhgit/blob/090e3586ee089f013659e02be16fd0472b629bc7/core/signatures.go
		core.Signatures = core.ConfigSigs.GetSignatures()
		core.Blacklisted_extensions = []string{".exe", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".psd", ".xcf", ".zip", ".tar.gz", ".ttf", ".lock"}
		if silent == true {
			core.PrintSecrets = false
		}
	}

	//if silent && debug {
	//	l.Log.Fatal("Enable silent mode or debug mode. Can't print debug information if silent mode is enabled.")
	//}

	//ensure output is being sent to console or outputfile.
	if silent && outputFile == "" {
		l.Log.Fatal("If you aren't saving the output with -o and you want the display silenced -S, what's the point of running JSubfinder?")
	}
}
