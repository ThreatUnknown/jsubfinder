package core

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
	l "github.com/hiddengearz/jsubfinder/core/logger"
)

var newSubdomains []string
var newSecrets []string
var SSHFolder string
var Certificate string
var Key string
var X509pair tls.Certificate
var subDomainlogger *log.Logger
var secretsLogger *log.Logger

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
	} else {
		proxy.Logger = log.New(ioutil.Discard, "", 0)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		//GetCertificate:     returnCert,
	}

	X509pair, err = tls.LoadX509KeyPair(Certificate, Key)
	if err != nil {
		log.Fatalf("Unable to load certificate %s: %v", Certificate, err)
		return errors.New(fmt.Sprintf("Unable to load certificate %s: %v", Certificate, err))
	}
	tlsConfig.Certificates = append(tlsConfig.Certificates, X509pair)

	// Not strictly required but appears to help with SNI
	tlsConfig.BuildNameToCertificate()

	goproxy.MitmConnect.TLSConfig = func(host string, ctx *goproxy.ProxyCtx) (*tls.Config, error) {
		return tlsConfig, nil
	}

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

		//fmt.Println("received request to", r.Request.URL.String())

		var result JavaScript

		//if !Greedy && !strings.HasSuffix(r.Request.URL.String(), ".js") {
		//	return r
		//}

		result.UrlAddr.string = r.Request.URL.String()

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			l.Log.Debug(errors.New("Failed to read body of " + result.UrlAddr.string))
			return nil
		}

		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(string(bodyBytes))))
		//fmt.Println(r.Body)

		//r.Body.Close()

		//fmt.Println(string(bodyBytes))

		result.Content = string(bodyBytes)
		if result.Content == "" {
			return r
		}

		contenType := r.Header.Get("Content-Type")

		if strings.Contains(contenType, "javascript") || strings.Contains(result.Content, "<script") ||
			strings.Contains(result.Content, "/script>") || strings.Contains(result.Content, "\"script\"") {
			//fmt.Println("bingo " + r.Request.URL.Hostname() + r.Request.URL.Path)
			go func() {
				ParseProxyResponse(result)
				//time.Sleep(2 * time.Second)
			}()
		}
		return r
	})

	if OutputFile != "" {
		f, err := os.OpenFile(OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		subDomainlogger = log.New(f, "", 0)

		if FindSecrets {
			f, err := os.OpenFile("secrets_"+OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			secretsLogger = log.New(f, "", 0)
		}
	}

	fmt.Println("Proxy started on", port)
	http.ListenAndServe(port, proxy)

	fmt.Println("Proxy stopped")
	return nil
}

func ParseProxyResponse(js JavaScript) {
	err := js.UrlAddr.GetRootDomain()
	if err != nil {
		l.Log.Debug(err)
		return
	}

	err = js.GetSubDomains()
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
			if !Silent {
				fmt.Println("Subdomain: " + subdomain)
			}
			if OutputFile != "" {
				subDomainlogger.Output(2, subdomain)
			}
			newSubdomains = append(newSubdomains, subdomain)

		}
	}
	for _, secret := range js.secrets {
		_, found := Find(newSecrets, secret)
		if !found {
			if PrintSecrets {
				fmt.Println(secret + " of " + js.UrlAddr.string)
			}
			if OutputFile != "" {
				secretsLogger.Output(2, secret)
			}
			newSecrets = append(newSecrets, secret+" of "+js.UrlAddr.string)
		}
	}

}
