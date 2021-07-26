package fsutil

import (
	"io/fs"
	"path/filepath"
)

type DirEntries []fs.DirEntry

func (de DirEntries) Names(dir string) []string {
	filenames := []string{}
	for _, entry := range de {
		if len(dir) > 0 {
			filenames = append(filenames, filepath.Join(dir, entry.Name()))
		} else {
			filenames = append(filenames, entry.Name())
		}
	}
	return filenames
}
