package osutil

import (
	"os"
	"time"
)

type FileInfoMore struct {
	FileInfo os.FileInfo
	ModAge   time.Duration
}

func NewFileInfoMoreFromPath(path string) (FileInfoMore, error) {
	fi, err := GetFileInfo(path)
	if err != nil {
		return FileInfoMore{}, err
	}
	fm := FileInfoMore{FileInfo: fi}
	fm.ModAge = time.Now().Sub(fi.ModTime())
	return fm, nil
}
