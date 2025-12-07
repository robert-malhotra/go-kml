package kml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
)

// Color represents a KML color in AABBGGRR format.
// Note: KML uses AABBGGRR (alpha, blue, green, red) format,
// which is different from the standard RGBA format.
type Color struct {
	A uint8 // Alpha
	B uint8 // Blue
	G uint8 // Green
	R uint8 // Red
}

// Predefined color constants in KML AABBGGRR format
var (
	White       = Color{255, 255, 255, 255}
	Black       = Color{255, 0, 0, 0}
	Red         = Color{255, 0, 0, 255}
	Green       = Color{255, 0, 255, 0}
	Blue        = Color{255, 255, 0, 0}
	Transparent = Color{0, 255, 255, 255}
)

// RGBA creates a Color from standard RGBA values.
// It converts from RGBA format to KML's AABBGGRR format.
func RGBA(r, g, b, a uint8) Color {
	return Color{
		A: a,
		B: b,
		G: g,
		R: r,
	}
}

// ParseColor parses a KML hex color string in AABBGGRR format.
// The string must be exactly 8 hexadecimal characters.
func ParseColor(s string) (Color, error) {
	if len(s) != 8 {
		return Color{}, errors.New("color string must be exactly 8 hexadecimal characters")
	}

	// Parse the hex string as a 32-bit unsigned integer
	val, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color string: %w", err)
	}

	// Extract AABBGGRR components
	return Color{
		A: uint8((val >> 24) & 0xFF),
		B: uint8((val >> 16) & 0xFF),
		G: uint8((val >> 8) & 0xFF),
		R: uint8(val & 0xFF),
	}, nil
}

// Hex returns the KML hex representation of the color in AABBGGRR format.
// The result is 8 lowercase hexadecimal characters.
func (c Color) Hex() string {
	return fmt.Sprintf("%02x%02x%02x%02x", c.A, c.B, c.G, c.R)
}

// MarshalXML implements xml.Marshaler for Color.
// It outputs the color as a hex string in AABBGGRR format.
func (c Color) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(c.Hex(), start)
}

// UnmarshalXML implements xml.Unmarshaler for Color.
// It parses a hex string in AABBGGRR format.
func (c *Color) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	parsed, err := ParseColor(s)
	if err != nil {
		return err
	}

	*c = parsed
	return nil
}
