package core

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	l "github.com/hiddengearz/jsubfinder/core/logger"
	tld "github.com/jpillora/go-tld"
)

type UrlAddr struct {
	string            //URL address
	rootDomain string //Top Level Domain of the URL
}

//GetContent retrieves the content of urls - #### MAYBE CHECK FOR redirects and follow them????
func (u *UrlAddr) GetContent(client *http.Client) (newContent string, isJS bool, err error) {
	//defer lock.RUnlock()

	var req *http.Request
	var resp *http.Response
	if Debug {
		defer TimeTrack(time.Now(), "GetContent "+u.string)
	}

	//If the provided URL starts with HTTP/s make a request
	if strings.HasPrefix(u.string, "https://") || strings.HasPrefix(u.string, "http://") {

		if IsUrlVisited(u.string) {
			err = errors.New("Url " + u.string + " was been scanned before")
			return
		}

		req, err = http.NewRequest(http.MethodGet, u.string, nil)
		if err != nil {
			l.Log.Debug("Client get failed: %s\n", err)
			return
		}

		req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")

		resp, err = client.Do(req)
		if err != nil {
			l.Log.Debug("Client get failed: %s\n", err)
			return
		}

	} else { //if the request doesn't start with HTTP/s, add it
		if IsUrlVisited("http://" + u.string) {
			err = errors.New("Url " + "http://" + u.string + " was been scanned before")
			return
		}

		req, err = http.NewRequest(http.MethodGet, "http://"+u.string, nil)
		if err != nil {
			l.Log.Debug("Client get failed: %s\n", err)
			return
		}

		req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")

		resp, err = client.Do(req)

		if err != nil && !strings.Contains(string(err.Error()), "no such host") { //if there is an error and its not due to dns
			l.Log.Debug("new err Client get failed: %s\n", err)

			if IsUrlVisited("https://" + u.string) {
				err = errors.New("Url " + u.string + " was been scanned before")
				return
			}

			req, err = http.NewRequest(http.MethodGet, "https://"+u.string, nil) //try a request with https
			if err != nil {
				l.Log.Debug("Client get failed: %s\n", err)
				return
			}

			resp, err = client.Do(req)
			if err != nil {
				l.Log.Debug("Client get failed: %s\n", err)
				return
			}
			u.string = "https://" + u.string
		} else if err != nil {
			l.Log.Debug("Client get failed: %s\n", err)
			return
		} else { //if no error
			u.string = "http://" + u.string

		}
	}

	contenType := resp.Header.Get("Content-Type")
	if contenType == "" {
		isJS = false
	} else if strings.Contains(contenType, "javascript") {
		isJS = true
		//fmt.Println("content type js" + u.string)
	}

	//read the body and return it
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	newContent = (string(bodyBytes))

	return
}

//Get the Top Level Domain of the URL
func (u *UrlAddr) GetRootDomain() (err error) {
	u2, err := tld.Parse(u.string)
	if err != nil {
		return
	}

	u.rootDomain = u2.Domain + "." + u2.TLD
	return
}
