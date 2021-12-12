package core

import (
	"regexp"
	"strings"
	"time"

	l "github.com/hiddengearz/jsubfinder/core/logger"
)

type JavaScript struct {
	UrlAddr
	Content    string          //Content of the JS file
	subdomains map[string]bool //Subdomains found in content of JavaScript
	secrets    map[string]bool //Secrets found content of JavaScript
}

//GetSubDomains uses regex to find subdomains in the content of JS files
func (js *JavaScript) GetSubDomains() error {
	if Debug {
		defer TimeTrack(time.Now(), "GetSubDomains "+js.UrlAddr.string)
	}

	//Regex for finding a subdomain
	subdomainRegex, err := regexp.Compile("([a-zA-Z0-9][a-zA-Z0-9-]*\\." + js.UrlAddr.rootDomain + ")")
	if err != nil {
		l.Log.Debug(err)
		return err
	}

	results := subdomainRegex.FindAllStringSubmatch(js.Content, -1)
	for _, result := range results { //for all found subdomains...
		tmp := result[1]
		tmp = strings.Replace(tmp, "u002F", "", -1)
		tmp = strings.Replace(tmp, "x2F", "", -1)
		js.subdomains[tmp] = true
		if Command != "proxy" {
			if IsNewSubdomain(tmp) {
				AddNewSubdomain(tmp)
			}
		}
	}
	return nil

}

//GetSecrets uses regex to find secrets in the content of JS files
func (js *JavaScript) GetSecrets() error {
	//js.secrets = make(map[string]bool)
	if Debug {
		defer TimeTrack(time.Now(), "GetSecrets "+js.UrlAddr.string)
	}
	for _, sig := range Signatures {
		for _, entry := range sig.Match(js) {
			js.secrets[entry] = true
			if Command != "proxy" {
				if IsNewSecret(entry) {
					AddNewSecret(entry)
				}
			}
		}

	}

	return nil
}
