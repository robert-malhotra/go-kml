package kml

import (
	"bytes"
	"reflect"
	"testing"
)

// TestRoundTripSimple tests that a simple document with one placemark
// can be written and parsed back correctly.
func TestRoundTripSimple(t *testing.T) {
	// Create KML with builder
	k := NewKMLBuilder().
		Document("Test Document").
		Placemark("Test Point").
		Point(-122.0, 37.0).
		Done().(*DocumentBuilder).
		Build()

	// Write to bytes
	var buf bytes.Buffer
	if err := k.Write(&buf); err != nil {
		t.Fatalf("Failed to write KML: %v", err)
	}

	// Parse bytes back
	k2, err := ParseBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to parse KML: %v", err)
	}

	// Verify structure matches
	doc1, ok := k.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Document feature")
	}

	doc2, ok := k2.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Document feature in parsed KML")
	}

	if doc1.Name != doc2.Name {
		t.Errorf("Document name mismatch: got %q, want %q", doc2.Name, doc1.Name)
	}

	if len(doc1.Features) != len(doc2.Features) {
		t.Fatalf("Feature count mismatch: got %d, want %d", len(doc2.Features), len(doc1.Features))
	}

	pm1, ok := doc1.Features[0].(*Placemark)
	if !ok {
		t.Fatal("Expected Placemark feature")
	}

	pm2, ok := doc2.Features[0].(*Placemark)
	if !ok {
		t.Fatal("Expected Placemark feature in parsed KML")
	}

	if pm1.Name != pm2.Name {
		t.Errorf("Placemark name mismatch: got %q, want %q", pm2.Name, pm1.Name)
	}

	// Verify Point geometry
	pt1, ok := pm1.Geometry.(*Point)
	if !ok {
		t.Fatal("Expected Point geometry")
	}

	pt2, ok := pm2.Geometry.(*Point)
	if !ok {
		t.Fatal("Expected Point geometry in parsed KML")
	}

	if !coordEqual(pt1.Coordinates, pt2.Coordinates) {
		t.Errorf("Coordinates mismatch: got %v, want %v", pt2.Coordinates, pt1.Coordinates)
	}
}

// TestRoundTripWithStyles tests that styles are preserved through round-trip.
func TestRoundTripWithStyles(t *testing.T) {
	// Create KML with LineStyle, PolyStyle, IconStyle
	k := NewKMLBuilder().
		Document("Styled Document").
		Style("lineStyle").
		LineStyle().
		Color(RGBA(255, 0, 0, 255)).
		Width(2.5).
		Done().
		Done().
		Style("polyStyle").
		PolyStyle().
		Color(RGBA(0, 255, 0, 128)).
		Fill(true).
		Outline(false).
		Done().
		Done().
		Style("iconStyle").
		IconStyle().
		Color(RGBA(0, 0, 255, 255)).
		Scale(1.5).
		Icon("http://example.com/icon.png").
		Done().
		Done().
		Placemark("Styled Line").
		StyleURL("#lineStyle").
		LineString(Coord(-122.0, 37.0), Coord(-121.0, 38.0)).
		Done().(*DocumentBuilder).
		Build()

	// Round-trip
	var buf bytes.Buffer
	if err := k.Write(&buf); err != nil {
		t.Fatalf("Failed to write KML: %v", err)
	}

	k2, err := ParseBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to parse KML: %v", err)
	}

	// Verify styles preserved
	doc1, _ := k.Feature.(*Document)
	doc2, _ := k2.Feature.(*Document)

	if len(doc1.Styles) != len(doc2.Styles) {
		t.Fatalf("Style count mismatch: got %d, want %d", len(doc2.Styles), len(doc1.Styles))
	}

	// Verify LineStyle
	ls1 := doc1.Styles[0].LineStyle
	ls2 := doc2.Styles[0].LineStyle
	if ls1 == nil || ls2 == nil {
		t.Fatal("LineStyle is nil")
	}
	if !colorEqual(ls1.Color, ls2.Color) {
		t.Errorf("LineStyle color mismatch: got %v, want %v", ls2.Color, ls1.Color)
	}
	if ls1.Width != ls2.Width {
		t.Errorf("LineStyle width mismatch: got %v, want %v", ls2.Width, ls1.Width)
	}

	// Verify PolyStyle
	ps1 := doc1.Styles[1].PolyStyle
	ps2 := doc2.Styles[1].PolyStyle
	if ps1 == nil || ps2 == nil {
		t.Fatal("PolyStyle is nil")
	}
	if !colorEqual(ps1.Color, ps2.Color) {
		t.Errorf("PolyStyle color mismatch: got %v, want %v", ps2.Color, ps1.Color)
	}
	if *ps1.Fill != *ps2.Fill {
		t.Errorf("PolyStyle fill mismatch: got %v, want %v", *ps2.Fill, *ps1.Fill)
	}
	if *ps1.Outline != *ps2.Outline {
		t.Errorf("PolyStyle outline mismatch: got %v, want %v", *ps2.Outline, *ps1.Outline)
	}

	// Verify IconStyle
	is1 := doc1.Styles[2].IconStyle
	is2 := doc2.Styles[2].IconStyle
	if is1 == nil || is2 == nil {
		t.Fatal("IconStyle is nil")
	}
	if !colorEqual(is1.Color, is2.Color) {
		t.Errorf("IconStyle color mismatch: got %v, want %v", is2.Color, is1.Color)
	}
	if is1.Scale != is2.Scale {
		t.Errorf("IconStyle scale mismatch: got %v, want %v", is2.Scale, is1.Scale)
	}
	if is1.Icon.Href != is2.Icon.Href {
		t.Errorf("IconStyle icon href mismatch: got %q, want %q", is2.Icon.Href, is1.Icon.Href)
	}

	// Verify placemark styleUrl
	pm1, _ := doc1.Features[0].(*Placemark)
	pm2, _ := doc2.Features[0].(*Placemark)
	if pm1.StyleURL != pm2.StyleURL {
		t.Errorf("StyleURL mismatch: got %q, want %q", pm2.StyleURL, pm1.StyleURL)
	}
}

// TestRoundTripWithFolders tests that nested folder hierarchy is preserved.
func TestRoundTripWithFolders(t *testing.T) {
	// Create KML with multiple nested folders
	k := NewKMLBuilder().
		Document("Document with Folders").
		Folder("Folder 1").
		Placemark("Point in Folder 1").
		Point(-120.0, 35.0).
		Done().(*FolderBuilder).
		Folder("Nested Folder 1.1").
		Placemark("Point in Nested 1.1").
		Point(-121.0, 36.0).
		Done().(*FolderBuilder).
		Done().(*FolderBuilder).
		Done().(*DocumentBuilder).
		Folder("Folder 2").
		Placemark("Point in Folder 2").
		Point(-122.0, 37.0).
		Done().(*FolderBuilder).
		Done().(*DocumentBuilder).
		Build()

	// Round-trip
	var buf bytes.Buffer
	if err := k.Write(&buf); err != nil {
		t.Fatalf("Failed to write KML: %v", err)
	}

	k2, err := ParseBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to parse KML: %v", err)
	}

	// Verify folder hierarchy
	doc1, _ := k.Feature.(*Document)
	doc2, _ := k2.Feature.(*Document)

	if len(doc1.Features) != len(doc2.Features) {
		t.Fatalf("Top-level feature count mismatch: got %d, want %d", len(doc2.Features), len(doc1.Features))
	}

	// Check Folder 1
	f1_1, ok := doc1.Features[0].(*Folder)
	if !ok {
		t.Fatal("Expected Folder feature")
	}
	f1_2, ok := doc2.Features[0].(*Folder)
	if !ok {
		t.Fatal("Expected Folder feature in parsed KML")
	}

	if f1_1.Name != f1_2.Name {
		t.Errorf("Folder 1 name mismatch: got %q, want %q", f1_2.Name, f1_1.Name)
	}

	if len(f1_1.Features) != len(f1_2.Features) {
		t.Fatalf("Folder 1 feature count mismatch: got %d, want %d", len(f1_2.Features), len(f1_1.Features))
	}

	// Check nested folder
	nf1_1, ok := f1_1.Features[1].(*Folder)
	if !ok {
		t.Fatal("Expected nested Folder")
	}
	nf1_2, ok := f1_2.Features[1].(*Folder)
	if !ok {
		t.Fatal("Expected nested Folder in parsed KML")
	}

	if nf1_1.Name != nf1_2.Name {
		t.Errorf("Nested folder name mismatch: got %q, want %q", nf1_2.Name, nf1_1.Name)
	}

	// Check Folder 2
	f2_1, ok := doc1.Features[1].(*Folder)
	if !ok {
		t.Fatal("Expected Folder 2 feature")
	}
	f2_2, ok := doc2.Features[1].(*Folder)
	if !ok {
		t.Fatal("Expected Folder 2 feature in parsed KML")
	}

	if f2_1.Name != f2_2.Name {
		t.Errorf("Folder 2 name mismatch: got %q, want %q", f2_2.Name, f2_1.Name)
	}
}

// TestRoundTripAllGeometries tests all geometry types (Point, LineString, Polygon, MultiGeometry).
func TestRoundTripAllGeometries(t *testing.T) {
	// Create polygons with inner boundaries
	outerRing := []Coordinate{
		Coord(-122.0, 37.0),
		Coord(-121.0, 37.0),
		Coord(-121.0, 38.0),
		Coord(-122.0, 38.0),
		Coord(-122.0, 37.0),
	}
	innerRing := []Coordinate{
		Coord(-121.8, 37.2),
		Coord(-121.2, 37.2),
		Coord(-121.2, 37.8),
		Coord(-121.8, 37.8),
		Coord(-121.8, 37.2),
	}

	k := NewKMLBuilder().
		Document("All Geometries").
		Placemark("Point Test").
		Point(-122.0, 37.0, 100.0).
		Done().(*DocumentBuilder).
		Placemark("LineString Test").
		LineString(
			Coord(-122.0, 37.0),
			Coord(-121.0, 38.0),
			Coord(-120.0, 39.0, 500.0),
		).
		Done().(*DocumentBuilder).
		Placemark("Polygon Test").
		Polygon(outerRing, innerRing).
		Done().(*DocumentBuilder).
		Build()

	// Add MultiGeometry manually (builder doesn't support it directly)
	doc := k.Feature.(*Document)
	multiGeom := &Placemark{
		Name: "MultiGeometry Test",
		Geometry: &MultiGeometry{
			Geometries: []Geometry{
				&Point{Coordinates: Coord(-122.0, 37.0)},
				&LineString{
					Coordinates: []Coordinate{
						Coord(-121.0, 36.0),
						Coord(-120.0, 35.0),
					},
				},
			},
		},
	}
	doc.Features = append(doc.Features, multiGeom)

	// Round-trip
	var buf bytes.Buffer
	if err := k.Write(&buf); err != nil {
		t.Fatalf("Failed to write KML: %v", err)
	}

	k2, err := ParseBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to parse KML: %v", err)
	}

	doc1, _ := k.Feature.(*Document)
	doc2, _ := k2.Feature.(*Document)

	// Verify Point
	pm1, _ := doc1.Features[0].(*Placemark)
	pm2, _ := doc2.Features[0].(*Placemark)
	pt1, _ := pm1.Geometry.(*Point)
	pt2, _ := pm2.Geometry.(*Point)
	if !coordEqual(pt1.Coordinates, pt2.Coordinates) {
		t.Errorf("Point coordinates mismatch: got %v, want %v", pt2.Coordinates, pt1.Coordinates)
	}

	// Verify LineString
	pm1 = doc1.Features[1].(*Placemark)
	pm2 = doc2.Features[1].(*Placemark)
	ls1, _ := pm1.Geometry.(*LineString)
	ls2, _ := pm2.Geometry.(*LineString)
	if !coordSliceEqual(ls1.Coordinates, ls2.Coordinates) {
		t.Errorf("LineString coordinates mismatch: got %v, want %v", ls2.Coordinates, ls1.Coordinates)
	}

	// Verify Polygon
	pm1 = doc1.Features[2].(*Placemark)
	pm2 = doc2.Features[2].(*Placemark)
	pg1, _ := pm1.Geometry.(*Polygon)
	pg2, _ := pm2.Geometry.(*Polygon)
	if !coordSliceEqual(pg1.OuterBoundary.Coordinates, pg2.OuterBoundary.Coordinates) {
		t.Errorf("Polygon outer boundary mismatch")
	}
	if len(pg1.InnerBoundaries) != len(pg2.InnerBoundaries) {
		t.Fatalf("Polygon inner boundary count mismatch: got %d, want %d",
			len(pg2.InnerBoundaries), len(pg1.InnerBoundaries))
	}
	if !coordSliceEqual(pg1.InnerBoundaries[0].Coordinates, pg2.InnerBoundaries[0].Coordinates) {
		t.Errorf("Polygon inner boundary mismatch")
	}

	// Verify MultiGeometry
	pm1 = doc1.Features[3].(*Placemark)
	pm2 = doc2.Features[3].(*Placemark)
	mg1, _ := pm1.Geometry.(*MultiGeometry)
	mg2, _ := pm2.Geometry.(*MultiGeometry)
	if len(mg1.Geometries) != len(mg2.Geometries) {
		t.Fatalf("MultiGeometry count mismatch: got %d, want %d", len(mg2.Geometries), len(mg1.Geometries))
	}

	// Verify MultiGeometry Point
	mgPt1, _ := mg1.Geometries[0].(*Point)
	mgPt2, _ := mg2.Geometries[0].(*Point)
	if !coordEqual(mgPt1.Coordinates, mgPt2.Coordinates) {
		t.Errorf("MultiGeometry point coordinates mismatch")
	}

	// Verify MultiGeometry LineString
	mgLs1, _ := mg1.Geometries[1].(*LineString)
	mgLs2, _ := mg2.Geometries[1].(*LineString)
	if !coordSliceEqual(mgLs1.Coordinates, mgLs2.Coordinates) {
		t.Errorf("MultiGeometry linestring coordinates mismatch")
	}
}

// TestRoundTripColors verifies that colors round-trip correctly.
func TestRoundTripColors(t *testing.T) {
	testColors := []struct {
		name  string
		color Color
	}{
		{"Red", RGBA(255, 0, 0, 255)},
		{"Green", RGBA(0, 255, 0, 255)},
		{"Blue", RGBA(0, 0, 255, 255)},
		{"Semi-transparent", RGBA(128, 128, 128, 128)},
		{"White", White},
		{"Black", Black},
	}

	for _, tc := range testColors {
		t.Run(tc.name, func(t *testing.T) {
			k := NewKMLBuilder().
				Document("Color Test").
				Style("testStyle").
				LineStyle().
				Color(tc.color).
				Width(2.0).
				Done().
				Done().
				Build()

			// Round-trip
			var buf bytes.Buffer
			if err := k.Write(&buf); err != nil {
				t.Fatalf("Failed to write KML: %v", err)
			}

			k2, err := ParseBytes(buf.Bytes())
			if err != nil {
				t.Fatalf("Failed to parse KML: %v", err)
			}

			// Verify color
			doc1, _ := k.Feature.(*Document)
			doc2, _ := k2.Feature.(*Document)

			color1 := doc1.Styles[0].LineStyle.Color
			color2 := doc2.Styles[0].LineStyle.Color

			if !colorEqual(color1, color2) {
				t.Errorf("Color mismatch: got %v (hex: %s), want %v (hex: %s)",
					color2, color2.Hex(), color1, color1.Hex())
			}
		})
	}
}

// TestRoundTripFromString parses real KML string, writes it, and parses again.
func TestRoundTripFromString(t *testing.T) {
	kmlString := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Round Trip Test</name>
    <Style id="testStyle">
      <LineStyle>
        <color>ff0000ff</color>
        <width>2</width>
      </LineStyle>
    </Style>
    <Folder>
      <name>Test Folder</name>
      <Placemark>
        <name>Test Line</name>
        <styleUrl>#testStyle</styleUrl>
        <LineString>
          <coordinates>-122.0,37.0 -121.0,38.0</coordinates>
        </LineString>
      </Placemark>
    </Folder>
  </Document>
</kml>
`

	// Parse original
	k1, err := ParseBytes([]byte(kmlString))
	if err != nil {
		t.Fatalf("Failed to parse original KML: %v", err)
	}

	// Write to bytes
	var buf bytes.Buffer
	if err := k1.Write(&buf); err != nil {
		t.Fatalf("Failed to write KML: %v", err)
	}

	// Parse again
	k2, err := ParseBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to parse written KML: %v", err)
	}

	// Verify structure
	doc1, _ := k1.Feature.(*Document)
	doc2, _ := k2.Feature.(*Document)

	if doc1.Name != doc2.Name {
		t.Errorf("Document name mismatch: got %q, want %q", doc2.Name, doc1.Name)
	}

	// Verify style
	if len(doc1.Styles) != len(doc2.Styles) {
		t.Fatalf("Style count mismatch: got %d, want %d", len(doc2.Styles), len(doc1.Styles))
	}

	style1 := doc1.Styles[0]
	style2 := doc2.Styles[0]

	if style1.ID != style2.ID {
		t.Errorf("Style ID mismatch: got %q, want %q", style2.ID, style1.ID)
	}

	if !colorEqual(style1.LineStyle.Color, style2.LineStyle.Color) {
		t.Errorf("LineStyle color mismatch: got %s, want %s",
			style2.LineStyle.Color.Hex(), style1.LineStyle.Color.Hex())
	}

	if style1.LineStyle.Width != style2.LineStyle.Width {
		t.Errorf("LineStyle width mismatch: got %v, want %v",
			style2.LineStyle.Width, style1.LineStyle.Width)
	}

	// Verify folder
	folder1, _ := doc1.Features[0].(*Folder)
	folder2, _ := doc2.Features[0].(*Folder)

	if folder1.Name != folder2.Name {
		t.Errorf("Folder name mismatch: got %q, want %q", folder2.Name, folder1.Name)
	}

	// Verify placemark
	pm1, _ := folder1.Features[0].(*Placemark)
	pm2, _ := folder2.Features[0].(*Placemark)

	if pm1.Name != pm2.Name {
		t.Errorf("Placemark name mismatch: got %q, want %q", pm2.Name, pm1.Name)
	}

	if pm1.StyleURL != pm2.StyleURL {
		t.Errorf("StyleURL mismatch: got %q, want %q", pm2.StyleURL, pm1.StyleURL)
	}

	// Verify geometry
	ls1, _ := pm1.Geometry.(*LineString)
	ls2, _ := pm2.Geometry.(*LineString)

	if !coordSliceEqual(ls1.Coordinates, ls2.Coordinates) {
		t.Errorf("LineString coordinates mismatch: got %v, want %v",
			ls2.Coordinates, ls1.Coordinates)
	}
}

// Helper functions

// coordEqual checks if two coordinates are equal within a small tolerance.
func coordEqual(a, b Coordinate) bool {
	const epsilon = 1e-9
	return floatEqual(a.Lon, b.Lon, epsilon) &&
		floatEqual(a.Lat, b.Lat, epsilon) &&
		floatEqual(a.Alt, b.Alt, epsilon)
}

// coordSliceEqual checks if two coordinate slices are equal.
func coordSliceEqual(a, b []Coordinate) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !coordEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// colorEqual checks if two colors are equal.
func colorEqual(a, b Color) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B && a.A == b.A
}

// floatEqual checks if two floats are equal within a tolerance.
func floatEqual(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

// compareKML compares two KML structures for equality.
// Returns true if they are structurally equivalent.
// This is a utility function available for more complex comparison tests.
var _ = compareKML // silence unused warning

func compareKML(a, b *KML) bool {
	if a.Xmlns != b.Xmlns {
		return false
	}

	// Compare features using reflection for deep equality
	return reflect.DeepEqual(a.Feature, b.Feature)
}
