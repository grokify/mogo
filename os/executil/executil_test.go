//go:build linux || darwin
// +build linux darwin

package executil

import (
	"strings"
	"testing"
)

var execTests = []struct {
	v    string
	want []string
}{
	{v: "ls -al .", want: []string{".", "..", "executil.go"}},
}

func TestExecSimple(t *testing.T) {
	for _, tt := range execTests {
		o, _, err := ExecSimple(tt.v)
		if err != nil {
			t.Errorf("executil.ExecSimple(\"%s\") error (%s)",
				tt.v, err.Error())
		}
		lines := strings.Split(o.String(), "\n")
		missing := []string{}
	WANT:
		for _, w := range tt.want {
			for _, line := range lines {
				if strings.Contains(line, w) {
					continue WANT
				}
			}
			missing = append(missing, w)
		}
		if len(missing) > 0 {
			t.Errorf("executil.ExecSimple(\"%s\") mismatch want (%s) not found",
				tt.v, strings.Join(missing, ","))
		}
	}
}
