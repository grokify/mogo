// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"errors"
	"fmt"
	"time"
)

type DurationPct struct {
	DurationStartTime time.Time
	DurationEndTime   time.Time
	DurationActive    time.Duration
	DurationRange     time.Duration
	DurationPct       float64
}

func (dp *DurationPct) Inflate() error {
	if dp.DurationRange.Seconds() == 0 {
		return errors.New("No Duration Range")
	}
	err := dp.InflateDurationActive()
	if err != nil {
		return err
	}
	dp.DurationPct = float64(dp.DurationActive.Seconds()) /
		float64(dp.DurationRange.Seconds())
	return nil
}

func (dp *DurationPct) InflateDurationActive() error {
	if dp.DurationActive.Seconds() == 0 {
		fmt.Println(dp.DurationStartTime)
		fmt.Println(dp.DurationEndTime)
		dp.DurationActive = dp.DurationEndTime.Sub(dp.DurationStartTime)
	}
	return nil
}

type ImpactPct struct {
	ImpactNum int
	TotalNum  int
	ImpactPct float64
}

func (ip *ImpactPct) Inflate() error {
	if ip.TotalNum == 0 {
		return errors.New("Total num is zero")
	}
	ip.ImpactPct = float64(ip.ImpactNum) / float64(ip.TotalNum)
	return nil
}

type Event struct {
	DurationPct    DurationPct
	ImpactPct      ImpactPct
	EventImpactPct float64
	EventUptimePct float64
}

func (e *Event) Inflate() error {
	err := e.DurationPct.Inflate()
	if err != nil {
		return err
	}
	err = e.ImpactPct.Inflate()
	if err != nil {
		return err
	}
	e.EventImpactPct = e.DurationPct.DurationPct * e.ImpactPct.ImpactPct
	e.EventUptimePct = 1 - e.EventImpactPct
	return nil
}
