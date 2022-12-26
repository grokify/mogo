package ioutil

/*
import (
	"os"
	"time"
)

type FileInfoMore struct {
	FileInfo os.FileInfo
	ModAge   time.Duration
}

func NewFileInfoMoreFromPath(path string) (FileInfoMore, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return FileInfoMore{}, err
	}
	fm := FileInfoMore{FileInfo: fi}
	modAge, err := FileModAge(fi)
	if err != nil {
		fm.ModAge = modAge
	}
	return fm, nil
}

func FileModAge(fi os.FileInfo) (time.Duration, error) {
	now := time.Now()
	age := now.Sub(fi.ModTime())
	return age, nil
}

func FilenameModAge(filename string) (time.Duration, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		dur, _ := time.ParseDuration("0s")
		return dur, err
	}
	return FileModAge(fi)
}

func FilenameModAgeGTE(filename string, s string) (bool, error) {
	ageCheck, err := time.ParseDuration(s)
	if err != nil {
		return false, err
	}
	fileAge, err := FilenameModAge(filename)
	if err != nil {
		return false, err
	}
	if fileAge.Hours() >= ageCheck.Hours() {
		return true, nil
	} else {
		return false, nil
	}
}

func FilenameModAgeLTE(filename string, s string) (bool, error) {
	ageCheck, err := time.ParseDuration(s)
	if err != nil {
		return false, err
	}
	fileAge, err := FilenameModAge(filename)
	if err != nil {
		return false, err
	}
	if fileAge.Hours() <= ageCheck.Hours() {
		return true, nil
	} else {
		return false, nil
	}
}
*/
