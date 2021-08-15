package codegen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grokify/simplego/os/osutil"
	"github.com/pkg/errors"
)

func ConvertFilesInPlaceNestedstructsToPointers(dir string, rx *regexp.Regexp) ([]string, error) {
	filenames := []string{}
	if rx == nil {
		rx = regexp.MustCompile(`.*\.go$`)
	}
	entries, err := osutil.ReadDirMore(dir, rx, false, true, false)
	// files, err := ioutilmore.DirEntriesRxSizeGt0(dir, ioutilmore.File, rx)
	if err != nil {
		return filenames, errors.Wrap(err, "codegen.ConvertFilesInPlace.ReadDirMore")
	}
	//filenames := osutil.DirEntries(entries).Names(dir, true)
	//for _, filename := range filenames {
	for _, entry := range entries {
		filename := filepath.Join(dir, entry.Name())
		fileinfo, err := entry.Info()
		if err != nil {
			return filenames, errors.Wrap(err, fmt.Sprintf("codegen.ConvertFilesInPlaceNestedstructsToPointers...entry.Info() [%s]", entry.Name()))
		}
		err = ConvertFileNestedstructsToPointers(filename, filename, fileinfo.Mode().Perm())
		if err != nil {
			return filenames, errors.Wrap(err, "codegen.ConvertFilesInPlace.ConvertFile")
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func ConvertFileNestedstructsToPointers(inFile, outFile string, perm os.FileMode) error {
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(
		outFile,
		[]byte(GoCodeNestedstructsToPointers(string(data))),
		perm)
}

var (
	rxParenOpen  = regexp.MustCompile(`^type\s+\S+\s+struct\s+{\s*$`)
	rxParenClose = regexp.MustCompile(`^\s*}\s*$`)
	rxCustomType = regexp.MustCompile(`^(\s*[0-9A-Za-z]+\s+(?:[0-9a-z\]\[]+\])?)([A-Z].*)$`)
)

// GoCodeNestedstructsToPointers is designed to convert
// nested structs to pointers.
func GoCodeNestedstructsToPointers(code string) string {
	oldLines := strings.Split(code, "\n")
	newLines := []string{}
	inParen := false
	for _, line := range oldLines {
		if rxParenOpen.MatchString(line) {
			inParen = true
			newLines = append(newLines, line)
			continue
		} else if rxParenClose.MatchString(line) {
			inParen = false
			newLines = append(newLines, line)
			continue
		} else if inParen {
			m := rxCustomType.FindStringSubmatch(line)
			if len(m) > 0 {
				newLines = append(newLines, m[1]+"*"+m[2])
				continue
			}
		}
		newLines = append(newLines, line)
	}
	return strings.Join(newLines, "\n")
}
