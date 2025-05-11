package ordered

import (
	"cmp"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Dedupe is like `slices.Dedupe()` but operates on elements conforming to `constraints.Ordered`.
func Dedupe[S ~[]E, E constraints.Ordered](s S) S {
	deduped := []E{}
	seen := map[E]int{}
	for _, val := range s {
		if _, ok := seen[val]; ok {
			continue
		}
		seen[val] = 1
		deduped = append(deduped, val)
	}
	return deduped
}

// dedupeOrderedMore is like `Dedupe()` but operates on elements conforming to `constraints.Ordered`.
func dedupeOrderedMore[S ~[]E, E constraints.Ordered](s S, seen map[E]int) (S, map[E]int) {
	deduped := []E{}
	for _, val := range s {
		if _, ok := seen[val]; ok {
			continue
		}
		seen[val] = 1
		deduped = append(deduped, val)
	}
	return deduped, seen
}

func AppendOrdered[S ~[]E, E constraints.Ordered](dedupe bool, s ...S) S {
	result := S{}
	if len(s) == 0 {
		return result
	}
	seen := map[E]int{}
	for _, si := range s {
		if dedupe {
			si, seen = dedupeOrderedMore(si, seen)
			result = append(result, si...)
		} else {
			result = append(result, si...)
		}
	}
	return result
}

func MinMax[T cmp.Ordered](s ...T) (min, max T) {
	for i, val := range s {
		if i == 0 {
			min = val
			max = val
		} else {
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
		}
	}
	return
}

// Union returns an ordered set of deduped elements.
func Union[S ~[]E, E constraints.Ordered](s ...S) S {
	union := S{}
	if len(s) == 0 {
		return union
	}
	for _, sx := range s {
		sx = Dedupe(sx)
		slices.Sort(sx)
		for _, e := range sx {
			unionAll := true
			for _, sy := range s {
				sy = Dedupe(sy)
				slices.Sort(sy)
				if slices.Index(sy, e) == -1 {
					unionAll = false
					break
				}
			}
			if unionAll {
				union = append(union, e)
			}
		}
	}
	union = Dedupe(union)
	slices.Sort(union)
	return union
}
