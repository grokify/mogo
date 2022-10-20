package filepathutil

// $ go get github.com/grokify/mogo/git/apps/gitremoteaddupstream

import (
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

// FilepathLeaf returns the last element of a path.
func FilepathLeaf(s string) string {
	_, file := filepath.Split(s)
	return file
}

// CurLeafDir returns the leaf dir of a nested directory.
func CurLeafDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dirParts := strings.Split(dir, string(os.PathSeparator))
	if len(dirParts) > 0 {
		return dirParts[len(dirParts)-1], nil
	}
	return "", nil
}

// UserToAbsolute converts ~ directories to absolute directories
// in filepaths. This is useful because ~ cannot be resolved by
// ioutil.ReadFile().
func UserToAbsolute(file string) (string, error) {
	file = strings.TrimSpace(file)
	parts := strings.Split(file, string(os.PathSeparator))
	if len(parts) == 0 {
		return file, nil
	} else if parts[0] != "~" {
		return file, nil
	}
	usr, err := user.Current()
	if err != nil {
		return file, err
	}
	if len(parts) == 1 {
		return usr.HomeDir, nil
	}

	return strings.Join(
		append([]string{usr.HomeDir}, parts[1:]...),
		string(os.PathSeparator)), nil
}

// Trim trims the provided filepath using `os.PathSeparator`
func Trim(file string) string { return strings.Trim(file, string(os.PathSeparator)) }

// TrimLeft left trims the provided filepath using `os.PathSeparator`
func TrimLeft(file string) string { return strings.TrimLeft(file, string(os.PathSeparator)) }

// TrimRight right trims the provided filepath using `os.PathSeparator`
func TrimRight(file string) string { return strings.TrimRight(file, string(os.PathSeparator)) }

var rxExt = regexp.MustCompile(`\.[^/.]*$`)

// TrimExt removes the extension, including period.
func TrimExt(path string) string {
	return rxExt.ReplaceAllString(path, "")
}

// FilterFilepaths filters a slice of filepaths using various options.
func FilterFilepaths(paths []string, inclExists, inclNotExists, inclFiles, inclDirs bool) []string {
	filtered := []string{}
	for _, path := range paths {
		exists := true
		fi, err := os.Stat(path)
		if os.IsNotExist(err) {
			exists = false
		} else if err != nil {
			continue
		}
		if !(inclExists && inclNotExists) &&
			((!inclExists && exists) || (!inclNotExists && !exists)) {
			continue
		}
		if !(inclFiles && inclDirs) {
			if (!inclFiles && fi.Mode().IsRegular()) ||
				(!inclDirs && fi.Mode().IsDir()) {
				continue
			}
		}
		filtered = append(filtered, path)
	}
	return filtered
}
