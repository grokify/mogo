package codegen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/io/ioutilmore"
	"github.com/pkg/errors"
)

func ConvertFilesInPlaceNestedstructsToPointers(dir string, rx *regexp.Regexp, perm os.FileMode) ([]string, error) {
	filepaths := []string{}
	if rx == nil {
		rx = regexp.MustCompile(`.*\.go$`)
	}
	files, err := ioutilmore.DirEntriesRxSizeGt0(dir, ioutilmore.File, rx)
	if err != nil {
		return filepaths, errors.Wrap(err, "codegen.ConvertFilesInPlace.DirEntriesReNotEmpty")
	}
	for _, file := range files {
		filepath := filepath.Join(dir, file.Name())
		err := ConvertFileNestedstructsToPointers(filepath, filepath, perm)
		if err != nil {
			return filepaths, errors.Wrap(err, "codegen.ConvertFilesInPlace.ConvertFile")
		}
		filepaths = append(filepaths, filepath)
	}
	return filepaths, nil
}

func ConvertFileNestedstructsToPointers(inFile, outFile string, perm os.FileMode) error {
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}
	newLines := GoCodeNestedstructsToPointers(string(data))
	return ioutil.WriteFile(outFile, []byte(newLines), perm)
}

var (
	rxParenOpen         = regexp.MustCompile(`{\s*$`)
	rxParenClose        = regexp.MustCompile(`^\s*}\s*$`)
	rxCustomTypeComplex = regexp.MustCompile(`^(\s*[0-9A-Za-z]+\s+[0-9a-z\]\[]+\])([A-Z].*)$`)
	rxCustomTypeSimple  = regexp.MustCompile(`^([\s\t]*[0-9A-Za-z]+[\s\t]+)([A-Z].*)$`)
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
			mc := rxCustomTypeComplex.FindAllStringSubmatch(line, -1)
			if len(mc) > 0 {
				newLines = append(newLines, mc[0][1]+"*"+mc[0][2])
				continue
			}
			ms := rxCustomTypeSimple.FindAllStringSubmatch(line, -1)
			if len(ms) > 0 {
				newLines = append(newLines, ms[0][1]+"*"+ms[0][2])
				continue
			}
		}
		newLines = append(newLines, line)
	}
	return strings.Join(newLines, "\n")
}
