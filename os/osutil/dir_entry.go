package osutil

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DirEntries provides utility functions for `[]os.DirEntry`. Use as
// `entries := osutil.DirEntries(slice)`.
type DirEntries []os.DirEntry

func (entries DirEntries) Len() int           { return len(entries) }
func (entries DirEntries) Less(i, j int) bool { return entries[i].Name() < entries[j].Name() }
func (entries DirEntries) Swap(i, j int)      { entries[i], entries[j] = entries[j], entries[i] }

// Sort sorts dir entries by name.
func (entries DirEntries) Sort() { sort.Sort(entries) }

// Names returns a slice of entry names. It can optionally
// add the directory path and sort the values.
func (entries DirEntries) Names(dir string, sortNames bool) []string {
	if len(strings.TrimSpace(dir)) == 0 {
		dir = ""
	}
	names := []string{}
	for _, item := range entries {
		if len(dir) == 0 {
			names = append(names, item.Name())
		} else {
			names = append(names, filepath.Join(dir, item.Name()))
		}
	}
	if sortNames {
		sort.Strings(names)
	}
	return names
}

// WriteFileNames writes a text file with filenames, one per line.
func (entries DirEntries) WriteFileNames(filename, dir string, sortNames bool, perm os.FileMode) error {
	if len(filename) == 0 {
		return errors.New("filename required")
	}
	names := entries.Names(dir, sortNames)
	return os.WriteFile(
		filename,
		[]byte(strings.Join(names, "\n")+"\n"),
		perm)
}

// Infos returns a `[]os.FileInfo` slice.
func (entries DirEntries) Infos() ([]os.FileInfo, error) {
	var infos []os.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return infos, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}
