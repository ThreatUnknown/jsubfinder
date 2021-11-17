package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	C "github.com/hiddengearz/jsubfinder/core"
	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/mitchellh/go-homedir"
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
		if !strings.HasPrefix(C.UpsteamProxy, "http://") {
			l.Log.Fatal(errors.New("Upsteam Proxy doesn't contain http://"))
		}
		//check if valid url...

		home, err := homedir.Dir()
		if err != nil {
			l.Log.Fatal(err)
		}
		C.SSHFolder = home + "/.ssh/"
		if !C.FolderExists(C.SSHFolder) {
			l.Log.Fatal("Folder " + C.SSHFolder + " doesnt exist. Please create it")
		}

		C.Certificate = C.SSHFolder + "jsubfinder.pub"
		C.Key = C.SSHFolder + "jsubfinder"

		if !C.FileExists(C.Certificate) || !C.FileExists(C.Key) {

			fmt.Println("creating cert!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			err = C.CreateAuthority(C.Certificate, C.Key)
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}
