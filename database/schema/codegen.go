package schema

import (
	"fmt"
	"sort"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
)

// SchemaCodegen generates Ent schema lines for enum values where the
// count of unique values is less than the provided `enumLenLTE` value.
// The keys to the map are expected to be column names.
func SchemaCodegen(v map[string][]string, enumLenLTE int, schemaLineFunc func(colName string, vals []string) (string, error)) ([]string, map[string]int, error) {
	var cols []string
	colsExcl := map[string]int{}
	colNames := maputil.Keys(v)
	for _, colName := range colNames {
		if vals, ok := v[colName]; !ok {
			continue
		} else if len(vals) > enumLenLTE {
			colsExcl[colName] = len(vals)
		} else if typeLine, err := schemaLineFunc(colName, vals); err != nil {
			return []string{}, nil, err
		} else {
			cols = append(cols, typeLine)
		}
	}
	return cols, colsExcl, nil
}

func SchemaCodegenEnt(v map[string][]string, enumLenLTE int) ([]string, map[string]int, error) {
	return SchemaCodegen(v, enumLenLTE, SchemaTypeEnt)
}

func SchemaCodegenMySQL(v map[string][]string, enumLenLTE int) ([]string, map[string]int, error) {
	return SchemaCodegen(v, enumLenLTE, SchemaTypeMySQL)
}

func SchemaTypeEnt(colName string, vals []string) (string, error) {
	vals = slicesutil.Dedupe(vals)
	sort.Strings(vals)
	j, err := jsonutil.MarshalSlice(slicesutil.Dedupe(vals), true)
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf(`field.Enum("%s").Values(%s)`, colName, string(j)), nil
	}
}

func SchemaTypeMySQL(colName string, vals []string) (string, error) {
	vals = slicesutil.Dedupe(vals)
	sort.Strings(vals)
	j, err := jsonutil.MarshalSlice(slicesutil.Dedupe(vals), true)
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf(`%s ENUM(%s)`, colName, string(j)), nil
	}
}
