package urlutil

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/grokify/gotilla/type/stringsutil"
)

// AppendURLValues appends one url.Values to another url.Values.
func AppendURLValues(v1, v2 url.Values) url.Values {
	for key, vals := range v2 {
		for _, val := range vals {
			v1.Add(key, val)
		}
	}
	return v1
}

// ToSlug creates a slug byte array from an input byte array.
// Slugs have words separated by a hyphen with no punctuation
// or spaces.
func ToSlug(slug []byte) []byte {
	// Convert punctuation and spaces to hyphens: string([]byte{45}) = "-"
	slug = regexp.MustCompile(`[[:punct:]\s_]+`).ReplaceAll(slug, []byte{45})
	slug = regexp.MustCompile(`["']+`).ReplaceAll(slug, []byte{})
	return regexp.MustCompile(`(^-+|-+$)`).ReplaceAll(slug, []byte{})
}

// ToSlugLowerString creates a lower-cased slug string
func ToSlugLowerString(s string) string {
	return string(ToSlug([]byte(strings.ToLower(s))))
}

// BuildURLFromMap returns a URL as a string from a base URL and a
// set of query parameters as a map[string]string{}
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

// BuildURL returns a URL string from a base URL and url.Values.
func BuildURL(baseUrl string, queryValues url.Values) string {
	qryString := queryValues.Encode()
	if len(qryString) > 0 {
		return baseUrl + "?" + qryString
	}
	return baseUrl
}

func BuildURLQueryString(baseUrl string, qry interface{}) string {
	v, _ := query.Values(qry)
	qryString := v.Encode()
	if len(qryString) > 0 {
		return baseUrl + "?" + qryString
	}
	return baseUrl
}

// GetURLBody returns an HTTP response byte array body from
// a URL.
func GetURLBody(absoluteUrl string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, absoluteUrl, nil)
	if err != nil {
		return []byte{}, err
	}
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// GetURLPostBody returns a HTTP post body as a byte array from a
// URL, body type and an io.Reader.
func GetURLPostBody(absoluteUrl string, bodyType string, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Post(absoluteUrl, bodyType, reqBody)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// JoinAbsolute performs a path.Join() while preserving two slashes after the scheme.
func JoinAbsolute(elem ...string) string {
	return regexp.MustCompile(`^([A-Za-z]+:/)`).ReplaceAllString(path.Join(elem...), "${1}/")
}

var (
	rxFwdSlashMore *regexp.Regexp = regexp.MustCompile(`/+`)
	rxUriScheme    *regexp.Regexp = regexp.MustCompile(`^([A-Za-z][0-9A-Za-z]*:/)`)
)

// CondenseUri trims spaces and condenses slashes.
func CondenseUri(uri string) string {
	return rxUriScheme.ReplaceAllString(
		rxFwdSlashMore.ReplaceAllString(strings.TrimSpace(uri), "/"),
		"${1}/")
}

// UrlValuesStringSorted returns and encoded string with sorting
func UrlValuesEncodeSorted(v url.Values, priorities []string) string {
	encoded := []string{}
	priorityKeys := map[string]int{}

	priorities = stringsutil.Dedupe(priorities)

	for _, key := range priorities {
		if vals, ok := v[key]; ok {
			sort.Strings(vals)
			for _, val := range vals {
				qry := url.QueryEscape(key) + "=" + url.QueryEscape(val)
				encoded = append(encoded, qry)
			}
		}

		priorityKeys[key] = 1
	}

	sortedKeys := []string{}
	for k, _ := range v {
		if _, ok := priorityKeys[k]; !ok {
			sortedKeys = append(sortedKeys, k)
		}
	}
	sort.Strings(sortedKeys)
	for _, key := range sortedKeys {
		if vals, ok := v[key]; ok {
			sort.Strings(vals)
			for _, val := range vals {
				qry := url.QueryEscape(key) + "=" + url.QueryEscape(val)
				encoded = append(encoded, qry)
			}
		}
	}
	return strings.Join(encoded, "&")
}
