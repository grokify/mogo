// Package osutil implements some OS utility functions.
package osutil

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// EmptyAll will delete all contents of a directory, leaving
// the provided directory. This is different from os.Remove
// which also removes the directory provided.
func EmptyAll(name string) error {
	aEntries, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}
	for _, f := range aEntries {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}
		err = os.Remove(name + "/" + f.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// FileModAge returns a time.Duration representing the age
// of the named file from FileInfo.ModTime().
func FileModAge(name string) (time.Duration, error) {
	stat, err := os.Stat(name)
	if err != nil {
		dur0, _ := time.ParseDuration("0s")
		return dur0, err
	}
	return time.Since(stat.ModTime()), nil
}

// FileModAgeFromInfo returns the file last modification
// age as a time.Duration.
func FileModAgeFromInfo(fi os.FileInfo) time.Duration {
	return time.Since(fi.ModTime())
}

// GetFileInfo returns an os.FileInfo from a filepath.
func GetFileInfo(path string) (os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

// AbsFilepath returns an absolute filepath, using
// the user's current / home directory if indicated in
// the filepath string.
func AbsFilepath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return path, nil
	} else if filepath.IsAbs(path) {
		return path, nil
	}
	parts := strings.Split(path, string(os.PathSeparator))
	if parts[0] != "~" {
		return filepath.Abs(path)
	}

	usr, err := user.Current()
	if err != nil {
		return path, err
	}
	if len(parts) == 1 {
		return usr.HomeDir, nil
	}

	return strings.Join(
		append([]string{usr.HomeDir}, parts[1:]...),
		string(os.PathSeparator)), nil
}

// FinfosToFilepaths returns a slice of string from a directory
// and sli=ce of `os.FileInfo`.
func FinfosToFilepaths(dir string, fis []os.FileInfo) []string {
	filepaths := []string{}
	for _, fi := range fis {
		filepaths = append(filepaths, filepath.Join(dir, fi.Name()))
	}
	return filepaths
}

// CreateFileWithLines creates a file and writes lines to it. It will
// optionally add a `lineSuffix` (e.g. `"\n"`) and use `bufio`.
func CreateFileWithLines(filename string, lines []string, lineSuffix string, useBuffer bool) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if !useBuffer {
		for _, line := range lines {
			_, err := f.WriteString(line + lineSuffix)
			if err != nil {
				return err
			}
		}
		return f.Sync()
	}
	w := bufio.NewWriter(f)
	for _, line := range lines {
		_, err := w.WriteString(line + lineSuffix)
		if err != nil {
			return err
		}
	}
	return w.Flush()
}
