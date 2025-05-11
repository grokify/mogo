package timeutil

import (
	"fmt"

	"github.com/grokify/mogo/type/ordered"
)

func MinInt32(ints []int32) int32 {
	min := int32(0)
	init := false
	for _, this := range ints {
		if !init {
			min = this
			init = true
			continue
		} else if this < min {
			min = this
		}
	}
	return min
}

func MaxInt32(ints []int32) int32 {
	max := int32(0)
	init := false
	for _, this := range ints {
		if !init {
			max = this
			init = true
			continue
		} else if this > max {
			max = this
		}
	}
	return max
}

func YearQuarterTimeline(yyyyqs []int) ([]int, error) {
	qtrs := []int{}
	min, max := ordered.MinMax(yyyyqs...)
	if min > 0 {
		qtrs = append(qtrs, min)
		this := min
		var err error
		for this < max {
			this, err = YearQuarterAdd(this, 1)
			if err != nil {
				return qtrs, err
			}
			qtrs = append(qtrs, this)
		}
	}
	return qtrs, nil
}

type QuarterTimeline struct {
	min         int
	max         int
	initialized bool
	timeline    []int
}

func (qt *QuarterTimeline) AddInit(yyyyq int) {
	qt.min = yyyyq
	qt.max = yyyyq
	qt.initialized = true
}

func (qt *QuarterTimeline) Add(yyyyq int) {
	if !qt.initialized {
		qt.AddInit(yyyyq)
		return
	}
	if yyyyq < qt.min {
		qt.min = yyyyq
	}
	if yyyyq > qt.max {
		qt.max = yyyyq
	}
}

func (qt *QuarterTimeline) Inflate() error {
	timeline, err := qt.Timeline()
	if err != nil {
		return err
	}
	qt.timeline = timeline
	return nil
}

func (qt *QuarterTimeline) Timeline() ([]int, error) {
	return YearQuarterTimeline([]int{qt.min, qt.max})
}

func (qt *QuarterTimeline) TimelineIndex(yyyyq int) (int, error) {
	if len(qt.timeline) == 0 {
		if err := qt.Inflate(); err != nil {
			return 0, err
		}
	}
	for i, try := range qt.timeline {
		if try == yyyyq {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Quarter not found [%v] MIN[%v] MAX[%v]", yyyyq, qt.min, qt.max)
}
