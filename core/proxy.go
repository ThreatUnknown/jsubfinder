package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/elazarl/goproxy"
	l "github.com/hiddengearz/jsubfinder/core/logger"
)

var newSubdomains []string
var newSecrets []string

func StartProxy(port string, upsteamProxySet bool) (err error) {
	proxy := goproxy.NewProxyHttpServer()

	if upsteamProxySet {
		proxy.Tr = &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(UpsteamProxy)
		}}
		proxy.ConnectDial = proxy.NewConnectDialToProxy(UpsteamProxy)
	}
	if Debug {
		proxy.Verbose = true
	}
	GenerateCert()
	if err != nil {
		l.Log.Fatal(err)
	}

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	//goproxy.GoproxyCa
	//goproxy.GoproxyCa.PrivateKey

	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		fmt.Printf("received request to", r.Request.URL.String())

		var result JavaScript
		result.UrlAddr.string = r.Request.URL.String()

		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			l.Log.Debug(errors.New("Failed to read body of " + result.UrlAddr.string))
			return r
		}

		result.Content = string(bodyBytes)
		go func() {
			ParseProxyResponse(result)
			time.Sleep(2 * time.Second)
		}()
		return r
	})

	fmt.Println("Proxy started on", port)
	http.ListenAndServe(port, proxy)

	fmt.Println("Proxy stopped")
	return nil
}

func ParseProxyResponse(js JavaScript) {
	err := js.GetSubDomains()
	if err != nil {
		l.Log.Debug(err)
		return
	}
	if FindSecrets {
		err := js.GetSecrets()
		if err != nil {
			l.Log.Debug(err)
			return
		}
	}

	for _, subdomain := range js.subdomains {
		_, found := Find(newSubdomains, subdomain)
		if !found {
			fmt.Println(subdomain)
			newSubdomains = append(newSubdomains, subdomain)
		}
	}
	for _, secret := range js.secrets {
		_, found := Find(newSecrets, secret)
		if !found {
			if PrintSecrets {
				fmt.Println(secret + " of " + js.UrlAddr.string)
			}
			newSecrets = append(newSecrets, secret+" of "+js.UrlAddr.string)
		}
	}

}
