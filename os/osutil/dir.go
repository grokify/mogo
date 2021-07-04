package osutil

import (
	"fmt"
	"os"
	"regexp"
	"sort"
)

func ReadDirMore(dir string, rx *regexp.Regexp, inclDirs, inclFiles, inclEmptyFiles bool) ([]os.DirEntry, error) {
	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sdirs := []os.DirEntry{}
	for _, item := range items {
		if item.IsDir() {
			if !inclDirs {
				continue
			}
		} else if !inclFiles {
			continue
		}
		if rx != nil && !rx.MatchString(item.Name()) {
			continue
		}
		if !inclEmptyFiles && !item.IsDir() {
			fi, err := item.Info()
			if err != nil {
				return nil, err
			}
			if fi.Size() <= 0 {
				continue
			}
		}
		sdirs = append(sdirs, item)
	}
	return sdirs, nil
}

func ReadSubdirMin(dir string, rx *regexp.Regexp) (os.DirEntry, error) {
	sdirs, err := ReadDirMore(dir, rx, true, false, false)
	if err != nil {
		return nil, err
	}
	if len(sdirs) == 0 {
		return nil, fmt.Errorf("no subdirectories for dir [%s]", dir)
	}
	sort.Sort(DirEntrySlice(sdirs))
	return sdirs[0], nil
}

func ReadSubdirMax(dir string, rx *regexp.Regexp) (os.DirEntry, error) {
	sdirs, err := ReadDirMore(dir, rx, true, false, false)
	if err != nil {
		return nil, err
	}
	if len(sdirs) == 0 {
		return nil, fmt.Errorf("no subdirectories for dir [%s]", dir)
	}
	sort.Sort(DirEntrySlice(sdirs))
	return sdirs[len(sdirs)-1], nil
}
