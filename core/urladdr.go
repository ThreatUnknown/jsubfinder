package core

import (
	"errors"
	"strings"
	"time"

	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/valyala/fasthttp"
)

type UrlAddr struct {
	string        //URL address
	tld    string //Top Level Domain of the URL
}

//GetContent retrieves the content of urls - #### MAYBE CHECK FOR redirects and follow them????
func (u *UrlAddr) GetContent(client *fasthttp.Client) (err error, newContent string) {
	if Debug {
		defer TimeTrack(time.Now(), "GetContent "+u.string)
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")
	if strings.HasPrefix(u.string, "https://") || strings.HasPrefix(u.string, "http://") {

		req.SetRequestURI(u.string)

		err = client.Do(req, resp)
		if err != nil {
			l.Log.Debug("Client get failed: %s\n", err)
			return
		}

	} else {
		req.SetRequestURI("https://" + u.string)

		err = client.Do(req, resp)
		if err != nil && !strings.Contains(string(err.Error()), "no such host") {
			l.Log.Debug("new err Client get failed: %s\n", err)
			req.SetRequestURI("http://" + u.string)

			err = client.Do(req, resp)
			if err != nil {
				l.Log.Debug("Client get failed: %s\n", err)
				return
			}
			err = errors.New("http")

		} else if err != nil {
			l.Log.Debug("Client get failed: %s\n", err)
			return
		} else {
			err = errors.New("https")
		}
	}

	bodyBytes := resp.Body()
	newContent = (string(bodyBytes))

	return err, newContent
}
