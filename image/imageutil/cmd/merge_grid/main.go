// Command merge_grid merges images into a grid and writes it as a JPEG.
// Images are placed left to right, top to bottom; each row is scaled to a
// common height and rows to a common width, so a short last row is scaled
// up to the full grid width.
//
//	merge_grid -o grid.jpg -c 2 -w 2000 front.jpg back.jpg detail1.jpg detail2.jpg
package main

import (
	"errors"
	"fmt"
	"image/jpeg"
	"log"

	"github.com/spf13/cobra"

	"github.com/grokify/mogo/image/imageutil"
)

func main() {
	var (
		out     string
		cols    int
		width   int
		quality int
	)
	cmd := &cobra.Command{
		Use:   "merge_grid [flags] image...",
		Short: "Merge images into a grid and write it as a JPEG",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if cols < 1 {
				return errors.New("--cols must be at least 1")
			}
			mat, err := imageutil.MatrixRead(gridRows(args, cols))
			if err != nil {
				return err
			}
			img := mat.Merge(true, true)
			if width > 0 {
				img = imageutil.Resize(width, 0, img, imageutil.ScalerBest())
			}
			im := imageutil.Image{Image: img}
			opts := &imageutil.JPEGEncodeOptions{Options: &jpeg.Options{Quality: quality}}
			if err := im.WriteJPEGFile(out, opts); err != nil {
				return err
			}
			b := img.Bounds()
			fmt.Printf("wrote %s (%dx%d)\n", out, b.Dx(), b.Dy())
			return nil
		},
	}
	cmd.Flags().StringVarP(&out, "out", "o", "merged.jpg", "output JPEG file")
	cmd.Flags().IntVarP(&cols, "cols", "c", 2, "images per row")
	cmd.Flags().IntVarP(&width, "width", "w", 0, "output width in pixels, preserving aspect ratio (0 keeps merged size)")
	cmd.Flags().IntVarP(&quality, "quality", "q", imageutil.JPEGQualityMax, "JPEG quality (1-100)")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// gridRows splits paths into rows of cols elements; the last row may be shorter.
func gridRows(paths []string, cols int) [][]string {
	var rows [][]string
	for len(paths) > cols {
		rows = append(rows, paths[:cols])
		paths = paths[cols:]
	}
	return append(rows, paths)
}
