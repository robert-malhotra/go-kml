package kml

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// AltitudeMode specifies how altitude values are interpreted.
type AltitudeMode string

const (
	AltitudeModeClampToGround      AltitudeMode = "clampToGround"
	AltitudeModeRelativeToGround   AltitudeMode = "relativeToGround"
	AltitudeModeAbsolute           AltitudeMode = "absolute"
	AltitudeModeClampToSeaFloor    AltitudeMode = "clampToSeaFloor"    // gx extension
	AltitudeModeRelativeToSeaFloor AltitudeMode = "relativeToSeaFloor" // gx extension
)

// Geometry is an interface implemented by all geometry types.
type Geometry interface {
	geometryType() string
}

// Point represents a single geographic location.
type Point struct {
	ID           string       `xml:"id,attr,omitempty"`
	Extrude      bool         `xml:"extrude,omitempty"`
	AltitudeMode AltitudeMode `xml:"altitudeMode,omitempty"`
	Coordinates  Coordinate   `xml:"-"`
}

func (p *Point) geometryType() string {
	return "Point"
}

// MarshalXML implements custom XML marshaling for Point.
func (p *Point) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Point"

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if p.ID != "" {
		if err := e.EncodeElement(p.ID, xml.StartElement{Name: xml.Name{Local: "id"}}); err != nil {
			return err
		}
	}

	if p.Extrude {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "extrude"}}); err != nil {
			return err
		}
	}

	if p.AltitudeMode != "" {
		if err := e.EncodeElement(p.AltitudeMode, xml.StartElement{Name: xml.Name{Local: "altitudeMode"}}); err != nil {
			return err
		}
	}

	// Encode coordinates as a single coordinate string
	coordStr := p.Coordinates.String()
	if err := e.EncodeElement(coordStr, xml.StartElement{Name: xml.Name{Local: "coordinates"}}); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML implements custom XML unmarshaling for Point.
func (p *Point) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			p.ID = attr.Value
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "extrude":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				p.Extrude = v != 0
			case "altitudeMode", "gx:altitudeMode":
				var mode string
				if err := d.DecodeElement(&mode, &el); err != nil {
					return err
				}
				p.AltitudeMode = AltitudeMode(mode)
			case "coordinates":
				var coordStr string
				if err := d.DecodeElement(&coordStr, &el); err != nil {
					return err
				}
				coords, err := parseCoordinates(coordStr)
				if err != nil {
					return err
				}
				if len(coords) > 0 {
					p.Coordinates = coords[0]
				}
			default:
				// Skip unknown elements
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// LineString represents a connected set of line segments.
type LineString struct {
	ID           string       `xml:"id,attr,omitempty"`
	Extrude      bool         `xml:"extrude,omitempty"`
	Tessellate   bool         `xml:"tessellate,omitempty"`
	AltitudeMode AltitudeMode `xml:"altitudeMode,omitempty"`
	Coordinates  []Coordinate `xml:"-"`
}

func (ls *LineString) geometryType() string {
	return "LineString"
}

// MarshalXML implements custom XML marshaling for LineString.
func (ls *LineString) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "LineString"

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if ls.ID != "" {
		if err := e.EncodeElement(ls.ID, xml.StartElement{Name: xml.Name{Local: "id"}}); err != nil {
			return err
		}
	}

	if ls.Extrude {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "extrude"}}); err != nil {
			return err
		}
	}

	if ls.Tessellate {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "tessellate"}}); err != nil {
			return err
		}
	}

	if ls.AltitudeMode != "" {
		if err := e.EncodeElement(ls.AltitudeMode, xml.StartElement{Name: xml.Name{Local: "altitudeMode"}}); err != nil {
			return err
		}
	}

	// Encode coordinates
	coordStr := coordinatesToString(ls.Coordinates)
	if err := e.EncodeElement(coordStr, xml.StartElement{Name: xml.Name{Local: "coordinates"}}); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML implements custom XML unmarshaling for LineString.
func (ls *LineString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			ls.ID = attr.Value
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "extrude":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				ls.Extrude = v != 0
			case "tessellate":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				ls.Tessellate = v != 0
			case "altitudeMode", "gx:altitudeMode":
				var mode string
				if err := d.DecodeElement(&mode, &el); err != nil {
					return err
				}
				ls.AltitudeMode = AltitudeMode(mode)
			case "coordinates":
				var coordStr string
				if err := d.DecodeElement(&coordStr, &el); err != nil {
					return err
				}
				coords, err := parseCoordinates(coordStr)
				if err != nil {
					return err
				}
				ls.Coordinates = coords
			default:
				// Skip unknown elements
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// LinearRing defines a closed line string (first and last points identical).
type LinearRing struct {
	ID           string       `xml:"id,attr,omitempty"`
	Extrude      bool         `xml:"extrude,omitempty"`
	Tessellate   bool         `xml:"tessellate,omitempty"`
	AltitudeMode AltitudeMode `xml:"altitudeMode,omitempty"`
	Coordinates  []Coordinate `xml:"-"`
}

func (lr *LinearRing) geometryType() string {
	return "LinearRing"
}

// MarshalXML implements custom XML marshaling for LinearRing.
func (lr *LinearRing) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "LinearRing"

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if lr.ID != "" {
		if err := e.EncodeElement(lr.ID, xml.StartElement{Name: xml.Name{Local: "id"}}); err != nil {
			return err
		}
	}

	if lr.Extrude {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "extrude"}}); err != nil {
			return err
		}
	}

	if lr.Tessellate {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "tessellate"}}); err != nil {
			return err
		}
	}

	if lr.AltitudeMode != "" {
		if err := e.EncodeElement(lr.AltitudeMode, xml.StartElement{Name: xml.Name{Local: "altitudeMode"}}); err != nil {
			return err
		}
	}

	// Encode coordinates
	coordStr := coordinatesToString(lr.Coordinates)
	if err := e.EncodeElement(coordStr, xml.StartElement{Name: xml.Name{Local: "coordinates"}}); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML implements custom XML unmarshaling for LinearRing.
func (lr *LinearRing) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			lr.ID = attr.Value
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "extrude":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				lr.Extrude = v != 0
			case "tessellate":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				lr.Tessellate = v != 0
			case "altitudeMode", "gx:altitudeMode":
				var mode string
				if err := d.DecodeElement(&mode, &el); err != nil {
					return err
				}
				lr.AltitudeMode = AltitudeMode(mode)
			case "coordinates":
				var coordStr string
				if err := d.DecodeElement(&coordStr, &el); err != nil {
					return err
				}
				coords, err := parseCoordinates(coordStr)
				if err != nil {
					return err
				}
				lr.Coordinates = coords
			default:
				// Skip unknown elements
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// Polygon represents a polygon with an outer boundary and optional inner boundaries (holes).
type Polygon struct {
	ID              string       `xml:"id,attr,omitempty"`
	Extrude         bool         `xml:"extrude,omitempty"`
	Tessellate      bool         `xml:"tessellate,omitempty"`
	AltitudeMode    AltitudeMode `xml:"altitudeMode,omitempty"`
	OuterBoundary   LinearRing   `xml:"-"`
	InnerBoundaries []LinearRing `xml:"-"`
}

func (p *Polygon) geometryType() string {
	return "Polygon"
}

// MarshalXML implements custom XML marshaling for Polygon.
func (p *Polygon) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Polygon"

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if p.ID != "" {
		if err := e.EncodeElement(p.ID, xml.StartElement{Name: xml.Name{Local: "id"}}); err != nil {
			return err
		}
	}

	if p.Extrude {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "extrude"}}); err != nil {
			return err
		}
	}

	if p.Tessellate {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "tessellate"}}); err != nil {
			return err
		}
	}

	if p.AltitudeMode != "" {
		if err := e.EncodeElement(p.AltitudeMode, xml.StartElement{Name: xml.Name{Local: "altitudeMode"}}); err != nil {
			return err
		}
	}

	// Encode outerBoundaryIs
	outerStart := xml.StartElement{Name: xml.Name{Local: "outerBoundaryIs"}}
	if err := e.EncodeToken(outerStart); err != nil {
		return err
	}
	if err := e.Encode(&p.OuterBoundary); err != nil {
		return err
	}
	if err := e.EncodeToken(outerStart.End()); err != nil {
		return err
	}

	// Encode innerBoundaryIs elements
	for _, inner := range p.InnerBoundaries {
		innerStart := xml.StartElement{Name: xml.Name{Local: "innerBoundaryIs"}}
		if err := e.EncodeToken(innerStart); err != nil {
			return err
		}
		if err := e.Encode(&inner); err != nil {
			return err
		}
		if err := e.EncodeToken(innerStart.End()); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML implements custom XML unmarshaling for Polygon.
func (p *Polygon) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			p.ID = attr.Value
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "extrude":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				p.Extrude = v != 0
			case "tessellate":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				p.Tessellate = v != 0
			case "altitudeMode", "gx:altitudeMode":
				var mode string
				if err := d.DecodeElement(&mode, &el); err != nil {
					return err
				}
				p.AltitudeMode = AltitudeMode(mode)
			case "outerBoundaryIs":
				// Look for LinearRing inside
			outerBoundaryLoop:
				for {
					token, err := d.Token()
					if err != nil {
						return err
					}
					switch inner := token.(type) {
					case xml.StartElement:
						if inner.Name.Local == "LinearRing" {
							var ring LinearRing
							if err := d.DecodeElement(&ring, &inner); err != nil {
								return err
							}
							p.OuterBoundary = ring
						} else {
							if err := d.Skip(); err != nil {
								return err
							}
						}
					case xml.EndElement:
						break outerBoundaryLoop
					}
				}
			case "innerBoundaryIs":
				// Look for LinearRing inside
			innerBoundaryLoop:
				for {
					token, err := d.Token()
					if err != nil {
						return err
					}
					switch inner := token.(type) {
					case xml.StartElement:
						if inner.Name.Local == "LinearRing" {
							var ring LinearRing
							if err := d.DecodeElement(&ring, &inner); err != nil {
								return err
							}
							p.InnerBoundaries = append(p.InnerBoundaries, ring)
						} else {
							if err := d.Skip(); err != nil {
								return err
							}
						}
					case xml.EndElement:
						break innerBoundaryLoop
					}
				}
			default:
				// Skip unknown elements
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// MultiGeometry contains multiple geometry objects.
type MultiGeometry struct {
	ID         string     `xml:"id,attr,omitempty"`
	Geometries []Geometry `xml:"-"`
}

func (mg *MultiGeometry) geometryType() string {
	return "MultiGeometry"
}

// MarshalXML implements custom XML marshaling for MultiGeometry.
func (mg *MultiGeometry) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "MultiGeometry"

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if mg.ID != "" {
		if err := e.EncodeElement(mg.ID, xml.StartElement{Name: xml.Name{Local: "id"}}); err != nil {
			return err
		}
	}

	// Encode each geometry
	for _, geom := range mg.Geometries {
		if err := e.Encode(geom); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML implements custom XML unmarshaling for MultiGeometry.
func (mg *MultiGeometry) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			mg.ID = attr.Value
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch el := token.(type) {
		case xml.StartElement:
			var geom Geometry
			switch el.Name.Local {
			case "Point":
				var p Point
				if err := d.DecodeElement(&p, &el); err != nil {
					return err
				}
				geom = &p
			case "LineString":
				var ls LineString
				if err := d.DecodeElement(&ls, &el); err != nil {
					return err
				}
				geom = &ls
			case "LinearRing":
				var lr LinearRing
				if err := d.DecodeElement(&lr, &el); err != nil {
					return err
				}
				geom = &lr
			case "Polygon":
				var p Polygon
				if err := d.DecodeElement(&p, &el); err != nil {
					return err
				}
				geom = &p
			case "MultiGeometry":
				var mg2 MultiGeometry
				if err := d.DecodeElement(&mg2, &el); err != nil {
					return err
				}
				geom = &mg2
			default:
				// Skip unknown elements
				if err := d.Skip(); err != nil {
					return err
				}
				continue
			}
			if geom != nil {
				mg.Geometries = append(mg.Geometries, geom)
			}
		case xml.EndElement:
			return nil
		}
	}
}

// parseCoordinates parses a KML coordinate string into a slice of Coordinates.
// KML format: lon,lat[,alt] with whitespace-separated tuples.
func parseCoordinates(s string) ([]Coordinate, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	var coords []Coordinate
	tuples := strings.Fields(s)

	for _, tuple := range tuples {
		parts := strings.Split(tuple, ",")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid coordinate tuple: %s", tuple)
		}

		lon, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude in tuple %s: %w", tuple, err)
		}

		lat, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude in tuple %s: %w", tuple, err)
		}

		coord := Coordinate{Lon: lon, Lat: lat}

		if len(parts) >= 3 && parts[2] != "" {
			alt, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid altitude in tuple %s: %w", tuple, err)
			}
			coord.Alt = alt
		}

		coords = append(coords, coord)
	}

	return coords, nil
}

// coordinatesToString converts a slice of Coordinates to KML coordinate string format.
func coordinatesToString(coords []Coordinate) string {
	if len(coords) == 0 {
		return ""
	}

	var parts []string
	for _, c := range coords {
		parts = append(parts, c.String())
	}

	return strings.Join(parts, " ")
}
