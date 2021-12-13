package strconvutil

import (
	"testing"

	"github.com/grokify/mogo/text/languageutil"
)

var strconvLangTests = []struct {
	lang   string
	strval string
	intval int
}{
	{languageutil.English, "1,234,567,890", 1234567890},
	{languageutil.Thai, "1,234,567,890", 1234567890},
	{languageutil.Danish, "1 234 567 890", 1234567890},
	{languageutil.Finnish, "1 234 567 890", 1234567890},
	{languageutil.German, "1 234 567 890", 1234567890},
	{languageutil.French, "1 234 567 890", 1234567890},
	{languageutil.Italian, "1.234.567.890", 1234567890},
	{languageutil.Norwegian, "1.234.567.890", 1234567890},
	{languageutil.Spanish, "1.234.567.890", 1234567890},
}

func TestAtoiLang(t *testing.T) {
	for _, tt := range strconvLangTests {
		tryInt, err := AtoiLang(tt.lang, tt.strval)
		if err != nil {
			t.Errorf("strconvutil.AtoiLang(\"%s\", \"%s\") Error: [%s]",
				tt.lang, tt.strval, err.Error())
		}
		if err == nil && tryInt != tt.intval {
			t.Errorf("strconvutil.AtoiLang(\"%s\", \"%s\") Error: want [%d], got [%d]",
				tt.lang, tt.strval, tt.intval, tryInt)
		}
	}
}
