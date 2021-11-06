package core

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"regexp/syntax"

	l "github.com/hiddengearz/jsubfinder/core/logger"
	"gopkg.in/yaml.v2"

	//"strconv"
	"strings"
)

var ConfigSigs ConfigSignature
var Signatures []Signature
var Blacklisted_extensions []string
var PrintSecrets bool = false

type ConfigSignature struct {
	Signatures []struct {
		Part    string `yaml:"part"`
		Match   string `yaml:"match,omitempty"`
		Name    string `yaml:"name"`
		Regex   string `yaml:"regex,omitempty"`
		Comment string `yaml:"comment,omitempty"`
	} `yaml:"signatures"`
}

func (config *ConfigSignature) ParseConfig(fileName string) {
	if fileExists(fileName) {
		//fmt.Println("Parsing " + fileName)

		yamlFile, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Error reading YAML file: %s\n", err)
			return
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
		}

		//fmt.Printf("Result: %v\n", config)

	} else {
		l.Log.Fatal(fileName + "doesn't exist")
		return
	}

}

type Signature interface {
	Name() string
	Match(data *JSData) []string
}

func (config *ConfigSignature) GetSignatures() []Signature {
	var signatures []Signature
	for _, signature := range config.Signatures {
		if signature.Match != "" {
			signatures = append(signatures, &SimpleSignature{
				name:  signature.Name,
				part:  signature.Part,
				match: signature.Match,
			})
		} else {
			if _, err := syntax.Parse(signature.Match, syntax.FoldCase); err == nil {
				signatures = append(signatures, &PatternSignature{
					name:  signature.Name,
					part:  signature.Part,
					match: regexp.MustCompile(signature.Regex),
				})
			}
		}
	}
	//fmt.Println("total signatures: " + strconv.Itoa(len(signatures)))
	return signatures
}

type SimpleSignature struct {
	part  string
	match string
	name  string
}

func (s *SimpleSignature) Name() string {
	return s.name
}
func (s *SimpleSignature) Match(data *JSData) (secrets []string) {
	/*
			switch s.part {
		case PartPath:
			haystack = &file.Path
			matchPart = PartPath
		case PartFilename:
			haystack = &file.Filename
			matchPart = PartPath
		case PartExtension:
			haystack = &file.Extension
			matchPart = PartPath
		default:
			return false, matchPart
		}


	*/

	if s.match != "" {

		if strings.Contains(data.UrlAddr.string, s.match) {
			secrets = append(secrets, s.name+" "+s.match+" found in URL")
			if Debug {
				fmt.Println(s.name + " " + s.match + " found within URL of " + data.UrlAddr.string)
			}
		}
		if strings.Contains(data.Content, s.match) {
			secrets = append(secrets, s.name+" "+s.match+" found in content")
			if Debug {
				fmt.Println(s.name + " " + s.match + " found in content of " + data.UrlAddr.string)
			}
		}

		return secrets

	}

	return secrets

}

type PatternSignature struct {
	part  string
	match *regexp.Regexp
	name  string
}

func (s *PatternSignature) Name() string {
	return s.name
}
func (s *PatternSignature) Match(data *JSData) (secrets []string) {
	if s.match != nil {
		if s.match.MatchString(data.UrlAddr.string) {
			tmp := s.match.FindAllString(data.UrlAddr.string, -1)
			for _, secret := range tmp {
				secrets = append(secrets, s.name+" "+secret+" found in URL")
				if Debug {
					fmt.Println(s.name + " " + secret + " found within URL of " + data.UrlAddr.string)
				}
			}
		}
		if s.match.MatchString(data.Content) {
			tmp := s.match.FindAllString(data.Content, -1)
			for _, secret := range tmp {
				secrets = append(secrets, s.name+" "+secret+" found in content")
				if Debug {
					fmt.Println(s.name + " " + secret + " found within content of " + data.UrlAddr.string)
				}
			}
		}
	}

	return secrets
}
