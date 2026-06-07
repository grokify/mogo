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

// var rxExt = regexp.MustCompile(`\.[^.]+$`)

// TrimExt removes the extension, including a trailing period.
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

func ChangeStripExtension(fp, newext string) string {
	if newext != "" {
		if strings.Index(newext, ".") != 0 {
			newext = "." + newext
		}
	}
	if strings.Contains(fp, ".") {
		fp = rxExt.ReplaceAllString(fp, newext)
	} else {
		fp += newext
	}
	return fp
}

// NormalizePath converts a file path to use forward slashes for consistent
// cross-platform comparison. This is useful when comparing paths that may
// come from different sources (e.g., Unix-style paths in config files on Windows).
//
// The path is first cleaned using filepath.Clean, then all backslashes are
// converted to forward slashes.
func NormalizePath(p string) string {
	p = filepath.Clean(p)
	return strings.ReplaceAll(p, "\\", "/")
}

// IsAbsAny returns true if the path is absolute on any platform.
// Unlike filepath.IsAbs, this function recognizes both Unix-style paths
// (starting with /) and Windows-style paths (starting with drive letter like C:\)
// regardless of the current operating system.
//
// This is useful when handling paths that may come from different platforms,
// such as in cross-platform tools, SARIF output, or configuration files.
//
// Examples:
//
//	IsAbsAny("/Users/test/file.txt")     // true (Unix absolute)
//	IsAbsAny("C:\\Users\\test\\file.txt") // true (Windows absolute)
//	IsAbsAny("C:/Users/test/file.txt")    // true (Windows with forward slashes)
//	IsAbsAny("relative/path.txt")         // false
//	IsAbsAny("./relative/path.txt")       // false
func IsAbsAny(path string) bool {
	if path == "" {
		return false
	}
	// Unix-style absolute path
	if path[0] == '/' {
		return true
	}
	// Windows-style absolute path (drive letter)
	if len(path) >= 2 && path[1] == ':' {
		c := path[0]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			return true
		}
	}
	return false
}
