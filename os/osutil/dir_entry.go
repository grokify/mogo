package osutil

import (
	"os"
	"sort"
)

type DirEntrySlice []os.DirEntry

func (p DirEntrySlice) Len() int           { return len(p) }
func (p DirEntrySlice) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p DirEntrySlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p DirEntrySlice) Sort()              { sort.Sort(p) }

func (p DirEntrySlice) Names() []string {
	names := []string{}
	for _, item := range p {
		names = append(names, item.Name())
	}
	return names
}
