package osutil

import (
	"path/filepath"
	"testing"
)

var existTests = []struct {
	filename string
	exists   bool
}{
	{"exist.txt", true},
	{"doesnotexist.txt", false},
}

func TestExist(t *testing.T) {
	for _, tt := range existTests {
		exists, err := Exists(filepath.Join("filepath_testdata", tt.filename))
		if err != nil {
			t.Errorf("osutil.Exists(\"%s\") Error [%s]", tt.filename, err.Error())
		}
		if exists != tt.exists {
			t.Errorf("osutil.Exists(\"%s\") Want [%v] Got [%v]", tt.filename, tt.exists, exists)
		}
	}
}
