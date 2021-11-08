package core

import (
	"crypto/tls"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	l "github.com/hiddengearz/jsubfinder/core/logger"
	tld "github.com/jpillora/go-tld"
	"github.com/valyala/fasthttp"
)

type WebPage struct {
	UrlAddr
	Content string //Contant of the webpage
	JSFiles []JavaScript
}

func GetResults(url string) (wp WebPage) {
	if Debug {
		defer TimeTrack(time.Now(), url)
	}
	var err error

	wp.UrlAddr.string = url
	client := &fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: SSL},
	}
	err, wp.Content = wp.UrlAddr.GetContent(client)
	if err != nil && err.Error() == "https" {
		wp.UrlAddr.string = "https://" + url
		l.Log.Debug("Adding https to " + url)
	} else if err != nil && err.Error() == "http" {
		wp.UrlAddr.string = "http://" + url
		l.Log.Debug("Adding http to " + url)
	} else if err != nil {
		return
	}

	u2, err2 := tld.Parse(wp.UrlAddr.string)
	if err2 != nil {
		log.Fatal(err)
	}
	wp.UrlAddr.tld = u2.Domain + "." + u2.TLD

	if Crawl {
		wp.JSFiles, err = wp.GetJSLinks()
		if err != nil {
			l.Log.Error(err)
		}
	}

	if strings.Contains(wp.Content, "<script") || strings.Contains(wp.Content, "/script") {
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

				_, tmp := js.GetContent(client)
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
				for _, sig := range Signatures {
					wp.JSFiles[i].secrets = append(wp.JSFiles[i].secrets, sig.Match(&wp.JSFiles[i])...)
				}
			}
		}
	}

	return wp
}

//GetJSLinks retrieves the links to JS files from the content of the url
func (wp *WebPage) GetJSLinks() (JSFile []JavaScript, err error) {
	if Debug {
		defer TimeTrack(time.Now(), "GetJSLinks "+wp.string)
	}
	var results [][]string
	//([-a-zA-Z0-9@:%._/\+~#=]{1,256}\.js)\b
	//jsRegex, err := regexp.Compile("src\\s?=\\s?\"(.*.js)\"")
	jsRegex, err := regexp.Compile("\\s?=\\s?\"([-a-z0-9\\/@:%._\\+~#=]+.js)\"")
	if err != nil {
		log.Fatal(err)
	}
	GreedyRegex, err := regexp.Compile("\\s?=\\s?\"([-a-z0-9\\/@:%._\\+~#=]+)\"")
	if err != nil {
		log.Fatal(err)
	}

	if !Greedy {
		results = jsRegex.FindAllStringSubmatch(wp.Content, -1)
	} else {
		results = GreedyRegex.FindAllStringSubmatch(wp.Content, -1)
	}

	var found bool = false
	for _, result := range results {

		for _, js := range wp.JSFiles {
			if result[1] == js.UrlAddr.string {
				found = true
			}
		}

		if !found {
			var protocol string
			if result[1] != "" {
				if strings.HasPrefix(result[1], "http://") || strings.HasPrefix(result[1], "https://") {
					JSFile = append(JSFile, JavaScript{UrlAddr{result[1], wp.UrlAddr.tld}, "", nil, nil})
				} else if strings.HasPrefix(result[1], "//") {
					protocol, err = GetHTTprotocol(wp.UrlAddr.string)
					if err != nil {
						return
					}
					link := result[1]
					link = protocol + link[2:]
					JSFile = append(JSFile, JavaScript{UrlAddr{link, wp.UrlAddr.tld}, "", nil, nil})
				} else {
					protocol, err = GetHTTprotocol(wp.UrlAddr.string)
					if err != nil {
						return
					}
					link := strings.Replace(wp.UrlAddr.string, protocol, "", 1) + "/" + result[1]
					link = protocol + strings.Replace(link, "//", "/", -1)
					JSFile = append(JSFile, JavaScript{UrlAddr{link, wp.UrlAddr.tld}, "", nil, nil})
				}
			}
			//JSFile append
		}
	}

	return
}
