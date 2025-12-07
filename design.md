# Design Document: `kmlgo` — A Go Library for KML Files

**Author:** [Your Name]  
**Date:** December 2024  
**Status:** Draft  

---

## Overview

`kmlgo` is a Go library for parsing, creating, and manipulating KML (Keyhole Markup Language) files. KML is an XML-based format for representing geographic data, widely used by Google Earth, Google Maps, and other GIS applications.

This library aims to provide a type-safe, idiomatic Go API for working with KML documents, supporting both reading existing files and programmatically generating new ones.

---

## Goals

- Provide complete coverage of the KML 2.2 specification (OGC standard)
- Offer a clean, idiomatic Go API with strong typing
- Support streaming parsing for large KML files
- Enable round-trip parsing (read → modify → write) without data loss
- Zero external dependencies for the core library
- Comprehensive error handling with meaningful error messages

## Non-Goals

- Real-time GPS/tracking integration
- Geometric calculations (distance, area, intersections)
- Format conversion (KML ↔ GeoJSON, Shapefile, etc.) — may be a separate package
- Network-based KML fetching (NetworkLink resolution)
- KMZ archive handling in core (available as subpackage)

---

## Core Types

### Document Structure

```go
package kml

// KML represents the root element of a KML document.
type KML struct {
    XMLName xml.Name `xml:"kml"`
    Xmlns   string   `xml:"xmlns,attr"`
    Feature Feature  // Document, Folder, or Placemark
}

// Document is a container for features and shared styles.
type Document struct {
    ID          string        `xml:"id,attr,omitempty"`
    Name        string        `xml:"name,omitempty"`
    Description string        `xml:"description,omitempty"`
    Open        bool          `xml:"open,omitempty"`
    Visibility  *bool         `xml:"visibility,omitempty"`
    Styles      []Style       `xml:"Style,omitempty"`
    StyleMaps   []StyleMap    `xml:"StyleMap,omitempty"`
    Features    []Feature     `xml:"-"` // Custom unmarshaling required
}

// Folder is a container for organizing features hierarchically.
type Folder struct {
    ID          string    `xml:"id,attr,omitempty"`
    Name        string    `xml:"name,omitempty"`
    Description string    `xml:"description,omitempty"`
    Open        bool      `xml:"open,omitempty"`
    Visibility  *bool     `xml:"visibility,omitempty"`
    Features    []Feature `xml:"-"`
}

// Placemark represents a geographic feature with geometry.
type Placemark struct {
    ID          string    `xml:"id,attr,omitempty"`
    Name        string    `xml:"name,omitempty"`
    Description string    `xml:"description,omitempty"`
    Visibility  *bool     `xml:"visibility,omitempty"`
    StyleURL    string    `xml:"styleUrl,omitempty"`
    Style       *Style    `xml:"Style,omitempty"`
    Geometry    Geometry  `xml:"-"` // Point, LineString, Polygon, etc.
    ExtendedData *ExtendedData `xml:"ExtendedData,omitempty"`
}
```

### Geometry Types

```go
// Geometry is an interface implemented by all geometry types.
type Geometry interface {
    geometryType() string
}

// Coordinate represents a single geographic coordinate.
type Coordinate struct {
    Lon float64 // Longitude in degrees
    Lat float64 // Latitude in degrees
    Alt float64 // Altitude in meters (optional)
}

// Point represents a single geographic location.
type Point struct {
    ID            string      `xml:"id,attr,omitempty"`
    Extrude       bool        `xml:"extrude,omitempty"`
    AltitudeMode  AltitudeMode `xml:"altitudeMode,omitempty"`
    Coordinates   Coordinate  `xml:"coordinates"`
}

// LineString represents a connected set of line segments.
type LineString struct {
    ID            string       `xml:"id,attr,omitempty"`
    Extrude       bool         `xml:"extrude,omitempty"`
    Tessellate    bool         `xml:"tessellate,omitempty"`
    AltitudeMode  AltitudeMode `xml:"altitudeMode,omitempty"`
    Coordinates   []Coordinate `xml:"coordinates"`
}

// LinearRing defines a closed line string (first and last points identical).
type LinearRing struct {
    ID            string       `xml:"id,attr,omitempty"`
    Extrude       bool         `xml:"extrude,omitempty"`
    Tessellate    bool         `xml:"tessellate,omitempty"`
    AltitudeMode  AltitudeMode `xml:"altitudeMode,omitempty"`
    Coordinates   []Coordinate `xml:"coordinates"`
}

// Polygon represents a polygon with an outer boundary and optional inner boundaries (holes).
type Polygon struct {
    ID              string       `xml:"id,attr,omitempty"`
    Extrude         bool         `xml:"extrude,omitempty"`
    Tessellate      bool         `xml:"tessellate,omitempty"`
    AltitudeMode    AltitudeMode `xml:"altitudeMode,omitempty"`
    OuterBoundary   LinearRing   `xml:"outerBoundaryIs>LinearRing"`
    InnerBoundaries []LinearRing `xml:"innerBoundaryIs>LinearRing,omitempty"`
}

// MultiGeometry contains multiple geometry objects.
type MultiGeometry struct {
    ID         string     `xml:"id,attr,omitempty"`
    Geometries []Geometry `xml:"-"`
}

// AltitudeMode specifies how altitude values are interpreted.
type AltitudeMode string

const (
    AltitudeModeClampToGround     AltitudeMode = "clampToGround"
    AltitudeModeRelativeToGround  AltitudeMode = "relativeToGround"
    AltitudeModeAbsolute          AltitudeMode = "absolute"
    AltitudeModeClampToSeaFloor   AltitudeMode = "clampToSeaFloor"   // gx extension
    AltitudeModeRelativeToSeaFloor AltitudeMode = "relativeToSeaFloor" // gx extension
)
```

### Styling Types

```go
// Style defines the appearance of features.
type Style struct {
    ID         string      `xml:"id,attr,omitempty"`
    IconStyle  *IconStyle  `xml:"IconStyle,omitempty"`
    LabelStyle *LabelStyle `xml:"LabelStyle,omitempty"`
    LineStyle  *LineStyle  `xml:"LineStyle,omitempty"`
    PolyStyle  *PolyStyle  `xml:"PolyStyle,omitempty"`
    BalloonStyle *BalloonStyle `xml:"BalloonStyle,omitempty"`
}

// Color represents a KML color (AABBGGRR format).
type Color struct {
    A, B, G, R uint8
}

// IconStyle specifies how icons are drawn.
type IconStyle struct {
    Color   Color   `xml:"color,omitempty"`
    Scale   float64 `xml:"scale,omitempty"`
    Heading float64 `xml:"heading,omitempty"`
    Icon    *Icon   `xml:"Icon,omitempty"`
    HotSpot *HotSpot `xml:"hotSpot,omitempty"`
}

// LineStyle specifies how lines are drawn.
type LineStyle struct {
    Color Color   `xml:"color,omitempty"`
    Width float64 `xml:"width,omitempty"`
}

// PolyStyle specifies how polygons are drawn.
type PolyStyle struct {
    Color   Color `xml:"color,omitempty"`
    Fill    *bool `xml:"fill,omitempty"`
    Outline *bool `xml:"outline,omitempty"`
}

// StyleMap maps between normal and highlight styles.
type StyleMap struct {
    ID    string `xml:"id,attr,omitempty"`
    Pairs []Pair `xml:"Pair"`
}

type Pair struct {
    Key      string `xml:"key"`      // "normal" or "highlight"
    StyleURL string `xml:"styleUrl"`
}
```

---

## Public API

### Parsing

```go
// Parse reads a KML document from an io.Reader.
func Parse(r io.Reader) (*KML, error)

// ParseFile reads a KML document from a file path.
func ParseFile(path string) (*KML, error)

// ParseBytes reads a KML document from a byte slice.
func ParseBytes(data []byte) (*KML, error)
```

### Writing

```go
// Write writes a KML document to an io.Writer.
func (k *KML) Write(w io.Writer) error

// WriteIndent writes a KML document with indentation.
func (k *KML) WriteIndent(w io.Writer, prefix, indent string) error

// WriteFile writes a KML document to a file.
func (k *KML) WriteFile(path string) error

// Bytes returns the KML document as a byte slice.
func (k *KML) Bytes() ([]byte, error)
```

### Builder API (Fluent Construction)

```go
// NewKML creates a new KML document.
func NewKML() *KMLBuilder

// Example usage:
kml := kml.NewKML().
    Document("My Map").
    Style("myStyle").
        IconStyle().Scale(1.2).Icon("http://...").Done().
        LineStyle().Color(kml.Red).Width(2).Done().
    Done().
    Folder("Points of Interest").
        Placemark("Coffee Shop").
            Point(-122.4194, 37.7749).
            StyleURL("#myStyle").
        Done().
        Placemark("Restaurant").
            Point(-122.4089, 37.7851).
        Done().
    Done().
    Build()
```

### Iteration and Query

```go
// Walk traverses all features in a KML document.
func (k *KML) Walk(fn func(Feature) error) error

// Placemarks returns all placemarks in the document.
func (k *KML) Placemarks() []*Placemark

// FindByID finds a feature by its ID.
func (k *KML) FindByID(id string) Feature

// Filter returns features matching a predicate.
func (k *KML) Filter(fn func(Feature) bool) []Feature

// Bounds returns the bounding box of all coordinates.
func (k *KML) Bounds() (sw, ne Coordinate)
```

### Coordinate Utilities

```go
// ParseCoordinates parses a KML coordinate string.
func ParseCoordinates(s string) ([]Coordinate, error)

// Coord creates a coordinate from lon, lat, and optional altitude.
func Coord(lon, lat float64, alt ...float64) Coordinate

// String returns the KML string representation of a coordinate.
func (c Coordinate) String() string
```

### Color Utilities

```go
// Predefined colors
var (
    White       = Color{255, 255, 255, 255}
    Black       = Color{255, 0, 0, 0}
    Red         = Color{255, 0, 0, 255}
    Green       = Color{255, 0, 255, 0}
    Blue        = Color{255, 255, 0, 0}
    Transparent = Color{0, 255, 255, 255}
)

// RGBA creates a color from standard RGBA values.
func RGBA(r, g, b, a uint8) Color

// ParseColor parses a KML color string (AABBGGRR).
func ParseColor(s string) (Color, error)

// Hex returns the KML hex representation.
func (c Color) Hex() string
```

---

## Package Structure

```
github.com/yourorg/kmlgo/
├── kml.go           # Root types and parsing
├── document.go      # Document, Folder types
├── placemark.go     # Placemark type
├── geometry.go      # All geometry types
├── style.go         # Style types
├── coordinate.go    # Coordinate parsing/formatting
├── color.go         # Color utilities
├── builder.go       # Fluent builder API
├── walk.go          # Traversal utilities
├── errors.go        # Error types
├── kmz/             # KMZ archive support (subpackage)
│   ├── kmz.go
│   └── reader.go
└── gx/              # Google Earth extensions (subpackage)
    └── gx.go
```

---

## Error Handling

```go
// ParseError represents an error during KML parsing.
type ParseError struct {
    Line    int
    Column  int
    Message string
    Cause   error
}

func (e *ParseError) Error() string
func (e *ParseError) Unwrap() error

// ValidationError represents invalid KML structure.
type ValidationError struct {
    Element string
    Field   string
    Message string
}

// Sentinel errors
var (
    ErrInvalidCoordinate = errors.New("kml: invalid coordinate format")
    ErrInvalidColor      = errors.New("kml: invalid color format")
    ErrEmptyDocument     = errors.New("kml: document contains no features")
    ErrMissingGeometry   = errors.New("kml: placemark has no geometry")
)
```

---

## Implementation Notes

### XML Unmarshaling Strategy

Due to KML's polymorphic nature (Features can be Document, Folder, or Placemark; Geometry can be Point, LineString, etc.), we implement custom `UnmarshalXML` methods:

```go
func (d *Document) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
    // Use xml.RawToken for polymorphic children
    // Dispatch based on element name
}
```

### Coordinate Parsing

KML coordinates use a specific format: `lon,lat[,alt]` with whitespace-separated tuples. The parser handles various edge cases including trailing whitespace and missing altitude values.

### Memory Efficiency

For large files, consider:
- Streaming parser option that yields features via channel
- Lazy loading of nested structures
- Coordinate pooling for high-volume scenarios

---

## Example Usage

### Reading and Modifying

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourorg/kmlgo"
)

func main() {
    doc, err := kml.ParseFile("places.kml")
    if err != nil {
        log.Fatal(err)
    }
    
    for _, pm := range doc.Placemarks() {
        if pt, ok := pm.Geometry.(*kml.Point); ok {
            fmt.Printf("%s: %.4f, %.4f\n", 
                pm.Name, pt.Coordinates.Lat, pt.Coordinates.Lon)
        }
    }
}
```

### Creating from Scratch

```go
package main

import (
    "os"
    
    "github.com/yourorg/kmlgo"
)

func main() {
    doc := kml.NewKML().
        Document("My Hiking Trails").
        Folder("Bay Area").
            Placemark("Mission Peak Trail").
                Description("Challenging hike with great views").
                LineString(
                    kml.Coord(-121.8814, 37.5126),
                    kml.Coord(-121.8806, 37.5098),
                    kml.Coord(-121.8753, 37.5073),
                ).
                Style().
                    LineStyle().Color(kml.Red).Width(3).Done().
                Done().
            Done().
        Done().
        Build()
    
    doc.WriteFile("trails.kml")
}
```

---

## Testing Strategy

- Unit tests for all parsing edge cases
- Round-trip tests (parse → write → parse) for data integrity
- Fuzz testing for coordinate and color parsing
- Benchmark tests for large file parsing
- Compatibility tests with real-world KML files from Google Earth

---

## Future Considerations

- **v1.1:** KMZ support via `kmlgo/kmz` subpackage
- **v1.2:** Google Earth extensions (`gx:` namespace) support  
- **v1.3:** Streaming parser for very large files
- **v2.0:** Potential GeoJSON conversion utilities

---

## References

- [OGC KML 2.2 Specification](https://www.ogc.org/standards/kml)
- [Google KML Documentation](https://developers.google.com/kml/documentation)
- [KML Reference](https://developers.google.com/kml/documentation/kmlreference)
