package secretboxutil

import (
	"encoding/hex"
	"testing"
)

var secretKey = "6368616e676520746869732070617373776f726420746f206120736563726574"

var roundTripTests = []struct {
	v    string
	want string
}{
	{"Hello World!", "Hello World!"}}

func TestRoundTrip(t *testing.T) {
	secretKeyBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		t.Errorf("TestRoundTrip cannot decode key: %v", err)
	}
	for _, tt := range roundTripTests {
		enc, err := SealBox([]byte(tt.v), secretKeyBytes)
		if err != nil {
			t.Errorf("TestRoundTrip cannot SealBox: %v", err)
		}

		got, err := OpenBox(enc, secretKeyBytes)
		if err != nil {
			t.Errorf("TestRoundTrip cannot OpenBox: %v", err)
		}

		if string(got) != tt.want {
			t.Errorf("OpenBox(%v): want %v, got %v", string(enc), tt.want, string(got))
		}
	}
}

func TestRoundTripBase32(t *testing.T) {
	secretKeyBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		t.Errorf("TestRoundTrip cannot decode key: %v", err)
	}
	for _, tt := range roundTripTests {
		enc, err := SealBase32String([]byte(tt.v), secretKeyBytes)
		if err != nil {
			t.Errorf("TestRoundTrip cannot SealBox: %v", err)
		}

		got, err := OpenBase32String(enc, secretKeyBytes)
		if err != nil {
			t.Errorf("TestRoundTrip cannot OpenBox: %v", err)
		}

		if string(got) != tt.want {
			t.Errorf("OpenBox(%v): want %v, got %v", enc, tt.want, string(got))
		}
	}
}

var openBase32Tests = []struct {
	v    string
	want string
}{
	{"CJJ4WMXAJ3E5ODQZS5BHOBXA7NQWBASAIHSNPLIKOI3GUQXJ7PHIVQUPI7ZSTOA4MTKSTVD43TGPJMYJJUCQ====",
		"Hello World!"}}

func TestOpenBase32(t *testing.T) {
	secretKeyBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		t.Errorf("TestRoundTrip cannot decode key: %v", err)
	}
	for _, tt := range openBase32Tests {
		got, err := OpenBase32String(tt.v, secretKeyBytes)
		if err != nil {
			t.Errorf("TestRoundTrip cannot OpenBox: %v", err)
		}

		if string(got) != tt.want {
			t.Errorf("OpenBox(%v): want %v, got %v", tt.v, tt.want, string(got))
		}
	}
}
