package har

import (
	"encoding/json"
	"io/ioutil"
)

func ReadLogFile(filename string) (Log, error) {
	h := Log{}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return h, err
	}

	return h, json.Unmarshal(bytes, &h)
}
