package osutil

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/grokify/simplego/type/maputil"
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

// ReadDirRxSubmatch takes a directory, regular expression and boolean to indicate
// whether to include zero size files and returns the greatest of a single match in the
// regular expression.
func ReadDirRxSubmatch(dir string, rx *regexp.Regexp, subMatchIdx uint, inclDirs, inclFiles, inclEmptyFiles bool) (map[string][]os.DirEntry, error) {
	entryMap := map[string][]os.DirEntry{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return entryMap, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if !inclDirs {
				continue
			}
		} else if !inclFiles {
			continue
		}
		if !inclEmptyFiles && !entry.IsDir() {
			fi, err := entry.Info()
			if err != nil {
				return entryMap, err
			}
			if fi.Size() <= 0 {
				continue
			}
		}
		m := rx.FindStringSubmatch(entry.Name())
		if len(m) <= 0 {
			continue
		}
		if int(subMatchIdx) >= len(m) {
			return entryMap, fmt.Errorf("index too large for matches. matches [%d] idx [%d]", len(m), subMatchIdx)
		}
		submatch := m[int(subMatchIdx)]
		if _, ok := entryMap[submatch]; !ok {
			entryMap[submatch] = []os.DirEntry{}
		}
		entryMap[submatch] = append(entryMap[submatch], entry)
	}

	return entryMap, nil
}

func ReadDirRxSubmatchGreatestEntries(dir string, rx *regexp.Regexp, subMatchIdx uint, inclDirs, inclFiles, inclEmptyFiles bool) ([]os.DirEntry, error) {
	entryMap, err := ReadDirRxSubmatch(dir, rx, subMatchIdx, inclDirs, inclFiles, inclEmptyFiles)
	if err != nil {
		return []os.DirEntry{}, err
	}
	if len(entryMap) == 0 {
		return []os.DirEntry{}, nil
	}
	keysSorted := maputil.StringKeysSorted(entryMap)
	greatest := keysSorted[len(keysSorted)-1]
	return entryMap[greatest], nil
}

// ReadDirRxSubmatchGreatestCapture takes a directory, regular expression and returns the greatest of a single submatch in the
// regular expression.
func ReadDirRxSubmatchGreatestCapture(dir string, rx *regexp.Regexp, subMatchIdx uint, inclDirs, inclFiles, inclEmptyFiles bool) (string, error) {
	entryMap, err := ReadDirRxSubmatch(dir, rx, subMatchIdx, inclDirs, inclFiles, inclEmptyFiles)
	if err != nil {
		return "", err
	}
	if len(entryMap) == 0 {
		return "", errors.New("no match for ReadDirRxSubmatchGreatestMatch")
	}
	keysSorted := maputil.StringKeysSorted(entryMap)
	greatest := keysSorted[len(keysSorted)-1]
	return greatest, nil
}
