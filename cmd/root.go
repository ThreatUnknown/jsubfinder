package cmd

import (
	"os"

	C "github.com/hiddengearz/jsubfinder/core"
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

			if C.Debug && C.Silent {
				l.Log.Fatal("Please choose Debug mode or silent mode. Not both.")
			} else if C.Debug {

				l.InitDetailedLogger(f)
				l.Log.SetLevel(logrus.DebugLevel)

				l.Log.Debug("Debug mode enabled")
			} else if C.Silent {
				l.Log.SetLevel(logrus.PanicLevel)
			}

		},
	}
)

func init() {
	l.Log = logrus.New()
	rootCmd.AddCommand(searchExec)
	rootCmd.AddCommand(proxyExec)
	rootCmd.PersistentFlags().IntVarP(&C.Threads, "threads", "t", 5, "Ammount of threads to be used")
	rootCmd.PersistentFlags().StringVarP(&C.OutputFile, "outputFile", "o", "", "name/location to store the file")
	rootCmd.PersistentFlags().BoolVarP(&C.Greedy, "greedy", "g", false, "Check all files for URL's not just Javascript")
	rootCmd.PersistentFlags().BoolVarP(&C.Debug, "debug", "d", false, "Enable debug mode. Logs are stored in log.info")
	rootCmd.PersistentFlags().BoolVarP(&C.Crawl, "crawl", "c", false, "Enable crawling")
	rootCmd.PersistentFlags().BoolVarP(&C.FindSecrets, "secrets", "s", false, "Check results for secrets e.g api keys")
	rootCmd.PersistentFlags().BoolVarP(&C.Silent, "silent", "S", false, "Enable Silent mode")
	rootCmd.PersistentFlags().StringVar(&C.Sig, "sig", "~/.jsf_signatures.yaml", "Location of signatures for finding secrets")

}

func Execute() error {
	return rootCmd.Execute()
}

//Things to check before running any code.
func safetyChecks() {

	//Check if we can write to the outputFile
	if C.OutputFile != "" {
		file, err := os.OpenFile(C.OutputFile, os.O_WRONLY, 0666)
		if err != nil {
			if os.IsPermission(err) {
				l.Log.Fatal(err)
			}

		}
		file.Close()
	}

	//Ensure you don't provide both url and input file
	if C.InputFile != "" && C.Url != "" {
		l.Log.Fatal("Provide either -f or -u, you can't provide both")
	}

	//ensure signature file exists
	if C.FindSecrets {
		C.ConfigSigs.ParseConfig(C.Sig) //https://github.com/eth0izzle/shhgit/blob/090e3586ee089f013659e02be16fd0472b629bc7/core/signatures.go
		C.Signatures = C.ConfigSigs.GetSignatures()
		C.Blacklisted_extensions = []string{".exe", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".psd", ".xcf", ".zip", ".tar.gz", ".ttf", ".lock"}
		if C.Silent == true {
			C.PrintSecrets = false
		}
	}

	//if silent && debug {
	//	l.Log.Fatal("Enable silent mode or debug mode. Can't print debug information if silent mode is enabled.")
	//}

	//ensure output is being sent to console or outputfile.
	if C.Silent && C.OutputFile == "" {
		l.Log.Fatal("If you aren't saving the output with -o and you want the display silenced -S, what's the point of running JSubfinder?")
	}
}
