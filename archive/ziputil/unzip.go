package ziputil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SecureUnzip extracts all files from a zip.Reader to the destination directory.
// It uses FindUnsafeZipPaths to prevent directory traversal or unsafe paths.
func SecureUnzip(zr *zip.Reader, dest string) error {
	// Security check
	if bad, err := FindUnsafeZipPaths(zr); err != nil {
		return fmt.Errorf("refusing to unzip archive: %v (unsafe: %v)", err, bad)
	}

	for _, f := range zr.File {
		fpath := filepath.Join(dest, f.Name)

		// Extra sanity check
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("zip-slip detected at extraction: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)

		out, err := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			out.Close()
			return err
		}

		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return err
		}

		out.Close()
		rc.Close()
	}

	return nil
}
