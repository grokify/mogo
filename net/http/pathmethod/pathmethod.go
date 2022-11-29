package pathmethod

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/net/httputilmore"
)

// PathMethod returns a path-method string which can be used as a unique identifier for an HTTP endpoint.
func PathMethod(htPath, htMethod string) string {
	htPath = strings.TrimSpace(htPath)
	htMethod = strings.ToUpper(strings.TrimSpace(htMethod))
	parts := []string{}
	if len(htPath) > 0 {
		parts = append(parts, htPath)
	}
	if len(htMethod) > 0 {
		parts = append(parts, htMethod)
	}
	return strings.Join(parts, " ")
}

var ErrPathMethodInvalid = errors.New("pathmethod string invalid")

func ParsePathMethod(pathmethod string) (string, string, error) {
	parts := strings.Split(pathmethod, " ")
	if len(parts) != 2 {
		return "", "", ErrPathMethodInvalid
	}
	method, err := httputilmore.ParseHTTPMethodString(parts[1])
	return strings.TrimSpace(parts[0]), method, err
}

type PathMethodSet struct {
	PathMethods map[string]int
}

func NewPathMethodSet() PathMethodSet {
	return PathMethodSet{
		PathMethods: map[string]int{}}
}

func (pms *PathMethodSet) init() {
	if pms.PathMethods == nil {
		pms.PathMethods = map[string]int{}
	}
}

// Add adds pathmethod strings
func (pms *PathMethodSet) Add(pathmethods ...string) error {
	pms.init()
	for _, pm := range pathmethods {
		htPath, htMethod, err := ParsePathMethod(pm)
		if err != nil {
			return err
		}
		pathMethod := PathMethod(htPath, htMethod)
		pms.PathMethods[pathMethod]++
	}
	return nil
}

func (pms *PathMethodSet) Count() int {
	pms.init()
	return len(pms.PathMethods)
}

func (pms *PathMethodSet) Exists(htPath, htMethod string) bool {
	pms.init()
	pm := PathMethod(htPath, htMethod)
	_, ok := pms.PathMethods[pm]
	return ok
}

func (pms *PathMethodSet) StringExists(pathMethod string) bool {
	pms.init()
	_, ok := pms.PathMethods[pathMethod]
	return ok
}
