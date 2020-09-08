// Package osutil implements some OS utility functions.
package osutil

import (
	"os"
	"regexp"
	"strings"
)

type EnvVar struct {
	Key   string
	Value string
}

func Env() []EnvVar {
	envs := []EnvVar{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if len(pair) > 0 {
			key := strings.TrimSpace(pair[0])
			if len(key) > 0 {
				env := EnvVar{Key: key}
				if len(pair) > 1 {
					env.Value = strings.Join(pair[1:], "=")
				}
				envs = append(envs, env)
			}
		}
	}
	return envs
}

// EnvFiltered returns a map[string]string of environment
// variables that match a regular expression.
func EnvFiltered(rx *regexp.Regexp) map[string]string {
	res := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if rx.MatchString(pair[0]) {
			res[pair[0]] = pair[1]
		}
	}
	return res
}

func EnvExists(fields ...string) (missing []string, haveAll bool) {
	missing = []string{}
	for _, field := range fields {
		val := strings.TrimSpace(os.Getenv(field))
		if len(val) == 0 {
			missing = append(missing, field)
		}
	}
	haveAll = true
	if len(missing) > 0 {
		haveAll = false
	}
	return
}
