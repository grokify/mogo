// Command crop crops, squares, and splits an image from a file or URL and
// writes the results as JPEG. Operations run in order: --width/--height
// crop, then --square, then --split.
//
//	crop -o out.jpg --width 800 --align-x left photo.jpg
//	crop -o out.jpg --square larger --bg white https://example.com/photo.jpg
//	crop -o scan.jpg --split x two-sided-scan.jpg   # writes scan_left.jpg and scan_right.jpg
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/mogo/image/colors"
	"github.com/grokify/mogo/image/imageutil"
)

type options struct {
	width, height  int
	alignX, alignY string
	square         string // "", "smaller" (crop), or "larger" (pad)
	bg             color.Color
	split          string // "", "x" (left/right), or "y" (top/bottom)
}

// output is a processed image and the suffix added to the output file name.
type output struct {
	suffix string
	img    image.Image
}

func main() {
	var (
		o       options
		out, bg string
		quality int
	)
	cmd := &cobra.Command{
		Use:   "crop [flags] <file-or-url>",
		Short: "Crop, square, and split an image, writing JPEG output",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			c, err := colors.Parse(bg)
			if err != nil {
				return fmt.Errorf("--bg: %w", err)
			}
			o.bg = c
			img, _, err := imageutil.ReadImage(args[0])
			if err != nil {
				return err
			}
			outs, err := process(img, o)
			if err != nil {
				return err
			}
			opts := &imageutil.JPEGEncodeOptions{Options: &jpeg.Options{Quality: quality}}
			for _, r := range outs {
				name := outputName(out, r.suffix)
				if err := (imageutil.Image{Image: r.img}).WriteJPEGFile(name, opts); err != nil {
					return err
				}
				fmt.Printf("wrote %s (%dx%d)\n", name, r.img.Bounds().Dx(), r.img.Bounds().Dy())
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&out, "out", "o", "cropped.jpg", "output JPEG file; --split adds _left/_right or _top/_bottom")
	f.IntVar(&o.width, "width", 0, "crop to this width in pixels (0 keeps the width)")
	f.IntVar(&o.height, "height", 0, "crop to this height in pixels (0 keeps the height)")
	f.StringVar(&o.alignX, "align-x", imageutil.AlignCenter, "part kept by --width: left, center, or right")
	f.StringVar(&o.alignY, "align-y", imageutil.AlignCenter, "part kept by --height: top, center, or bottom")
	f.StringVar(&o.square, "square", "", `make square: "smaller" crops to the shorter side, "larger" pads to the longer side`)
	f.StringVar(&bg, "bg", "white", `background for --square larger padding and transparent areas (name or hex, e.g. "#0000ff")`)
	f.StringVar(&o.split, "split", "", `split into halves: "x" for left/right, "y" for top/bottom`)
	f.IntVarP(&quality, "quality", "q", imageutil.JPEGQualityMax, "JPEG quality (1-100)")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// process applies the crop, square, and split operations in that order, then
// flattens each result onto the background color since JPEG has no alpha.
func process(img image.Image, o options) ([]output, error) {
	outs, err := transform(img, o)
	if err != nil {
		return nil, err
	}
	bg := o.bg
	if bg == nil {
		bg = color.White
	}
	for i, r := range outs {
		flat := image.NewRGBA(r.img.Bounds())
		draw.Draw(flat, flat.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)
		draw.Draw(flat, flat.Bounds(), r.img, r.img.Bounds().Min, draw.Over)
		outs[i].img = flat
	}
	return outs, nil
}

func transform(img image.Image, o options) ([]output, error) {
	if o.width < 0 || o.height < 0 {
		return nil, fmt.Errorf("--width and --height must not be negative")
	}
	if err := checkChoice("--align-x", o.alignX, imageutil.AlignLeft, imageutil.AlignCenter, imageutil.AlignRight); err != nil {
		return nil, err
	}
	if err := checkChoice("--align-y", o.alignY, imageutil.AlignTop, imageutil.AlignCenter, imageutil.AlignBottom); err != nil {
		return nil, err
	}
	if o.width > 0 {
		img = imageutil.CropX(img, o.width, o.alignX)
	}
	if o.height > 0 {
		img = imageutil.CropY(img, o.height, o.alignY)
	}

	switch o.square {
	case "":
	case "smaller":
		img = imageutil.Image{Image: img}.SquareSmaller()
	case "larger":
		img = imageutil.Image{Image: img}.SquareLarger(color.Transparent) // filled by process
	default:
		return nil, fmt.Errorf(`--square must be "smaller" or "larger", got %q`, o.square)
	}

	switch o.split {
	case "":
		return []output{{"", img}}, nil
	case "x":
		half := img.Bounds().Dx() / 2
		return []output{
			{"_left", imageutil.CropX(img, half, imageutil.AlignLeft)},
			{"_right", imageutil.CropX(img, half, imageutil.AlignRight)},
		}, nil
	case "y":
		half := img.Bounds().Dy() / 2
		return []output{
			{"_top", imageutil.CropY(img, half, imageutil.AlignTop)},
			{"_bottom", imageutil.CropY(img, half, imageutil.AlignBottom)},
		}, nil
	default:
		return nil, fmt.Errorf(`--split must be "x" or "y", got %q`, o.split)
	}
}

func checkChoice(flag, val string, choices ...string) error {
	for _, c := range choices {
		if val == c {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %s, got %q", flag, strings.Join(choices, ", "), val)
}

// outputName inserts suffix before the extension of name.
func outputName(name, suffix string) string {
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext) + suffix + ext
}
