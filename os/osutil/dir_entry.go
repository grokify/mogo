package osutil

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DirEntrySlice []os.DirEntry

func (p DirEntrySlice) Len() int           { return len(p) }
func (p DirEntrySlice) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p DirEntrySlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p DirEntrySlice) Sort()              { sort.Sort(p) }

func (p DirEntrySlice) Names(dir string, sortNames bool) []string {
	if len(strings.TrimSpace(dir)) == 0 {
		dir = ""
	}
	names := []string{}
	for _, item := range p {
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
