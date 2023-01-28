package interfaceutil

import "github.com/grokify/mogo/type/stringsutil"

func SplitSliceInterface(items []any, max int) [][]any {
	slices := [][]any{}
	current := []any{}

	for _, item := range items {
		current = append(current, item)
		if len(current) >= max {
			slices = append(slices, current)
			current = []any{}
		}
	}
	if len(current) > 0 {
		slices = append(slices, current)
	}

	return slices
}

func ToBool(v any) bool {
	if v == nil {
		return false
	} else if valBool, ok := v.(bool); ok {
		return valBool
	} else if valString, ok := v.(string); ok {
		return stringsutil.ToBool(valString)
	} else if valInt, ok := v.(int); ok {
		return valInt != 0
	} else if valFloat, ok := v.(float64); ok {
		return valFloat != 0.0
	}
	return false
}

func ToBoolFlip(v any) bool {
	return !ToBool(v)
}

func ToBoolInt(v any) int {
	if ToBool(v) {
		return 1
	}
	return 0
}

func ToInt(v any, defaultValue int) int {
	if v == nil {
		return defaultValue
	} else if valBool, ok := v.(bool); ok {
		if valBool {
			return 1
		}
		return 0
	} else if valString, ok := v.(string); ok {
		if stringsutil.ToBool(valString) {
			return 1
		}
		return 0
	} else if valInt, ok := v.(int); ok {
		return valInt
	} else if valFloat, ok := v.(float64); ok {
		return int(valFloat)
	}
	return defaultValue
}
