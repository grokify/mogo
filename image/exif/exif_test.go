package exif_test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
	"time"

	"github.com/grokify/mogo/image/imageutil"
)

func createTestImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 255 / width),
				G: uint8(y * 255 / height),
				B: 0,
				A: 255,
			})
		}
	}
	return img
}

func TestExifToBytes(t *testing.T) {
	exif := New()
	exif.Make = "Test Camera"
	exif.Model = "Test Model"
	exif.DateTime = time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	exif.ExposureTime = "1/125"
	exif.FNumber = 2.8
	exif.ISOSpeed = 100
	exif.FocalLength = 50.0
	exif.Software = "Test Software"
	exif.Copyright = "Test Copyright"
	exif.Description = "Test Description"

	data, err := exif.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes failed: %v", err)
	}

	// Verify Exif header
	if len(data) < 8 {
		t.Fatal("Exif data too short")
	}
	if string(data[0:6]) != "Exif\x00\x00" {
		t.Error("Invalid Exif header")
	}
	if data[6] != 'I' || data[7] != 'I' {
		t.Error("Invalid byte order")
	}

	// Verify TIFF header
	if len(data) < 16 {
		t.Fatal("TIFF header too short")
	}
	if data[8] != 0x2A || data[9] != 0x00 {
		t.Error("Invalid TIFF identifier")
	}

	// Verify content
	if !bytes.Contains(data, []byte("Test Camera")) {
		t.Error("Make not found in Exif data")
	}
	if !bytes.Contains(data, []byte("Test Model")) {
		t.Error("Model not found in Exif data")
	}
	if !bytes.Contains(data, []byte("2024:03:15 12:00:00")) {
		t.Error("DateTime not found in Exif data")
	}
}

func TestExifFromBytes(t *testing.T) {
	// Create test Exif data
	exif := New()
	exif.Make = "Test Camera"
	exif.Model = "Test Model"
	exif.DateTime = time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	data, err := exif.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes failed: %v", err)
	}

	// Parse back
	parsed, err := FromBytes(data)
	if err != nil {
		t.Fatalf("FromBytes failed: %v", err)
	}

	// Verify parsed data
	if parsed.Make != exif.Make {
		t.Errorf("Make mismatch: got %q, want %q", parsed.Make, exif.Make)
	}
	if parsed.Model != exif.Model {
		t.Errorf("Model mismatch: got %q, want %q", parsed.Model, exif.Model)
	}
	if !parsed.DateTime.Equal(exif.DateTime) {
		t.Errorf("DateTime mismatch: got %v, want %v", parsed.DateTime, exif.DateTime)
	}
}

func TestExifWithJPEG(t *testing.T) {
	// Create test image
	img := createTestImage(100, 100)

	// Create Exif data
	exf := New()
	exf.Make = "Test Camera"
	exf.Model = "Test Model"
	exf.DateTime = time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	exf.Software = "Test Software"

	// Convert Exif to bytes
	exifData, err := exf.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes failed: %v", err)
	}

	// Encode JPEG with Exif
	buf := &bytes.Buffer{}
	if err := imageutil.EncodeJPEGWithExif(buf, img, &jpeg.Options{Quality: 90}, exifData); err != nil {
		t.Fatalf("EncodeJPEGWithExif failed: %v", err)
	}

	// Verify output
	data := buf.Bytes()
	if len(data) < 4 {
		t.Fatal("Output too short")
	}
	if data[0] != 0xFF || data[1] != 0xD8 {
		t.Error("Invalid JPEG header")
	}
	if data[2] != 0xFF || data[3] != 0xE1 {
		t.Error("Invalid Exif marker")
	}

	// Verify Exif content
	if !bytes.Contains(data, []byte("Test Camera")) {
		t.Error("Make not found in output")
	}
	if !bytes.Contains(data, []byte("Test Model")) {
		t.Error("Model not found in output")
	}
	if !bytes.Contains(data, []byte("2024:03:15 12:00:00")) {
		t.Error("DateTime not found in output")
	}
}
