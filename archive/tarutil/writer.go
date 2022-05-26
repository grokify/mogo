package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"sort"
)

// CreateArchiveGzipFile creates a gzipped tarball with the name `tgzFilename`. `files`
// is a `map[string]string` where the keys are local file paths and the values are
// archive filepaths. If the archive filepath is empty, the local filepath is used.
func CreateArchiveGzipFile(tgzFilename string, files map[string]string) error {
	// Create output file
	out, err := os.Create(tgzFilename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create the archive and write the output to the "out" Writer
	return CreateArchiveGzipWriter(out, files)
	// Adapted from: https://www.arthurkoziel.com/writing-tar-gz-files-in-go/ by Arthur Koziel.
	// Converted to reusable functions with archive internal filename support.
}

func CreateArchiveGzipWriter(buf io.Writer, files map[string]string) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	filenames := []string{}
	for filename := range files {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)
	// Iterate over files and add them to the tar archive
	for _, filename := range filenames {
		var err error
		archivename, ok := files[filename]
		if !ok {
			panic("map key not found")
		}
		if len(archivename) > 0 {
			err = AddFileToArchive(tw, filename, archivename)
		} else {
			err = AddFileToArchive(tw, filename, "")
		}
		if err != nil {
			return err
		}
	}

	return nil
	// Adapted from: https://www.arthurkoziel.com/writing-tar-gz-files-in-go/ by Arthur Koziel.
	// Converted to reusable functions with archive internal filename support.
}

func AddFileToArchive(tw *tar.Writer, filename, archivename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	if len(archivename) > 0 {
		header.Name = archivename
	} else {
		header.Name = filename
	}

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	return tw.Flush()
	// Adapted from: https://www.arthurkoziel.com/writing-tar-gz-files-in-go/ by Arthur Koziel.
	// Converted to reusable functions with archive internal filename support.
}
