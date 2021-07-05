package convertutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/simplego/image/imageutil"
	"github.com/grokify/simplego/io/ioutilmore"
	"github.com/grokify/simplego/os/osutil"
	"github.com/pkg/errors"
)

const (
	PDFSpecs    = "1950pxw300dpi"
	KindleSpecs = "600pxw300dpi"
)

type CopyType int

const (
	PDFFormat    CopyType = iota // 0 convert cli value
	KindleFormat                 // 1 convert cli value
)

// ReformatImages converts images in one dir to another using default
// formats for Kindle and PDF.
func ReformatImages(baseSrcDir, baseOutDir string, copyType CopyType, rewrite bool) error {
	var err error
	baseSrcDir, err = filepath.Abs(strings.TrimSpace(baseSrcDir))
	if err != nil {
		return err
	}
	baseOutDir, err = filepath.Abs(strings.TrimSpace(baseOutDir))
	if err != nil {
		return err
	}
	return reformatImagesSubdir(baseSrcDir, baseOutDir, "", copyType, rewrite)
}

func reformatImagesSubdir(baseSrcDir, baseOutDir, dirPart string, copyType CopyType, rewrite bool) error {
	thisSrcDir := baseSrcDir
	thisOutDir := baseOutDir
	dirPart = strings.TrimSpace(dirPart)
	if len(dirPart) > 0 {
		thisSrcDir = filepath.Join(thisSrcDir, dirPart)
		thisOutDir = filepath.Join(thisOutDir, dirPart)

		isDir, err := osutil.IsDir(thisSrcDir)
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
		err := reformatImagesSubdirFile(
			filepath.Join(thisSrcDir, file.Name()),
			filepath.Join(thisOutDir, file.Name()),
			copyType, rewrite)
		if err != nil {
			return err
		}
		/*
			thisSrcFile := filepath.Join(thisSrcDir, file.Name())
			thisOutFile := filepath.Join(thisOutDir, file.Name())

			if !imageutil.IsImageExt(thisSrcFile) {
				continue
			}
			if !rewrite {
				isFile, err := ioutilmore.IsFile(thisOutFile)
				if err == nil && isFile {
					continue
				}
			}

			switch copyType {
			case PDFFormat:
				_, _, err := ConvertToPDF(thisSrcFile, thisOutFile)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("ConvertToPDF failed for [%s]", thisSrcFile))
				}
			case KindleFormat:
				_, stderr, err := ConvertToKindle(thisSrcFile, thisOutFile)
				err = CheckError(err, stderr)
				if err != nil {
					return err
				}
			}
		*/
	}
	for _, sdir := range sdirs {
		subDir := sdir.Name()
		if len(dirPart) > 0 {
			subDir = filepath.Join(dirPart, subDir)
		}
		err := reformatImagesSubdir(baseSrcDir, baseOutDir, subDir, copyType, rewrite)
		if err != nil {
			return err
		}
	}
	return nil
}

func reformatImagesSubdirFile(thisSrcFile, thisOutFile string, copyType CopyType, rewrite bool) error {
	if !imageutil.IsImageExt(thisSrcFile) {
		return nil
	}
	if !rewrite {
		isFile, err := osutil.IsFile(thisOutFile)
		if err == nil && isFile {
			return nil
		}
	}

	switch copyType {
	case PDFFormat:
		_, _, err := ConvertToPDF(thisSrcFile, thisOutFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("ConvertToPDF failed for [%s]", thisSrcFile))
		}
	case KindleFormat:
		_, stderr, err := ConvertToKindle(thisSrcFile, thisOutFile)
		err = CheckError(err, stderr)
		if err != nil {
			return err
		}
	}
	return nil
}
