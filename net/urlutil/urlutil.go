package urlutil

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func ToSlug(src []byte) []byte {
	rxp := regexp.MustCompile(`[\*\s]+`)
	out := rxp.ReplaceAll(src, []byte("-"))
	rxm := regexp.MustCompile(`-+`)
	out = rxm.ReplaceAll(out, []byte("-"))
	rxb := regexp.MustCompile(`^[\s-]+`)
	out = rxb.ReplaceAll(out, []byte{})
	rxe := regexp.MustCompile(`[\s-]+$`)
	out = rxe.ReplaceAll(out, []byte{})
	return out
}

func ToSlugLowerString(s string) string {
	return string(ToSlug([]byte(strings.ToLower(s))))
}

func BuildURLFromMap(baseUrl string, queryParams map[string]string) string {
	if len(queryParams) < 1 {
		return baseUrl
	}
	queryValues := url.Values{}
	for key, val := range queryParams {
		queryValues.Set(key, val)
	}
	return BuildURL(baseUrl, queryValues)
}

func BuildURL(baseUrl string, queryValues url.Values) string {
	qryString := queryValues.Encode()
	if len(qryString) > 0 {
		return baseUrl + "?" + qryString
	}
	return baseUrl
}

func GetURLBody(absoluteUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", absoluteUrl, nil)
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func GetURLPostBody(absoluteUrl string, bodyType string, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Post(absoluteUrl, bodyType, reqBody)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
