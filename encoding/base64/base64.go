// base64 supports Base36 encoding and decoding.
package base64

func StripPadding(str string) string {
	return strings.Replace(str, "=", "", -1)
}
