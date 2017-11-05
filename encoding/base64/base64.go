package base64

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
)

// RFC7617UserPass base64 encodes a user-id and password per:
// https://tools.ietf.org/html/rfc7617#section-2
func RFC7617UserPass(userid, password string) (string, error) {
	if strings.Index(userid, ":") > -1 {
		return "", fmt.Errorf(
			"RFC7617 user-id cannot include a colon (':') [%v]", userid)
	}
	userpass := strings.Join([]string{userid, password}, ":")
	return b64.StdEncoding.EncodeToString([]byte(userpass)), nil
}
