package urlutil

import (
	"net/url"
	"regexp"
	"strconv"
)

// UrlMore provides additional URL parsing and reconstruction capabilties
// above and beyond URL. Specifically it can parse out the port number and
// return URLs that strip off the target fragment as well as the query
// string.
type URLMore struct {
	URL                  *url.URL
	URLWoQueryWoFragment string
	URLWoFragment        string
	Port                 int
}

func NewURLMore() URLMore {
	urlMore := URLMore{
		URLWoQueryWoFragment: "",
		URLWoFragment:        ""}
	return urlMore
}

// Parse uses `url.Parse()` to create a URL object. When using an already
// created URL object, simply set the `Url` property and then call `Inflate`.
func (urlMore *URLMore) Parse(rawurl string) error {
	myURL, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	urlMore.URL = myURL
	urlMore.Inflate()
	return nil
}

func (urlMore *URLMore) Inflate() {
	myURL := urlMore.URL
	urlMore.URLWoQueryWoFragment = myURL.Scheme + "://" + myURL.Host + myURL.Path
	if len(myURL.RawQuery) > 0 {
		urlMore.URLWoFragment = urlMore.URLWoQueryWoFragment + "?" + myURL.RawQuery
	} else {
		urlMore.URLWoFragment = urlMore.URLWoQueryWoFragment
	}
	rx := regexp.MustCompile(`:([0-9]+)$`)
	rs := rx.FindStringSubmatch(myURL.Host)
	if len(rs) > 0 {
		port, err := strconv.Atoi(rs[1])
		if err == nil {
			urlMore.Port = port
		}
	}
}
