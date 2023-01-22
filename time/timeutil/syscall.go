package timeutil

import (
	"syscall"
	"time"
)

func Timespec(t syscall.Timespec) time.Time {
	return time.Unix(t.Sec, t.Nsec)
}

func Timeval(t syscall.Timeval) time.Time {
	return time.Unix(0, t.Nano())
}
