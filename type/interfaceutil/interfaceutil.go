package interfaceutil

func SplitSliceInterface(items []interface{}, max int) [][]interface{} {
	slices := [][]interface{}{}
	current := []interface{}{}

	for _, item := range items {
		current = append(current, item)
		if len(current) >= max {
			slices = append(slices, current)
			current = []interface{}{}
		}
	}
	if len(current) > 0 {
		slices = append(slices, current)
	}

	return slices
}
