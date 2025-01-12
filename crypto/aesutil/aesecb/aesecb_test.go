package aesecb

import (
	"fmt"
	"testing"
)

var aesecbTests = []struct {
	plaintext string
	key       string
}{
	{"My Secret Password", "0123456789abcdef"},
	{"My Secret Password", "1234567890123456"},
}

func TestEncryptDecrypt(t *testing.T) {
	for _, tt := range aesecbTests {
		enc, err := EncryptBase64(tt.plaintext, tt.key)
		if err != nil {
			t.Errorf("aesecb.EncryptBase64 error(%s)", err.Error())
		}
		dec, err := DecryptBase64(enc, tt.key)
		if err != nil {
			t.Errorf("aesecb.DecryptBase64 error(%s)", err.Error())
		}
		if dec != tt.plaintext {
			fmt.Printf("[%v]\n", []byte(dec))
			fmt.Printf("[%v]\n", []byte(tt.plaintext))
			t.Errorf("encrypt/decrypt AES ECB error: want decrypted (%s), got (%s)", tt.plaintext, dec)
		}
	}
}
