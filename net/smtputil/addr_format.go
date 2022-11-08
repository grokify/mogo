package smtputil

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/sort/sortutil"
)

// www.example.com@@info
// com.example.www@@@info

var (
	ErrEmailAddressEmpty             = errors.New("empty email address")
	ErrEmailAddressInvalidTooManyAts = errors.New("invalid email address. Too many ats")
)

// EmailAddrToSwapped converts an email address to "swapped" format for sorting like:
// "www.example.com@@info"
func EmailAddrToSwapped(e string) (string, error) {
	user, host, err := ParseEmailAddress(e)
	if err != nil {
		return "", err
	}
	if len(host) == 0 {
		return user, nil
	}
	return EmailUserHostToSwapped(user, host), nil
}

// EmailAddrToSwapped converts an email address to "reverse" format for sorting like:
// "com.example.www@@info"
func EmailAddrToReverse(e string) (string, error) {
	user, host, err := ParseEmailAddress(e)
	if err != nil {
		return "", err
	}
	if len(host) == 0 {
		return user, nil
	}
	return EmailUserHostToReverse(user, host), nil
}

func ParseEmailAddress(e string) (string, string, error) {
	parts := strings.Split(e, "@")
	if len(parts) == 0 {
		return "", "", ErrEmailAddressEmpty
	} else if len(parts) == 1 {
		return parts[0], "", nil
	} else if len(parts) == 2 {
		return parts[0], parts[1], nil
	}
	return "", "", ErrEmailAddressInvalidTooManyAts
}

func EmailUserHostToSwapped(user, host string) string {
	user = strings.TrimSpace(user)
	host = strings.TrimSpace(host)
	return host + "@@" + user
}

func EmailUserHostToReverse(user, host string) string {
	user = strings.TrimSpace(user)
	hostParts := strings.Split(strings.TrimSpace(host), ".")
	sortutil.ReverseSlice(hostParts)
	return strings.Join(hostParts, ".") + "@@@" + user
}
