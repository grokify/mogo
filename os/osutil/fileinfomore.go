package osutil

import (
	"os"
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
	fm.ModAge = time.Now().Sub(fi.ModTime())
	return fm, nil
}
