package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	l "github.com/hiddengearz/jsubfinder/core/logger"
)

type WebPage struct {
	UrlAddr
	Content string //Contant of the webpage
	JSFiles []JavaScript
}

//Get Subdomains and secrets from URL's
func GetResults(url string) (wp WebPage) {
	if Debug {
		defer TimeTrack(time.Now(), url)
	}
	var err error
	var contenTypeJS bool = false

	wp.UrlAddr.string = url //set the URL, needed for GetContent()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: SSL},
		},
	}

	err, wp.Content, contenTypeJS = wp.UrlAddr.GetContent(client) //retrieve the content of the webpage
	if err != nil {
		return
	}

	err = wp.UrlAddr.GetRootDomain() //Set the The root domain, e.g www.google.com > google.com
	if err != nil {
		return
	}

	if Crawl {
		wp.JSFiles, err = wp.GetJSLinks()
		if err != nil {
			l.Log.Error(err)
		}
	}

	//If the base url happens to be a JS via its content-type header, add it to the list
	if contenTypeJS {
		wp.JSFiles = append(wp.JSFiles, JavaScript{wp.UrlAddr, wp.Content, nil, nil}) //Add the base url to JSFiles as there may be inline JS
		fmt.Println("GetResults content type JS")
	}

	if strings.Contains(wp.Content, "<script") || strings.Contains(wp.Content, "/script>") || strings.Contains(wp.Content, "\"script\"") {
		wp.JSFiles = append(wp.JSFiles, JavaScript{wp.UrlAddr, wp.Content, nil, nil}) //Add the base url to JSFiles as there may be inline JS
	} else {
		l.Log.Debug("no script tags in: " + wp.UrlAddr.string)
	}

	var wg = sync.WaitGroup{}
	maxGoroutines := 2
	guard := make(chan struct{}, maxGoroutines)

	type result struct {
		int
		string
	}

	results := make(chan result, len(wp.JSFiles))

	for i, js := range wp.JSFiles {
		if js.Content == "" {
			guard <- struct{}{}
			wg.Add(1)
			go func(i int, js JavaScript) {

				_, tmp, _ := js.GetContent(client)
				tmpResult := result{
					int:    i,
					string: tmp,
				}
				results <- tmpResult
				<-guard
				wg.Done()
			}(i, js)
		}

	}

	wg.Wait()
	close(guard)
	close(results)

	var i int = 0
	for entry := range results {

		wp.JSFiles[entry.int].Content = entry.string

	}
	for i, _ = range wp.JSFiles {
		if wp.JSFiles[i].Content != "" {
			err := wp.JSFiles[i].GetSubDomains()
			if err != nil {
				l.Log.Error(err)
			}
		}
	}

	if FindSecrets {
		for i, _ = range wp.JSFiles {
			if wp.JSFiles[i].Content != "" {
				err := wp.JSFiles[i].GetSecrets()
				if err != nil {
					l.Log.Debug(err)
				}
			}
		}
	}

	return wp
}

//GetJSLinks retrieves the links to JS files from the content of the url
func (wp *WebPage) GetJSLinks() (JSFile []JavaScript, err error) {
	var found bool
	var results [][]string

	if Debug {
		defer TimeTrack(time.Now(), "GetJSLinks "+wp.string)
	}

	//([-a-zA-Z0-9@:%._/\+~#=]{1,256}\.js)\b
	//jsRegex, err := regexp.Compile("src\\s?=\\s?\"(.*.js)\"")
	jsRegex, err := regexp.Compile("\\s?=\\s?\"([-a-z0-9\\/@:%.-_\\+~#=]+.js)\"")
	if err != nil {
		log.Fatal(err)
	}
	GreedyRegex, err := regexp.Compile("\\s?=\\s?\"([-a-z0-9\\/@:%.-_\\+~#=]+)\"") //shold prob just remove this tbh
	if err != nil {
		log.Fatal(err)
	}

	if !Greedy {
		results = jsRegex.FindAllStringSubmatch(wp.Content, -1)
	} else {
		results = GreedyRegex.FindAllStringSubmatch(wp.Content, -1)
	}

	for _, result := range results {
		found = false

		for _, js := range wp.JSFiles {
			if result[1] == js.UrlAddr.string {
				found = true
			}
		}

		if !found {
			var protocol string
			if result[1] != "" {
				if strings.HasPrefix(result[1], "http://") || strings.HasPrefix(result[1], "https://") {
					JSFile = append(JSFile, JavaScript{UrlAddr{result[1], wp.UrlAddr.rootDomain}, "", nil, nil})
				} else if strings.HasPrefix(result[1], "//") {
					protocol, err = GetHTTprotocol(wp.UrlAddr.string)
					if err != nil {
						return
					}
					link := result[1]
					link = protocol + link[2:]
					JSFile = append(JSFile, JavaScript{UrlAddr{link, wp.UrlAddr.rootDomain}, "", nil, nil})
				} else {
					protocol, err = GetHTTprotocol(wp.UrlAddr.string)
					if err != nil {
						return
					}
					link := strings.Replace(wp.UrlAddr.string, protocol, "", 1) + "/" + result[1]
					link = protocol + strings.Replace(link, "//", "/", -1)
					JSFile = append(JSFile, JavaScript{UrlAddr{link, wp.UrlAddr.rootDomain}, "", nil, nil})
				}
			}
			//JSFile append
		}
	}

	return
}
