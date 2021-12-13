package smtputil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

// GmailHost is the Google mail smtp hostname.
const GmailHost = "gmail.com"

// GmailUserPlusTime creates a Gmail SMTP address appending
// a time in "DT14" time format.
func GmailAddressPlusTime(smtpUser string) (string, error) {
	smtpUser = strings.TrimSpace(smtpUser)
	if len(smtpUser) == 0 {
		return "", errors.New("No SMTP User address provided.")
	}
	t := time.Now()
	return fmt.Sprintf(`%s+%s@%s`, smtpUser, t.Format(timeutil.DT14), GmailHost), nil
}
