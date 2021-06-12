package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/time/timezone"
)

func main() {
	zones := timezone.ZonesSystem(timezone.DefaultZoneDirs())
	if 1 == 1 {
		zones = timezone.ZonesPortable()
	}
	fmtutil.PrintJSON(zones)
	fmtutil.PrintJSONMin(zones)

	tz := "America/New_York"
	offset, err := timezone.ZoneOffsetSeconds(time.Now(), tz)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("NAME [%v] OFFSET [%v]\n", tz, offset)
}
