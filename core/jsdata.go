package core

import (
	"regexp"
	"strings"
	"time"

	l "github.com/hiddengearz/jsubfinder/core/logger"
)

type JavaScript struct {
	UrlAddr
	Content    string   //Content of the JS file
	subdomains []string //Subdomains found in content of JavaScript
	secrets    []string //Secrets found content of JavaScript
}

//GetSubDomains uses regex to find subdomains in the content of JS files
func (js *JavaScript) GetSubDomains() error {
	if Debug {
		defer TimeTrack(time.Now(), "GetSubDomains "+js.UrlAddr.string)
	}
	domainRegex, err := regexp.Compile("([a-zA-Z0-9][a-zA-Z0-9-]*\\." + js.UrlAddr.tld + ")")
	if err != nil {
		l.Log.Debug(err)
		return err
	}
	results := domainRegex.FindAllStringSubmatch(js.Content, -1)
	for _, result := range results {
		_, found := Find(js.subdomains, result[1])
		if !found {
			tmp := result[1]
			tmp = strings.Replace(tmp, "u002F", "", -1)

			//l.Log.Debug("Replacing " + result[1] + " with " + tmp + " due to likely false positive u002F in url")
			js.subdomains = append(js.subdomains, tmp)
		}
	}
	return nil

}

func (js *JavaScript) GetSecrets() error {
	if Debug {
		defer TimeTrack(time.Now(), "GetSecrets "+js.UrlAddr.string)
	}
	for _, sig := range Signatures {
		js.secrets = append(js.secrets, sig.Match(js)...)
	}
	return nil
}
