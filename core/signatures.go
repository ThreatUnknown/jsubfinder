package core

import (
	"errors"
	"io/ioutil"
	"regexp"
	"regexp/syntax"

	l "github.com/ThreatUnkown/jsubfinder/core/logger"
	"gopkg.in/yaml.v3"

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

func (config *ConfigSignature) ParseConfig(fileName string) error {
	if FileExists(fileName) {
		//fmt.Println("Parsing " + fileName)

		yamlFile, err := ioutil.ReadFile(fileName)
		if err != nil {
			return err

		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return err
		}

		//fmt.Printf("Result: %v\n", config)

	} else {
		return errors.New(fileName + " doesn't exist")
	}
	return nil

}

type Signature interface {
	Name() string
	Match(js *JavaScript) []string
}

func (config *ConfigSignature) GetSignatures() ([]Signature, error) {
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
			} else {
				return signatures, err
			}
		}
	}
	//fmt.Println("total signatures: " + strconv.Itoa(len(signatures)))
	return signatures, nil
}

type SimpleSignature struct {
	part  string
	match string
	name  string
}

func (s *SimpleSignature) Name() string {
	return s.name
}
func (s *SimpleSignature) Match(js *JavaScript) (secrets []string) {
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

		if strings.Contains(js.UrlAddr.string, s.match) {
			secrets = append(secrets, s.name+" "+s.match+" found in URL")
			l.Log.Debug(s.name + " " + s.match + " found within URL of " + js.UrlAddr.string)
		}
		if strings.Contains(js.Content, s.match) {
			secrets = append(secrets, s.name+" "+s.match+" found in content")
			l.Log.Debug(s.name + " " + s.match + " found in content of " + js.UrlAddr.string)
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
func (s *PatternSignature) Match(js *JavaScript) (secrets []string) {
	if s.match != nil {
		if s.match.MatchString(js.UrlAddr.string) {
			tmp := s.match.FindAllString(js.UrlAddr.string, -1)
			for _, secret := range tmp {
				secrets = append(secrets, s.name+" "+secret+" found within URL of "+js.UrlAddr.string)
				l.Log.Debug(s.name + " " + secret + " found within URL of " + js.UrlAddr.string)
			}
		}
		if s.match.MatchString(js.Content) {
			tmp := s.match.FindAllString(js.Content, -1)
			for _, secret := range tmp {
				secrets = append(secrets, s.name+" "+secret+" found within content of "+js.UrlAddr.string)
				l.Log.Debug(s.name + " " + secret + " found within content of " + js.UrlAddr.string)
			}
		}
	}

	return secrets
}
