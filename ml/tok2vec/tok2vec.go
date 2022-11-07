package tok2vec

import "errors"

type Tok2Vec struct {
	MapToks2Vecs map[string]int
	MapVecs2Toks map[int]string
}

func NewTok2Vec() Tok2Vec {
	return Tok2Vec{
		MapToks2Vecs: map[string]int{},
		MapVecs2Toks: map[int]string{}}
}

func (t2v *Tok2Vec) Toks2Vecs(tokens []string) []int {
	vecs := []int{}
	for _, tok := range tokens {
		if vec, ok := t2v.MapToks2Vecs[tok]; ok {
			vecs = append(vecs, vec)
			continue
		}
		vec := len(t2v.MapToks2Vecs) + 1
		t2v.MapToks2Vecs[tok] = vec
		t2v.MapVecs2Toks[vec] = tok
		vecs = append(vecs, vec)
	}
	return vecs
}

func (t2v *Tok2Vec) Doc2VecFloat64(tokens ...string) []float64 {
	vecs := []float64{}
	for _, tok := range tokens {
		if vec, ok := t2v.MapToks2Vecs[tok]; ok {
			vecs = append(vecs, float64(vec))
			continue
		}
		vec := len(t2v.MapToks2Vecs) + 1
		t2v.MapToks2Vecs[tok] = vec
		t2v.MapVecs2Toks[vec] = tok
		vecs = append(vecs, float64(vec))
	}
	return vecs
}

var (
	ErrTokenNotFound  = errors.New("token not found")
	ErrVectorNotFound = errors.New("vector not found")
)

func (t2v *Tok2Vec) Tok(vec int) (string, error) {
	if tok, ok := t2v.MapVecs2Toks[vec]; ok {
		return tok, nil
	}
	return "", ErrVectorNotFound
}

func (t2v *Tok2Vec) Vec(tok string) (int, error) {
	if vec, ok := t2v.MapToks2Vecs[tok]; ok {
		return vec, nil
	}
	return -1, ErrTokenNotFound
}
