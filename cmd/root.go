package cmd

import "github.com/spf13/cobra"

var (
	rootCmd     = &cobra.Command{Use: "app"}
	threads     int
	inputFile   string
	url         string
	outputFile  string
	greedy      bool
	debug       bool
	crawl       bool
	findSecrets bool
	sig         string
)

func Execute() error {
	return rootCmd.Execute()
}
