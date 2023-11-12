package timeutil

import (
	"syscall"
	"testing"
	"time"
)

var syscallTests = []struct {
	nsec        int64
	rfc3339Nano string
}{
	{1674382312 * int64(time.Second), "2023-01-22T10:11:52Z"},
	{1674382312111222333, "2023-01-22T10:11:52.111222333Z"},
}

func TestSyscall(t *testing.T) {
	for _, tt := range syscallTests {
		ts := syscall.NsecToTimespec(tt.nsec)
		tsTime := Timespec(ts)
		tsTry := tsTime.UTC().Format(time.RFC3339Nano)

		if tsTry != tt.rfc3339Nano {
			t.Errorf("Timespec(): nsec [%d] want [%s], got [%s]", tt.nsec, tt.rfc3339Nano, tsTry)
		}

		ts2 := syscall.Timespec{Nsec: tt.nsec}
		ts2Time := Timespec(ts2)
		ts2Try := ts2Time.UTC().Format(time.RFC3339Nano)

		if ts2Try != tt.rfc3339Nano {
			t.Errorf("Timespec(): nsec [%d] want [%s], got [%s]", tt.nsec, tt.rfc3339Nano, ts2Try)
		}

		tv := syscall.NsecToTimeval(tt.nsec)
		tvTime := Timeval(tv)
		tvTry := tvTime.UTC().Format(time.RFC3339Nano)

		if tsTry != tt.rfc3339Nano {
			t.Errorf("Timeval(): nsec [%d] want [%s], got [%s]", tt.nsec, tt.rfc3339Nano, tvTry)
		}
	}
}
