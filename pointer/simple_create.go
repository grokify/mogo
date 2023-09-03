// pointer package provides some pointer shortcuts. See:
// 1: https://stackoverflow.com/questions/50830676/set-int-pointer-to-int-value-golang
// 2: https://github.com/openlyinc/pointy
package pointer

func Int64(num int64) *int64 { return &num }
func Pointer[E any](e E) *E  { return &e }

func ToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Dereference[E any](e *E) E {
	if e == nil {
		return *new(E)
	} else {
		return *e
	}
}
