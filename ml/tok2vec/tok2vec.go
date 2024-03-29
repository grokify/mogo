package tok2vec

import (
	"errors"

	"github.com/grokify/mogo/type/number"
)

// Tok2Vec is designed to convert tokens to vectors for multiple documents so that
// tokens that exist in multiple documents will be assigned a single vector.
type Tok2Vec struct {
	toks2VecsMap map[string]int
	vecs2ToksMap map[int]string
}

func NewTok2Vec() Tok2Vec {
	return Tok2Vec{
		toks2VecsMap: map[string]int{},
		vecs2ToksMap: map[int]string{}}
}

func (t2v *Tok2Vec) Toks2Vecs(tokens []string, deduplicate bool) []int {
	vecs := []int{}
	for _, tok := range tokens {
		if vec, ok := t2v.toks2VecsMap[tok]; ok {
			if !deduplicate {
				vecs = append(vecs, vec)
			}
			continue
		}
		vec := len(t2v.toks2VecsMap) + 1
		t2v.toks2VecsMap[tok] = vec
		t2v.vecs2ToksMap[vec] = tok
		vecs = append(vecs, vec)
	}
	return vecs
}

func (t2v *Tok2Vec) Toks2VecsFloat64(tokens []string, deduplicate bool) []float64 {
	return number.SliceToFloat64(t2v.Toks2Vecs(tokens, deduplicate))
}

var (
	ErrTokenNotFound  = errors.New("token not found")
	ErrVectorNotFound = errors.New("vector not found")
)

func (t2v *Tok2Vec) Tok(vec int) (string, error) {
	if tok, ok := t2v.vecs2ToksMap[vec]; ok {
		return tok, nil
	}
	return "", ErrVectorNotFound
}

func (t2v *Tok2Vec) Vec(tok string) (int, error) {
	if vec, ok := t2v.toks2VecsMap[tok]; ok {
		return vec, nil
	}
	return -1, ErrTokenNotFound
}
