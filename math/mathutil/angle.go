package mathutil

import "math"

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func RadiansToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}
