package main

import (
	"fmt"

	"github.com/grokify/gotilla/strconv/phonenumber"
	"github.com/kellydunn/golang-geo"
)

const (
	AREACODE_USNYC = 212
	AREACODE_USSFO = 415
	USNYC_LAT_GOOG = 40.6976684
	USNYC_LON_GOOG = -74.2605588
	USSFO_LAT_GOOG = 37.7578149
	USSFO_LON_GOOG = -122.5078121
)

func GcdGoogle() {
	p1 := geo.NewPoint(USNYC_LAT_GOOG, USNYC_LON_GOOG)
	p2 := geo.NewPoint(USSFO_LAT_GOOG, USSFO_LON_GOOG)

	dist := p1.GreatCircleDistance(p2)
	fmt.Printf("Great circle distance NYC to SFO: %v\n", dist)
}

func GcdAreaCode() {
	a2g := phonenumber.NewAreaCodeToGeo()
	a2g.ReadData()

	acNYC, ok := a2g.AreaCodeInfos[AREACODE_USNYC]
	if !ok {
		panic(fmt.Sprintf("AreaCode %v Not Found.", AREACODE_USNYC))
	}
	acSFO, ok := a2g.AreaCodeInfos[AREACODE_USSFO]
	if !ok {
		panic(fmt.Sprintf("AreaCode %v Not Found.", AREACODE_USSFO))
	}

	dist2 := acNYC.Point.GreatCircleDistance(acSFO.Point)
	fmt.Printf("Great circle distance %v to %v: %v\n", AREACODE_USNYC, AREACODE_USSFO, dist2)
}

func main() {
	GcdGoogle()
	GcdAreaCode()
	fmt.Println("DONE")
}
