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
)

func Copy(src string, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
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

// http://stackoverflow.com/questions/8824571/golang-determining-whether-file-points-to-file-or-directory
func DirAndFileFromPath(path string) (string, string, error) {
	path = strings.TrimRight(path, "/\\")
	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", "", err
	}
	isFile := false
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return "", "", errors.New("Path is Dir")
	case mode.IsRegular():
		isFile = true
	}
	if isFile == false {
		return "", "", errors.New("Is not a file")
	}
	rx1 := regexp.MustCompile(`^(.+)[/\\]([^/\\]+)`)
	rs1 := rx1.FindStringSubmatch(path)
	dir := ""
	file := ""
	if len(rs1) > 1 {
		dir = rs1[1]
		file = rs1[2]
	}
	return dir, file, nil
}

// DirFilesSubmatchGreatest takes a directory, regular expression and boolean to indicate
// whether to include zero size files and returns the greatest of a single match in the
// regular expression.
func DirFilesSubmatchGreatest(dir string, rx1 *regexp.Regexp, nonZeroFilesOnly bool) (string, error) {
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

func GetFileInfo(path string) (os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

func IsDir(path string) (bool, error) {
	fi, err := GetFileInfo(path)
	if err != nil {
		return false, err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

func IsFile(path string) (bool, error) {
	fi, err := GetFileInfo(path)
	if err != nil {
		return false, err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return false, nil
	case mode.IsRegular():
		return true, nil
	}
	return false, nil
}

func IsFileWithSizeGtZero(path string) (bool, error) {
	fi, err := GetFileInfo(path)
	if err != nil {
		return false, err
	}
	if fi.Mode().IsRegular() == false {
		err = errors.New("400: file path is not a file")
		return false, err
	} else if fi.Size() <= 0 {
		return false, nil
	}
	return true, nil
}

func FilterFilenamesSizeGtZero(filepaths ...string) []string {
	filepathsExist := []string{}

	for _, envPathVal := range filepaths {
		envPathVals := strings.Split(envPathVal, ",")
		for _, envPath := range envPathVals {
			envPath = strings.TrimSpace(envPath)

			good, err := IsFileWithSizeGtZero(envPath)
			if err == nil && good {
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

func WriteJSON(filepath string, data interface{}, perm os.FileMode, wantPretty bool) error {
	bytes := []byte{}
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
	}
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

	w := bufio.NewWriter(file)

	fw.File = file
	fw.Writer = w

	return fw, nil
}

func (f *FileWriter) Close() {
	f.Writer.Flush()
	f.File.Close()
}
