package main

import (
	"fmt"

	"github.com/grokify/gotilla/strconv/phonenumber"
	"github.com/kellydunn/golang-geo"
)

const (
	USNYC_AREACODE = 212
	USNYC_LAT_GOOG = 40.6976684
	USNYC_LON_GOOG = -74.2605588

	USSFO_AREACODE = 415
	USSFO_LAT_GOOG = 37.7578149
	USSFO_LON_GOOG = -122.5078121
)

func GcdGoogle() {
	p1 := geo.NewPoint(USNYC_LAT_GOOG, USNYC_LON_GOOG)
	p2 := geo.NewPoint(USSFO_LAT_GOOG, USSFO_LON_GOOG)

	dist := p1.GreatCircleDistance(p2)
	fmt.Printf("Great circle distance NYC to SFO: %v\n", dist)
}

func main() {
	GcdGoogle()

	a2g := phonenumber.NewAreaCodeToGeo()
	a2g.ReadData()

	dist, err := a2g.GcdAreaCodes(USNYC_AREACODE, USSFO_AREACODE)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Great circle distance %v to %v: %v\n", USNYC_AREACODE, USSFO_AREACODE, dist)
	fmt.Println("DONE")
}
