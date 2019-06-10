package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/time/timeutil"
)

func main() {
	zones := timeutil.ZonesSystem(timeutil.DefaultZoneDirs())
	if 1 == 1 {
		zones = timeutil.ZonesPortable()
	}
	fmtutil.PrintJSON(zones)
	fmtutil.PrintJSONMin(zones)

	tz := "America/New_York"
	offset, err := timeutil.ZoneOffsetSeconds(time.Now(), tz)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("NAME [%v] OFFSET [%v]\n", tz, offset)
}
