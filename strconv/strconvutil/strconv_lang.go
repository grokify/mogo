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
		s = strings.Replace(s, sepThousands, "", -1)
	}
	return strconv.Atoi(s)
}

func DecimalSeparator(lang string) (string, error) {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == languageutil.English ||
		lang == languageutil.Thai {
		return ".", nil
	} else if lang == languageutil.Italian ||
		lang == languageutil.Norwegian ||
		lang == languageutil.Spanish {
		return ",", nil
	} else if lang == languageutil.Danish ||
		lang == languageutil.Finnish ||
		lang == languageutil.French ||
		lang == languageutil.German {
		return ",", nil
	}
	return "", fmt.Errorf("language not found [%s]", lang)
}

func ThousandsSeparator(lang string) (string, error) {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == languageutil.English ||
		lang == languageutil.Thai {
		return ",", nil
	} else if lang == languageutil.Italian ||
		lang == languageutil.Norwegian ||
		lang == languageutil.Spanish {
		return ".", nil
	} else if lang == languageutil.Danish ||
		lang == languageutil.Finnish ||
		lang == languageutil.French ||
		lang == languageutil.German {
		return " ", nil
	}
	return "", fmt.Errorf("language not found [%s]", lang)
}
