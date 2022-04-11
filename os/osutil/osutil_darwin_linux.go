//go:build linux || darwin
// +build linux darwin

package osutil

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
