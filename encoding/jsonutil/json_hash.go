package jsonutil

import (
	"bytes"
	"encoding/json"

	"github.com/grokify/mogo/crypto/shautil"
	"github.com/grokify/mogo/type/stringsutil"
)

func SHA512d256Base32(v any, padding rune) (string, error) {
	if b, err := json.Marshal(v); err != nil {
		return "", err
	} else if sha, err := shautil.Sum512d256Base32(bytes.NewReader(b), padding); err != nil {
		return "", err
	} else {
		return stringsutil.RemoveNonPrintable(sha), nil
	}
}
