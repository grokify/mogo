package osutil

import (
	"sort"
)

// SortDirEntriesModTime sorts `DirEntries` by last modified time. It will panic if
// an entry cannot retrieve `FileInfo` information.
func SortDirEntriesModTime(files DirEntries) {
	sort.Slice(files, func(i, j int) bool {
		iF, err := files[i].Info()
		if err != nil {
			panic(err)
		}
		jF, err := files[j].Info()
		if err != nil {
			panic(err)
		}
		return iF.ModTime().Unix() < jF.ModTime().Unix()
	})
}
