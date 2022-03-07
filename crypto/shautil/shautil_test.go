package shautil

import (
	"testing"
)

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
