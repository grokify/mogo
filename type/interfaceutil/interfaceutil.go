package interfaceutil

import "github.com/grokify/mogo/type/stringsutil"

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

func ToBool(value interface{}) bool {
	if value == nil {
		return false
	} else if valBool, ok := value.(bool); ok {
		return valBool
	} else if valString, ok := value.(string); ok {
		return stringsutil.ToBool(valString)
	} else if valInt, ok := value.(int); ok {
		if valInt == 0 {
			return false
		}
		return true
	} else if valFloat, ok := value.(float64); ok {
		if valFloat == 0.0 {
			return false
		}
		return true
	}
	return false
}

func ToBoolFlip(value interface{}) bool {
	if ToBool(value) {
		return false
	}
	return true
}

func ToBoolInt(value interface{}) int {
	if ToBool(value) {
		return 1
	}
	return 0
}

func ToInt(value interface{}, defaultValue int) int {
	if value == nil {
		return defaultValue
	} else if valBool, ok := value.(bool); ok {
		if valBool {
			return 1
		}
		return 0
	} else if valString, ok := value.(string); ok {
		if stringsutil.ToBool(valString) {
			return 1
		}
		return 0
	} else if valInt, ok := value.(int); ok {
		return valInt
	} else if valFloat, ok := value.(float64); ok {
		return int(valFloat)
	}
	return defaultValue
}
