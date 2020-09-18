package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	tld "github.com/jpillora/go-tld"
	"github.com/valyala/fasthttp"
)

type UrlData struct {
	UrlAddr
	Content string
	JSFiles []JSData
}

func NewURLData(u string) (data UrlData) {
	if Debug {
		defer TimeTrack(time.Now(), u)
	}
	var err error

	data.UrlAddr.string = u
	client := &fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	err, data.Content = data.UrlAddr.GetContent(client)
	if err != nil && err.Error() == "https" {
		data.UrlAddr.string = "https://" + u
		if Debug {
			fmt.Println("Adding https to " + u)
		}
	} else if err != nil && err.Error() == "http" {
		data.UrlAddr.string = "http://" + u
		if Debug {
			fmt.Println("Adding http to " + u)
		}
	} else if err != nil {
		return
	}

	u2, err2 := tld.Parse(data.UrlAddr.string)
	if err2 != nil {
		log.Fatal(err)
	}
	data.UrlAddr.tld = u2.Domain + "." + u2.TLD

	if Crawl {
		data.JSFiles = data.GetJSLinks()
	}

	if strings.Contains(data.Content, "<script") || strings.Contains(data.Content, "/script") {
		data.JSFiles = append(data.JSFiles, JSData{data.UrlAddr, data.Content, nil, nil}) //Add the base url to JSFiles as there may be inline JS
	} else {
		if Debug {
			fmt.Println("no script tags in: " + data.UrlAddr.string)
		}
	}

	var wg = sync.WaitGroup{}
	maxGoroutines := 5
	guard := make(chan struct{}, maxGoroutines)

	type result struct {
		int
		string
	}

	results := make(chan result, len(data.JSFiles))

	for i, js := range data.JSFiles {
		if js.Content == "" {
			guard <- struct{}{}
			wg.Add(1)
			go func(i int, js JSData) {

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

		data.JSFiles[entry.int].Content = entry.string

	}
	for i, _ = range data.JSFiles {
		if data.JSFiles[i].Content != "" {
			data.JSFiles[i].GetSubDomains()
		}
	}

	if FindSecrets {
		for i, _ = range data.JSFiles {
			if data.JSFiles[i].Content != "" {
				for _, sig := range Signatures {
					data.JSFiles[i].secrets = append(data.JSFiles[i].secrets, sig.Match(&data.JSFiles[i])...)
				}
			}
		}
	}

	return data
}

//GetJSLinks retrieves the links to JS files from the content of the url
func (u *UrlData) GetJSLinks() (JSFile []JSData) {
	if Debug {
		defer TimeTrack(time.Now(), "GetJSLinks "+u.string)
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
		results = jsRegex.FindAllStringSubmatch(u.Content, -1)
	} else {
		results = GreedyRegex.FindAllStringSubmatch(u.Content, -1)
	}

	var found bool = false
	for _, result := range results {

		for _, js := range u.JSFiles {
			if result[1] == js.UrlAddr.string {
				found = true
			}
		}

		if !found {
			if result[1] != "" {
				if strings.HasPrefix(result[1], "http://") || strings.HasPrefix(result[1], "https://") {
					JSFile = append(JSFile, JSData{UrlAddr{result[1], u.UrlAddr.tld}, "", nil, nil})
				} else if strings.HasPrefix(result[1], "//") {
					protocol, err := GetHTTprotocol(u.UrlAddr.string)
					if err == nil {
						link := result[1]
						link = protocol + link[2:]
						JSFile = append(JSFile, JSData{UrlAddr{link, u.UrlAddr.tld}, "", nil, nil})
					}
				} else {
					protocol, err := GetHTTprotocol(u.UrlAddr.string)
					if err == nil {
						link := strings.Replace(u.UrlAddr.string, protocol, "", 1) + "/" + result[1]
						link = protocol + strings.Replace(link, "//", "/", -1)
						JSFile = append(JSFile, JSData{UrlAddr{link, u.UrlAddr.tld}, "", nil, nil})
					}
				}
			}
			//JSFile append
		}
	}

	return JSFile
}
