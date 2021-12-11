package core

import (
	"fmt"
	"strconv"
	"sync"

	l "github.com/hiddengearz/jsubfinder/core/logger"
)

var (
	InputURLs    []string
	Threads      int
	InputFile    string
	OutputFile   string
	Greedy       bool
	Debug        bool = false
	Crawl        bool
	FindSecrets  bool
	Sig          string
	Silent       bool
	SSL          bool = false
	LocalPort    int
	UpsteamProxy string
)

func ExecSearch(concurrency int, outputFile string) error {

	//fmt.Print(Urls)
	var data []WebPage
	var wg = sync.WaitGroup{}
	var newSubdomains []string
	var newSecrets []string
	maxGoroutines := concurrency
	guard := make(chan struct{}, maxGoroutines)

	//Start a go routine and start fetching results for each URL provided
	results := make(chan WebPage, len(InputURLs))
	for _, url := range InputURLs {
		guard <- struct{}{}
		wg.Add(1)
		go func(url string) {

			results <- GetResults(url) //fetch results and return them to a channel
			<-guard
			wg.Done()
		}(url)
	}

	wg.Wait()
	close(guard)
	close(results)

	//Take results from the channel and add them to []webpage
	for result := range results {
		if result.Content != "" { //the urladdr will be blank if the page can't be reached. Thus don't add it.
			data = append(data, result)
		}
	}

	//If Debug mode, print results
	if Debug {
		for _, url := range data { //For each URL the user provided
			fmt.Println("url: " + url.UrlAddr.string)              //print the url
			fmt.Println("\trootDomain: " + url.UrlAddr.rootDomain) //print the root domain
			for _, js := range url.JSFiles {                       //For each URL with JS
				fmt.Println("\tjs: " + js.UrlAddr.string)                           //Print the URL
				fmt.Println("\t\tcontent length: " + strconv.Itoa(len(js.Content))) // Print the content length
				for _, subdomain := range js.subdomains {                           //print the subdomain found in the js
					fmt.Println("\t\tsubdomain: " + subdomain)
					_, found := Find(newSubdomains, subdomain)
					if !found { //add the subdomain to the list of new subdomains if not in the list
						newSubdomains = append(newSubdomains, subdomain)
					}
				}
				for _, secret := range js.secrets {
					fmt.Println("\t\tsecret: " + secret)
					_, found := Find(newSecrets, secret)
					if !found {
						newSecrets = append(newSecrets, secret+" of "+js.UrlAddr.string)
					}
				}
			}
		}
	} else {
		for _, url := range data {
			for _, js := range url.JSFiles {
				for _, subdomain := range js.subdomains {
					_, found := Find(newSubdomains, subdomain)
					if !found {
						if !Silent {
							fmt.Println(subdomain)
						}
						newSubdomains = append(newSubdomains, subdomain)
					}
				}
				for _, secret := range js.secrets {
					_, found := Find(newSecrets, secret)
					if !found {

						newSecrets = append(newSecrets, secret+" of "+js.UrlAddr.string)
					}
				}
			}
		}
	}

	if PrintSecrets {
		for _, secret := range newSecrets {
			fmt.Println(secret)
		}
	}

	if outputFile != "" {
		err := SaveResults(outputFile, newSubdomains)
		if err != nil {
			l.Log.Error(err)
		}
		if FindSecrets {
			SaveResults("secrets_"+outputFile, newSecrets)
			if err != nil {
				l.Log.Error(err)
			}
		}
	}
	return nil
}
