package shautil

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"strings"
	"testing"
)

var sumStringTests = []struct {
	sum256 string
	input  string
}{
	{"09ca7e4eaa6e8ae9c7d261167129184883644d07dfba7cbfbc4c8a2e08360d5b", "hello, world"},
	{"03675ac53ff9cd1535ccc7dfcdfa2c458c5218371f418dc136f2d19ac1fbe8a5", "Hello, World"},
}

func TestSumString(t *testing.T) {
	for _, tt := range sumStringTests {
		try256 := Sum256String(tt.input)
		if try256 != tt.sum256 {
			t.Errorf("shautil.Sum256String(\"%s\") Want [%s], Got [%v]", tt.input, tt.sum256, try256)
		}
	}
}

var readImageFileTests = []struct {
	filename                string
	sha1Hex                 string
	sha512d224Base32        string
	sha512d224Base32PadNone string
}{
	{"testdata/gopher_color.png", "1aeefe9e60eb95e3415edbeacb80273774cce060",
		"YTDXGKS6RAAED7YE5YHLB5OZQXT6GWNPUNVJD6WKATBU2===",
		"YTDXGKS6RAAED7YE5YHLB5OZQXT6GWNPUNVJD6WKATBU2"},
}

func TestShaSumFile(t *testing.T) {
	for _, tt := range readImageFileTests {
		/*
			trySha1, err := Sum1HexFile(tt.filename)
			if err != nil {
				t.Errorf("shautil.Sum1HexFile(\"%s\") Error: [%v]", tt.filename, err.Error())
			}
			if trySha1 != tt.sha1Hex {
				t.Errorf("shautil.Sum1HexFile(\"%s\") Format Want [%s], Got [%v]", tt.filename, tt.sha1Hex, trySha1)
			}
		*/
		trySha512d224Base32, err := Sum512d224Base32File(tt.filename, '=')
		if err != nil {
			t.Errorf("shautil.Sum512d224Base32File(\"%s\") Error: [%v]", tt.filename, err.Error())
		}
		if trySha512d224Base32 != tt.sha512d224Base32 {
			t.Errorf("shautil.Sum512d224Base32File(\"%s\") Format Want [%s], Got [%v]", tt.filename, tt.sha512d224Base32, trySha512d224Base32)
		}
		trySha512d224Base32PadNone, err := Sum512d224Base32File(tt.filename, -1)
		if err != nil {
			t.Errorf("shautil.Sum512d224Base32File(\"%s\") Error: [%v]", tt.filename, err.Error())
		}
		if trySha512d224Base32PadNone != tt.sha512d224Base32PadNone {
			t.Errorf("shautil.Sum512d224Base32File(\"%s\") Format Want [%s], Got [%v]", tt.filename, tt.sha512d224Base32PadNone, trySha512d224Base32PadNone)
		}
	}
}

var suha512d224CustomPadTests = []struct {
	input        string
	raw          []byte
	base32PadStd string
	base32Pad1   string
}{
	{"foobar",
		[]byte{57, 231, 201, 91, 243, 249, 45, 205, 23, 29, 69, 45, 6, 10, 61, 195, 183, 202, 151, 158, 4, 87, 241, 12, 165, 176, 228, 179},
		"HHT4SW7T7EW42FY5IUWQMCR5YO34VF46ARL7CDFFWDSLG===",
		"HHT4SW7T7EW42FY5IUWQMCR5YO34VF46ARL7CDFFWDSLG111",
	},
}

func TestSha512d224CustomPad(t *testing.T) {
	for _, tt := range suha512d224CustomPadTests {
		d := sha512.Sum512_224([]byte(tt.input))
		d2 := d[:]

		if !bytes.Equal(d2, tt.raw) {
			t.Errorf("sha512.Sum512_224 output unexpected: want [%v] got [%v]", tt.raw, d2)
		}

		b32Std := base32.StdEncoding.EncodeToString(d2)

		if b32Std != tt.base32PadStd {
			t.Errorf("base32.StdEncoding.EncodeToString output unexpected: want [%v] got [%v]", tt.base32PadStd, b32Std)
		}

		b32P1 := Sum512d224Base32String(tt.input, '1')
		if b32P1 != tt.base32Pad1 {
			t.Errorf("shautil.Sum512d224Base32String output unexpected: want [%v] got [%v]", tt.base32Pad1, b32P1)
		}

		b32P1Reader, err := Sum512d224Base32(strings.NewReader(tt.input), '1')
		if err != nil {
			t.Errorf("shautil.Sum512d224Base32(\"%s\") Error: [%v]", tt.input, err.Error())
		}
		if b32P1Reader != tt.base32Pad1 {
			t.Errorf("shautil.Sum512d224Base32 output unexpected: want [%v] got [%v]", tt.base32Pad1, b32P1Reader)
		}
	}
}
