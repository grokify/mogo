// Command rotate rotates an image from a file or URL by any number of
// degrees and writes the result as JPEG, or PNG when the output name ends
// in .png. Positive degrees turn clockwise and negative degrees turn
// counterclockwise. Multiples of 90° are lossless.
//
//	rotate -o out.jpg -d 90 photo.jpg              # 90° clockwise
//	rotate -o out.jpg -d -90 photo.jpg             # 90° counterclockwise
//	rotate -o out.jpg -d 7.5 --bg black photo.jpg  # larger canvas, black corners
//	rotate -o out.jpg -d 7.5 --crop photo.jpg      # crop away the corners
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
	xdraw "golang.org/x/image/draw"

	"github.com/grokify/mogo/image/colors"
	"github.com/grokify/mogo/image/imageutil"
)

type options struct {
	degrees float64
	crop    bool
	bg      color.Color
	interp  xdraw.Interpolator
	opaque  bool // flatten onto bg, needed for formats without alpha
}

func main() {
	var (
		o            options
		out, bg, alg string
		quality      int
	)
	cmd := &cobra.Command{
		Use:   "rotate [flags] <file-or-url>",
		Short: "Rotate an image clockwise or counterclockwise by any number of degrees",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("degrees") {
				return fmt.Errorf("--degrees is required")
			}
			c, err := colors.Parse(bg)
			if err != nil {
				return fmt.Errorf("--bg: %w", err)
			}
			o.bg = c
			if o.interp, err = parseInterpolator(alg); err != nil {
				return fmt.Errorf("--interp: %w", err)
			}
			isPNG := strings.EqualFold(filepath.Ext(out), ".png")
			o.opaque = !isPNG

			img, _, err := imageutil.ReadImage(args[0])
			if err != nil {
				return err
			}
			res := process(img, o)

			if isPNG {
				err = imageutil.Image{Image: res}.WritePNGFile(out)
			} else {
				err = imageutil.Image{Image: res}.WriteJPEGFile(out,
					&imageutil.JPEGEncodeOptions{Options: &jpeg.Options{Quality: quality}})
			}
			if err != nil {
				return err
			}
			fmt.Printf("wrote %s (%dx%d)\n", out, res.Bounds().Dx(), res.Bounds().Dy())
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&out, "out", "o", "rotated.jpg", "output file; JPEG unless the name ends in .png")
	f.Float64VarP(&o.degrees, "degrees", "d", 0, "degrees to rotate: positive is clockwise, negative is counterclockwise (required)")
	f.BoolVar(&o.crop, "crop", false, "crop to the largest rectangle with no background instead of growing the canvas")
	f.StringVar(&bg, "bg", "white", `color for the exposed corners and transparent areas (name or hex, e.g. "#0000ff")`)
	f.StringVar(&alg, "interp", imageutil.AlgCatmullRom, "pixel sampling: nearestneighbor, approxbilinear, bilinear, or catmullrom")
	f.IntVarP(&quality, "quality", "q", imageutil.JPEGQualityMax, "JPEG quality (1-100)")
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// process rotates img, then flattens it onto the background color when the
// output format has no alpha.
func process(img image.Image, o options) image.Image {
	var res *image.RGBA
	if o.crop {
		res = imageutil.RotateCrop(img, o.degrees, o.interp)
	} else {
		res = imageutil.Rotate(img, o.degrees, o.bg, o.interp)
	}
	if !o.opaque || o.bg == nil {
		return res
	}
	flat := image.NewRGBA(res.Bounds())
	draw.Draw(flat, flat.Bounds(), image.NewUniform(o.bg), image.Point{}, draw.Src)
	draw.Draw(flat, flat.Bounds(), res, res.Bounds().Min, draw.Over)
	return flat
}

func parseInterpolator(alg string) (xdraw.Interpolator, error) {
	s, err := imageutil.ParseScaler(alg)
	if err != nil {
		return nil, err
	}
	interp, ok := s.(xdraw.Interpolator)
	if !ok {
		return nil, fmt.Errorf("%q cannot rotate images", alg)
	}
	return interp, nil
}
