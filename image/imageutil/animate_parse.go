package imageutil

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ReadGIF decodes all frames of an animated GIF. Frames are returned as stored
// in the file, i.e. each may be a partial "delta" image; use `GIFFrames` to get
// fully composited frames.
func ReadGIF(r io.Reader) (*gif.GIF, error) {
	return gif.DecodeAll(r)
}

// ReadGIFFile decodes all frames of an animated GIF file.
func ReadGIFFile(filename string) (*gif.GIF, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	g, err := ReadGIF(f)
	if err != nil {
		_ = f.Close() // decode error takes precedence
		return nil, err
	}
	return g, f.Close()
}

// GIFFrameInfo describes a single frame of a GIF.
type GIFFrameInfo struct {
	Index    int             // zero-based
	Bounds   image.Rectangle // region of the canvas this frame updates
	Delay    time.Duration   // display time
	Start    time.Duration   // offset from start of animation
	Disposal byte            // gif.DisposalNone, DisposalBackground, DisposalPrevious
}

// GIFInfo summarizes an animated GIF.
type GIFInfo struct {
	Width, Height int
	FrameCount    int
	LoopCount     int // 0 = loop forever, -1 = play once, n = n+1 plays
	Duration      time.Duration
	Frames        []GIFFrameInfo
}

// GIFInfoFromGIF returns metadata for a decoded GIF.
func GIFInfoFromGIF(g *gif.GIF) (GIFInfo, error) {
	if g == nil || len(g.Image) == 0 {
		return GIFInfo{}, errors.New("gif has no frames")
	}
	w, h := gifCanvasSize(g)
	info := GIFInfo{Width: w, Height: h, FrameCount: len(g.Image), LoopCount: g.LoopCount}
	var at time.Duration
	for i, img := range g.Image {
		fi := GIFFrameInfo{Index: i, Bounds: img.Bounds(), Start: at, Disposal: gifDisposal(g, i)}
		if i < len(g.Delay) {
			fi.Delay = time.Duration(g.Delay[i]) * 10 * time.Millisecond // GIF delay is in 1/100 s
		}
		at += fi.Delay
		info.Frames = append(info.Frames, fi)
	}
	info.Duration = at
	return info, nil
}

// GIFFrames composites an animated GIF into full-canvas frames, honoring each
// frame's disposal method so that every returned image is what a viewer would
// see at that point in the animation.
func GIFFrames(g *gif.GIF) ([]*image.RGBA, error) {
	if g == nil || len(g.Image) == 0 {
		return nil, errors.New("gif has no frames")
	}
	w, h := gifCanvasSize(g)
	canvasRect := image.Rect(0, 0, w, h)
	canvas := image.NewRGBA(canvasRect) // starts fully transparent

	frames := make([]*image.RGBA, 0, len(g.Image))
	for i, src := range g.Image {
		var prev *image.RGBA
		disposal := gifDisposal(g, i)
		if disposal == gif.DisposalPrevious {
			prev = cloneRGBA(canvas)
		}

		draw.Draw(canvas, src.Bounds(), src, src.Bounds().Min, draw.Over)
		frames = append(frames, cloneRGBA(canvas))

		switch disposal {
		case gif.DisposalBackground:
			// Renderers overwhelmingly treat this as transparent, not the
			// palette's background color.
			draw.Draw(canvas, src.Bounds(), image.Transparent, image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			canvas = prev
		}
	}
	return frames, nil
}

// WriteGIFFramesPNG writes each composited frame to `dir` as
// `<prefix>_NNN.png` and returns the file paths in order. The directory is
// created if needed.
func WriteGIFFramesPNG(g *gif.GIF, dir, prefix string) ([]string, error) {
	frames, err := GIFFrames(g)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(frames))
	for i, fr := range frames {
		name := filepath.Join(dir, fmt.Sprintf("%s_%03d.png", prefix, i))
		if err := writeGIFFramePNG(name, fr); err != nil {
			return names, err
		}
		names = append(names, name)
	}
	return names, nil
}

func writeGIFFramePNG(name string, img image.Image) error {
	f, err := os.Create(name) // #nosec G304 -- caller-supplied output path
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		_ = f.Close() // encode error takes precedence
		return err
	}
	return f.Close()
}

// gifCanvasSize returns the logical screen size, falling back to the union of
// frame bounds when the GIF has no (or a zero) Config.
func gifCanvasSize(g *gif.GIF) (int, int) {
	if g.Config.Width > 0 && g.Config.Height > 0 {
		return g.Config.Width, g.Config.Height
	}
	var r image.Rectangle
	for _, img := range g.Image {
		r = r.Union(img.Bounds())
	}
	return r.Max.X, r.Max.Y
}

func gifDisposal(g *gif.GIF, i int) byte {
	if i < len(g.Disposal) {
		return g.Disposal[i]
	}
	return gif.DisposalNone
}

func cloneRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Rect)
	copy(dst.Pix, src.Pix)
	return dst
}
