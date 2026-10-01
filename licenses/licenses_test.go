package licenses

import (
	"os"
	"strings"
	"testing"
)

func TestIsMIT(t *testing.T) {
	reflowed := strings.ToUpper(strings.Join(strings.Fields(mitTerms), "\n"))
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"terms only", mitTerms, true},
		{"title and copyright", "MIT License\n\nCopyright (c) 2026 Example Author\n\n" + mitTerms, true},
		{"(MIT) title, multiple copyrights", "The MIT License (MIT)\n\nCopyright 2020 A\n(c) 2021 B\n\n" + mitTerms + "\n", true},
		{"reflowed and uppercased", reflowed, true},
		{"empty", "", false},
		{"header only", "MIT License\nCopyright 2026 A", false},
		{"extra clause appended", mitTerms + "\n\nCommons Clause: the Software may not be sold.", false},
		{"altered terms", strings.Replace(mitTerms, "free of charge", "for a fee", 1), false},
		{"unexpected header line", "MIT License\nSee also NOTICE\n\n" + mitTerms, false},
		{"other license", "Apache License\nVersion 2.0, January 2004", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMIT(tt.s); got != tt.want {
				t.Errorf("IsMIT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMITRepoLicense(t *testing.T) {
	b, err := os.ReadFile("../LICENSE")
	if err != nil {
		t.Fatal(err)
	}
	if !IsMIT(string(b)) {
		t.Error("IsMIT(mogo LICENSE) = false, want true")
	}
}
