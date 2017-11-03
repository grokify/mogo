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
	return fng.RandomLineNumberMinMax(fakeLineNumberMin, fakeLineNumberMax)
}

// RandomLineNumber generates a random line number
func (fng *FakeNumberGenerator) RandomLineNumberMinMax(min, max uint16) uint16 {
	return uint16(fng.Rand.Intn(int(max)-int(min))) + min
}

// RandomLocalNumberUS returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUS() uint64 {
	return fng.LocalNumberUS(fng.RandomAreaCode(), fng.RandomLineNumber())
}

// LocalNumberUS returns a US E.164 number given an areacode and line number
func (fng *FakeNumberGenerator) LocalNumberUS(ac uint16, ln uint16) uint64 {
	return 10000000000 + (uint64(ac) * 10000000) + (5550000) + uint64(ln)
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
