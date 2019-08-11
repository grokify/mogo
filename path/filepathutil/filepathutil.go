package filepathutil

// $ go get github.com/grokify/gotilla/git/apps/gitremoteaddupstream

import (
	"os"
	"os/user"
	"strings"
)

func CurLeafDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dirParts := strings.Split(dir, string(os.PathSeparator))
	if len(dirParts) > 0 {
		return dirParts[len(dirParts)-1], nil
	}
	return "", nil
}

// UserToAbsolute converts ~ directories to absolute directories
// in filepaths. This is useful because ~ cannot be resolved by
// ioutil.ReadFile().
func UserToAbsolute(file string) (string, error) {
	file = strings.TrimSpace(file)
	parts := strings.Split(file, string(os.PathSeparator))
	if len(parts) == 0 {
		return file, nil
	} else if parts[0] != "~" {
		return file, nil
	}
	usr, err := user.Current()
	if err != nil {
		return file, err
	}
	if len(parts) == 1 {
		return usr.HomeDir, nil
	}

	return strings.Join(
		append([]string{usr.HomeDir}, parts[1:]...),
		string(os.PathSeparator)), nil
}
