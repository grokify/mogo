package osutil

import (
	"io/ioutil"
	"os"
)

// checks whether a given filepath exists, file or dir

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

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
