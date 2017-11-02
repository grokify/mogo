package phonenumber

import (
	"math/rand"
	"time"
)

const (
	fakeLineNumberMin = uint16(100)
	fakeLineNumberMax = uint16(199)
)

type FakeNumberGenerator struct {
	AreaCodes []uint16
	Rand      *rand.Rand
}

func NewFakeNumberGenerator(areacodes []uint16) FakeNumberGenerator {
	fng := FakeNumberGenerator{
		AreaCodes: areacodes,
		Rand:      rand.New(rand.NewSource(time.Now().Unix())),
	}
	return fng
}

// RandomAreaCode generates a random area code.
func (fng *FakeNumberGenerator) RandomAreaCode() uint16 {
	return fng.AreaCodes[fng.Rand.Intn(len(fng.AreaCodes))]
}

// RandomLineNumber generates a random line number
func (fng *FakeNumberGenerator) RandomLineNumber() uint16 {
	return uint16(fng.Rand.Intn(int(fakeLineNumberMax)-int(fakeLineNumberMin))) +
		fakeLineNumberMin
}

// RandomLocalNumberUS returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUS() uint64 {
	ac := uint64(fng.RandomAreaCode())
	ln := uint64(fng.RandomLineNumber())
	return 10000000000 + (ac * 10000000) + (5550000) + ln
}

// RandomLocalNumberUS returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUSUnique(set map[uint64]int8) (uint64, map[uint64]int8) {
	try := fng.RandomLocalNumberUS()
	_, ok := set[try]
	for ok {
		try := fng.RandomLocalNumberUS()
		_, ok = set[try]
	}
	set[try] = 1
	return try, set
}
