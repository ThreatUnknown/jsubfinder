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
		defer TimeTrack(time.Now(), "GetResults "+url)
	}
	var err error
	var contenTypeJS bool = false

	wp.UrlAddr.string = url //set the URL, needed for GetContent()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: SSL},
		},
	}

	wp.Content, contenTypeJS, err = wp.UrlAddr.GetContent(client) //retrieve the content of the webpage
	if err != nil {
		l.Log.Debug(err)
		return
	}
	if contenTypeJS {
		AddUrlVisited(wp.UrlAddr.string)
	}

	err = wp.UrlAddr.GetRootDomain() //Set the The root domain, e.g www.google.com > google.com
	if err != nil {
		l.Log.Debug(err)
		return
	}

	if Crawl {
		wp.JSFiles, err = wp.GetJSLinks()
		if err != nil {
			l.Log.Debug(err)
		}
	}

	//If the base url happens to be a JS via its content-type header, add it to the list
	if contenTypeJS {
		wp.JSFiles = append(wp.JSFiles, JavaScript{wp.UrlAddr, wp.Content, make(map[string]bool), make(map[string]bool)}) //Add the base url to JSFiles as there may be inline JS
		fmt.Println("GetResults content type JS")
	}

	if strings.Contains(wp.Content, "<script") || strings.Contains(wp.Content, "/script>") || strings.Contains(wp.Content, "\"script\"") {
		wp.JSFiles = append(wp.JSFiles, JavaScript{wp.UrlAddr, wp.Content, make(map[string]bool), make(map[string]bool)}) //Add the base url to JSFiles as there may be inline JS
	} else {
		l.Log.Debug("no script tags in: " + wp.UrlAddr.string)
	}

	//setup go routings
	var wg = sync.WaitGroup{}
	maxGoroutines := 2
	guard := make(chan struct{}, maxGoroutines)

	type result struct {
		int
		string
	}

	results := make(chan result, len(wp.JSFiles))

	//for each JSFile, get content
	for i, js := range wp.JSFiles {
		if js.Content == "" {
			guard <- struct{}{}
			wg.Add(1)
			go func(i int, js JavaScript) {

				tmp, contenTypeJS, err := js.GetContent(client)
				if err != nil {
					l.Log.Debug(err)
				}
				if contenTypeJS { //if the page was a JS file or has js, add it to urls visited
					AddUrlVisited(wp.UrlAddr.string)
				}
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
	for i = range wp.JSFiles {
		if wp.JSFiles[i].Content != "" {
			err := wp.JSFiles[i].GetSubDomains()
			if err != nil {
				l.Log.Error(err)
			}
		}
	}

	if FindSecrets {
		for i = range wp.JSFiles {
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

		var protocol string
		if result[1] != "" {
			if strings.HasPrefix(result[1], "http://") || strings.HasPrefix(result[1], "https://") {
				if IsUrlVisited(result[1]) {
					continue
				}
				JSFile = append(JSFile, JavaScript{UrlAddr{result[1], wp.UrlAddr.rootDomain}, "", make(map[string]bool), make(map[string]bool)})
			} else if strings.HasPrefix(result[1], "//") {
				protocol, err = GetHTTprotocol(wp.UrlAddr.string) //assumption that the JS file will be hosted on same protocol as web server
				if err != nil {
					return
				}
				link := result[1]
				link = protocol + link[2:]
				if IsUrlVisited(link) {
					continue
				}
				JSFile = append(JSFile, JavaScript{UrlAddr{link, wp.UrlAddr.rootDomain}, "", make(map[string]bool), make(map[string]bool)})
			} else {
				protocol, err = GetHTTprotocol(wp.UrlAddr.string) //assumption that the JS file will be hosted on same protocol as web server
				if err != nil {
					return
				}
				link := strings.Replace(wp.UrlAddr.string, protocol, "", 1) + "/" + result[1]
				link = protocol + strings.Replace(link, "//", "/", -1)
				if IsUrlVisited(link) {
					continue
				}
				JSFile = append(JSFile, JavaScript{UrlAddr{link, wp.UrlAddr.rootDomain}, "", make(map[string]bool), make(map[string]bool)})
			}
			//JSFile append
		}
	}
	//time.Sleep(300)
	return
}
