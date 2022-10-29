package har

import (
	"encoding/json"
	"os"
)

func ReadLogFile(filename string) (Log, error) {
	h := Log{}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return h, err
	}

	return h, json.Unmarshal(bytes, &h)
}
