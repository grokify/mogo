package duration

import (
	"fmt"
	"time"
)

type Spec struct {
	Value int64 `json:"value" yaml:"value"`
	Unit  Unit  `json:"unit" yaml:"unit"`
}

func (s Spec) Duration() (time.Duration, error) {
	switch s.Unit {
	case UnitYear:
		return time.Duration(s.Value) * Year, nil
	case UnitMonth:
		return time.Duration(s.Value) * 30 * Day, nil
	case UnitDay:
		return time.Duration(s.Value) * Day, nil
	case UnitHour:
		return time.Duration(s.Value) * time.Hour, nil
	case UnitMinute:
		return time.Duration(s.Value) * time.Minute, nil
	case UnitSecond:
		return time.Duration(s.Value) * time.Second, nil
	case UnitMillisecond:
		return time.Duration(s.Value) * time.Millisecond, nil
	case UnitMicrosecond:
		return time.Duration(s.Value) * time.Microsecond, nil
	case UnitNanosecond:
		return time.Duration(s.Value), nil
	default:
		return 0, fmt.Errorf("unknown time duration unit (%s)", string(s.Unit))
	}
}
