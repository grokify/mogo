package tok2vec

import (
	"strings"
	"testing"

	"github.com/grokify/mogo/text/languageutil"
)

var t2vSingleTests = []struct {
	tok string
	vec int
}{
	{"A", 1},
	{"C", 3},
	{"X", 24},
	{"Z", 26},
}

func TestVec(t *testing.T) {
	alphabet := strings.Split(languageutil.AlphabetEN, "")
	t2v := NewTok2Vec()
	t2v.Toks2Vecs(alphabet, false)

	for _, tt := range t2vSingleTests {
		tryVec, err := t2v.Vec(tt.tok)
		if err != nil {
			panic(err)
		}
		if tryVec != tt.vec {
			t.Errorf("Tok2Vec.Vec(\"%s\") want [%d], got [%d]", tt.tok, tt.vec, tryVec)
		}
	}
}

var t2vMultiTests = []struct {
	toks              []string
	vecsMulti         []int
	vecsDeduped       []int
	vecMultiWantErr   bool
	vecDedupedWantErr bool
}{
	{[]string{"Foo", "Bar", "Baz"}, []int{1, 2, 3}, []int{1, 2, 3}, false, false},
	{[]string{"Foo", "Bar", "Bar", "Qux"}, []int{1, 2, 2, 3}, []int{1, 2, 3}, false, false},
}

func TestToks2Vecs(t *testing.T) {
	for _, tt := range t2vMultiTests {
		t2vMulti := NewTok2Vec()
		vecsMultiTry := t2vMulti.Toks2Vecs(tt.toks, false)
		if !tt.vecMultiWantErr && len(vecsMultiTry) != len(tt.vecsMulti) {
			t.Errorf("Tok2Vec.Toks2Vecs() wrong output count: params [%v, %v] want [%v] got [%v]", tt.toks, false, tt.vecsMulti, vecsMultiTry)
		} else {
			for i, vecWant := range tt.vecsMulti {
				if !tt.vecMultiWantErr && vecWant != vecsMultiTry[i] {
					t.Errorf("Tok2Vec.Toks2Vecs() wrong output vecs: params [%v, %v] want [%v] got [%v]", tt.toks, false, tt.vecsMulti, vecsMultiTry)
				}
			}
		}

		t2vDeduped := NewTok2Vec()
		vecsDedupedTry := t2vDeduped.Toks2Vecs(tt.toks, true)
		if !tt.vecDedupedWantErr && len(vecsDedupedTry) != len(tt.vecsDeduped) {
			t.Errorf("Tok2Vec.Toks2Vecs() wrong output count: params [%v, %v] want [%v] got [%v]", tt.toks, true, tt.vecsDeduped, vecsDedupedTry)
		} else {
			for i, vecWant := range tt.vecsDeduped {
				if !tt.vecMultiWantErr && vecWant != vecsDedupedTry[i] {
					t.Errorf("Tok2Vec.Toks2Vecs() wrong output vecs: params [%v, %v] want [%v] got [%v]", tt.toks, true, tt.vecsDeduped, vecsDedupedTry)
				}
			}
		}
	}
}
