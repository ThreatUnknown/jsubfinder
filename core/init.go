package core

import (
	"fmt"
	"strconv"
	"sync"
)

func Exec(urls []string, concurrency int, outputFile string) {
	var data []UrlData
	var wg = sync.WaitGroup{}
	maxGoroutines := concurrency
	guard := make(chan struct{}, maxGoroutines)

	results := make(chan UrlData, len(urls))
	for _, url := range urls {
		guard <- struct{}{}
		wg.Add(1)
		go func(url string) {

			results <- NewURLData(url)
			<-guard
			wg.Done()
		}(url)
	}

	wg.Wait()
	close(guard)
	close(results)

	for result := range results {
		if result.Content != "" { //the urladdr will be blank if the page can't be reached. Thus don't add it.
			data = append(data, result)
		}
	}

	saveresults := IsFlagPassed("o")
	var newSubdomains []string
	var newSecrets []string
	if Debug {
		for _, url := range data {
			fmt.Println("url: " + url.UrlAddr.string)
			fmt.Println("\ttld: " + url.UrlAddr.tld)
			for _, js := range url.JSFiles {
				fmt.Println("\tjs: " + js.UrlAddr.string)
				fmt.Println("\t\tcontent length: " + strconv.Itoa(len(js.Content)))
				for _, subdomain := range js.subdomains {
					fmt.Println("\t\tsubdomain: " + subdomain)
					_, found := Find(newSubdomains, subdomain)
					if !found {
						newSubdomains = append(newSubdomains, subdomain)
					}
				}
				for _, secret := range js.secrets {
					fmt.Println("\t\tsecret: " + secret)
					_, found := Find(newSecrets, secret)
					if !found {
						newSecrets = append(newSecrets, secret)
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
						fmt.Println(subdomain)
						newSubdomains = append(newSubdomains, subdomain)
					}
				}
				for _, secret := range js.secrets {
					_, found := Find(newSecrets, secret)
					if !found {
						newSecrets = append(newSecrets, secret)
					}
				}
			}
		}
	}

	if saveresults {
		SaveResults(outputFile, newSubdomains)
	}
	if saveresults {
		SaveResults("secrets_"+outputFile, newSecrets)
	}
}
