package osutil

import (
	"io"
	"os"
)

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

// CopyFileSecure copies a file from src to dst after validating the destination
// path does not contain path traversal sequences (".."). This is the recommended
// function for library code that receives destination paths from callers.
//
// For CLI entry points where the user explicitly provides paths, use
// CopyFile directly with a //nolint:gosec comment instead.
func CopyFileSecure(src, dst string) error {
	cleanDst, err := CleanPathSecure(dst)
	if err != nil {
		return err
	}
	return CopyFile(src, cleanDst)
}
