package timeutil

import (
	"syscall"
	"time"
)

func TimespecToTime(ts syscall.Timespec) time.Time { return time.Unix(ts.Sec, ts.Nsec) }
