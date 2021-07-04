package osutil

import (
	"fmt"
	"os"
	"regexp"
	"sort"
)

func ReadSubdirs(dir string, regex *regexp.Regexp) ([]os.DirEntry, error) {
	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sdirs := []os.DirEntry{}
	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		if regex != nil && !regex.MatchString(item.Name()) {
			continue
		}
		sdirs = append(sdirs, item)
	}
	return sdirs, nil
}

func ReadSubdirMin(dir string, regex *regexp.Regexp) (os.DirEntry, error) {
	sdirs, err := ReadSubdirs(dir, regex)
	if err != nil {
		return nil, err
	}
	if len(sdirs) == 0 {
		return nil, fmt.Errorf("no subdirectories for dir [%s]", dir)
	}
	sort.Sort(DirEntrySlice(sdirs))
	return sdirs[0], nil
}

func ReadSubdirMax(dir string, regex *regexp.Regexp) (os.DirEntry, error) {
	sdirs, err := ReadSubdirs(dir, regex)
	if err != nil {
		return nil, err
	}
	if len(sdirs) == 0 {
		return nil, fmt.Errorf("no subdirectories for dir [%s]", dir)
	}
	sort.Sort(DirEntrySlice(sdirs))
	return sdirs[len(sdirs)-1], nil
}
