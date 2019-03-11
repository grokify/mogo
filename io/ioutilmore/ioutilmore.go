package ioutilmore

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/grokify/gotilla/type/maputil"
)

func CopyFile(src, dst string) (err error) {
	r, err := os.Open(src)
	if err != nil {
		return
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := w.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(w, r)
	if err != nil {
		return
	}

	err = w.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func ReadDirSplit(dirname string, skipDotDirs bool) ([]os.FileInfo, []os.FileInfo, error) {
	all, err := ioutil.ReadDir(dirname)
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}
	subdirs, regular := FileInfosSplit(all, skipDotDirs)
	return subdirs, regular, nil
}

func FileInfosSplit(all []os.FileInfo, skipDotDirs bool) ([]os.FileInfo, []os.FileInfo) {
	subdirs := []os.FileInfo{}
	regular := []os.FileInfo{}
	for _, f := range all {
		if f.Mode().IsDir() {
			if !skipDotDirs || (f.Name() != "." && f.Name() != "..") {
				subdirs = append(subdirs, f)
			}
		} else {
			regular = append(regular, f)
		}
	}
	return subdirs, regular
}

func DirEntriesReSizeGt0(dir string, rx1 *regexp.Regexp) ([]os.FileInfo, error) {
	filesMatch := []os.FileInfo{}
	filesAll, e := ioutil.ReadDir(dir)
	if e != nil {
		return filesMatch, e
	}
	for _, f := range filesAll {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}
		if f.Size() > int64(0) {
			rs1 := rx1.FindStringSubmatch(f.Name())
			if len(rs1) > 0 {
				filesMatch = append(filesMatch, f)
			}
		}
	}
	return filesMatch, nil
}

// DirEntriesRegexpGreatest takes a directory, regular expression and boolean to indicate
// whether to include zero size files and returns the greatest of a single match in the
// regular expression.
func DirFilesRegexpSubmatchGreatest(dir string, rx1 *regexp.Regexp, nonZeroFilesOnly bool) ([]os.FileInfo, error) {
	files := map[string][]os.FileInfo{}

	filesAll, e := ioutil.ReadDir(dir)
	if e != nil {
		return []os.FileInfo{}, e
	}
	for _, f := range filesAll {
		if f.Name() == "." || f.Name() == ".." ||
			(nonZeroFilesOnly && f.Size() <= int64(0)) {
			continue
		}

		if rs1 := rx1.FindStringSubmatch(f.Name()); len(rs1) > 1 {
			extract := rs1[1]
			if _, ok := files[extract]; !ok {
				files[extract] = []os.FileInfo{}
			}
			files[extract] = append(files[extract], f)
		}
	}
	keysSorted := maputil.StringKeysSorted(files)
	greatest := keysSorted[len(keysSorted)-1]
	return files[greatest], nil
}

// DirFilesRegexpSubmatchGreatestSubmatch takes a directory, regular expression and boolean to indicate
// whether to include zero size files and returns the greatest of a single match in the
// regular expression.
func DirFilesRegexpSubmatchGreatestSubmatch(dir string, rx1 *regexp.Regexp, nonZeroFilesOnly bool) (string, error) {
	filesAll, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	strs := []string{}
	for _, f := range filesAll {
		if nonZeroFilesOnly && f.Size() <= int64(0) {
			continue
		}
		rs1 := rx1.FindStringSubmatch(f.Name())
		if len(rs1) > 1 {
			strs = append(strs, rs1[1])
		}
	}
	sort.Strings(strs)
	if len(strs) == 0 {
		return "", fmt.Errorf("No matches found")
	}
	return strs[len(strs)-1], nil
}

func DirFromPath(path string) (string, error) {
	path = strings.TrimRight(path, "/\\")
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	isFile := false
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return path, nil
	case mode.IsRegular():
		isFile = true
	}
	if isFile == false {
		return "", nil
	}
	rx1 := regexp.MustCompile(`^(.+)[/\\][^/\\]+`)
	rs1 := rx1.FindStringSubmatch(path)
	dir := ""
	if len(rs1) > 1 {
		dir = rs1[1]
	}
	return dir, nil
}

func IsDir(name string) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsDir() {
		return false, nil
	}
	return true, nil
}

func IsFile(name string) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsRegular() {
		return false, nil
	}
	return true, nil
}

// IsFileWithSizeGtZero verifies a path exists, is a file and is not empty,
// returning an error otherwise. An os file not exists check can be done
// with os.IsNotExist(err) which acts on error from os.Stat()
func IsFileWithSizeGtZero(name string) (bool, error) {
	if fi, err := os.Stat(name); err != nil {
		return false, err
	} else if !fi.Mode().IsRegular() {
		return false, nil
		// return fmt.Errorf("Filepath [%v] exists but is not a file.", name)
	} else if fi.Size() <= 0 {
		return false, nil
		// return fmt.Errorf("Filepath [%v] exists but is empty with size [%v].", name, fi.Size())
	}
	return true, nil
}

func FilterFilenamesSizeGtZero(filepaths ...string) []string {
	filepathsExist := []string{}

	for _, envPathVal := range filepaths {
		envPathVals := strings.Split(envPathVal, ",")
		for _, envPath := range envPathVals {
			envPath = strings.TrimSpace(envPath)

			if isFile, err := IsFileWithSizeGtZero(envPath); isFile && err == nil {
				filepathsExist = append(filepathsExist, envPath)
			}
		}
	}
	return filepathsExist
}

func RemoveAllChildren(dir string) error {
	isDir, err := IsDir(dir)
	if err != nil {
		return err
	}
	if isDir == false {
		err = errors.New("400: Path Is Not Directory")
		return err
	}
	filesAll, err := ioutil.ReadDir(dir)
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

// ReaderToBytes reads from an io.Reader, e.g. io.ReadCloser
func ReaderToBytes(ior io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(ior)
	return buf.Bytes()
}

// ReadFileJSON reads and unmarshals a file.
func ReadFileJSON(file string, v interface{}) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func WriteFileJSON(filepath string, data interface{}, perm os.FileMode, prefix, indent string) error {
	var bytes []byte
	bytes, err := jsonutil.MarshalSimple(data, prefix, indent)
	if err != nil {
		return err
	}
	/*
		if wantPretty {
			bytesTry, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return err
			}
			bytes = bytesTry
		} else {
			bytesTry, err := json.Marshal(data)
			if err != nil {
				return err
			}
			bytes = bytesTry
		}*/
	return ioutil.WriteFile(filepath, bytes, perm)
}

type FileWriter struct {
	File   *os.File
	Writer *bufio.Writer
}

func NewFileWriter(path string) (FileWriter, error) {
	fw := FileWriter{}
	file, err := os.Create(path)
	if err != nil {
		return fw, err
	}

	fw.File = file
	fw.Writer = bufio.NewWriter(file)

	return fw, nil
}

func (f *FileWriter) Close() {
	f.Writer.Flush()
	f.File.Close()
}
