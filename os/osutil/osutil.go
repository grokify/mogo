package osutil

import (
	"io/ioutil"
	"os"
)

func EmptyAll(path string) error {
	aEntries, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range aEntries {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}
		err = os.Remove(path + "/" + f.Name())
		if err != nil {
			return err
		}
	}
	return nil
}
