package maputil

import (
	"fmt"
	"strings"
)

type MapStringFloat32 map[string]float32

func (m MapStringFloat32) String(layout, sep1, sep2 string) string {
	var parts []string
	for k, v := range m {
		parts = append(parts, strings.Join(
			[]string{
				k,
				fmt.Sprintf(layout, v),
			}, sep1))
	}
	return strings.Join(parts, sep2)
}
