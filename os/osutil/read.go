package osutil

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/grokify/mogo/errors/errorsutil"
)

func ReadFileByLine(name string, lineFunc func(idx uint, line string) error) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	i := uint(0)
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		txt := scanner.Text()
		err := lineFunc(i, txt)
		if err != nil {
			return err
		}
		i++
	}

	return scanner.Err()
}

// ReadFileJSON reads and unmarshals a file.
func ReadFileJSON(file string, v any) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func CloseFileWithError(file *os.File, err error) error {
	errFile := file.Close()
	if err != nil {
		return errorsutil.Wrap(err, errFile.Error())
	}
	return err
}
