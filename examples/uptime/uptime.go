package main

import (
	"fmt"
	"time"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/time/timeutil"
)

func getIncident(t0s, t1s, format string, durRange time.Duration, impactNum, totalNum int) {
	durPct := timeutil.DurationPct{DurationRange: durRange}
	t0, err := time.Parse(format, t0s)
	if err != nil {
		panic(err)
	}
	t1, err := time.Parse(format, t1s)
	if err != nil {
		panic(err)
	}
	durPct.DurationStartTime = t0
	durPct.DurationEndTime = t1

	fmt.Println(durPct.DurationRange.String())

	impPct := timeutil.ImpactPct{
		ImpactNum: impactNum,
		TotalNum:  totalNum}

	inc := timeutil.Event{
		DurationPct: durPct,
		ImpactPct:   impPct}
	err = inc.Inflate()
	if err != nil {
		panic(err)
	}

	fmtutil.PrintJSON(inc)

}

func main() {
	dt, err := time.Parse(time.RFC3339, "2017-08-08T00:00:00Z")
	if err != nil {
		panic(err)
	}
	rangeDur := timeutil.QuarterDuration(dt)
	fmt.Println(rangeDur.String())

	getIncident(
		"2017-09-04T09:15:00-0700",
		"2017-09-05T12:11:04-0600",
		timeutil.ISO8601, rangeDur, 140, 6577)

	getIncident(
		"2017-09-29T17:09:00Z",
		"2017-10-02T07:02:00Z",
		time.RFC3339, rangeDur, 180, 6577)

	fmt.Println("DONE")
}

/*
func getInc1(t0, t1 string, format string, rangeDur time.Duration, impactNum int) (timeutil.Incident, error) {
	t10, err := time.Parse(format, t0)
	if err != nil {
		return timeutil.Incident{}, err
	}
	t11, err := time.Parse(format, t1)
	if err != nil {
		return timeutil.Incident{}, err
	}
	inc1 := timeutil.Incident{
		StartTime:     t10,
		EndTime:       t11,
		ImpactNum:     impactNum,
		TotalNum:      6577,
		RangeDuration: rangeDur,
	}

	err = inc1.Inflate()
	if err != nil {
		return timeutil.Incident{}, err
	}

	return inc1, nil
}

func main() {
	dt, err := time.Parse(time.RFC3339, "2017-08-08T00:00:00Z")
	if err != nil {
		panic(err)
	}
	rangeDur, err := timeutil.QuarterDuration(dt)
	if err != nil {
		panic(err)
	}
	fmt.Println(rangeDur.String())

	inc1, err := getInc1(
		"2017-09-04T09:15:00-0700",
		"2017-09-05T12:11:04-0600",
		timeutil.ISO8601Z4, rangeDur, 142)
	if err != nil {
		panic(err)
	}
	fmtutil.PrintJSON(inc1)
	fmt.Println(inc1.Duration.String())

	inc2, err := getInc1(
		"2017-09-29T17:09:00Z",
		"2017-10-02T07:02:00Z",
		time.RFC3339, rangeDur, 180)
	if err != nil {
		panic(err)
	}
	fmtutil.PrintJSON(inc2)
	fmt.Println(inc2.Duration.String())

	fmt.Println("DONE")
}

*/
