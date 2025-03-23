package jsonraw

import (
	"encoding/json"
	"os"

	"github.com/grokify/mogo/type/maputil"
)

func ObjectKeys(b []byte) ([]string, error) {
	msa := map[string]any{}
	if err := json.Unmarshal(b, &msa); err != nil {
		return []string{}, err
	} else {
		return maputil.Keys(msa), nil
	}
}

func ObjectKeysFile(filename string) ([]string, error) {
	if b, err := os.ReadFile(filename); err != nil {
		return []string{}, err
	} else {
		return ObjectKeys(b)
	}
}

// ObjectModify creates a new byte slice, from an existing byte slice, with only the supplied field names.
func ObjectModify(b []byte, fieldNamesInclCopy []string, fieldNameInclUpsertValues map[string]any) ([]byte, error) {
	msa := map[string]any{}
	err := json.Unmarshal(b, &msa)
	if err != nil {
		return []byte{}, err
	}
	if len(msa) == 0 {
		return b, nil
	}
	incl := map[string]int{}
	for _, f := range fieldNamesInclCopy {
		incl[f]++
	}
	out := map[string]any{}
	for k, v := range msa {
		if _, ok := incl[k]; ok {
			out[k] = v
		}
	}
	for k, v := range fieldNameInclUpsertValues {
		out[k] = v
	}
	return json.Marshal(out)
}
