// Package osutil implements some OS utility functions.
package osutil

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// EmptyAll will delete all contents of a directory, leaving
// the provided directory. This is different from os.Remove
// which also removes the directory provided.
func EmptyAll(name string) error {
	aEntries, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}
	for _, f := range aEntries {
		if f.Name() == "." || f.Name() == ".." {
			continue
		}
		err = os.Remove(name + "/" + f.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// Exists checks whether the named filepath exists or not for
// a file or directory.
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// FileModAge returns a time.Duration representing the age
// of the named file from FileInfo.ModTime().
func FileModAge(name string) (time.Duration, error) {
	stat, err := os.Stat(name)
	if err != nil {
		dur0, _ := time.ParseDuration("0s")
		return dur0, err
	}
	return time.Now().Sub(stat.ModTime()), nil
}

// FileModAgeFromInfo returns the file last modification
// age as a time.Duration.
func FileModAgeFromInfo(fi os.FileInfo) time.Duration {
	return time.Now().Sub(fi.ModTime())
}

// GetFileInfo returns an os.FileInfo from a filepath.
func GetFileInfo(path string) (os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

type EnvVar struct {
	Key   string
	Value string
}

func Env() []EnvVar {
	envs := []EnvVar{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if len(pair) > 0 {
			key := strings.TrimSpace(pair[0])
			if len(key) > 0 {
				env := EnvVar{Key: key}
				if len(pair) > 1 {
					env.Value = strings.Join(pair[1:], "=")
				}
				envs = append(envs, env)
			}
		}
	}
	return envs
}

// EnvFiltered returns a map[string]string of environment
// variables that match a regular expression.
func EnvFiltered(rx *regexp.Regexp) map[string]string {
	res := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if rx.MatchString(pair[0]) {
			res[pair[0]] = pair[1]
		}
	}
	return res
}

func AbsFilepath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return path, nil
	} else if filepath.IsAbs(path) {
		return path, nil
	}
	parts := strings.Split(path, string(os.PathSeparator))
	if parts[0] != "~" {
		return filepath.Abs(path)
	}

	usr, err := user.Current()
	if err != nil {
		return path, err
	}
	if len(parts) == 1 {
		return usr.HomeDir, nil
	}

	return strings.Join(
		append([]string{usr.HomeDir}, parts[1:]...),
		string(os.PathSeparator)), nil
}
