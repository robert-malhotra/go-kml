# kmlgo

A Go library for parsing, creating, and manipulating KML (Keyhole Markup Language) files.

[![Go Reference](https://pkg.go.dev/badge/github.com/robert-malhotra/go-kml.svg)](https://pkg.go.dev/github.com/robert-malhotra/go-kml)
[![Go Report Card](https://goreportcard.com/badge/github.com/robert-malhotra/go-kml)](https://goreportcard.com/report/github.com/robert-malhotra/go-kml)

## Features

- **Complete KML 2.2 Support** - Parse and generate KML documents following the OGC standard
- **Type-Safe API** - Strongly typed Go structs for all KML elements
- **Fluent Builder** - Chainable builder API for creating KML documents programmatically
- **Round-Trip Safe** - Parse, modify, and write KML without data loss
- **Zero Dependencies** - Uses only the Go standard library
- **Traversal Utilities** - Walk, filter, and query KML document trees
- **Real-World Tested** - Validated against Google's official KML samples

## Installation

```bash
go get github.com/robert-malhotra/go-kml
```

## Quick Start

### Parsing a KML File

```go
package main

import (
    "fmt"
    "log"

    "github.com/robert-malhotra/go-kml"
)

func main() {
    // Parse from file
    doc, err := kml.ParseFile("places.kml")
    if err != nil {
        log.Fatal(err)
    }

    // Iterate through all placemarks
    for _, pm := range doc.Placemarks() {
        if pt, ok := pm.Geometry.(*kml.Point); ok {
            fmt.Printf("%s: %.4f, %.4f\n", pm.Name, pt.Coordinates.Lat, pt.Coordinates.Lon)
        }
    }
}
```

### Creating a KML Document

```go
package main

import (
    "github.com/robert-malhotra/go-kml"
)

func main() {
    // Using the fluent builder API
    doc := kml.NewKMLBuilder().
        Document("My Places").
        Style("redLine").
            LineStyle().Color(kml.Red).Width(3).Done().
        Done().
        Folder("Favorites").
            Placemark("Coffee Shop").
                Description("Best coffee in town").
                Point(-122.4194, 37.7749).
                StyleURL("#redLine").
            Done().
            Placemark("Home").
                Point(-122.4089, 37.7851).
            Done().
        Done().
        Build()

    // Write to file
    doc.WriteFile("places.kml")
}
```

### Parsing from Different Sources

```go
// From file
doc, err := kml.ParseFile("document.kml")

// From io.Reader
doc, err := kml.Parse(reader)

// From byte slice
doc, err := kml.ParseBytes(data)
```

### Writing KML

```go
// To file
err := doc.WriteFile("output.kml")

// To io.Writer
err := doc.Write(writer)

// To io.Writer with indentation
err := doc.WriteIndent(writer, "", "  ")

// To byte slice
bytes, err := doc.Bytes()
```

## Supported KML Elements

### Geometry Types

| Type | Description |
|------|-------------|
| `Point` | Single geographic location |
| `LineString` | Connected line segments |
| `LinearRing` | Closed line string |
| `Polygon` | Polygon with optional holes |
| `MultiGeometry` | Collection of geometries |

### Container Types

| Type | Description |
|------|-------------|
| `KML` | Root element |
| `Document` | Container with shared styles |
| `Folder` | Hierarchical organization |
| `Placemark` | Geographic feature |

### Style Types

| Type | Description |
|------|-------------|
| `Style` | Style definition |
| `StyleMap` | Normal/highlight style pairs |
| `IconStyle` | Point icon appearance |
| `LabelStyle` | Label appearance |
| `LineStyle` | Line appearance |
| `PolyStyle` | Polygon appearance |
| `BalloonStyle` | Info balloon appearance |

### Other Types

| Type | Description |
|------|-------------|
| `ExtendedData` | Custom data fields |
| `Data` | Name-value pairs |
| `Coordinate` | Geographic coordinate |
| `Color` | KML color (AABBGGRR format) |

## API Reference

### Parsing Functions

```go
// Parse reads a KML document from an io.Reader
func Parse(r io.Reader) (*KML, error)

// ParseFile reads a KML document from a file path
func ParseFile(path string) (*KML, error)

// ParseBytes reads a KML document from a byte slice
func ParseBytes(data []byte) (*KML, error)
```

### KML Methods

```go
// Write writes a KML document to an io.Writer
func (k *KML) Write(w io.Writer) error

// WriteIndent writes with custom indentation
func (k *KML) WriteIndent(w io.Writer, prefix, indent string) error

// WriteFile writes to a file
func (k *KML) WriteFile(path string) error

// Bytes returns the KML as a byte slice
func (k *KML) Bytes() ([]byte, error)

// Walk traverses all features depth-first
func (k *KML) Walk(fn func(Feature) error) error

// Placemarks returns all placemarks in the document
func (k *KML) Placemarks() []*Placemark

// FindByID finds a feature by its ID attribute
func (k *KML) FindByID(id string) Feature

// Filter returns features matching a predicate
func (k *KML) Filter(fn func(Feature) bool) []Feature

// Bounds returns the bounding box of all coordinates
func (k *KML) Bounds() (sw, ne Coordinate)
```

### Coordinate Utilities

```go
// Create a coordinate
coord := kml.Coord(-122.4194, 37.7749)           // lon, lat
coord := kml.Coord(-122.4194, 37.7749, 100.0)    // lon, lat, altitude

// Parse coordinate string
coords, err := kml.ParseCoordinates("-122.0,37.0 -121.0,38.0")

// Convert to string
str := coord.String()  // "-122.4194,37.7749"
```

### Color Utilities

```go
// Predefined colors
kml.White       // Opaque white
kml.Black       // Opaque black
kml.Red         // Opaque red
kml.Green       // Opaque green
kml.Blue        // Opaque blue
kml.Transparent // Fully transparent

// Create from RGBA (converts to KML's AABBGGRR format)
color := kml.RGBA(255, 0, 0, 128)  // Semi-transparent red

// Parse KML hex color
color, err := kml.ParseColor("ff0000ff")  // Opaque red

// Get hex representation
hex := color.Hex()  // "ff0000ff"
```

## Builder API

The builder API provides a fluent interface for constructing KML documents:

```go
doc := kml.NewKMLBuilder().
    Document("Document Name").
        Description("Document description").
        Open(true).

        // Add styles
        Style("styleId").
            IconStyle().
                Scale(1.2).
                Icon("http://maps.google.com/mapfiles/kml/shapes/icon.png").
            Done().
            LineStyle().
                Color(kml.RGBA(255, 0, 0, 255)).
                Width(2.0).
            Done().
            PolyStyle().
                Color(kml.RGBA(0, 255, 0, 128)).
                Fill(true).
                Outline(true).
            Done().
        Done().

        // Add folders
        Folder("Folder Name").
            Description("Folder description").

            // Add placemarks
            Placemark("Point Name").
                Description("A point of interest").
                StyleURL("#styleId").
                Point(-122.0, 37.0).
            Done().

            Placemark("Line Name").
                LineString(
                    kml.Coord(-122.0, 37.0),
                    kml.Coord(-121.0, 38.0),
                    kml.Coord(-120.0, 37.5),
                ).
            Done().

            Placemark("Polygon Name").
                Polygon(
                    []kml.Coordinate{  // Outer boundary
                        kml.Coord(-122.0, 37.0),
                        kml.Coord(-121.0, 37.0),
                        kml.Coord(-121.0, 38.0),
                        kml.Coord(-122.0, 38.0),
                        kml.Coord(-122.0, 37.0),
                    },
                    []kml.Coordinate{  // Inner boundary (hole)
                        kml.Coord(-121.8, 37.2),
                        kml.Coord(-121.2, 37.2),
                        kml.Coord(-121.2, 37.8),
                        kml.Coord(-121.8, 37.8),
                        kml.Coord(-121.8, 37.2),
                    },
                ).
            Done().

            // Nested folder
            Folder("Subfolder").
                Placemark("Nested Point").
                    Point(-119.0, 36.0).
                Done().
            Done().

        Done().
    Done().
    Build()
```

## Traversing Documents

### Walk All Features

```go
err := doc.Walk(func(f kml.Feature) error {
    switch feat := f.(type) {
    case *kml.Document:
        fmt.Printf("Document: %s\n", feat.Name)
    case *kml.Folder:
        fmt.Printf("Folder: %s\n", feat.Name)
    case *kml.Placemark:
        fmt.Printf("Placemark: %s\n", feat.Name)
    }
    return nil
})
```

### Get All Placemarks

```go
placemarks := doc.Placemarks()
for _, pm := range placemarks {
    fmt.Printf("- %s\n", pm.Name)
}
```

### Filter Features

```go
// Find all placemarks with "Park" in the name
parks := doc.Filter(func(f kml.Feature) bool {
    if pm, ok := f.(*kml.Placemark); ok {
        return strings.Contains(pm.Name, "Park")
    }
    return false
})
```

### Find by ID

```go
feature := doc.FindByID("myPlacemark")
if pm, ok := feature.(*kml.Placemark); ok {
    fmt.Printf("Found: %s\n", pm.Name)
}
```

### Calculate Bounds

```go
sw, ne := doc.Bounds()
fmt.Printf("Southwest: %.4f, %.4f\n", sw.Lat, sw.Lon)
fmt.Printf("Northeast: %.4f, %.4f\n", ne.Lat, ne.Lon)
```

## Working with Geometry

### Point

```go
point := &kml.Point{
    Coordinates:  kml.Coord(-122.0, 37.0, 100.0),
    AltitudeMode: kml.AltitudeModeAbsolute,
    Extrude:      true,
}
```

### LineString

```go
line := &kml.LineString{
    Coordinates: []kml.Coordinate{
        kml.Coord(-122.0, 37.0),
        kml.Coord(-121.0, 38.0),
        kml.Coord(-120.0, 37.5),
    },
    Tessellate:   true,
    AltitudeMode: kml.AltitudeModeClampToGround,
}
```

### Polygon with Holes

```go
polygon := &kml.Polygon{
    OuterBoundary: kml.LinearRing{
        Coordinates: []kml.Coordinate{
            kml.Coord(-122.0, 37.0),
            kml.Coord(-121.0, 37.0),
            kml.Coord(-121.0, 38.0),
            kml.Coord(-122.0, 38.0),
            kml.Coord(-122.0, 37.0),  // Close the ring
        },
    },
    InnerBoundaries: []kml.LinearRing{
        {
            Coordinates: []kml.Coordinate{
                kml.Coord(-121.8, 37.2),
                kml.Coord(-121.2, 37.2),
                kml.Coord(-121.2, 37.8),
                kml.Coord(-121.8, 37.8),
                kml.Coord(-121.8, 37.2),
            },
        },
    },
    Extrude:      true,
    AltitudeMode: kml.AltitudeModeRelativeToGround,
}
```

### MultiGeometry

```go
multi := &kml.MultiGeometry{
    Geometries: []kml.Geometry{
        &kml.Point{Coordinates: kml.Coord(-122.0, 37.0)},
        &kml.LineString{
            Coordinates: []kml.Coordinate{
                kml.Coord(-122.0, 37.0),
                kml.Coord(-121.0, 38.0),
            },
        },
    },
}
```

## Altitude Modes

```go
kml.AltitudeModeClampToGround      // Default, ignore altitude
kml.AltitudeModeRelativeToGround   // Altitude relative to ground
kml.AltitudeModeAbsolute           // Altitude relative to sea level
kml.AltitudeModeClampToSeaFloor    // Clamp to sea floor (gx extension)
kml.AltitudeModeRelativeToSeaFloor // Relative to sea floor (gx extension)
```

## Extended Data

```go
placemark := &kml.Placemark{
    Name: "Golf Course",
    ExtendedData: &kml.ExtendedData{
        Data: []kml.Data{
            {Name: "holeNumber", Value: "1"},
            {Name: "par", Value: "4"},
            {Name: "yardage", Value: "380"},
        },
    },
    Geometry: &kml.Point{
        Coordinates: kml.Coord(-122.0, 37.0),
    },
}
```

## Error Handling

The library provides detailed error types:

```go
// ParseError for XML parsing issues
type ParseError struct {
    Line    int
    Column  int
    Message string
    Cause   error
}

// ValidationError for invalid KML structure
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

## Package Structure

```
github.com/robert-malhotra/go-kml/
├── kml.go           # Root KML type and parsing
├── document.go      # Document, Folder, Feature interface
├── placemark.go     # Placemark, ExtendedData
├── geometry.go      # Point, LineString, Polygon, etc.
├── style.go         # Style types
├── coordinate.go    # Coordinate parsing/formatting
├── color.go         # Color utilities
├── builder.go       # Fluent builder API
├── walk.go          # Traversal utilities
├── errors.go        # Error types
└── testdata/        # Real-world KML test samples
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

The library includes:
- Unit tests for all core types
- Round-trip tests (parse → write → parse)
- Real-world validation tests using [Google's KML samples](https://github.com/googlearchive/kml-samples)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## References

- [OGC KML 2.2 Specification](https://www.ogc.org/standards/kml)
- [Google KML Documentation](https://developers.google.com/kml/documentation)
- [KML Reference](https://developers.google.com/kml/documentation/kmlreference)
