package imageutil

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestNewWriterExif(t *testing.T) {
	tests := []struct {
		name    string
		exif    []byte
		wantErr bool
	}{
		{
			name:    "nil exif",
			exif:    nil,
			wantErr: false,
		},
		{
			name:    "empty exif",
			exif:    []byte{},
			wantErr: false,
		},
		{
			name:    "small exif",
			exif:    []byte{1, 2, 3, 4},
			wantErr: false,
		},
		{
			name:    "medium exif",
			exif:    bytes.Repeat([]byte{1}, 1000),
			wantErr: false,
		},
		{
			name:    "large exif",
			exif:    bytes.Repeat([]byte{1}, 65000),
			wantErr: false,
		},
		{
			name:    "too large exif",
			exif:    bytes.Repeat([]byte{1}, 70000),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing case: %s", tt.name)
			if tt.exif != nil {
				t.Logf("Exif length: %d bytes", len(tt.exif))
			}

			// Test with new implementation
			buf1 := &bytes.Buffer{}
			writer1, err1 := newWriterExif(buf1, tt.exif)
			if (err1 != nil) != tt.wantErr {
				t.Errorf("newWriterExif() error = %v, wantErr %v", err1, tt.wantErr)
				return
			}
			if err1 != nil {
				t.Logf("Got expected error: %v", err1)
				return // Skip further testing if we got an expected error
			}

			// Write some test data to verify the writer works
			testData := []byte("test data")
			if _, err := writer1.Write(testData); err != nil {
				t.Errorf("writer1.Write() error = %v", err)
				return
			}

			// Verify the output format
			data := buf1.Bytes()

			// Print first few bytes for debugging
			if len(data) > 0 {
				t.Logf("Output length: %d bytes", len(data))
				t.Logf("First 10 bytes: %v", data[:min(10, len(data))])
				t.Logf("First 10 bytes (hex): %x", data[:min(10, len(data))])
			}

			// Check SOI marker
			if len(data) < 2 {
				t.Errorf("Data too short: got %d bytes, want at least 2", len(data))
				return
			}
			if data[0] != JPEGMarkerPrefix {
				t.Errorf("Invalid JPEG marker prefix: got 0x%02x, want 0x%02x", data[0], JPEGMarkerPrefix)
			}
			if data[1] != JPEGMarkerSOI {
				t.Errorf("Invalid SOI marker: got 0x%02x, want 0x%02x", data[1], JPEGMarkerSOI)
			}

			// If we have Exif data, verify its format
			if tt.exif != nil {
				// Check Exif marker
				if len(data) < 4 {
					t.Errorf("Data too short for Exif: got %d bytes, want at least 4", len(data))
					return
				}
				if data[2] != JPEGMarkerPrefix {
					t.Errorf("Invalid Exif marker prefix: got 0x%02x, want 0x%02x", data[2], JPEGMarkerPrefix)
				}
				if data[3] != JPEGMarkerExif {
					t.Errorf("Invalid Exif marker: got 0x%02x, want 0x%02x", data[3], JPEGMarkerExif)
				}

				// Check length bytes
				markerLen := 2 + len(tt.exif)
				if len(data) < 6 {
					t.Errorf("Data too short for length bytes: got %d bytes, want at least 6", len(data))
					return
				}
				gotHigh := data[4]
				gotLow := data[5]
				wantHigh := byte(markerLen >> 8)
				wantLow := byte(markerLen & 0xFF)
				if gotHigh != wantHigh || gotLow != wantLow {
					t.Errorf("Invalid length bytes: got [0x%02x 0x%02x], want [0x%02x 0x%02x] (length=%d)",
						gotHigh, gotLow, wantHigh, wantLow, markerLen)
				}

				// Check Exif data
				exifStart := 6
				exifEnd := exifStart + len(tt.exif)
				if len(data) < exifEnd {
					t.Errorf("Data too short for Exif content: got %d bytes, want at least %d", len(data), exifEnd)
					return
				}
				if !bytes.Equal(data[exifStart:exifEnd], tt.exif) {
					t.Logf("Exif data mismatch at position %d", exifStart)
					t.Logf("Expected: %x", tt.exif)
					t.Logf("Got:      %x", data[exifStart:exifEnd])
					t.Errorf("Exif data mismatch")
				}

				// Check test data - it should be written after the Exif data
				if !bytes.Equal(data[exifEnd:], testData) {
					t.Logf("Test data mismatch at position %d", exifEnd)
					t.Logf("Expected: %x", testData)
					t.Logf("Got:      %x", data[exifEnd:])
					t.Errorf("Test data mismatch")
				}
			} else {
				// For nil Exif, the test data should be written after the SOI marker
				if !bytes.Equal(data[2:], testData) {
					t.Logf("Test data mismatch with nil Exif at position 2")
					t.Logf("Expected: %x", testData)
					t.Logf("Got:      %x", data[2:])
					t.Errorf("Test data mismatch with nil Exif")
				}
			}
		})
	}
}

func TestExifEndToEnd(t *testing.T) {
	// Create a test JPEG file
	testImg := createTestJPEG(t)
	defer os.Remove(testImg)

	// Test Exif data
	exifData := []byte{
		0x45, 0x78, 0x69, 0x66, 0x00, 0x00, // Exif header
		0x49, 0x49, // Intel byte order
		0x2A, 0x00, // TIFF identifier
		0x08, 0x00, 0x00, 0x00, // Offset to first IFD
		0x01, 0x00, // Number of entries
		0x01, 0x00, // Tag: ImageDescription
		0x02, 0x00, // Type: ASCII
		0x0A, 0x00, 0x00, 0x00, // Count
		0x0A, 0x00, 0x00, 0x00, // Value offset
		'T', 'e', 's', 't', ' ', 'E', 'x', 'i', 'f', 0x00, // "Test Exif\0"
	}

	// Read the test image
	img, err := readJPEGFile(testImg)
	if err != nil {
		t.Fatalf("Failed to read test image: %v", err)
	}

	// Encode JPEG to buffer
	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("Failed to encode JPEG: %v", err)
	}
	jpegData := buf.Bytes()
	if len(jpegData) < 2 || jpegData[0] != 0xFF || jpegData[1] != 0xD8 {
		t.Fatalf("Not a valid JPEG SOI")
	}

	// Build Exif segment
	markerLen := 2 + len(exifData)
	if markerLen > 0xFFFF {
		t.Fatalf("Exif too large")
	}
	exifSegment := []byte{0xFF, 0xE1, byte(markerLen >> 8), byte(markerLen & 0xFF)}
	exifSegment = append(exifSegment, exifData...)

	// Insert Exif after SOI
	final := append([]byte{}, jpegData[:2]...)
	final = append(final, exifSegment...)
	final = append(final, jpegData[2:]...)

	// Debug output
	t.Logf("Final data size: %d bytes", len(final))
	if len(final) > 0 {
		t.Logf("First 20 bytes: %x", final[:min(20, len(final))])
	}

	// Verify JPEG structure
	if len(final) < 2 {
		t.Fatal("File too short")
	}
	if final[0] != JPEGMarkerPrefix || final[1] != JPEGMarkerSOI {
		t.Error("Invalid JPEG header")
		t.Logf("Got header: %x %x", final[0], final[1])
		t.Logf("Expected: %x %x", JPEGMarkerPrefix, JPEGMarkerSOI)
	}

	// Verify Exif marker
	if len(final) < 4 {
		t.Fatal("File too short for Exif marker")
	}
	if final[2] != JPEGMarkerPrefix || final[3] != JPEGMarkerExif {
		t.Error("Invalid Exif marker")
		t.Logf("Got marker: %x %x", final[2], final[3])
		t.Logf("Expected: %x %x", JPEGMarkerPrefix, JPEGMarkerExif)
	}

	// Verify Exif data
	if !bytes.Contains(final, exifData) {
		t.Error("Exif data not found in output")
		t.Logf("Expected Exif data: %x", exifData)
		t.Logf("File contents: %x", final[:min(100, len(final))])
	}

	// Now try with a file
	outFile := filepath.Join(filepath.Dir(testImg), "test_with_exif.jpg")
	defer os.Remove(outFile)

	out, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer out.Close()

	// Write the data to the file
	if _, err := out.Write(final); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}
}

func createTestJPEG(t *testing.T) string {
	// Create a simple test image
	img := createTestImage(100, 100)

	// Create test file
	testFile := filepath.Join("testdata", "test.jpg")
	if err := os.MkdirAll(filepath.Dir(testFile), 0755); err != nil {
		t.Fatalf("Failed to create testdata directory: %v", err)
	}

	// Save as JPEG
	out, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer out.Close()

	if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	return testFile
}

func readJPEGFile(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return jpeg.Decode(f)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(float64(x) * 255.0 / float64(width)),
				G: uint8(float64(y) * 255.0 / float64(height)),
				B: 0,
				A: 255,
			})
		}
	}
	return img
}

func TestEncodeJPEGWithExif_Unit(t *testing.T) {
	img := createTestImage(10, 10)
	exif := []byte("ExifUnitTestData")
	buf := &bytes.Buffer{}
	err := EncodeJPEGWithExif(buf, img, &jpeg.Options{Quality: 80}, exif)
	if err != nil {
		t.Fatalf("EncodeJPEGWithExif failed: %v", err)
	}
	data := buf.Bytes()
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 || data[2] != 0xFF || data[3] != 0xE1 {
		t.Errorf("Output does not start with SOI+Exif marker: %x", data[:4])
	}
	if !bytes.Contains(data, exif) {
		t.Errorf("Exif data not found in output")
	}
}

func TestEncodeJPEGWithExif_EndToEnd(t *testing.T) {
	img := createTestImage(50, 50)
	exif := []byte("ExifEndToEndTestData")
	outFile := "testdata/endtoend_with_exif.jpg"
	defer os.Remove(outFile)
	f, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()
	if err := EncodeJPEGWithExif(f, img, &jpeg.Options{Quality: 90}, exif); err != nil {
		t.Fatalf("EncodeJPEGWithExif failed: %v", err)
	}
	f.Close()
	// Read back and check
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 || data[2] != 0xFF || data[3] != 0xE1 {
		t.Errorf("Output does not start with SOI+Exif marker: %x", data[:4])
	}
	if !bytes.Contains(data, exif) {
		t.Errorf("Exif data not found in output")
	}
}
