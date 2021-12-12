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

var SSHFolder string
var Certificate string
var Key string
var X509pair tls.Certificate
var subDomainlogger *log.Logger
var secretsLogger *log.Logger
var Scope []string
var inScope bool

//start the proxy server
func StartProxy(port string, upsteamProxySet bool) (err error) {
	proxy := goproxy.NewProxyHttpServer()

	//if upstream proxy set, proxy all requests
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
	/*
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
	*/

	//always intercepthttp requests
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		var result JavaScript

		//if there is a scope, check it
		if Scope != nil {
			inScope = false
			for _, host := range Scope {
				if host+":"+r.Request.URL.Port() == r.Request.URL.Host {
					inScope = true
					break
				}

			}
		}

		//if provided url isnt in scope, then return
		if !inScope {
			//l.Log.Debug(r.Request.URL.String() + " not in scope")
			return r
		}

		result.subdomains = make(map[string]bool)
		result.secrets = make(map[string]bool)
		result.UrlAddr.string = r.Request.URL.String()

		//read request
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			l.Log.Debug(errors.New("Failed to read body of " + result.UrlAddr.string))
			return nil
		}

		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(string(bodyBytes))))

		result.Content = string(bodyBytes)
		if result.Content == "" { //if no content, then there is no JS, return
			return r
		}

		contenType := r.Header.Get("Content-Type")

		//if the header or page contains javascript
		if strings.Contains(contenType, "javascript") || strings.Contains(result.Content, "<script") ||
			strings.Contains(result.Content, "/script>") || strings.Contains(result.Content, "\"script\"") {
			if IsUrlVisited(result.UrlAddr.string) { //if the url has been visited return
				l.Log.Debug(result.UrlAddr.string + " has already been visited")
				return r
			}
			go func() { //process page
				ParseProxyResponse(result)
				//time.Sleep(2 * time.Second)
			}()
			AddUrlVisited(result.UrlAddr.string)
		}
		return r
	})

	//if outputFile is set, setup the output files
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

//Process requests, print them to console and to file
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
	for subdomain, _ := range js.subdomains {
		if IsNewSubdomain(subdomain) {
			if !Silent {
				fmt.Println("Subdomain: " + subdomain)
			}
			AddNewSubdomain(subdomain)
			if OutputFile != "" {
				subDomainlogger.Output(2, subdomain)
			}
		}
	}
	for secret, _ := range js.secrets {
		if IsNewSecret(secret) {
			if PrintSecrets {
				fmt.Println(secret + " of " + js.UrlAddr.string)
			}
			AddNewSecret(secret)
			if OutputFile != "" {
				secretsLogger.Output(2, secret)
			}
		}
	}

}
