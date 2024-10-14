package durationbin

import (
	"math/big"
	"slices"
	"strconv"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/slicesutil"
)

var (
	binDays7   = Bin{Name: "7 Days", Duration: 7 * timeutil.Day}
	binDays14  = Bin{Name: "14 Days", Duration: 14 * timeutil.Day}
	binDays30  = Bin{Name: "30 Days", Duration: 14 * timeutil.Day}
	binDays60  = Bin{Name: "60 Days", Duration: 60 * timeutil.Day}
	binDays90  = Bin{Name: "90 Days", Duration: 90 * timeutil.Day}
	binDays120 = Bin{Name: "120 Days", Duration: 120 * timeutil.Day}
	binDays180 = Bin{Name: "180 Days", Duration: 180 * timeutil.Day}
	binYears1  = Bin{Name: "1 Year", Duration: 1 * timeutil.Year}
	binYears2  = Bin{Name: "2 Years", Duration: 2 * timeutil.Year}
)

type Bins []Bin

type Bin struct {
	Name     string
	Duration time.Duration
}

func BinsDefault() Bins {
	return Bins{
		binDays7,
		binDays14,
		binDays30,
		binDays60,
		binDays90,
		binDays120,
		binDays180,
		binYears1,
		binYears2}
}

func BinFromDuration(d time.Duration) Bin {
	// cycle through default bins with non-standard durations.
	defBins := BinsDefault()
	for i, bin := range defBins { // asc
		if d == bin.Duration {
			return bin
		} else if i == 0 && d < bin.Duration {
			return bin
		} else if i > 0 && d > defBins[i-1].Duration && d <= bin.Duration { // use previous duration
			return bin
		}
	}
	// calculate standard year durations.
	bg := big.NewInt(int64(d))
	md := big.NewInt(int64(0))
	yr := big.NewInt(int64(timeutil.Year))
	bg, md = bg.DivMod(bg, yr, md)
	if md.Int64() == 0 {
		return YearsBin(uint(bg.Uint64()))
	} else {
		return YearsBin(1 + uint(bg.Uint64()))
	}
}

func YearsBin(i uint) Bin {
	if i == 1 {
		return Bin{
			Name:     "1 Year",
			Duration: time.Duration(i) * timeutil.Year}
	} else {
		return Bin{
			Name:     strconv.Itoa(int(i)) + " Years",
			Duration: time.Duration(i) * timeutil.Year}
	}
}

func BuildBinsDurationLeftRight(durLeft, durRight time.Duration) Bins {
	wantReverse := false
	if durRight < durLeft {
		tmpDur := durRight
		durRight = durLeft
		durLeft = tmpDur
		wantReverse = true
	}
	var out Bins
	minBin := BinFromDuration(durLeft)
	maxBin := BinFromDuration(durRight)
	out = append(out, minBin)
	if minBin.Duration == maxBin.Duration {
		return out
	}
	defBins := BinsDefault()
	isActive := false
	for _, defBin := range defBins {
		if minBin.Duration <= defBin.Duration {
			isActive = true
			continue
		} else if !isActive {
			continue
		} else if maxBin.Duration == defBin.Duration {
			out = append(out, defBin)
			break
		} else if isActive {
			out = append(out, defBin)
		}
	}
	yr := uint(3)
	for out[len(out)-1].Duration < maxBin.Duration {
		out = append(out, YearsBin(yr))
		yr++
	}
	if wantReverse {
		slicesutil.Reverse(out)
	}
	return out
}

func BuildBinsDuration(d, extLeft, extRight time.Duration) Bins {
	if extLeft == extRight {
		if d < extRight {
			return BuildBinsDurationLeftRight(d, extRight)
		} else {
			return BuildBinsDurationLeftRight(extLeft, d)
		}
	} else if extLeft < extRight { // sort=ascemd
		if d < extLeft {
			return BuildBinsDurationLeftRight(d, extRight)
		} else if d > extRight {
			return BuildBinsDurationLeftRight(extLeft, d)
		} else {
			return BuildBinsDurationLeftRight(extLeft, extRight)
		}
	} else { // extLeft > extRight (reverse or sort=descend)
		if d > extLeft {
			return BuildBinsDurationLeftRight(d, extRight)
		} else if d < extRight {
			return BuildBinsDurationLeftRight(extLeft, d)
		} else {
			return BuildBinsDurationLeftRight(extLeft, extRight)
		}
	}
}

func BuildBinsDefault(durs []time.Duration, extLeft, extRight *time.Duration) Bins {
	durs = slicesutil.Dedupe(durs)
	if len(durs) == 0 {
		if extLeft == nil && extRight == nil {
			return Bins{}
		} else if extLeft != nil && extRight != nil {
			return BuildBinsDurationLeftRight(*extLeft, *extRight)
		} else if extLeft != nil {
			return Bins{BinFromDuration(*extLeft)}
		} else {
			return Bins{BinFromDuration(*extRight)}
		}
	} else if len(durs) == 1 {
		if extLeft == nil && extRight == nil {
			return Bins{BinFromDuration(durs[0])}
		} else if extLeft != nil && extRight != nil {
			return BuildBinsDuration(durs[0], *extLeft, *extRight)
		} else if extLeft != nil {
			return BuildBinsDurationLeftRight(*extLeft, durs[0])
		} else {
			return BuildBinsDurationLeftRight(durs[0], *extRight)
		}
	}
	slices.Sort(durs)
	durMin := durs[0]
	durMax := durs[len(durs)-1]
	if extLeft == nil && extRight == nil {
		return BuildBinsDurationLeftRight(durMin, durMax)
	} else if extLeft != nil && extRight != nil {
		if *extLeft <= *extRight { // asc
			curLeft := durMin
			curRight := durMax
			if *extLeft < curLeft {
				curLeft = *extLeft
			}
			if *extRight > curRight {
				curRight = *extRight
			}
			return BuildBinsDurationLeftRight(curLeft, curRight)
		} else { // desc, left = high, right = low
			curLeft := durMax
			curRight := durMin
			if *extLeft > curLeft {
				curLeft = *extLeft
			}
			if *extRight < curRight {
				curRight = *extRight
			}
			return BuildBinsDurationLeftRight(curLeft, curRight)
		}
	} else {
		panic("not implemented")
	}
}
