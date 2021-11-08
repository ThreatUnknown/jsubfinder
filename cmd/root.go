package cmd

import (
	"errors"
	"os"

	C "github.com/hiddengearz/jsubfinder/core"
	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "JSubFinder",
		Short: "Agent for the [redacted] project",
		Long:  `[redacted]`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

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
	rootCmd.PersistentFlags().BoolVarP(&C.SSL, "nossl", "K", true, "Skip SSL cert verification")

}

func Execute() error {
	return rootCmd.Execute()
}

//Things to check before running any code.
func safetyChecks() error {
	f, err := os.OpenFile("log.info", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return (err)
	}

	if C.Debug && C.Silent {
		return errors.New("Please choose Debug mode or silent mode. Not both.")
	} else if C.Debug { //Setup logging

		l.InitDetailedLogger(f)
		l.Log.SetLevel(logrus.DebugLevel)

		l.Log.Debug("Debug mode enabled")
	} else if C.Silent {
		l.Log.SetLevel(logrus.PanicLevel)
	}

	//Check if we can write to the outputFile
	if C.OutputFile != "" {
		file, err := os.OpenFile(C.OutputFile, os.O_WRONLY, 0666)
		if err != nil {
			if os.IsPermission(err) {
				return (err)
			}

		}
		file.Close()
	}

	//Ensure you don't provide both url and input file
	if C.InputFile != "" && len(C.InputURLs) != 0 {
		return errors.New("Provide either -f or -u, you can't provide both")
	}

	//ensure signature file exists
	if C.FindSecrets {
		err := C.ConfigSigs.ParseConfig(C.Sig) //https://github.com/eth0izzle/shhgit/blob/090e3586ee089f013659e02be16fd0472b629bc7/core/signatures.go
		if err != nil {
			return err
		}
		C.Signatures, err = C.ConfigSigs.GetSignatures()
		if err != nil {
			return err
		}
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
		return errors.New("If you aren't saving the output with -o and you want the display silenced -S, what's the point of running JSubfinder?")
	}

	return nil
}
