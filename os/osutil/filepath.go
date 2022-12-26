package osutil

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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
	if len(parts) > 0 {
		if strings.TrimSpace(parts[0]) == "." {
			parts = parts[1:]
		}
		partsPath = filepath.Join(parts...)
		if strings.TrimSpace(partsPath) == "." {
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

// Filenames returns a list of filenames for files only (no directories). If a directory
// is provided it will return a list of filenames in that directory. If a `Regexp` or
// `inclEmptyFiles` params are provided, those will be use to filter the output.
func Filenames(name string, rx *regexp.Regexp, inclEmptyFiles, absPath bool) ([]string, error) {
	isFile, err := IsFile(name, !inclEmptyFiles)
	if err != nil {
		return []string{}, err
	}
	if isFile {
		if rx == nil || rx.MatchString(name) {
			if absPath {
				nameAbs, err := AbsFilepath(name)
				if err != nil {
					return []string{}, err
				}
				name = nameAbs
			}
			return []string{name}, nil
		} else {
			return []string{}, nil
		}
	}
	filenames := []string{}
	entries, err := ReadDirMore(name, rx, false, true, inclEmptyFiles)
	if err != nil {
		return []string{}, nil
	}
	for _, entry := range entries {
		filename := filepath.Join(name, entry.Name())
		if absPath {
			filenameAbs, err := AbsFilepath(filename)
			if err != nil {
				return []string{}, err
			}
			filename = filenameAbs
		}
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)
	return filenames, nil
}

func FilenamesFilterSizeGTZero(filepaths ...string) []string {
	filepathsExist := []string{}

	for _, envPathVal := range filepaths {
		envPathVals := strings.Split(envPathVal, ",")
		for _, envPath := range envPathVals {
			envPath = strings.TrimSpace(envPath)

			if isFile, err := IsFile(envPath, true); isFile && err == nil {
				filepathsExist = append(filepathsExist, envPath)
			}
		}
	}
	return filepathsExist
}

func SplitBetter(path string) (dir, file string) {
	isDir, err := IsDir(path)
	if err != nil && isDir {
		return dir, ""
	}
	return filepath.Split(path)
}

func SplitBest(path string) (dir, file string, err error) {
	isDir, err := IsDir(path)
	if err != nil {
		return "", "", err
	} else if isDir {
		return path, "", nil
	}
	isFile, err := IsFile(path, false)
	if err != nil {
		return "", "", err
	} else if isFile {
		dir, file := filepath.Split(path)
		return dir, file, nil
	}
	return "", "", fmt.Errorf("path is valid but not file or directory: [%v]", path)
}

// ReadDirSplit returnsa slides of `os.FileInfo` for directories and files.
// Note: this isn't as necessary any more since `os.ReadDir()` returns a slice of
// `os.DirEntry{}` which has a `IsDir()` func.
func ReadDirSplit(dirname string, inclDotDirs bool) ([]os.FileInfo, []os.FileInfo, error) {
	allDEs, err := os.ReadDir(dirname)
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}
	allFIs, err := DirEntriesToFileInfos(allDEs)
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}
	subdirs, regular := FileInfosSplit(allFIs, inclDotDirs)
	return subdirs, regular, nil
}

func FileInfosSplit(all []os.FileInfo, inclDotDirs bool) ([]os.FileInfo, []os.FileInfo) {
	subdirs := []os.FileInfo{}
	regular := []os.FileInfo{}
	for _, f := range all {
		if f.Mode().IsDir() {
			if f.Name() == "." && f.Name() == ".." {
				if inclDotDirs {
					subdirs = append(subdirs, f)
				}
			} else {
				subdirs = append(subdirs, f)
			}
		} else {
			regular = append(regular, f)
		}
	}
	return subdirs, regular
}

func DirFromPath(path string) (string, error) {
	path = strings.TrimRight(path, "/\\")
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	isFile := false
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return path, nil
	case mode.IsRegular():
		isFile = true
	}
	if !isFile {
		return "", nil
	}
	rx1 := regexp.MustCompile(`^(.+)[/\\][^/\\]+`)
	rs1 := rx1.FindStringSubmatch(path)
	dir := ""
	if len(rs1) > 1 {
		dir = rs1[1]
	}
	return dir, nil
}
