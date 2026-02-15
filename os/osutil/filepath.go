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
	} else {
		return true, nil
	}
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
	if userhomedir, err := os.UserHomeDir(); err != nil {
		panic(err)
	} else if len(subdirs) > 0 {
		return filepath.Join(userhomedir, filepath.Join(subdirs...))
	} else {
		return userhomedir
	}
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
	var filenames []string
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
	var filepathsExist []string

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
	// allFIs, err := DirEntriesToFileInfos(allDEs)
	allDEs2 := DirEntries(allDEs)
	allFIs, err := allDEs2.Infos()
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}
	subdirs, regular := FileInfosSplit(allFIs, inclDotDirs)
	return subdirs, regular, nil
}

func FileInfosSplit(all []os.FileInfo, inclDotDirs bool) ([]os.FileInfo, []os.FileInfo) {
	var subdirs []os.FileInfo
	var regular []os.FileInfo
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

// SanitizeOpts configures path sanitization behavior.
type SanitizeOpts struct {
	// MustExist requires the path to exist on the filesystem.
	MustExist bool
	// MustBeFile requires the path to be a regular file (not a directory).
	MustBeFile bool
	// MustBeDir requires the path to be a directory.
	MustBeDir bool
	// AllowedExts restricts the path to files with these extensions (e.g., ".mp4", ".srt").
	// Extensions should include the leading dot. If empty, any extension is allowed.
	AllowedExts []string
}

// SanitizePath cleans and validates a file path for safe use in subprocesses.
// It expands ~ to the user's home directory, cleans the path, resolves it to
// an absolute path, and optionally validates existence and type constraints.
//
// This function is useful for sanitizing user-provided paths before passing
// them to exec.Command or similar functions.
func SanitizePath(path string, opts *SanitizeOpts) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	// Expand ~ and convert to absolute path
	absPath, err := AbsFilepath(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Clean the path to remove . and .. components
	absPath = filepath.Clean(absPath)

	// If no options provided, just return the cleaned absolute path
	if opts == nil {
		return absPath, nil
	}

	// Validate existence
	if opts.MustExist {
		exists, err := Exists(absPath)
		if err != nil {
			return "", fmt.Errorf("failed to check path existence: %w", err)
		}
		if !exists {
			return "", fmt.Errorf("path does not exist: %s", absPath)
		}
	}

	// Validate file vs directory
	if opts.MustBeFile {
		isFile, err := IsFile(absPath, false)
		if err != nil {
			return "", fmt.Errorf("failed to check if path is file: %w", err)
		}
		if !isFile {
			return "", fmt.Errorf("path is not a file: %s", absPath)
		}
	}

	if opts.MustBeDir {
		isDir, err := IsDir(absPath)
		if err != nil {
			return "", fmt.Errorf("failed to check if path is directory: %w", err)
		}
		if !isDir {
			return "", fmt.Errorf("path is not a directory: %s", absPath)
		}
	}

	// Validate extension
	if len(opts.AllowedExts) > 0 {
		ext := strings.ToLower(filepath.Ext(absPath))
		allowed := false
		for _, allowedExt := range opts.AllowedExts {
			if strings.ToLower(allowedExt) == ext {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("invalid file extension %q, allowed: %v", ext, opts.AllowedExts)
		}
	}

	return absPath, nil
}
