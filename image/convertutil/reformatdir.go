package convertutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/gotilla/image/imageutil"
	"github.com/grokify/gotilla/io/ioutilmore"
)

type CopyType int

const (
	PDFFormat    CopyType = iota // 0 convert cli value
	KindleFormat                 // 1 convert cli value
)

// ReformatImages converts images in one dir to another using default
// formats for Kindle and PDF.
func ReformatImages(baseSrcDir, baseOutDir string, copyType CopyType) error {
	var err error
	baseSrcDir, err = filepath.Abs(strings.TrimSpace(baseSrcDir))
	if err != nil {
		return err
	}
	baseOutDir, err = filepath.Abs(strings.TrimSpace(baseOutDir))
	if err != nil {
		return err
	}
	return reformatImagesSubdir(baseSrcDir, baseOutDir, "", copyType)
}

func reformatImagesSubdir(baseSrcDir, baseOutDir, dirPart string, copyType CopyType) error {
	thisSrcDir := baseSrcDir
	thisOutDir := baseOutDir
	dirPart = strings.TrimSpace(dirPart)
	if len(dirPart) > 0 {
		thisSrcDir = filepath.Join(thisSrcDir, dirPart)
		thisOutDir = filepath.Join(thisOutDir, dirPart)

		isDir, err := ioutilmore.IsDir(thisSrcDir)
		if err != nil {
			return err
		}
		if !isDir {
			return fmt.Errorf("Need Dir [%s]", thisSrcDir)
		}
	}

	if err := os.MkdirAll(thisOutDir, 0755); err != nil {
		return err
	}

	sdirs, files, err := ioutilmore.ReadDirSplit(thisSrcDir, false)
	if err != nil {
		return err
	}
	for _, file := range files {
		thisSrcFile := filepath.Join(thisSrcDir, file.Name())
		thisOutFile := filepath.Join(thisOutDir, file.Name())
		if !imageutil.IsImageExt(thisSrcFile) {
			continue
		}
		switch copyType {
		case PDFFormat:
			_, stderr, err := ConvertToPDF(thisSrcFile, thisOutFile)
			err = CheckError(err, stderr)
			if err != nil {
				return err
			}
		case KindleFormat:
			_, stderr, err := ConvertToKindle(thisSrcFile, thisOutFile)
			err = CheckError(err, stderr)
			if err != nil {
				return err
			}
		}
	}
	for _, sdir := range sdirs {
		subDir := sdir.Name()
		if len(dirPart) > 0 {
			subDir = filepath.Join(dirPart, subDir)
		}
		err := reformatImagesSubdir(baseSrcDir, baseOutDir, subDir, copyType)
		if err != nil {
			return err
		}
	}
	return nil
}
