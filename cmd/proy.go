package cmd

import (
	"errors"
	"strconv"
	"strings"

	C "github.com/hiddengearz/jsubfinder/core"
	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/spf13/cobra"
)

func init() {
	//rootCmd.AddCommand(cmdExec)
	proxyExec.PersistentFlags().IntVarP(&C.LocalPort, "port", "p", 8444, "Port for the proxy to listen on")
	proxyExec.PersistentFlags().StringVarP(&C.UpsteamProxy, "upstream-proxy", "u", "http://127.0.0.1:8888", "Adress of upsteam proxy e.g http://127.0.0.1:8888")

}

//Start JSubFiner in proxy mode
var proxyExec = &cobra.Command{
	Use:   "proxy",
	Short: "Run JSubfinder as a proxy",
	Long:  `Execute the command specified`,
	Run: func(cmd *cobra.Command, arguments []string) {
		var upsteamProxySet bool = false
		if cmd.Flags().Changed("upstream-proxy") {
			upsteamProxySet = true
		}
		C.StartProxy(":"+strconv.Itoa(C.LocalPort), upsteamProxySet)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := safetyChecks()
		if err != nil {
			l.Log.Fatal(err)
		}
		if !strings.Contains(C.UpsteamProxy, "http://") {
			l.Log.Fatal(errors.New("Upsteam Proxy doesn't contain http://"))
		}
	},
}
