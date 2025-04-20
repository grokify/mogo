package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/time/timezone"
)

func main() {
	zones := timezone.ZonesSystem(timezone.DefaultZoneDirs())
	if 1 == 1 {
		zones = timezone.ZonesPortable()
	}
	fmtutil.MustPrintJSON(zones)
	if err := fmtutil.PrintJSONMin(zones); err != nil {
		log.Fatal(err)
	}

	tz := "America/New_York"
	offset, err := timezone.ZoneOffsetSeconds(time.Now(), tz)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("NAME [%v] OFFSET [%v]\n", tz, offset)
}
