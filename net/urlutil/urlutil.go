package urlutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/grokify/mogo/type/slicesutil"
)

// AppendValues appends one `url.Values` to another `url.Values`.
func AppendValues(v1, v2 url.Values, inclDuplicates bool) url.Values {
	out := url.Values{}
	if inclDuplicates {
		for key, vals := range v1 {
			out[key] = append(out[key], vals...)
		}
		for key, vals := range v2 {
			out[key] = append(out[key], vals...)
		}
		return out
	}
	exists := map[string]map[string]int{}
	for k, vals := range v1 {
		if _, ok := exists[k]; !ok {
			exists[k] = map[string]int{}
		}
		for _, v := range vals {
			out.Add(k, v)
			exists[k][v]++
		}
	}
	for k, vals := range v2 {
		for _, v := range vals {
			if srcVals, ok := exists[k]; !ok {
				out.Add(k, v)
			} else if _, ok := srcVals[v]; !ok {
				out.Add(k, v)
			}
		}
	}
	return out
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

/*
// URLAddQueryMap returns a URL as a string from a base URL and a
// set of query parameters as a map[string]string{}
func URLAddQueryMapString(baseUrl string, queryParams map[string]string) (string, error) {
	if len(queryParams) < 1 {
		return baseUrl, nil
	}
	queryValues := map[string][]string{}
	for key, val := range queryParams {
		queryValues[key] = []string{val}
	}
	curUrl, err := URLAddQueryValues(baseUrl, queryValues)
	if err != nil {
		return baseUrl, err
	}
	return curUrl.String(), nil
}

// BuildURL returns a URL string from a base URL and url.Values.
func BuildURL(baseUrl string, queryValues url.Values) string {
	qryString := queryValues.Encode()
	if len(qryString) > 0 {
		return baseUrl + "?" + qryString
	}
	return baseUrl
}

// BuildURLQueryString to be deprecated in favor of URLAddQueryString
func BuildURLQueryString(baseUrl string, qry any) string {
	v, _ := query.Values(qry)
	qryString := v.Encode()
	if len(qryString) > 0 {
		return baseUrl + "?" + qryString
	}
	return baseUrl
}
*/

func URLAddQuery(inputURL *url.URL, qry url.Values, inclDuplicates bool) *url.URL {
	if len(qry) == 0 {
		return inputURL
	}
	allQS := AppendValues(inputURL.Query(), qry, inclDuplicates)
	inputURL.RawQuery = allQS.Encode()
	return inputURL
}

func URLStringAddQuery(inputURL string, qry url.Values, inclDuplicates bool) (*url.URL, error) {
	if goURL, err := url.Parse(inputURL); err != nil {
		return nil, err
	} else {
		return URLAddQuery(goURL, qry, inclDuplicates), nil
	}
}

// GetURLBody returns an HTTP response byte array body from a URL.
func GetURLBody(absoluteURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, absoluteURL, nil)
	if err != nil {
		return []byte{}, err
	}
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

// GetURLPostBody returns a HTTP post body as a byte array from a
// URL, body type and an io.Reader.
func GetURLPostBody(absoluteURL string, bodyType string, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Post(absoluteURL, bodyType, reqBody)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

// JoinAbsolute performs a path.Join() while preserving two slashes after the scheme.
func JoinAbsolute(elem ...string) string {
	return regexp.MustCompile(`^([A-Za-z]+:/)`).ReplaceAllString(path.Join(elem...), "${1}/")
}

// Join returns joining URL paths parts.
func Join(elem ...string) string {
	return rxFwdSlashMore.ReplaceAllString(
		strings.Join(elem, "/"), "/")
}

var (
	rxURIScheme    *regexp.Regexp = regexp.MustCompile(`^([A-Za-z][0-9A-Za-z]*:/)`)
	rxFwdSlashMore *regexp.Regexp = regexp.MustCompile(`/+`)
)

// CondenseURI  trims spaces and condenses slashes.
func CondenseURI(uri string) string {
	return rxURIScheme.ReplaceAllString(
		rxFwdSlashMore.ReplaceAllString(strings.TrimSpace(uri), "/"),
		"${1}/")
}

// URLValuesEncodeSorted returns and encoded string with sorting
func URLValuesEncodeSorted(v url.Values, priorities []string) string {
	encoded := []string{}
	priorityKeys := map[string]int{}

	priorities = slicesutil.Dedupe(priorities)

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
	for k := range v {
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

// URLToAddress converts a string like `https://localhost` to `localhost:https`.
func URLToAddress(s string) string {
	delimit := "://"
	parts := strings.Split(s, delimit)
	if len(parts) > 0 {
		return strings.Join(parts[1:], delimit) + ":" + parts[0]
	}
	return s
}

func MSABytesToValues(b []byte) (url.Values, error) {
	msa := map[string]any{}
	if err := json.Unmarshal(b, &msa); err != nil {
		return url.Values{}, err
	} else {
		return MSAToValues(msa), nil
	}
}

func MSAToValues(msa map[string]any) url.Values {
	vs := url.Values{}
	for k, v := range msa {
		vs.Add(k, fmt.Sprintf("%v", v))
	}
	return vs
}
