// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	YearSeconds            = (365 * 24 * 60 * 60) + (6 * 60 * 60)
	WeekSeconds            = 7 * 24 * 60 * 60
	DaySeconds             = 24 * 60 * 60
	MonthsEN               = `["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]`
	MillisToNanoMultiplier = 1000000
)

var days = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

type Interval int

const (
	Decade Interval = iota
	Year
	Quarter
	Month
	Week
	Day
	Hour
	Minute
	Second
	Millisecond
	Microsecond
	Nanosecond
)

var intervals = [...]string{
	"decade",
	"year",
	"quarter",
	"month",
	"week",
	"day",
	"hour",
	"minute",
	"second",
	"millisecond",
	"microsecond",
	"nanosecond",
}

func (i Interval) String() string { return intervals[i] }

func ParseInterval(src string) (Interval, error) {
	canonical := strings.ToLower(strings.TrimSpace(src))
	for i, try := range intervals {
		if canonical == try {
			return Interval(i), nil
		}
	}
	return Year, fmt.Errorf("Interval [%v] not found.", src)
}

// ParseDuration adds days (d), weeks (w), years (y).
func ParseDuration(s string) (time.Duration, error) {
	rx := regexp.MustCompile(`(?i)^\s*(-?\d+)(d|w|y)\s*$`)
	rs := rx.FindStringSubmatch(s)

	if len(rs) > 0 {
		zeroDuration, _ := time.ParseDuration("0s")
		quantity := rs[1]
		units := strings.ToLower(rs[2])
		i, err := strconv.Atoi(quantity)
		if err != nil {
			return zeroDuration, err
		}
		if units == "d" {
			s = fmt.Sprintf("%vs", i*DaySeconds)
		} else if units == "w" {
			s = fmt.Sprintf("%vs", i*WeekSeconds)
		} else if units == "y" {
			s = fmt.Sprintf("%vs", i*YearSeconds)
		} else {
			return zeroDuration, errors.New("timeutil.ParseDuration Parse Error")
		}
	}
	return time.ParseDuration(s)
}

func MustParseDuration(s string) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return dur
}

func NowDeltaDuration(d time.Duration) time.Time {
	return time.Now().Add(d)
}

func NowDeltaParseDuration(s string) (time.Time, error) {
	d, err := ParseDuration(s)
	if err != nil {
		return time.Now(), err
	}
	return time.Now().Add(d), nil
}

// TimeRange represents a time range with a max and min value.
type TimeRange struct {
	HaveMax bool
	HaveMin bool
	Max     time.Time
	Min     time.Time
}

// Insert updates a time range min and max values for a given time.
func (tr *TimeRange) Insert(t time.Time) {
	tr.InsertMax(t)
	tr.InsertMin(t)
}

// InsertMax updates a time range max value for a given time.
func (tr *TimeRange) InsertMax(t time.Time) {
	if !tr.HaveMax {
		tr.Max = t
		tr.HaveMax = true
	} else if IsGreaterThan(t, tr.Max, false) {
		tr.Max = t
	}
}

// InsertMin updates a time range min value for a given time.
func (tr *TimeRange) InsertMin(t time.Time) {
	if !tr.HaveMin {
		tr.Min = t
		tr.HaveMin = true
	} else if IsLessThan(t, tr.Min, false) {
		tr.Min = t
	}
}

// TimeForEpochMillis returns the time.Time value for an epoch
// in milliseconds
func UnixMillis(epochMillis int64) time.Time {
	return time.Unix(0, epochMillis*MillisToNanoMultiplier)
}

// Dt6ForTime returns the Dt6 value for time.Time.
func Dt6ForTime(dt time.Time) int32 {
	dt = dt.UTC()
	return int32(dt.Year()*100 + int(dt.Month()))
}

// Dt6ForDt14 returns the Dt6 value for Dt14.
func Dt6ForDt14(dt14 int64) int32 {
	dt16f := float64(dt14) / float64(1000000)
	return int32(dt16f)
}

// TimeForDt6 returns a time.Time value given a Dt6 value.
func TimeForDt6(dt6 int32) (time.Time, error) {
	return time.Parse(DT6, strconv.FormatInt(int64(dt6), 10))
}

func ParseDt6(dt6 int32) (int16, int8) {
	year := dt6 / 100
	month := int(dt6) - (int(year) * 100)
	return int16(year), int8(month)
}

func PrevDt6(dt6 int32) int32 {
	year, month := ParseDt6(dt6)
	if month == 1 {
		month = 12
		year = year - 1
	} else {
		month = month - 1
	}
	return int32(year)*100 + int32(month)
}

func NextDt6(dt6 int32) int32 {
	year, month := ParseDt6(dt6)
	if month == 12 {
		month = 1
		year += 1
	} else {
		month += 1
	}
	return int32(year)*100 + int32(month)
}

func TimeDt6AddNMonths(dt time.Time, numMonths int) time.Time {
	dt6 := Dt6ForTime(dt)
	for i := 0; i < numMonths; i++ {
		dt6 = NextDt6(dt6)
	}
	dt6NextMonth, err := TimeForDt6(dt6)
	if err != nil {
		panic(fmt.Sprintf("Cannot find next month for time: %v\n", dt.Format(time.RFC3339)))
	}
	return dt6NextMonth
}

func TimeDt6SubNMonths(dt time.Time, numMonths int) time.Time {
	dt6 := Dt6ForTime(dt)
	for i := 0; i < numMonths; i++ {
		dt6 = PrevDt6(dt6)
	}
	dt6NextMonth, err := TimeForDt6(dt6)
	if err != nil {
		panic(fmt.Sprintf("Cannot find next month for time: %v\n", dt.Format(time.RFC3339)))
	}
	return dt6NextMonth
}

func TimeDt4AddNYears(dt time.Time, numYears int) time.Time {
	return time.Date(dt.UTC().Year()+numYears, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func Dt6MinMaxSlice(minDt6 int32, maxDt6 int32) []int32 {
	if maxDt6 < minDt6 {
		tmpDt6 := maxDt6
		maxDt6 = minDt6
		minDt6 = tmpDt6
	}
	dt6Range := []int32{}
	curDt6 := minDt6
	for curDt6 < maxDt6+1 {
		dt6Range = append(dt6Range, curDt6)
		curDt6 = NextDt6(curDt6)
	}
	return dt6Range
}

// Dt8Now returns Dt8 value for the current time.
func Dt8Now() int32 {
	return Dt8ForTime(time.Now())
}

// Dt8ForString returns a Dt8 value given a layout and value to parse to time.Parse.
func Dt8ForString(layout, value string) (int32, error) {
	dt8 := int32(0)
	t, err := time.Parse(layout, value)
	if err == nil {
		dt8 = Dt8ForTime(t)
	}
	return dt8, err
}

// Dt8ForInts returns a Dt8 value for year, month, and day.
func Dt8ForInts(yyyy int, mm int, dd int) int32 {
	sDt8 := fmt.Sprintf("%04d%02d%02d", yyyy, mm, dd)
	iDt8, _ := strconv.ParseInt(sDt8, 10, 32)
	return int32(iDt8)
}

// Dt8ForTime returns a Dt8 value given a time struct.
func Dt8ForTime(t time.Time) int32 {
	u := t.UTC()
	s := u.Format(DT8)
	iDt8, _ := strconv.ParseInt(s, 10, 32)
	return int32(iDt8)
}

// TimeForDt8 returns a time.Time value given a Dt8 value.
func TimeForDt8(dt8 int32) (time.Time, error) {
	return time.Parse(DT8, strconv.FormatInt(int64(dt8), 10))
}

// DurationForNowSubDt8 returns a duartion struct between a Dt8 value and the current time.
func DurationForNowSubDt8(dt8 int32) (time.Duration, error) {
	t, err := TimeForDt8(dt8)
	if err != nil {
		var d time.Duration
		return d, err
	}
	now := time.Now()
	return now.Sub(t), nil
}

// Dt14Now returns a Dt14 value for the current time.
func Dt14Now() int64 {
	return Dt14ForTime(time.Now())
}

// Dt14ForString returns a Dt14 value given a layout and value to parse to time.Parse.
func Dt14ForString(layout, value string) (int64, error) {
	dt14 := int64(0)
	t, err := time.Parse(layout, value)
	if err == nil {
		dt14 = Dt14ForTime(t)
	}
	return dt14, err
}

// Dt8ForInts returns a Dt8 value for a UTC year, month, day, hour, minute and second.
func Dt14ForInts(yyyy int, mm int, dd int, hr int, mn int, dy int) int64 {
	sDt14 := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", yyyy, mm, dd, hr, mn, dy)
	iDt14, _ := strconv.ParseInt(sDt14, 10, 64)
	return int64(iDt14)
}

// Dt14ForTime returns a Dt14 value given a time.Time struct.
func Dt14ForTime(t time.Time) int64 {
	u := t.UTC()
	s := u.Format(DT14)
	iDt14, _ := strconv.ParseInt(s, 10, 64)
	return int64(iDt14)
}

// TimeForDt14 returns a time.Time value given a Dt14 value.
func TimeForDt14(dt14 int64) (time.Time, error) {
	return time.Parse(DT14, strconv.FormatInt(dt14, 10))
}

// Reformat a time string from one format to another
func FromTo(timeStringSrc string, fromFormat string, toFormat string) (string, error) {
	t, err := time.Parse(fromFormat, timeStringSrc)
	if err != nil {
		return "", err
	}
	timeStringOut := t.Format(toFormat)
	return timeStringOut, nil
}

func DurationStringMinutesSeconds(durationSeconds int64) (string, error) {
	if durationSeconds <= 0 {
		return "0 sec", nil
	}
	dur, err := time.ParseDuration(fmt.Sprintf("%vs", durationSeconds))
	if err != nil {
		return "", err
	}
	modSeconds := math.Mod(float64(durationSeconds), float64(60))
	if dur.Minutes() < 1 {
		return fmt.Sprintf("%v sec", modSeconds), nil
	}
	return fmt.Sprintf("%v min %v sec", int(dur.Minutes()), modSeconds), nil
}

func MonthNames() []string {
	data := []string{}
	json.Unmarshal([]byte(MonthsEN), &data)
	return data
}

// WeekStart takes a time.Time object and a week start day
// in the time.Weekday format.
func WeekStart(dt time.Time, dow time.Weekday) (time.Time, error) {
	return TimeDeltaDowInt(dt.UTC(), int(dow), -1, true, true)
}

// MonthStart returns a time.Time for the beginning of the
// month in UTC time.
func MonthStart(dt time.Time) time.Time {
	dt = dt.UTC()
	return time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// QuarterStart returns a time.Time for the beginning of the
// quarter in UTC time.
func QuarterStart(dt time.Time) time.Time {
	dt = dt.UTC()
	qm := QuarterToMonth(MonthToQuarter(uint8(dt.Month())))
	return time.Date(dt.Year(), time.Month(qm), 1, 0, 0, 0, 0, time.UTC)
}

// QuarterEnd returns a time.Time for the end of the
// quarter by second in UTC time.
func QuarterEnd(dt time.Time) time.Time {
	qs := QuarterStart(dt.UTC())
	qn := TimeDt6AddNMonths(qs, 3)
	return time.Date(qn.Year(), qn.Month(), 0, 23, 59, 59, 0, time.UTC)
}

// YearStart returns a a time.Time for the beginning of the year
// in UTC time.
func YearStart(dt time.Time) time.Time {
	return time.Date(dt.UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

func QuarterStartString(dt time.Time) string {
	dtStart := QuarterStart(dt)
	return fmt.Sprintf("%v Q%v", dtStart.Year(), MonthToQuarter(uint8(dtStart.Month())))
}

func PrevQuarter(dt time.Time) time.Time {
	return TimeDt6SubNMonths(QuarterStart(dt), 3)
}

func NextQuarter(dt time.Time) time.Time {
	return TimeDt6AddNMonths(QuarterStart(dt), 3)
}

// QuarterDuration returns a time.Duration representing the
// calendar quarter for the time provided.
func QuarterDuration(dt time.Time) time.Duration {
	start := QuarterStart(dt)
	end := NextQuarter(start)
	return end.Sub(start)
}

func IntervalStart(dt time.Time, interval Interval, dow time.Weekday) (time.Time, error) {
	switch interval.String() {
	case "year":
		return YearStart(dt), nil
	case "quarter":
		return QuarterStart(dt), nil
	case "month":
		return MonthStart(dt), nil
	case "week":
		return WeekStart(dt, dow)
	default:
		return time.Time{}, fmt.Errorf("Interval [%v] not supported in timeutil.IntervalStart.", interval)
	}
}

// MonthToQuarter converts a month to a calendar quarter.
func MonthToQuarter(month uint8) uint8 {
	return uint8(math.Ceil(float64(month) / 3))
}

// QuarterToMonth converts a calendar quarter to a month.
func QuarterToMonth(quarter uint8) uint8 {
	return quarter*3 - 2
}

func ParseWeekday(s string) (time.Weekday, error) {
	for i, day := range days {
		if strings.ToLower(strings.TrimSpace(s)) == strings.ToLower(day) {
			return time.Weekday(i), nil
		}
	}
	return time.Weekday(0), fmt.Errorf("Cannot parse weekday: %s", s)
}

// TimeMeta is a struct for holding various times related
// to a current time, including year start, quarter start,
// month start, and week start.
type TimeMeta struct {
	This         time.Time
	YearStart    time.Time
	QuarterStart time.Time
	MonthStart   time.Time
	WeekStart    time.Time
}

// NewTimeMeta returns a TimeMeta struct given `time.Time`
// and `time.Weekday` parameters.
func NewTimeMeta(dt time.Time, dow time.Weekday) (TimeMeta, error) {
	dt = dt.UTC()
	meta := TimeMeta{
		This:         dt,
		YearStart:    YearStart(dt),
		QuarterStart: QuarterStart(dt),
		MonthStart:   MonthStart(dt)}

	week, err := WeekStart(dt, dow)
	if err != nil {
		return meta, err
	}
	meta.WeekStart = week
	return meta, nil
}
