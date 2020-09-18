package core

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type JSData struct {
	UrlAddr
	Content    string
	subdomains []string
	secrets    []string
}

//GetSubDomains uses regex to find subdomains in the content of JS files
func (js *JSData) GetSubDomains() {
	if Debug {
		defer TimeTrack(time.Now(), "GetSubDomains "+js.UrlAddr.string)
	}
	domainRegex, err := regexp.Compile("([a-zA-Z0-9][a-zA-Z0-9-]*\\." + js.UrlAddr.tld + ")")
	if err != nil {
		log.Fatal(err)
	}
	results := domainRegex.FindAllStringSubmatch(js.Content, -1)
	for _, result := range results {
		_, found := Find(js.subdomains, result[1])
		if !found {
			tmp := result[1]
			tmp = strings.Replace(tmp, "u002F", "", -1)
			if Debug {
				fmt.Println("Replacing " + result[1] + " with " + tmp + " due to likely false positive u002F in url")
			}
			js.subdomains = append(js.subdomains, tmp)
		}
	}

}

//GetSecrets retrieves secret keys from Java script files wip
func GetSecrets() {

}
