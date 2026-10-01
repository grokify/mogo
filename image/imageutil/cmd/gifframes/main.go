// Command gifframes splits an animated GIF into composited PNG frames and
// prints frame timing metadata.
//
//	gifframes -in anim.gif -out ./frames
package main

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/image/imageutil"
	"github.com/spf13/cobra"
)

func main() {
	var in, out, prefix string
	cmd := &cobra.Command{
		Use:   "gifframes",
		Short: "Split an animated GIF into PNG frames and print metadata",
		RunE: func(_ *cobra.Command, _ []string) error {
			g, err := imageutil.ReadGIFFile(in)
			if err != nil {
				return err
			}
			info, err := imageutil.GIFInfoFromGIF(g)
			if err != nil {
				return err
			}
			fmt.Printf("size=%dx%d frames=%d loop=%d duration=%v\n",
				info.Width, info.Height, info.FrameCount, info.LoopCount, info.Duration)
			for _, f := range info.Frames {
				fmt.Printf("frame %03d start=%v delay=%v disposal=%d bounds=%v\n",
					f.Index, f.Start, f.Delay, f.Disposal, f.Bounds)
			}
			if out == "" {
				return nil
			}
			names, err := imageutil.WriteGIFFramesPNG(g, out, prefix)
			if err != nil {
				return err
			}
			fmt.Printf("wrote %d frames to %s\n", len(names), out)
			return nil
		},
	}
	cmd.Flags().StringVar(&in, "in", "", "input GIF file (required)")
	cmd.Flags().StringVar(&out, "out", "", "output directory for PNG frames (omit for metadata only)")
	cmd.Flags().StringVar(&prefix, "prefix", "frame", "output filename prefix")
	if err := cmd.MarkFlagRequired("in"); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
