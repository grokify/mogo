// pointer package provides some pointer shortcuts. See:
// 1: https://stackoverflow.com/questions/50830676/set-int-pointer-to-int-value-golang
// 2: https://github.com/openlyinc/pointy
package pointer

func Dereference[E any](e *E) E {
	if e == nil {
		return *new(E)
	} else {
		return *e
	}
}

func DereferenceSlice[S ~[]*E, E any](s S) []E {
	out := []E{}
	for _, e := range s {
		out = append(out, *e)
	}
	return out
}

func Pointer[E any](e E) *E { return &e }

func PointerSlice[S ~[]E, E any](s S) []*E {
	out := []*E{}
	for i := range s {
		out = append(out, &s[i])
	}
	return out
}
