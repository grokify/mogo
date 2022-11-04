package osutil

import (
	"bufio"
	"os"
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

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
