package ioutilmore

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/os/osutil"
	"github.com/pkg/errors"
)

type FileType int

const (
	File FileType = iota
	Directory
	Any
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

func ReadDirSplit(dirname string, inclDotDirs bool) ([]os.FileInfo, []os.FileInfo, error) {
	all, err := ioutil.ReadDir(dirname)
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}
	subdirs, regular := FileInfosSplit(all, inclDotDirs)
	return subdirs, regular, nil
}

func FileInfosSplit(all []os.FileInfo, inclDotDirs bool) ([]os.FileInfo, []os.FileInfo) {
	subdirs := []os.FileInfo{}
	regular := []os.FileInfo{}
	for _, f := range all {
		if f.Mode().IsDir() {
			if f.Name() == "." && f.Name() == ".." {
				if inclDotDirs {
					subdirs = append(subdirs, f)
				}
			} else {
				subdirs = append(subdirs, f)
			}
		} else {
			regular = append(regular, f)
		}
	}
	return subdirs, regular
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
	if !isFile {
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

func SplitBetter(path string) (dir, file string) {
	isDir, err := osutil.IsDir(path)
	if err != nil && isDir {
		return dir, ""
	}
	return filepath.Split(path)
}

func SplitBest(path string) (dir, file string, err error) {
	isDir, err := osutil.IsDir(path)
	if err != nil {
		return "", "", err
	} else if isDir {
		return path, "", nil
	}
	isFile, err := osutil.IsFile(path, false)
	if err != nil {
		return "", "", err
	} else if isFile {
		dir, file := filepath.Split(path)
		return dir, file, nil
	}
	return "", "", fmt.Errorf("path is valid but not file or directory: [%v]", path)
}

func FilterFilenamesSizeGtZero(filepaths ...string) []string {
	filepathsExist := []string{}

	for _, envPathVal := range filepaths {
		envPathVals := strings.Split(envPathVal, ",")
		for _, envPath := range envPathVals {
			envPath = strings.TrimSpace(envPath)

			if isFile, err := osutil.IsFile(envPath, true); isFile && err == nil {
				filepathsExist = append(filepathsExist, envPath)
			}
		}
	}
	return filepathsExist
}

func RemoveAllChildren(dir string) error {
	isDir, err := osutil.IsDir(dir)
	if err != nil {
		return err
	}
	if !isDir {
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

func FileinfosNames(fis []os.FileInfo) []string {
	s := []string{}
	for _, e := range fis {
		s = append(s, e.Name())
	}
	return s
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
	bytes, err := jsonutil.MarshalSimple(data, prefix, indent)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, bytes, perm)
}

func CloseFileWithError(file *os.File, err error) error {
	errFile := file.Close()
	if err != nil {
		return errors.Wrap(err, errFile.Error())
	}
	return err
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

// ReadAllOrError will successfully return the data
// or return the error in the value return value.
// This is useful to simply test scripts where the
// data is printed for debugging or testing.
func ReadAllOrError(r io.Reader) []byte {
	data, err := io.ReadAll(r)
	if err != nil {
		return []byte(err.Error())
	}
	return data
}
