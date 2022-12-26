package osutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/grokify/mogo/type/maputil"
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
	sort.Sort(DirEntries(sdirs))
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
	sort.Sort(DirEntries(sdirs))
	return sdirs[len(sdirs)-1], nil
}

// ReadDirRxSubmatch takes a directory, regular expression and boolean to
// indicate whether to include zero size files and returns the greatest of
// a single match in the regular expression.
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

func ReadDirRxSubmatchEntriesGreatest(dir string, rx *regexp.Regexp, subMatchIdx uint, inclDirs, inclFiles, inclEmptyFiles bool) ([]os.DirEntry, error) {
	entryMap, err := ReadDirRxSubmatch(dir, rx, subMatchIdx, inclDirs, inclFiles, inclEmptyFiles)
	if err != nil {
		return []os.DirEntry{}, err
	}
	if len(entryMap) == 0 {
		return []os.DirEntry{}, nil
	}
	keysSorted := maputil.StringKeys(entryMap, nil)
	greatest := keysSorted[len(keysSorted)-1]
	return entryMap[greatest], nil
}

// ReadDirRxSubmatchCaptureGreatest takes a directory, regular expression
// and returns the greatest of a single submatch in the regular expression.
func ReadDirRxSubmatchCaptureGreatest(dir string, rx *regexp.Regexp, subMatchIdx uint, inclDirs, inclFiles, inclEmptyFiles bool) (string, error) {
	keysSorted, err := ReadDirRxSubmatchCaptures(dir, rx, subMatchIdx, inclDirs, inclFiles, inclEmptyFiles)
	if err != nil {
		return "", err
	}
	greatest := keysSorted[len(keysSorted)-1]
	return greatest, nil
}

// ReadDirRxSubmatchCaptures takes a directory, regular expression and
// returns the greatest of captures from the regular expression.
func ReadDirRxSubmatchCaptures(dir string, rx *regexp.Regexp, subMatchIdx uint, inclDirs, inclFiles, inclEmptyFiles bool) ([]string, error) {
	entryMap, err := ReadDirRxSubmatch(dir, rx, subMatchIdx, inclDirs, inclFiles, inclEmptyFiles)
	if err != nil {
		return nil, err
	}
	if len(entryMap) == 0 {
		return nil, errors.New("no match for ReadDirRxSubmatchGreatestMatch")
	}
	keysSorted := maputil.StringKeys(entryMap, nil)
	return keysSorted, nil
}

// VisitPath visit a directory and all subdirectories, executing the supplied `visitFunc` on each.
func VisitPath(dir string, rx *regexp.Regexp, inclDirs, inclFiles, inclEmptyFiles bool, visitFunc func(dir string) error) error {
	err := visitFunc(dir)
	if err != nil {
		return err
	}
	entries, err := ReadDirMore(dir, rx, inclDirs, inclFiles, inclEmptyFiles)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err := VisitPath(filepath.Join(dir, entry.Name()), rx, inclDirs, inclFiles, inclEmptyFiles, visitFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func VisitFiles(name string, visitFunc func(dir string, info fs.FileInfo) error) error {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return err
	}
	dir, _ := filepath.Split(name)
	err = visitFunc(dir, fileInfo)
	if err != nil {
		return err
	}
	fileMode := fileInfo.Mode()
	if !fileMode.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(name)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err := VisitFiles(filepath.Join(name, entry.Name()), visitFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func DirRemoveAllChildren(dir string) error {
	isDir, err := IsDir(dir)
	if err != nil {
		return err
	}
	if !isDir {
		err = errors.New("400: Path Is Not Directory")
		return err
	}
	filesAll, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range filesAll {
		if fi.Name() == "." || fi.Name() == ".." {
			continue
		}
		filepath := path.Join(dir, fi.Name())
		if fi.IsDir() {
			err = os.RemoveAll(filepath)
			if err != nil {
				return err
			}
		} else {
			err = os.Remove(filepath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
