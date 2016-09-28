package ioutilmore

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
	modAge, err := GetFileModAge(fi)
	if err != nil {
		fm.ModAge = modAge
	}
	return fm, nil
}

func GetFileModAge(fi os.FileInfo) (time.Duration, error) {
	now := time.Now()
	age := now.Sub(fi.ModTime())
	return age, nil
}

func GetFilepathModAge(path string) (time.Duration, error) {
	fi, err := GetFileInfo(path)
	if err != nil {
		dur, _ := time.ParseDuration("0s")
		return dur, err
	}
	return GetFileModAge(fi)
}

func FilepathModAgeGTE(path string, s string) (bool, error) {
	ageCheck, err := time.ParseDuration(s)
	if err != nil {
		return false, err
	}
	fileAge, err := GetFilepathModAge(path)
	if err != nil {
		return false, err
	}
	if fileAge.Hours() >= ageCheck.Hours() {
		return true, nil
	} else {
		return false, nil
	}
}

func FilepathModAgeLTE(path string, s string) (bool, error) {
	ageCheck, err := time.ParseDuration(s)
	if err != nil {
		return false, err
	}
	fileAge, err := GetFilepathModAge(path)
	if err != nil {
		return false, err
	}
	if fileAge.Hours() <= ageCheck.Hours() {
		return true, nil
	} else {
		return false, nil
	}
}
