// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"time"
)

type XOXTimes struct {
	CurrentTime   time.Time
	CurrentStart  time.Time
	PreviousTime  time.Time
	PreviousStart time.Time
}

func QOQTimes(thisTime time.Time) XOXTimes {
	xox := XOXTimes{CurrentTime: thisTime.UTC()}
	xox.CurrentStart = QuarterStart(xox.CurrentTime)
	xox.PreviousStart = PrevQuarter(xox.CurrentStart)

	dur := xox.CurrentTime.Sub(xox.CurrentStart)
	xox.PreviousTime = xox.PreviousStart.Add(dur)
	return xox
}

func YOYTimes(thisTime time.Time) XOXTimes {
	thisTime = thisTime.UTC()
	xox := XOXTimes{
		CurrentTime:  thisTime,
		CurrentStart: YearStart(thisTime)}
	xox.PreviousStart = TimeDt4AddNYears(xox.CurrentStart, -1)

	dur := xox.CurrentTime.Sub(xox.CurrentStart)
	xox.PreviousTime = xox.PreviousStart.Add(dur)
	return xox
}
