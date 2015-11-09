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

type UrlMore struct {
	Url                  *url.URL
	UrlWoQueryWoFragment string
	UrlWoFragment        string
	Port                 int
}

func NewUrlMore() UrlMore {
	urlMore := UrlMore{
		UrlWoQueryWoFragment: "",
		UrlWoFragment:        ""}
	return urlMore
}

// Parse uses `url.Parse()` to create a URL object. When using an already
// created URL object, simply set the `Url` property and then call `Inflate`.

func (urlMore *UrlMore) Parse(rawurl string) error {
	myUrl, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	urlMore.Url = myUrl
	urlMore.Inflate()
	return nil
}

func (urlMore *UrlMore) Inflate() {
	myUrl := urlMore.Url
	urlMore.UrlWoQueryWoFragment = myUrl.Scheme + "://" + myUrl.Host + myUrl.Path
	if len(myUrl.RawQuery) > 0 {
		urlMore.UrlWoFragment = urlMore.UrlWoQueryWoFragment + "?" + myUrl.RawQuery
	} else {
		urlMore.UrlWoFragment = urlMore.UrlWoQueryWoFragment
	}
	rx := regexp.MustCompile(`:([0-9]+)$`)
	rs := rx.FindStringSubmatch(myUrl.Host)
	if len(rs) > 0 {
		port, err := strconv.Atoi(rs[1])
		if err == nil {
			urlMore.Port = port
		}
	}
}
