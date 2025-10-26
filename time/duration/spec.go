package duration

import (
	"fmt"
	"time"
)

type Unit string

const (
	SlugYear   Unit = "year"
	SlugMonth  Unit = "month"
	SlugWeek   Unit = "week"
	SlugDay    Unit = "day"
	SlugHour   Unit = "hour"
	SlugMinute Unit = "minute"
	SlugSecond Unit = "second"

	SlugBusinessDay = "businessday"
)

type Spec struct {
	Value int64 `json:"value" yaml:"value"`
	Unit  Unit  `json:"unit" yaml:"unit"`
}

func (s Spec) Duration() (time.Duration, error) {
	switch s.Unit {
	case SlugYear:
		return time.Duration(s.Value) * Year, nil
	case SlugMonth:
		return time.Duration(s.Value) * 30 * Day, nil
	case SlugDay:
		return time.Duration(s.Value) * Day, nil
	case SlugHour:
		return time.Duration(s.Value) * time.Hour, nil
	case SlugMinute:
		return time.Duration(s.Value) * time.Minute, nil
	case SlugSecond:
		return time.Duration(s.Value) * time.Second, nil
	default:
		return 0, fmt.Errorf("unknown time duration unit (%s)", string(s.Unit))
	}
}
