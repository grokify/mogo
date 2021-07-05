package osutil

import (
	"os"
)

func IsDir(name string) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsDir() {
		return false, nil
	}
	return true, nil
}

func IsFile(name string) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsRegular() {
		return false, nil
	}
	return true, nil
}

// Exists checks whether the named filepath exists or not for
// a file or directory.
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// IsFileWithSizeGtZero verifies a path exists, is a file and is not empty,
// returning an error otherwise. An os file not exists check can be done
// with os.IsNotExist(err) which acts on error from os.Stat()
func IsFileWithSizeGtZero(name string) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsRegular() {
		return false, nil
		// return fmt.Errorf("Filepath [%v] exists but is not a file.", name)
	} else if fi.Size() <= 0 {
		return false, nil
		// return fmt.Errorf("Filepath [%v] exists but is empty with size [%v].", name, fi.Size())
	}
	return true, nil
}
