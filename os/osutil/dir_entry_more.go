package osutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DirEntryMore struct {
	Dir      string
	DirEntry fs.DirEntry
}

type DirEntriesMore []DirEntryMore

func (entries DirEntriesMore) Rows(inclDirs, inclFiles bool, cols ...string) ([][]string, error) {
	rows := [][]string{}
	for _, em := range entries {
		if em.DirEntry.IsDir() && !inclDirs {
			continue
		} else if !em.DirEntry.IsDir() && !inclFiles {
			continue
		}
		row, err := em.Row(cols...)
		if err != nil {
			return rows, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

const (
	ColDir  = "dir"
	ColName = "name"
	ColPath = "path"
	ColSize = "size"
)

func (em DirEntryMore) Row(cols ...string) ([]string, error) {
	row := []string{}
	for _, col := range cols {
		switch strings.TrimSpace(strings.ToLower(col)) {
		case ColDir:
			row = append(row, em.Dir)
		case ColName:
			row = append(row, em.DirEntry.Name())
		case ColPath:
			row = append(row, filepath.Join(em.Dir, em.DirEntry.Name()))
		case ColSize:
			info, err := em.DirEntry.Info()
			if err != nil {
				return row, err
			}
			row = append(row, strconv.Itoa(int(info.Size())))
		}
	}
	return row, nil
}

func ReadDirFiles(dir string, inclDirs, inclFiles, recursive bool) (DirEntriesMore, error) {
	entriesMore := DirEntriesMore{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return entriesMore, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if recursive {
				subDirEntries, err := ReadDirFiles(
					filepath.Join(dir, entry.Name()),
					inclDirs, inclFiles, recursive)
				if err != nil {
					return entriesMore, err
				}
				entriesMore = append(entriesMore, subDirEntries...)
			}
			if !inclDirs {
				continue
			}
		} else if !entry.IsDir() && !inclFiles {
			continue
		}
		entriesMore = append(entriesMore,
			DirEntryMore{
				Dir:      dir,
				DirEntry: entry})
	}
	return entriesMore, nil
}
