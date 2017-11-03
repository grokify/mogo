// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"time"
)

// QuarterProjection takes a time and numeric value, estimating the
// value at the end of the quarter using a straight-line projection.
func QuarterProjection(dt time.Time, current float64) float64 {
	qtStart := QuarterStart(dt)
	durQ2D := dt.Sub(qtStart)
	qtNext := TimeDt6AddNMonths(qtStart, 3)
	durQtr := qtNext.Sub(qtStart)

	projection := current / durQ2D.Seconds() * durQtr.Seconds()
	return projection
}

/*

Gap
Actual
Run Rate
Target

Target
Actual
Run Rate
Gap

Target
Shortfall


Current



*/
