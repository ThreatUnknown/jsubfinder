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
	string        //URL address
	tld    string //Top Level Domain of the URL
}

//GetContent retrieves the content of urls - #### MAYBE CHECK FOR redirects and follow them????
func (u *UrlAddr) GetContent(client *http.Client) (err error, newContent string) {
	var req *http.Request
	var resp *http.Response
	if Debug {
		defer TimeTrack(time.Now(), "GetContent "+u.string)
	}

	if strings.HasPrefix(u.string, "https://") || strings.HasPrefix(u.string, "http://") {

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

	} else {

		//Do request with HTTP://
		req, err = http.NewRequest(http.MethodGet, "http://"+u.string, nil)
		if err != nil {
			l.Log.Debug("Client get failed: %s\n", err)
			return
		}

		req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")

		resp, err = client.Do(req)
		if err != nil && !strings.Contains(string(err.Error()), "no such host") { //if there is an error and its not due to dns
			l.Log.Debug("new err Client get failed: %s\n", err)

			req, err = http.NewRequest(http.MethodGet, "https://"+u.string, nil) //try with https
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
		} else {
			u.string = "http://" + u.string
			err = errors.New("https")
		}
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	newContent = (string(bodyBytes))

	return
}

func (u *UrlAddr) GetTLD() (err error) {
	u2, err := tld.Parse(u.string)
	if err != nil {
		return
	}

	u.tld = u2.Domain + "." + u2.TLD
	return
}
