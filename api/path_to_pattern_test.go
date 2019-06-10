package api

import (
	"testing"
)

var urlToPatternTests = []struct {
	v    string
	want string
}{
	{"/pets", "/pets"},
	{"/pets?q=Fido", "/pets"},
	{"/pets/dog", "/pets/{petId}"},
	{"/pets/dog/vacinations", "/pets/{petId}/vacinations"},
	{"/pets/dog/vacinations/123", "/pets/{petId}/vacinations/{vacinationId}"}}

func TestURLToPattern(t *testing.T) {
	ut := NewURLTransformer()
	ut.LoadPaths([]string{
		"/pets",
		"/pets/{petId}",
		"/pets/{petId}/vacinations",
		"/pets/{petId}/vacinations/{vacinationId}"})

	for _, tt := range urlToPatternTests {
		got := ut.URLActualToPattern(tt.v)
		if got != tt.want {
			t.Errorf(`URLActualToPattern("%v") Failed: want [%v], got [%v]`, tt.v, tt.want, got)
		}
	}
}
