package duration

import (
	"strconv"
	"strings"
	"time"
)

// DurationMilliseconds wraps time.Duration for JSON marshaling as milliseconds.
// This preserves time.Duration semantics in Go while serializing as
// human-readable millisecond integers in JSON.
//
// Example usage:
//
//	type Segment struct {
//	    Start DurationMilliseconds `json:"start_ms"`
//	    End   DurationMilliseconds `json:"end_ms"`
//	}
//
// JSON output: {"start_ms": 1500, "end_ms": 3200}
type DurationMilliseconds time.Duration

// MarshalJSON implements json.Marshaler, encoding the duration as milliseconds.
func (d DurationMilliseconds) MarshalJSON() ([]byte, error) {
	ms := int64(time.Duration(d) / time.Millisecond)
	return []byte(strconv.FormatInt(ms, 10)), nil
}

// UnmarshalJSON implements json.Unmarshaler, decoding milliseconds to duration.
func (d *DurationMilliseconds) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "null" || s == "" {
		*d = 0
		return nil
	}
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*d = DurationMilliseconds(time.Duration(ms) * time.Millisecond)
	return nil
}

// Duration returns the underlying time.Duration value.
func (d DurationMilliseconds) Duration() time.Duration {
	return time.Duration(d)
}

// Milliseconds returns the duration as milliseconds (int64).
func (d DurationMilliseconds) Milliseconds() int64 {
	return time.Duration(d).Milliseconds()
}

// Seconds returns the duration as seconds (float64).
func (d DurationMilliseconds) Seconds() float64 {
	return time.Duration(d).Seconds()
}

// String returns the string representation of the duration.
func (d DurationMilliseconds) String() string {
	return time.Duration(d).String()
}

// FromDuration creates a DurationMilliseconds from a time.Duration.
func FromDuration(d time.Duration) DurationMilliseconds {
	return DurationMilliseconds(d)
}

// FromMilliseconds creates a DurationMilliseconds from milliseconds.
func FromMilliseconds(ms int64) DurationMilliseconds {
	return DurationMilliseconds(time.Duration(ms) * time.Millisecond)
}
