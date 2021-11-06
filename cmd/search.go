package cmd

func init() {
	//rootCmd.AddCommand(cmdExec)
	cmdExec.PersistentFlags().StringVarP(&inputFile, "inputFile", "f", "", "File containing domains")
	cmdExec.PersistentFlags().StringVarP(&url, "url", "u", "", "url to check")

}

func safetyChecks() (bool, error) {
	//Can i write to outputFile
	//if sig enables, does signature file exist
	//maybe is url valid?
	return false, nil
}
