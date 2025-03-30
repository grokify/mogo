package strconvutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/mogo/text/languageutil"
)

// var rxComma = regexp.MustCompile(`,`)

// AtoiLang provides language parsing to handle
// thousands separators.
// Number formats: https://docs.oracle.com/cd/E19455-01/806-0169/overview-9/index.html
func AtoiLang(lang, s string) (int, error) {
	s = strings.TrimSpace(s)
	sepThousands, err := ThousandsSeparator(lang)
	if err != nil {
		return 0, err
	}
	if len(sepThousands) > 0 {
		s = strings.ReplaceAll(s, sepThousands, "")
	}
	return strconv.Atoi(s)
}

func DecimalSeparator(lang string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case languageutil.English,
		languageutil.Thai:
		return ".", nil
	case languageutil.Danish,
		languageutil.Finnish,
		languageutil.French,
		languageutil.German,
		languageutil.Italian,
		languageutil.Norwegian,
		languageutil.Spanish:
		return ",", nil
	default:
		return "", fmt.Errorf("language not found [%s]", lang)
	}
}

func ThousandsSeparator(lang string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case languageutil.English,
		languageutil.Thai:
		return ",", nil
	case languageutil.Italian,
		languageutil.Norwegian,
		languageutil.Spanish:
		return ".", nil
	case languageutil.Danish,
		languageutil.Finnish,
		languageutil.French,
		languageutil.German:
		return " ", nil
	default:
		return "", fmt.Errorf("language not found [%s]", lang)
	}
}
