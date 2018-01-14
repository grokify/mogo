package urlutil

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type URLValidator struct {
	RequiredSchemes map[string]int
}

func (uv *URLValidator) SchemesToLower() {
	newSchemes := map[string]int{}
	for scheme, v := range uv.RequiredSchemes {
		newSchemes[strings.ToLower(scheme)] = v
	}
	uv.RequiredSchemes = newSchemes
}

func (uv *URLValidator) ValidateURLString(s string) (*url.URL, error) {
	if len(strings.TrimSpace(s)) < 1 {
		return nil, fmt.Errorf("URL is empty.")
	}
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return u, err
	}
	return uv.ValidateURL(u)
}

func (uv *URLValidator) ValidateURL(u *url.URL) (*url.URL, error) {
	if len(uv.RequiredSchemes) > 0 {
		if _, ok := uv.RequiredSchemes[u.Scheme]; !ok {
			return u,
				fmt.Errorf("Scheme `%v` is not in list of required schemes: %v",
					u.Scheme, uv.RequiredSchemesSortedString())
		}
	}
	return u, nil
}

func (uv *URLValidator) RequiredSchemesSorted() []string {
	schemes := []string{}
	for scheme, _ := range uv.RequiredSchemes {
		schemes = append(schemes, scheme)
	}
	sort.Strings(schemes)
	return schemes
}

func (uv *URLValidator) RequiredSchemesSortedString() string {
	return strings.Join(uv.RequiredSchemesSorted(), ",")
}
