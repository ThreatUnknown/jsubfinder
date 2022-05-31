package core

import (
	"fmt"
	"strconv"
	"sync"

	l "github.com/ThreatUnkown/jsubfinder/core/logger"
)

var (
	InputURLs         []string
	Threads           int
	InputFile         string
	OutputFile        string
	SecretsOutputFile string
	Greedy            bool
	Debug             bool = false
	Crawl             bool
	FindSecrets       bool = false
	Sig               string
	Silent            bool
	SSL               bool = false
	LocalPort         int
	UpsteamProxy      string

	allPages []WebPage

	newSubdomains map[string]bool = make(map[string]bool)
	newSecrets    map[string]bool = make(map[string]bool)
	urlsVisited   map[string]bool = make(map[string]bool)

	lock sync.RWMutex

	Command string
)

func ExecSearch() error {

	//setup go routine
	var allPages []WebPage
	var wg = sync.WaitGroup{}

	guard := make(chan struct{}, Threads)

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
			allPages = append(allPages, result)
		}
	}

	//If Debug mode, print results in debug format
	if Debug {
		for _, url := range allPages { //For each URL the user provided
			fmt.Println("url: " + url.UrlAddr.string)              //print the url
			fmt.Println("\trootDomain: " + url.UrlAddr.rootDomain) //print the root domain
			for _, js := range url.JSFiles {                       //For each URL with JS
				fmt.Println("\tjs: " + js.UrlAddr.string)                           //Print the URL
				fmt.Println("\t\tcontent length: " + strconv.Itoa(len(js.Content))) // Print the content length
				for subdomain := range js.subdomains {                              //print the subdomain found in the js
					fmt.Println("\t\tsubdomain: " + subdomain)

				}
				for secret := range js.secrets {
					fmt.Println("\t\tsecret: " + secret)

				}
			}
		}
	} else { //if not debug mode print subdomains & secrets
		var tmp []string
		for subdomain := range newSubdomains {
			if !Silent {
				fmt.Println(subdomain)
			}
			if OutputFile != "" {
				tmp = append(tmp, subdomain)
			}

		}
		if OutputFile != "" {
			err := SaveResults(OutputFile, tmp)
			if err != nil {
				l.Log.Error(err)
			}
		}
		if FindSecrets {
			var tmp []string
			for secret := range newSecrets {
				if !Silent {
					fmt.Println(secret)
				}
				if SecretsOutputFile != "" {
					tmp = append(tmp, secret)
				}
			}
			if SecretsOutputFile != "" {
				err := SaveResults(SecretsOutputFile, tmp)
				if err != nil {
					l.Log.Error(err)
				}
			}
		}

	}

	return nil
}

//Add string to urlVisited, thread safe
func AddUrlVisited(url string) {
	lock.Lock()
	urlsVisited[url] = true
	lock.Unlock()
}

//is string in urlVisited, thread safe
func IsUrlVisited(url string) bool {
	lock.RLock()
	if urlsVisited[url] {
		lock.RUnlock()
		return true
	}
	lock.RUnlock()
	return false
}

//add string to newSubdomains, thread safe
func AddNewSubdomain(url string) {
	lock.Lock()
	newSubdomains[url] = true
	lock.Unlock()
}

//is string in newSubdomains, thread safe
func IsNewSubdomain(url string) bool {
	lock.RLock()
	if !newSubdomains[url] {
		lock.RUnlock()
		return true
	}
	lock.RUnlock()
	return false
}

//add string to newSecrets
func AddNewSecret(url string) {
	lock.Lock()
	newSecrets[url] = true
	lock.Unlock()
}

//is string in newSecrets
func IsNewSecret(url string) bool {
	lock.RLock()
	if !newSecrets[url] {
		lock.RUnlock()
		return true
	}
	lock.RUnlock()
	return false
}
