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

func main() {
	GcdGoogle()

	a2g := phonenumber.NewAreaCodeToGeo()
	a2g.ReadData()

	dist, err := a2g.GcdAreaCodes(AREACODE_USNYC, AREACODE_USSFO)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Great circle distance %v to %v: %v\n", AREACODE_USNYC, AREACODE_USSFO, dist)
	fmt.Println("DONE")
}
