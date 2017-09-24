package phonenumber

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/grokify/gotilla/encoding/csvutil"
	"github.com/kellydunn/golang-geo"
)

const (
	A2gCsvRelPath = "github.com/grokify/gotilla/strconv/phonenumber/us-area-code-geo.csv"
)

type AreaCodeInfo struct {
	AreaCode int
	Point    *geo.Point
}

// NewAreaCodeInfoStrings returns an AreaCodeInfo based on string area code,
// lat and lon values.
func NewAreaCodeInfoStrings(ac, lat, lon string) (AreaCodeInfo, error) {
	aci := AreaCodeInfo{}
	i, err := strconv.Atoi(ac)
	if err != nil {
		return aci, err
	}
	if i < 100 || i > 999 {
		return aci, errors.New(fmt.Sprintf("Invalid Area Code %v", i))
	}
	aci.AreaCode = i
	geo, err := NewPointString(lat, lon)
	if err != nil {
		return aci, err
	}
	aci.Point = geo
	return aci, nil
}

// NewPointString returns a *geo.Point based on string lat and lon values.
func NewPointString(lat string, lon string) (*geo.Point, error) {
	f1, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return geo.NewPoint(0, 0), err
	}
	f2, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return geo.NewPoint(0, 0), err
	}
	return geo.NewPoint(f1, f2), nil
}

type AreaCodeToGeo struct {
	AreaCodeInfos  map[int]AreaCodeInfo
	DistanceMatrix map[int]map[int]float64
}

func NewAreaCodeToGeo() AreaCodeToGeo {
	return AreaCodeToGeo{AreaCodeInfos: map[int]AreaCodeInfo{}}
}

func (a2g *AreaCodeToGeo) ReadData() error {
	return a2g.ReadCsvPath(A2gCsvFullPath())
}

func (a2g *AreaCodeToGeo) ReadCsvPath(csvpath string) error {
	csv, file, err := csvutil.NewReader(A2gCsvFullPath(), ',', false)
	if err != nil {
		return err
	}

	for {
		rec, errx := csv.Read()
		if errx == io.EOF {
			break
		} else if errx != nil {
			err = errx
			break
		} else if len(rec) != 3 {
			err = errors.New(fmt.Sprintf("Bad LatLon Data: %v\n", rec))
			break
		}
		aci, errx := NewAreaCodeInfoStrings(rec[0], rec[1], rec[2])
		if errx != nil {
			err = errx
			break
		}
		a2g.AreaCodeInfos[aci.AreaCode] = aci
	}
	file.Close()
	if err != nil {
		return err
	}
	a2g.Inflate()
	return nil
}

func (a2g *AreaCodeToGeo) AreaCodeSlice() []AreaCodeInfo {
	acSlice := []AreaCodeInfo{}
	for _, aci := range a2g.AreaCodeInfos {
		acSlice = append(acSlice, aci)
	}
	return acSlice
}

func (a2g *AreaCodeToGeo) Inflate() {
	a2g.DistanceMatrix = a2g.GetDistanceMatrix()
}

func (a2g *AreaCodeToGeo) GetDistanceMatrix() map[int]map[int]float64 {
	acis := a2g.AreaCodeSlice()
	distanceMatrix := map[int]map[int]float64{}

	l := len(acis)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			ac1 := acis[i]
			ac2 := acis[j]
			gcd := ac1.Point.GreatCircleDistance(ac2.Point)
			if _, ok := distanceMatrix[ac1.AreaCode]; !ok {
				distanceMatrix[ac1.AreaCode] = map[int]float64{}
			}
			distanceMatrix[ac1.AreaCode][ac2.AreaCode] = gcd
			if _, ok := distanceMatrix[ac2.AreaCode]; !ok {
				distanceMatrix[ac2.AreaCode] = map[int]float64{}
			}
			distanceMatrix[ac2.AreaCode][ac1.AreaCode] = gcd
		}
	}
	return distanceMatrix
}

// A2gCsvFullPath reads data from:
// https://github.com/ravisorg/Area-Code-Geolocation-Database
func A2gCsvFullPath() string {
	return path.Join(os.Getenv("GOPATH"), "src", A2gCsvRelPath)
}
