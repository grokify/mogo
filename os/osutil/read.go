package osutil

import (
	"bufio"
	"encoding/json"
	"io/fs"
	"os"

	"github.com/grokify/mogo/errors/errorsutil"
)

func ReadFileByLine(name string, lineFunc func(idx uint, line string) error) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	i := uint(0)
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		txt := scanner.Text()
		err := lineFunc(i, txt)
		if err != nil {
			return err
		}
		i++
	}

	return scanner.Err()
}

// ReadFileJSON reads and unmarshals a file.
func ReadFileJSON(file string, v any) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func CloseFileWithError(file *os.File, err error) error {
	errFile := file.Close()
	if err != nil {
		return errorsutil.Wrap(err, errFile.Error())
	}
	return err
}

// ReadDirFilesSecure reads all files from a directory recursively using os.Root
// for symlink-safe operations. This prevents TOCTOU race conditions where symlinks
// could be swapped during traversal (gosec G122).
//
// Returns a map of relative paths to file contents. Paths use forward slashes
// regardless of OS.
//
// Requires Go 1.24+ for os.Root support.
func ReadDirFilesSecure(dir string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := WalkFilesSecure(dir, nil, func(path string, content []byte) error {
		files[path] = content
		return nil
	})
	return files, err
}

// WalkFilesSecure walks dir recursively using os.Root for symlink-safe
// operations, preventing TOCTOU race conditions where a symlink could be
// swapped between traversal and read (gosec G122). Unlike ReadDirFilesSecure,
// it doesn't buffer every file in memory at once — fn is called with each
// regular file's path (relative to dir, forward-slashed) and content as the
// walk proceeds.
//
// skipDir, if non-nil, is called with each directory's base name (not its
// full path); returning true skips that directory and its contents entirely
// — use this to exclude vendor/, node_modules/, .git/, testdata/, and
// similar. A nil skipDir walks every directory.
//
// If fn returns an error, the walk stops immediately and that error is
// returned to the caller.
//
// Requires Go 1.24+ for os.Root support.
func WalkFilesSecure(dir string, skipDir func(name string) bool, fn func(path string, content []byte) error) error {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()

	return fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != "." && skipDir != nil && skipDir(d.Name()) {
				return fs.SkipDir
			}
			return nil
		}
		content, err := root.ReadFile(path)
		if err != nil {
			return err
		}
		return fn(path, content)
	})
}
