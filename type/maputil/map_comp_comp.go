package maputil

type MapCompComp[T, U comparable] map[T]U

func (m MapCompComp[T, U]) Equal(n map[T]U) bool {
	if len(m) != len(n) {
		return false
	}
	for k, vM := range m {
		vN, ok := n[k]
		if !ok || vM != vN {
			return false
		}
	}
	return true
}

// FilterMergeByMap returns key/values from m where keys exist in n.
// Additionally, add optional default for keys in n that don't exist in m.
func (m MapCompComp[T, U]) FilterMergeByMap(n map[T]U, useNOnlyVal bool, nOnlyDefault *U) map[T]U {
	out := make(map[T]U)
	for k, v := range m {
		if _, ok := n[k]; ok {
			out[k] = v
		}
	}
	if useNOnlyVal || nOnlyDefault != nil {
		for k, v := range n {
			if _, ok := m[k]; !ok {
				if useNOnlyVal {
					out[k] = v
				} else if nOnlyDefault != nil {
					out[k] = *nOnlyDefault
				}
			}
		}
	}
	return out
}

func (m MapCompComp[T, U]) ValueKeyCounts() map[U]int {
	rev := make(map[U]map[T]int)
	for k, v := range m {
		if _, ok := rev[v]; !ok {
			rev[v] = map[T]int{}
		}
		rev[v][k]++
	}
	out := make(map[U]int)
	for k, v := range rev {
		out[k] = len(v)
	}
	return out
}
