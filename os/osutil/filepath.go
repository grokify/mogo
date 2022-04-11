package osutil

import (
	"go/build"
	"os"
	"path/filepath"
)

func IsDir(name string) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsDir() {
		return false, nil
	}
	return true, nil
}

// IsFile verifies a path exists and is a file. It will optionally
// check if a file is not empty. An os file not exists check can be
// done with os.IsNotExist(err) which acts on error from os.Stat().
func IsFile(name string, sizeGtZero bool) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsRegular() {
		return false, nil
	} else if sizeGtZero && fi.Size() <= 0 {
		return false, nil
	}
	return true, nil
}

// Exists checks whether the named filepath exists or not for a file or directory.
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func MustUserHomeDir(subdirs ...string) string {
	userhomedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	if len(subdirs) > 0 {
		return filepath.Join(userhomedir, filepath.Join(subdirs...))
	}
	return userhomedir
}

func GoPath(parts ...string) string {
	partsPath := ""
	if len(partsPath) > 0 {
		if parts[0] == "." {
			parts = parts[1:]
		}
		partsPath = filepath.Join(parts...)
		if partsPath == "." {
			partsPath = ""
		}
	}
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		if len(partsPath) > 0 {
			return filepath.Join(gopath, partsPath)
		}
		return gopath
	}
	gopath = build.Default.GOPATH
	if len(partsPath) > 0 {
		return filepath.Join(gopath, partsPath)
	}
	return gopath
}

/*
func GoPathSrc(pkgparts ...string) string {
	if len(pkgparts) > 0 {
		packagePath := filepath.Join(pkgparts...)
		if len(packagePath) > 0 && packagePath != "." {
			return filepath.Join(GoPath(), "src", packagePath)
		}
	}
	return filepath.Join(GoPath(), "src")
}
*/
