package ioutil

import (
	"bufio"
	"bytes"
	"io"
)

/*
type FileType int

const (
	File FileType = iota
	Directory
	Any
)
*/

/*
func CopyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := w.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	err = w.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, si.Mode())
}
*/

/*
// ReadDirSplit returnsa slides of `os.FileInfo` for directories and files.
// Note: this isn't as necessary any more since `os.ReadDir()` returns a slice of
// `os.DirEntry{}` which has a `IsDir()` func.
func ReadDirSplit(dirname string, inclDotDirs bool) ([]os.FileInfo, []os.FileInfo, error) {
	allDEs, err := os.ReadDir(dirname)
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}
	allFIs, err := osutil.DirEntriesToFileInfos(allDEs)
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}
	subdirs, regular := FileInfosSplit(allFIs, inclDotDirs)
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
*/

/*
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
*/

/*
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
*/

/*
func RemoveAllChildren(dir string) error {
	isDir, err := osutil.IsDir(dir)
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
*/

/*
func FileinfosNames(fis []os.FileInfo) []string {
	s := []string{}
	for _, e := range fis {
		s = append(s, e.Name())
	}
	return s
}
*/

// ReaderToBytes reads from an io.Reader, e.g. io.ReadCloser
func ReaderToBytes(r io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

/*
// ReadFileJSON reads and unmarshals a file.
func ReadFileJSON(file string, v interface{}) error {
	bytes, err := os.ReadFile(file)
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
	return os.WriteFile(filepath, bytes, perm)
}

func CloseFileWithError(file *os.File, err error) error {
	errFile := file.Close()
	if err != nil {
		return errorsutil.Wrap(err, errFile.Error())
	}
	return err
}
*/

/*
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
*/

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

// Write writes from `Writer` to a `Reader`. See `osutil.WriteFileReader()`.
func Write(w *bufio.Writer, r io.Reader) error {
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}
	return w.Flush()
}
