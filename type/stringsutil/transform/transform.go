package transform

import (
	"errors"

	"github.com/grokify/mogo/type/maputil"
)

type TransformFunc func(s string) string

func TransformMap(xf func(string) string, s []string) (map[string]string, map[string][]string, error) {
	out := map[string]string{}
	for _, si := range s {
		out[si] = xf(si)
	}
	dupes := maputil.DuplicateValues(out)
	if !maputil.UniqueValues(out) {
		return out, dupes, errors.New("strcase collisions")
	}
	return out, dupes, nil
}
