package timeutil

import (
	"errors"
	"fmt"
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

func MinMaxInt32(ints []int32) (int32, int32) {
	min := int32(0)
	max := int32(0)
	init := false
	for _, this := range ints {
		if !init {
			min = this
			max = this
			init = true
			continue
		}
		if this < min {
			min = this
		}
		if this > max {
			max = this
		}
	}
	return min, max
}

func QuarterInt32Timeline(ints []int32) ([]int32, error) {
	qtrs := []int32{}
	min, max := MinMaxInt32(ints)
	if min > 0 {
		qtrs = append(qtrs, min)
		this := min
		err := errors.New("")
		for this < max {
			this, err = NextQuarterInt32(this)
			if err != nil {
				return qtrs, err
			}
			qtrs = append(qtrs, this)
		}
	}
	return qtrs, nil
}

type QuarterTimeline struct {
	min         int32
	max         int32
	initialized bool
	timeline    []int32
}

func (qt *QuarterTimeline) AddInit(yyyyq int32) {
	qt.min = yyyyq
	qt.max = yyyyq
	qt.initialized = true
}

func (qt *QuarterTimeline) Add(yyyyq int32) {
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

func (qt *QuarterTimeline) Timeline() ([]int32, error) {
	return QuarterInt32Timeline([]int32{qt.min, qt.max})
}

func (qt *QuarterTimeline) TimelineIndex(yyyyq int32) (int, error) {
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
