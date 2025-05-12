// Package exif provides functionality for working with Exif metadata in images.
package exif

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/grokify/mogo/type/number"
)

// ExifTag represents a single Exif tag with its value
type ExifTag struct {
	ID    uint16
	Type  uint16
	Value any
}

// Exif represents Exif metadata in a human-readable format
type Exif struct {
	// Common Exif tags
	Make         string    // Camera manufacturer
	Model        string    // Camera model
	DateTime     time.Time // Date and time when the image was taken
	ExposureTime string    // Exposure time (e.g., "1/125")
	FNumber      float64   // F-number (aperture)
	ISOSpeed     uint16    // ISO speed rating
	FocalLength  float64   // Focal length in mm
	Flash        uint16    // Flash status
	Orientation  uint16    // Image orientation
	GPSLatitude  float64   // GPS latitude
	GPSLongitude float64   // GPS longitude
	GPSAltitude  float64   // GPS altitude
	Software     string    // Software used to create the image
	Copyright    string    // Copyright information
	Artist       string    // Artist/author
	Description  string    // Image description
	UserComment  string    // User comment
}

// Common Exif tag IDs
const (
	TagMake         uint16 = 0x010F
	TagModel        uint16 = 0x0110
	TagDateTime     uint16 = 0x0132
	TagExposureTime uint16 = 0x829A
	TagFNumber      uint16 = 0x829D
	TagISOSpeed     uint16 = 0x8827
	TagFocalLength  uint16 = 0x920A
	TagFlash        uint16 = 0x9209
	TagOrientation  uint16 = 0x0112
	TagGPSLatitude  uint16 = 0x0002
	TagGPSLongitude uint16 = 0x0004
	TagGPSAltitude  uint16 = 0x0006
	TagSoftware     uint16 = 0x0131
	TagCopyright    uint16 = 0x8298
	TagArtist       uint16 = 0x013B
	TagDescription  uint16 = 0x010E
	TagUserComment  uint16 = 0x9286
)

// Common Exif data types
const (
	TypeByte      uint16 = 1
	TypeASCII     uint16 = 2
	TypeShort     uint16 = 3
	TypeLong      uint16 = 4
	TypeRational  uint16 = 5
	TypeUndefined uint16 = 7
	TypeSLong     uint16 = 9
	TypeSRational uint16 = 10
)

// New creates a new Exif struct with default values
func New() *Exif {
	return &Exif{
		DateTime: time.Now(),
	}
}

// ToBytes converts the Exif struct to a byte array in the correct format
func (e *Exif) ToBytes() ([]byte, error) {
	var tags []ExifTag

	// Add tags if they have non-zero values
	if e.Make != "" {
		tags = append(tags, ExifTag{TagMake, TypeASCII, e.Make})
	}
	if e.Model != "" {
		tags = append(tags, ExifTag{TagModel, TypeASCII, e.Model})
	}
	if !e.DateTime.IsZero() {
		tags = append(tags, ExifTag{TagDateTime, TypeASCII, e.DateTime.Format("2006:01:02 15:04:05")})
	}
	if e.ExposureTime != "" {
		tags = append(tags, ExifTag{TagExposureTime, TypeASCII, e.ExposureTime})
	}
	if e.FNumber != 0 {
		tags = append(tags, ExifTag{TagFNumber, TypeRational, e.FNumber})
	}
	if e.ISOSpeed != 0 {
		tags = append(tags, ExifTag{TagISOSpeed, TypeShort, e.ISOSpeed})
	}
	if e.FocalLength != 0 {
		tags = append(tags, ExifTag{TagFocalLength, TypeRational, e.FocalLength})
	}
	if e.Flash != 0 {
		tags = append(tags, ExifTag{TagFlash, TypeShort, e.Flash})
	}
	if e.Orientation != 0 {
		tags = append(tags, ExifTag{TagOrientation, TypeShort, e.Orientation})
	}
	if e.GPSLatitude != 0 {
		tags = append(tags, ExifTag{TagGPSLatitude, TypeRational, e.GPSLatitude})
	}
	if e.GPSLongitude != 0 {
		tags = append(tags, ExifTag{TagGPSLongitude, TypeRational, e.GPSLongitude})
	}
	if e.GPSAltitude != 0 {
		tags = append(tags, ExifTag{TagGPSAltitude, TypeRational, e.GPSAltitude})
	}
	if e.Software != "" {
		tags = append(tags, ExifTag{TagSoftware, TypeASCII, e.Software})
	}
	if e.Copyright != "" {
		tags = append(tags, ExifTag{TagCopyright, TypeASCII, e.Copyright})
	}
	if e.Artist != "" {
		tags = append(tags, ExifTag{TagArtist, TypeASCII, e.Artist})
	}
	if e.Description != "" {
		tags = append(tags, ExifTag{TagDescription, TypeASCII, e.Description})
	}
	if e.UserComment != "" {
		tags = append(tags, ExifTag{TagUserComment, TypeASCII, e.UserComment})
	}

	// Build the Exif data
	buf := &bytes.Buffer{}

	// Write Exif header
	buf.WriteString("Exif\x00\x00")
	buf.Write([]byte("II")) // Intel byte order

	// Check for tag count overflow before writing as uint16
	if len(tags) > int(^uint16(0)) {
		return nil, fmt.Errorf("too many tags for Exif")
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(0x2A)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(8)); err != nil {
		return nil, err
	}
	if len(tags) > int(^uint16(0)) {
		return nil, fmt.Errorf("too many tags for Exif")
	} else if lenTagsU16, err := number.Itou16(len(tags)); err != nil {
		return nil, err
	} else if err := binary.Write(buf, binary.LittleEndian, lenTagsU16); err != nil {
		return nil, err
	}

	// Write tag entries (without values)
	tagEntries := make([]byte, len(tags)*12)
	for i, tag := range tags {
		offset := i * 12
		binary.LittleEndian.PutUint16(tagEntries[offset:offset+2], tag.ID)
		binary.LittleEndian.PutUint16(tagEntries[offset+2:offset+4], tag.Type)
		// Count and value offset will be filled in later
	}
	buf.Write(tagEntries)

	// Write values and update tag entries with correct value offsets
	valueOffset := buf.Len()
	for i, tag := range tags {
		offset := i * 12
		var count uint32
		var value uint32
		switch tag.Type {
		case TypeASCII:
			str := tag.Value.(string)
			strLen := len(str) + 1
			if strLen > int(^uint32(0)) {
				return nil, fmt.Errorf("string too long for Exif tag")
			}
			var err error
			count, err = number.Itou32(strLen)
			//count = uint32(strLen)
			if err != nil {
				return nil, err
			}
			if valueOffset > int(^uint32(0)) {
				return nil, fmt.Errorf("offset too large for Exif tag")
			}
			value, err = number.Itou32(valueOffset)
			//value = uint32(valueOffset)
			if err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			buf.WriteString(str)
			buf.WriteByte(0) // Null terminator
			valueOffset = buf.Len()
		case TypeShort:
			count = 1
			value = uint32(tag.Value.(uint16))
		case TypeLong:
			count = 1
			value = tag.Value.(uint32)
		case TypeRational:
			count = 1
			if valueOffset > int(^uint32(0)) {
				return nil, fmt.Errorf("offset too large for Exif tag")
			}
			var err error
			value, err = number.Itou32(valueOffset)
			//value = uint32(valueOffset)
			if err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, uint32(tag.Value.(float64)*100)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, uint32(100)); err != nil {
				return nil, err
			}
			valueOffset = buf.Len()
		default:
			return nil, fmt.Errorf("unsupported Exif type: %d", tag.Type)
		}
		binary.LittleEndian.PutUint32(tagEntries[offset+4:offset+8], count)
		binary.LittleEndian.PutUint32(tagEntries[offset+8:offset+12], value)
	}

	// Update tag entries in the buffer
	buf.Reset()
	buf.WriteString("Exif\x00\x00")
	buf.Write([]byte("II"))
	if err := binary.Write(buf, binary.LittleEndian, uint16(0x2A)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(8)); err != nil {
		return nil, err
	}
	if len(tags) > int(^uint16(0)) {
		return nil, fmt.Errorf("too many tags for Exif")
	} else if lenTagsU16, err := number.Itou16(len(tags)); err != nil {
		return nil, err
	} else if err := binary.Write(buf, binary.LittleEndian, lenTagsU16); err != nil {
		return nil, err
	}
	buf.Write(tagEntries)

	// Write values
	for _, tag := range tags {
		switch tag.Type {
		case TypeASCII:
			str := tag.Value.(string)
			buf.WriteString(str)
			buf.WriteByte(0) // Null terminator
		case TypeRational:
			if err := binary.Write(buf, binary.LittleEndian, uint32(tag.Value.(float64)*100)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, uint32(100)); err != nil {
				return nil, err
			}
		}
	}

	data := buf.Bytes()
	return data, nil
}

// FromBytes parses Exif data from a byte array into an Exif struct
func FromBytes(data []byte) (*Exif, error) {
	if len(data) < 14 {
		return nil, fmt.Errorf("data too short for Exif header and TIFF header")
	}
	/*
		fmt.Printf("Debug FromBytes:\n")
		fmt.Printf("  Data length: %d\n", len(data))
		fmt.Printf("  First 20 bytes: % x\n", data[:min(20, len(data))])
	*/
	// Check Exif header
	if string(data[0:6]) != "Exif\x00\x00" {
		return nil, fmt.Errorf("invalid Exif header")
	}

	// Check byte order
	isIntel := data[6] == 'I' && data[7] == 'I'
	if !isIntel {
		return nil, fmt.Errorf("only Intel byte order is supported")
	}

	// Parse TIFF header
	tiffHeaderStart := 6
	if len(data) < tiffHeaderStart+8 {
		return nil, fmt.Errorf("data too short for TIFF header")
	}
	/*
		fmt.Printf("  TIFF header start: %d\n", tiffHeaderStart)
		fmt.Printf("  TIFF bytes: % x\n", data[tiffHeaderStart:tiffHeaderStart+8])
	*/
	// Verify TIFF identifier
	if data[tiffHeaderStart+2] != 0x2A || data[tiffHeaderStart+3] != 0x00 {
		return nil, fmt.Errorf("invalid TIFF identifier")
	}

	// Get offset to first IFD (relative to TIFF header)
	ifdOffset := binary.LittleEndian.Uint32(data[tiffHeaderStart+4 : tiffHeaderStart+8])
	//fmt.Printf("  Raw IFD offset: %d\n", ifdOffset)
	actualIfdOffset := tiffHeaderStart + int(ifdOffset)
	// fmt.Printf("  Actual IFD offset: %d\n", actualIfdOffset)

	if actualIfdOffset >= len(data) {
		return nil, fmt.Errorf("invalid IFD offset: %d (data length: %d)", actualIfdOffset, len(data))
	}

	// Parse number of entries
	if len(data) < actualIfdOffset+2 {
		return nil, fmt.Errorf("data too short for IFD count")
	}
	numEntries := binary.LittleEndian.Uint16(data[actualIfdOffset : actualIfdOffset+2])
	// fmt.Printf("  Number of entries: %d\n", numEntries)

	// Create new Exif struct
	exif := New()

	// Parse each tag
	for i := uint16(0); i < numEntries; i++ {
		tagOffset := actualIfdOffset + 2 + int(i)*12
		if len(data) < tagOffset+12 {
			return nil, fmt.Errorf("data too short for tag")
		}

		tagID := binary.LittleEndian.Uint16(data[tagOffset : tagOffset+2])
		tagType := binary.LittleEndian.Uint16(data[tagOffset+2 : tagOffset+4])
		count := binary.LittleEndian.Uint32(data[tagOffset+4 : tagOffset+8])
		valueOffset := binary.LittleEndian.Uint32(data[tagOffset+8 : tagOffset+12])

		// fmt.Printf("  Tag %d: ID=%04x Type=%d Count=%d ValueOffset=%d\n", i, tagID, tagType, count, valueOffset)

		// Parse value based on tag type
		switch tagType {
		case TypeASCII:
			if int(valueOffset) >= len(data) {
				continue
			}
			// Print the next 20 bytes at the value offset for debugging
			// fmt.Printf("    ASCII bytes at offset %d: % x\n", valueOffset, data[valueOffset:min(int(valueOffset)+20, len(data))])
			// Find null terminator
			end := int(valueOffset)
			for end < len(data) && data[end] != 0 {
				end++
			}
			value := string(data[valueOffset:end])

			// Set field based on tag ID
			switch tagID {
			case TagMake:
				exif.Make = value
			case TagModel:
				exif.Model = value
			case TagDateTime:
				if t, err := time.Parse("2006:01:02 15:04:05", value); err == nil {
					exif.DateTime = t
				}
			case TagSoftware:
				exif.Software = value
			case TagCopyright:
				exif.Copyright = value
			case TagArtist:
				exif.Artist = value
			case TagDescription:
				exif.Description = value
			case TagUserComment:
				exif.UserComment = value
			}
		case TypeShort:
			if count == 1 {
				value := binary.LittleEndian.Uint16(data[tagOffset+8 : tagOffset+10])
				switch tagID {
				case TagISOSpeed:
					exif.ISOSpeed = value
				case TagFlash:
					exif.Flash = value
				case TagOrientation:
					exif.Orientation = value
				}
			}
		case TypeRational:
			if count == 1 && int(valueOffset)+8 <= len(data) {
				numerator := binary.LittleEndian.Uint32(data[valueOffset : valueOffset+4])
				denominator := binary.LittleEndian.Uint32(data[valueOffset+4 : valueOffset+8])
				if denominator != 0 {
					value := float64(numerator) / float64(denominator)
					switch tagID {
					case TagFNumber:
						exif.FNumber = value
					case TagFocalLength:
						exif.FocalLength = value
					case TagGPSLatitude:
						exif.GPSLatitude = value
					case TagGPSLongitude:
						exif.GPSLongitude = value
					case TagGPSAltitude:
						exif.GPSAltitude = value
					}
				}
			}
		}
	}

	return exif, nil
}

/*
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
*/
