package timeutil

import (
	"syscall"
	"testing"
	"time"
)

var syscallTests = []struct {
	nsec    int64
	rfc3339 string
}{
	{1674382312 * NanosPerSecond, "2023-01-22T10:11:52Z"},
}

func TestSyscall(t *testing.T) {
	for _, tt := range syscallTests {
		ts := syscall.NsecToTimespec(tt.nsec)
		tsTime := Timespec(ts)
		tsTry := tsTime.UTC().Format(time.RFC3339)

		if tsTry != tt.rfc3339 {
			t.Errorf("Timespec(): nsec [%d] want [%s], got [%s]", tt.nsec, tt.rfc3339, tsTry)
		}

		tv := syscall.NsecToTimeval(tt.nsec)
		tvTime := Timeval(tv)
		tvTry := tvTime.UTC().Format(time.RFC3339)

		if tsTry != tt.rfc3339 {
			t.Errorf("Timeval(): nsec [%d] want [%s], got [%s]", tt.nsec, tt.rfc3339, tvTry)
		}
	}
}
