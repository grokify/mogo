package timeutil

import "time"

func TimeUpdateLocation(t time.Time, z string) (time.Time, error) {
	if loc, err := time.LoadLocation(z); err != nil {
		return t, err
	} else {
		return t.In(loc), nil
	}
}
