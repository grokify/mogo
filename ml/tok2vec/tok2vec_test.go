package tok2vec

import (
	"strings"
	"testing"

	"github.com/grokify/mogo/text/languageutil"
)

var t2vTests = []struct {
	tok string
	vec int
}{
	{"A", 1},
	{"C", 3},
	{"X", 24},
	{"Z", 26},
}

func TestMapStringIntSort(t *testing.T) {
	alphabet := strings.Split(languageutil.AlphabetEN, "")
	t2v := NewTok2Vec()
	t2v.Toks2Vecs(alphabet)

	for _, tt := range t2vTests {
		tryVec, err := t2v.Vec(tt.tok)
		if err != nil {
			panic(err)
		}
		if tryVec != tt.vec {
			t.Errorf("Tok2Vec.Vec(\"%s\") want [%d], got [%d]", tt.tok, tt.vec, tryVec)
		}
	}
}
