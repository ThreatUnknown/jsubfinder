package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	l "github.com/hiddengearz/jsubfinder/core/logger"
)

//RadFile reads the content of a file and returns it in a slice
func ReadFile(filePath string) (content []string, err error) {
	if Debug {
		defer TimeTrack(time.Now(), "ReadFile "+filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		l.Log.Debug(err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return content, err
	}

	return
}

//Find searches a []string for a substring and return it's position in the array and a bool for if it's in the array
func Find(slice []string, val string) (int, bool) {
	//if Debug {
	//	defer TimeTrack(time.Now(), "Find "+val)
	//}
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

//GetHTTprotocol parses a given URL to get it's protocol e.g http:// or https://
func GetHTTprotocol(url string) (protocol string, err error) {
	//if Debug {
	//	defer TimeTrack(time.Now(), "GetHTTprotocol "+url)
	//}
	if strings.HasPrefix(url, "http://") {
		protocol = "http://"
		return
	} else if strings.HasPrefix(url, "https://") {
		protocol = "https://"
		return
	} else {
		err = errors.New("No prefix")
		return
	}
}

//TimeTrack times a function just for testing
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(name + " took " + strconv.FormatFloat(elapsed.Seconds(), 'f', 3, 64) + "s")
}

//SaveResults saves the content to the spcified file
func SaveResults(fileLocation string, newContent []string) error {
	newFile, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(newFile)

	for _, data := range newContent {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	newFile.Close()
	return nil
}

//fileExists returns a bool if the file exists or not
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		l.Log.Debug(err)
		return false
	}
	return !info.IsDir()
}
