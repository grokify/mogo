package maputil

import "errors"

var (
	ErrKeyNotExist = errors.New("key not found")
	ErrNotString   = errors.New("value not string")
)

// MapStrAny represents a `map[string]any`
type MapStrAny map[string]any

func (msa MapStrAny) ValueString(k string, errOnNotExist bool) (string, error) {
	v, ok := msa[k]
	if !ok {
		if errOnNotExist {
			return "", ErrKeyNotExist
		} else {
			return "", nil
		}
	}
	s, ok := v.(string)
	if !ok {
		return "", ErrNotString
	}
	return s, nil
}
