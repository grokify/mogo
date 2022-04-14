package accounting

import "math"

// AnnualToMonthly uses compounding where uses 100% = 1.
func AnnualToMonthly(annual float64) float64 {
	return math.Pow(1+annual, (1.0/12.0)) - 1
}

// AnnualToQuarterly uses compounding where 100% = 1.
func AnnualToQuarterly(annual float64) float64 {
	return math.Pow(1+annual, (1.0/4.0)) - 1
}
