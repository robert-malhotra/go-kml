package kml

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// Coordinate represents a single geographic coordinate.
type Coordinate struct {
	Lon float64 // Longitude in degrees
	Lat float64 // Latitude in degrees
	Alt float64 // Altitude in meters (optional, 0 if not specified)
}

// Coord creates a coordinate from longitude, latitude, and optional altitude.
// If altitude is not provided, it defaults to 0.
func Coord(lon, lat float64, alt ...float64) Coordinate {
	c := Coordinate{
		Lon: lon,
		Lat: lat,
		Alt: 0,
	}
	if len(alt) > 0 {
		c.Alt = alt[0]
	}
	return c
}

// String returns the KML string representation of a coordinate.
// Returns "lon,lat,alt" if altitude is non-zero, otherwise "lon,lat".
func (c Coordinate) String() string {
	if c.Alt == 0 {
		return fmt.Sprintf("%g,%g", c.Lon, c.Lat)
	}
	return fmt.Sprintf("%g,%g,%g", c.Lon, c.Lat, c.Alt)
}

// ParseCoordinates parses a KML coordinate string into a slice of Coordinates.
// KML format is: "lon,lat[,alt] lon,lat[,alt] ..." (whitespace-separated tuples,
// comma-separated values within each tuple).
//
// The function handles edge cases including:
// - Leading/trailing whitespace
// - Multiple spaces between coordinates
// - Missing altitude values (defaults to 0)
// - Empty strings
func ParseCoordinates(s string) ([]Coordinate, error) {
	// Trim leading and trailing whitespace
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("%w: empty coordinate string", ErrInvalidCoordinate)
	}

	// Split by whitespace (handles multiple spaces, tabs, newlines)
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return nil, fmt.Errorf("%w: empty coordinate string", ErrInvalidCoordinate)
	}

	coords := make([]Coordinate, 0, len(fields))

	for _, field := range fields {
		// Each field should be "lon,lat[,alt]"
		parts := strings.Split(field, ",")
		if len(parts) < 2 || len(parts) > 3 {
			return nil, fmt.Errorf("%w: expected 2 or 3 values, got %d in %q",
				ErrInvalidCoordinate, len(parts), field)
		}

		// Parse longitude
		lon, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid longitude %q: %v",
				ErrInvalidCoordinate, parts[0], err)
		}

		// Parse latitude
		lat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid latitude %q: %v",
				ErrInvalidCoordinate, parts[1], err)
		}

		// Parse altitude (optional)
		alt := 0.0
		if len(parts) == 3 {
			alt, err = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid altitude %q: %v",
					ErrInvalidCoordinate, parts[2], err)
			}
		}

		coords = append(coords, Coordinate{
			Lon: lon,
			Lat: lat,
			Alt: alt,
		})
	}

	return coords, nil
}

// Coordinates is a slice of Coordinate values that implements custom XML
// marshaling/unmarshaling for the KML coordinate string format.
type Coordinates []Coordinate

// MarshalXML encodes Coordinates into KML coordinate string format.
func (c Coordinates) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(c) == 0 {
		return e.EncodeElement("", start)
	}

	var sb strings.Builder
	for i, coord := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(coord.String())
	}

	return e.EncodeElement(sb.String(), start)
}

// UnmarshalXML decodes KML coordinate string format into Coordinates.
func (c *Coordinates) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	coords, err := ParseCoordinates(s)
	if err != nil {
		return err
	}

	*c = coords
	return nil
}
