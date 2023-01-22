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
	xox.CurrentStart = quarterStart(xox.CurrentTime)
	xox.PreviousStart = QuarterAdd(xox.CurrentStart, -1)

	dur := xox.CurrentTime.Sub(xox.CurrentStart)
	xox.PreviousTime = xox.PreviousStart.Add(dur)
	return xox
}

func YOYTimes(t time.Time) XOXTimes {
	t = t.UTC()
	tm := NewTimeMore(t, 0)
	xox := XOXTimes{
		CurrentTime:  t,
		CurrentStart: tm.YearStart()}
	xox.PreviousStart = TimeDt4AddNYears(xox.CurrentStart, -1)

	dur := xox.CurrentTime.Sub(xox.CurrentStart)
	xox.PreviousTime = xox.PreviousStart.Add(dur)
	return xox
}
