//go:build linux || darwin
// +build linux darwin

package osutil

import (
	"path/filepath"
	"sort"
	"syscall"
)

func FileStatT(filename string) (syscall.Stat_t, error) {
	var stat syscall.Stat_t
	return stat, syscall.Stat(filename, &stat)
}

func MustFileStatT(filename string) syscall.Stat_t {
	stat, err := FileStatT(filename)
	if err != nil {
		panic(err)
	}
	return stat
}

// SortDirEntriesBirthtimeSec sorts `DirEntries` by `Stat_t.Birthtimespec` which
// is available on OS-X but not all systems. It will panic if a entry is not found.
func SortDirEntriesBirthtimeSec(dir string, entries DirEntries) {
	sort.Slice(entries, func(i, j int) bool {
		return MustFileStatT(filepath.Join(dir, entries[i].Name())).Birthtimespec.Sec <
			MustFileStatT(filepath.Join(dir, entries[j].Name())).Birthtimespec.Sec
	})
}
