package convertutil

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/gotilla/cmd/cmdutil"
)

/*
https://stackoverflow.com/questions/9586048/imagemagick-no-decode-delegate

brew uninstall graphicsmagick jpeg libtiff jasper
brew install imagemagick

*/

const (
	PdfWidth                      = 1950
	KindleWidth                   = 600
	PressDpi                      = 300
	WebDpi                        = 72
	ResolutionPixelsPerInchString = "PixelsPerInch"
	ConvertFormat                 = `convert %s -resize %sX%s %s %s %s`
	//ConvertFormat                 = `convert %s -resize %sx%s -set units %s -density %d %s`
)

type ResolutionUnits int

const (
	ResolutionUndefined           ResolutionUnits = iota // 0 convert cli value
	ResolutionPixelsPerInch                              // 1 convert cli value
	ResolutionPixelsPerCentimeter                        // 2 convert cli value
)

type ConvertParams struct {
	SourcePath        string
	OutputPath        string
	OutputWidth       int
	OutputHeight      int
	ResolutionDensity int
	ResolutionUnits   ResolutionUnits
}

/*
	cmdPDF := convertutil.ConvertCommand(
		filepathOrg, filepathPDF, convertutil.PdfWidth, 0, convertutil.PressDpi,
		convertutil.ResolutionPixelsPerInch)
*/

//func ConvertCommand(srcPath, outPath string, width, height, resolutionDensity, resolutionUnits int) string {

func ConvertCommand(params ConvertParams) string {
	widthStr := ""
	heightStr := ""
	densityStr := ""
	unitsStr := ""
	if params.OutputWidth > 0 {
		widthStr = strconv.Itoa(params.OutputWidth)
	}
	if params.OutputHeight > 0 {
		heightStr = strconv.Itoa(params.OutputHeight)
	}
	if params.ResolutionDensity > 0 {
		densityStr = fmt.Sprintf(" -density %d", params.ResolutionDensity)
	} /*
		if len(resolutionUnits) > 0 {
			unitsStr = fmt.Sprintf(` -set units %s`, resolutionUnits)
		}*/
	if params.ResolutionUnits == ResolutionPixelsPerInch {
		unitsStr = fmt.Sprintf(` -set units %s`, ResolutionPixelsPerInchString)
	}
	if 1 == 0 {
		quote := `"`
		params.SourcePath = quote + strings.TrimSpace(params.SourcePath) + quote
		params.OutputPath = quote + strings.TrimSpace(params.OutputPath) + quote
	}
	cmd := fmt.Sprintf(
		ConvertFormat,
		params.SourcePath,
		widthStr,
		heightStr,
		unitsStr,
		densityStr,
		params.OutputPath)

	return cmd

	//convert input.jpg -resize 1950x -set units PixelsPerInch -density 300 output.jpg
}

func ConvertToKindle(sourcePath, outputPath string) (bytes.Buffer, bytes.Buffer, error) {
	command := ConvertCommand(ConvertParams{
		SourcePath:        sourcePath,
		OutputPath:        outputPath,
		OutputWidth:       KindleWidth,
		OutputHeight:      0,
		ResolutionDensity: PressDpi,
		ResolutionUnits:   ResolutionPixelsPerInch})
	return cmdutil.ExecSimple(command)
}

func ConvertToPDF(sourcePath, outputPath string) (bytes.Buffer, bytes.Buffer, error) {
	command := ConvertCommand(ConvertParams{
		SourcePath:        sourcePath,
		OutputPath:        outputPath,
		OutputWidth:       KindleWidth,
		OutputHeight:      0,
		ResolutionDensity: PressDpi,
		ResolutionUnits:   ResolutionPixelsPerInch})
	fmt.Println(command)
	return cmdutil.ExecSimple(command)
}

const (
	convertErrorStatus = "exit status 1"
	convertError       = "convert: no decode delegate for this image format `'"
)

// CheckError removes the error if ReadImage doesn't work.
// exit status 1: convert: no decode delegate for this image format `' @ error/constitute.c/ReadImage/501.
func CheckError(err error, stderr bytes.Buffer) error {
	if err.Error() == convertErrorStatus &&
		strings.Index(stderr.String(), convertError) > -1 {
		return nil
	}
	return err
}
