package osutil

import (
	"os"
	"syscall"
	"time"
)

// FileInfoMore provides a struct hold FileInfo with
// additional information.
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

func FileStatT(filename string) (syscall.Stat_t, error) {
	var stat syscall.Stat_t
	return stat, syscall.Stat(filename, &stat)
}

func MustFileStatT(filename string) syscall.Stat_t {
	stat, err := FileStatT(filename)
	if err != nil {
		panic(err)
	}
	return stat
}
