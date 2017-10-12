package phonenumber

import (
	"math/rand"
	"time"
)

const (
	fakeLineNumberMin = 100
	fakeLineNumberMax = 199
)

type FakeNumberGenerator struct {
	AreaCodeToGeo AreaCodeToGeo
	AreaCodes     []int
	Rand          *rand.Rand
}

func NewFakeNumberGenerator(a2g AreaCodeToGeo) FakeNumberGenerator {
	fng := FakeNumberGenerator{
		AreaCodeToGeo: a2g,
		AreaCodes:     a2g.AreaCodes(),
		Rand:          rand.New(rand.NewSource(time.Now().Unix())),
	}
	return fng
}

// RandomAreaCode generates a random area code.
func (fng *FakeNumberGenerator) RandomAreaCode() int {
	return fng.AreaCodes[fng.Rand.Intn(len(fng.AreaCodes))]
}

// RandomLineNumber generates a random line number
func (fng *FakeNumberGenerator) RandomLineNumber() int {
	return fng.Rand.Intn(fakeLineNumberMax-fakeLineNumberMin) + fakeLineNumberMin
}

// RandomLocalNumberUS returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUS() int {
	ac := fng.RandomAreaCode()
	ln := fng.RandomLineNumber()
	return 10000000000 + (ac * 10000000) + (5550000) + ln
}

// RandomLocalNumberUS returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUSUnique(set map[int]int) (int, map[int]int) {
	try := fng.RandomLocalNumberUS()
	_, ok := set[try]
	for ok {
		try := fng.RandomLocalNumberUS()
		_, ok = set[try]
	}
	set[try] = 1
	return try, set
}
