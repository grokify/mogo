package osutil

//go:build cgo && (linux || darwin)
import (
	"syscall"
)

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
