package ziputil

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/grokify/mogo/archive/archivesecure"
)

// FindUnsafeZipPaths scans all entries in a zip.Reader and returns a slice of unsafe paths.
// Unsafe paths include:
//   - directory traversal (../ or ..\)
//   - absolute paths (/foo, C:\foo)
//   - drive letters
//   - null-byte injection
//
// Returns a non-nil error if any unsafe paths are detected.
func FindUnsafeZipPaths(zr *zip.Reader) ([]string, error) {
	var bad []string
	for _, f := range zr.File {
		if archivesecure.IsUnsafePath(f.Name, archivesecure.PathCheckOptions{}) {
			bad = append(bad, f.Name)
		}
	}
	if len(bad) > 0 {
		return bad, errors.New("unsafe paths detected in ZIP")
	}
	return nil, nil
}

// StreamScanZipPaths streams a ZIP from io.Reader to a temp file and validates it using FindUnsafeZipPaths.
func StreamScanZipPaths(r io.Reader) ([]string, error) {
	tmp, err := os.CreateTemp("", "zipstream-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err := io.Copy(tmp, r); err != nil {
		return nil, err
	}

	stat, err := tmp.Stat()
	if err != nil {
		return nil, err
	}

	readerAt, err := os.Open(tmp.Name())
	if err != nil {
		return nil, err
	}
	defer readerAt.Close()

	zr, err := zip.NewReader(readerAt, stat.Size())
	if err != nil {
		return nil, err
	}

	bad, err := FindUnsafeZipPaths(zr)
	if err != nil {
		return bad, fmt.Errorf("unsafe file paths detected in streamed ZIP: %v", bad)
	}

	return nil, nil
}
