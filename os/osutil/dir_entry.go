package osutil

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DirEntrySlice []os.DirEntry

func (slice DirEntrySlice) Len() int           { return len(slice) }
func (slice DirEntrySlice) Less(i, j int) bool { return slice[i].Name() < slice[j].Name() }
func (slice DirEntrySlice) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }
func (slice DirEntrySlice) Sort()              { sort.Sort(slice) }

// Names returns a slice of entry names. It can optionally
// add the directory path and sort the values.

func (slice DirEntrySlice) Names(dir string, sortNames bool) []string {
	if len(strings.TrimSpace(dir)) == 0 {
		dir = ""
	}
	names := []string{}
	for _, item := range slice {
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

// Infos returns a `[]os.FileInfo` slice.
func (slice DirEntrySlice) Infos() ([]os.FileInfo, error) {
	var infos []os.FileInfo
	for _, entry := range slice {
		info, err := entry.Info()
		if err != nil {
			return infos, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}
