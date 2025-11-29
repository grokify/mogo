package archivesecure

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Options for archive path validation
type PathCheckOptions struct {
	CheckSymlink bool // tar only
	CheckDevice  bool // tar only
}

// IsUnsafePath returns true if the path is unsafe
func IsUnsafePath(name string, opts PathCheckOptions) bool {
	// Null byte
	if strings.Contains(name, "\x00") {
		return true
	}

	// Absolute paths
	if filepath.IsAbs(name) || strings.HasPrefix(name, "/") {
		return true
	}

	// Windows drive letters
	if len(name) > 2 && name[1] == ':' {
		return true
	}

	if runtime.GOOS != "windows" && len(name) > 2 && name[1] == ':' {
		return true
	}

	clean := filepath.Clean(name)

	// Traversal
	if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return true
	}
	if strings.Contains(name, "../") || strings.Contains(name, `..\`) {
		return true
	}

	// Symlink or device checks can be done externally by passing opts
	return false
}
