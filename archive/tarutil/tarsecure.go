package tarutil

import (
	"archive/tar"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/grokify/mogo/archive/archivesecure"
)

func FindUnsafeTarPaths(tr *tar.Reader) ([]string, error) {
	var bad []string
	for {
		hdr, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if archivesecure.IsUnsafePath(hdr.Name, archivesecure.PathCheckOptions{
			CheckSymlink: hdr.Typeflag == tar.TypeSymlink || hdr.Typeflag == tar.TypeLink,
			CheckDevice:  hdr.Typeflag == tar.TypeChar || hdr.Typeflag == tar.TypeBlock,
		}) {
			bad = append(bad, hdr.Name)
		}

		if hdr.Typeflag == tar.TypeSymlink || hdr.Typeflag == tar.TypeLink {
			target := hdr.Linkname
			if strings.HasPrefix(filepath.Clean(target), "..") {
				bad = append(bad, hdr.Name+" -> "+target)
			}
		}

		if hdr.Typeflag == tar.TypeChar || hdr.Typeflag == tar.TypeBlock {
			bad = append(bad, hdr.Name+" (device file)")
		}
	}

	if len(bad) > 0 {
		return bad, errors.New("unsafe paths detected in tar")
	}
	return nil, nil
}
