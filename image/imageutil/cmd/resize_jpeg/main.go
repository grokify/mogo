package main

import (
	"image/jpeg"
	"log"

	"github.com/grokify/mogo/image/imageutil"
	flags "github.com/jessevdk/go-flags"
)

type cliOptions struct {
	Input   string `short:"i" long:"input dir/file" description:"A dir or file" value-name:"FILE" required:"true"`
	Output  string `short:"o" long:"output dir/filefile" description:"A dir or file" required:"true"`
	Height  uint   `short:"h" long:"height" description:"Height"`
	Width   uint   `short:"w" long:"width" description:"Width"`
	Quality int    `short:"q" long:"quality" description:"Quality"`
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	err = imageutil.ResizePathJPEG(opts.Input, opts.Output, opts.Width, opts.Height,
		&imageutil.JPEGEncodeOptions{
			Options: &jpeg.Options{Quality: opts.Quality},
		})
	if err != nil {
		log.Fatal(err)
	}

	/*
		isDirSrc, err := osutil.IsDir(opts.Input)
		logutil.FatalErr(err)

		if isDirSrc {
			isDirOut, err := osutil.IsDir(opts.Output)
			logutil.FatalErr(err)
			if !isDirOut {
				logutil.FatalErr(errors.New("output must be a directory"))
			}
			// write regexp to match .jpg or .jpeg file extensions
			files, err := osutil.ReadDirMore(opts.Input, imageutil.RxFileExtensionJPG, false, true, false)
			logutil.FatalErr(err)
			n := len(files)
			for i, e := range files {
				fmt.Printf("Processing %d of %d: %s\n", i+1, n, e.Name())
				srcPath := filepath.Join(opts.Input, e.Name())
				outPath := filepath.Join(opts.Output, e.Name())
				err := imageutil.ResizeFileJPEG(srcPath, outPath, opts.Width, opts.Height, jopts)
				logutil.FatalErr(err)
			}
		} else {
			err := imageutil.ResizeFileJPEG(opts.Input, opts.Output, opts.Width, opts.Height, jopts)
			logutil.FatalErr(err)
		}
	*/
}
