package osutil

import (
	"os"
	"time"
)

// FileInfoMore provides a struct hold FileInfo with additional information.
type FileInfoMore struct {
	FileInfo os.FileInfo
	ModAge   time.Duration
}

// NewFileInfoMoreFromPath returns a FileInfoMore struct
// populatig both FileInfo and ModAge (last modification time).
func NewFileInfoMoreFromPath(path string) (FileInfoMore, error) {
	fi, err := GetFileInfo(path)
	if err != nil {
		return FileInfoMore{}, err
	}
	fm := FileInfoMore{FileInfo: fi}
	fm.ModAge = time.Since(fi.ModTime())
	return fm, nil
}

// FileInfoModAge returns the file last modification age as a time.Duration.
func FileInfoModAge(fi os.FileInfo) time.Duration {
	return time.Since(fi.ModTime())
}

// FileModAge returns a time.Duration representing the age
// of the named file from FileInfo.ModTime().
func FileModAge(filename string) (time.Duration, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		dur, _ := time.ParseDuration("0s")
		return dur, err
	}
	return FileInfoModAge(fi), nil
}

func FilenModAgeGTE(filename string, s string) (bool, error) {
	ageCheck, err := time.ParseDuration(s)
	if err != nil {
		return false, err
	}
	fileAge, err := FileModAge(filename)
	if err != nil {
		return false, err
	}
	if fileAge.Hours() >= ageCheck.Hours() {
		return true, nil
	} else {
		return false, nil
	}
}

func FileModAgeLTE(filename string, s string) (bool, error) {
	ageCheck, err := time.ParseDuration(s)
	if err != nil {
		return false, err
	}
	fileAge, err := FileModAge(filename)
	if err != nil {
		return false, err
	}
	if fileAge.Hours() <= ageCheck.Hours() {
		return true, nil
	} else {
		return false, nil
	}
}

func FileInfosNames(fis []os.FileInfo) []string {
	s := []string{}
	for _, e := range fis {
		s = append(s, e.Name())
	}
	return s
}

// MustFileSize returns value of `FileInfo.Size()` which is length in bytes for
// regular files; system-dependent for others.
// It returns `-1` if an error is encountered.
func MustFileSize(filename string) int64 {
	fi, err := os.Lstat(filename)
	if err != nil {
		return -1
	}
	return fi.Size()
}
