package bytesutil

import (
	"bytes"

	"github.com/grokify/mogo/math/mathutil"
)

const UTF8BOM = "\xef\xbb\xbf"

func TrimUTF8BOM(data []byte) []byte {
	return bytes.TrimPrefix(data, []byte(UTF8BOM))
}

func BytesToInt(s []byte) int {
	var res int
	for _, v := range s {
		res <<= 8
		res |= int(v)
	}
	return res
}

// https://stackoverflow.com/questions/48178008/convert-byte-slice-to-int-slice
func BytesToInts(bytes []byte) []int {
	ints := []int{}
	for _, b := range bytes {
		ints = append(ints, int(b))
	}
	return ints
}

// https://stackoverflow.com/questions/48178008/convert-byte-slice-to-int-slice
func BytesToIntsMore(bytes []byte, intLength int) []int {
	//ints := make([]int, len(bytes))
	ints := []int{}
	curNum := []byte{}
	for i, b := range bytes {
		curNum = append(curNum, b)
		if intLength > 0 {
			if mod := mathutil.ModPyInt(i, intLength); mod == intLength-1 {
				ints = append(ints, BytesToInt(curNum))
				curNum = []byte{}
			}
		}
	}
	if len(curNum) > 0 {
		ints = append(ints, BytesToInt(curNum))
	}
	return ints
}
