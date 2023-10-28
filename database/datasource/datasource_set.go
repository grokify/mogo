package datasource

import (
	"fmt"

	"github.com/grokify/mogo/type/maputil"
)

type DataSourceSet struct {
	Data map[string]DataSource
}

func (dss DataSourceSet) GetDataSource(key string) (DataSource, error) {
	if ds, ok := dss.Data[key]; ok {
		return ds, nil
	}
	return DataSource{}, fmt.Errorf("key not found (%s)", key)
}

func (dss DataSourceSet) Keys() []string {
	return maputil.Keys(dss.Data)
}
