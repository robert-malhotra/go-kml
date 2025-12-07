package kml

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

// TestParseSimpleKML tests parsing a simple KML document with one placemark
func TestParseSimpleKML(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Test</name>
    <Placemark>
      <name>Test Point</name>
      <Point>
        <coordinates>-122.0,37.0</coordinates>
      </Point>
    </Placemark>
  </Document>
</kml>`

	k, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("Failed to parse KML: %v", err)
	}

	// Verify the document
	doc, ok := k.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document, got %T", k.Feature)
	}

	if doc.Name != "Test" {
		t.Errorf("Expected document name 'Test', got '%s'", doc.Name)
	}

	// Verify there's one placemark
	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature, got %d", len(doc.Features))
	}

	placemark, ok := doc.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("Expected Placemark, got %T", doc.Features[0])
	}

	if placemark.Name != "Test Point" {
		t.Errorf("Expected placemark name 'Test Point', got '%s'", placemark.Name)
	}

	// Verify the point geometry
	point, ok := placemark.Geometry.(*Point)
	if !ok {
		t.Fatalf("Expected Point geometry, got %T", placemark.Geometry)
	}

	if point.Coordinates.Lon != -122.0 {
		t.Errorf("Expected longitude -122.0, got %f", point.Coordinates.Lon)
	}

	if point.Coordinates.Lat != 37.0 {
		t.Errorf("Expected latitude 37.0, got %f", point.Coordinates.Lat)
	}
}

// TestParseWithFolder tests parsing KML with nested folders
func TestParseWithFolder(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Test Document</name>
    <Folder>
      <name>Folder 1</name>
      <Placemark>
        <name>Point in Folder</name>
        <Point>
          <coordinates>-122.0,37.0</coordinates>
        </Point>
      </Placemark>
      <Folder>
        <name>Nested Folder</name>
        <Placemark>
          <name>Nested Point</name>
          <Point>
            <coordinates>-121.0,36.0</coordinates>
          </Point>
        </Placemark>
      </Folder>
    </Folder>
  </Document>
</kml>`

	k, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("Failed to parse KML: %v", err)
	}

	doc, ok := k.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document, got %T", k.Feature)
	}

	if doc.Name != "Test Document" {
		t.Errorf("Expected document name 'Test Document', got '%s'", doc.Name)
	}

	// Verify the folder
	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature in document, got %d", len(doc.Features))
	}

	folder, ok := doc.Features[0].(*Folder)
	if !ok {
		t.Fatalf("Expected Folder, got %T", doc.Features[0])
	}

	if folder.Name != "Folder 1" {
		t.Errorf("Expected folder name 'Folder 1', got '%s'", folder.Name)
	}

	// Verify folder has 2 features (1 placemark + 1 nested folder)
	if len(folder.Features) != 2 {
		t.Fatalf("Expected 2 features in folder, got %d", len(folder.Features))
	}

	// Check first feature is a placemark
	placemark1, ok := folder.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("Expected Placemark as first feature, got %T", folder.Features[0])
	}

	if placemark1.Name != "Point in Folder" {
		t.Errorf("Expected placemark name 'Point in Folder', got '%s'", placemark1.Name)
	}

	// Check second feature is a nested folder
	nestedFolder, ok := folder.Features[1].(*Folder)
	if !ok {
		t.Fatalf("Expected Folder as second feature, got %T", folder.Features[1])
	}

	if nestedFolder.Name != "Nested Folder" {
		t.Errorf("Expected nested folder name 'Nested Folder', got '%s'", nestedFolder.Name)
	}

	// Verify nested placemark
	if len(nestedFolder.Features) != 1 {
		t.Fatalf("Expected 1 feature in nested folder, got %d", len(nestedFolder.Features))
	}

	placemark2, ok := nestedFolder.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("Expected Placemark in nested folder, got %T", nestedFolder.Features[0])
	}

	if placemark2.Name != "Nested Point" {
		t.Errorf("Expected placemark name 'Nested Point', got '%s'", placemark2.Name)
	}
}

// TestParseWithStyles tests parsing KML with Style definitions and styleUrl references
func TestParseWithStyles(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Styled Document</name>
    <Style id="redIcon">
      <IconStyle>
        <color>ff0000ff</color>
        <scale>1.5</scale>
      </IconStyle>
      <LineStyle>
        <color>ff0000ff</color>
        <width>2.0</width>
      </LineStyle>
    </Style>
    <Placemark>
      <name>Red Point</name>
      <styleUrl>#redIcon</styleUrl>
      <Point>
        <coordinates>-122.0,37.0</coordinates>
      </Point>
    </Placemark>
  </Document>
</kml>`

	k, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("Failed to parse KML: %v", err)
	}

	doc, ok := k.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document, got %T", k.Feature)
	}

	// Verify styles
	if len(doc.Styles) != 1 {
		t.Fatalf("Expected 1 style, got %d", len(doc.Styles))
	}

	style := doc.Styles[0]
	if style.ID != "redIcon" {
		t.Errorf("Expected style ID 'redIcon', got '%s'", style.ID)
	}

	// Verify icon style
	if style.IconStyle == nil {
		t.Fatal("Expected IconStyle to be present")
	}

	if style.IconStyle.Scale != 1.5 {
		t.Errorf("Expected scale 1.5, got %f", style.IconStyle.Scale)
	}

	expectedColor := "ff0000ff"
	if style.IconStyle.Color.Hex() != expectedColor {
		t.Errorf("Expected color %s, got %s", expectedColor, style.IconStyle.Color.Hex())
	}

	// Verify line style
	if style.LineStyle == nil {
		t.Fatal("Expected LineStyle to be present")
	}

	if style.LineStyle.Width != 2.0 {
		t.Errorf("Expected width 2.0, got %f", style.LineStyle.Width)
	}

	// Verify placemark styleUrl reference
	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature, got %d", len(doc.Features))
	}

	placemark, ok := doc.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("Expected Placemark, got %T", doc.Features[0])
	}

	if placemark.StyleURL != "#redIcon" {
		t.Errorf("Expected styleUrl '#redIcon', got '%s'", placemark.StyleURL)
	}
}

// TestParseBytes tests the ParseBytes function
func TestParseBytes(t *testing.T) {
	kmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Byte Test</name>
    <Placemark>
      <name>Point</name>
      <Point>
        <coordinates>0.0,0.0</coordinates>
      </Point>
    </Placemark>
  </Document>
</kml>`)

	k, err := ParseBytes(kmlData)
	if err != nil {
		t.Fatalf("ParseBytes failed: %v", err)
	}

	doc, ok := k.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document, got %T", k.Feature)
	}

	if doc.Name != "Byte Test" {
		t.Errorf("Expected document name 'Byte Test', got '%s'", doc.Name)
	}

	// Test with invalid data
	invalidData := []byte("not valid xml")
	_, err = ParseBytes(invalidData)
	if err == nil {
		t.Error("Expected error when parsing invalid data, got nil")
	}
}

// TestWriteKML tests the Write method produces valid XML
func TestWriteKML(t *testing.T) {
	// Create KML programmatically
	k := createTestKML()

	// Write to buffer
	var buf bytes.Buffer
	err := k.Write(&buf)
	if err != nil {
		t.Fatalf("Failed to write KML: %v", err)
	}

	output := buf.String()

	// Verify XML declaration is present
	if !strings.HasPrefix(output, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("Expected XML declaration at start of output")
	}

	// Verify namespace is present
	if !strings.Contains(output, "xmlns=\"http://www.opengis.net/kml/2.2\"") {
		t.Error("Expected KML namespace in output")
	}

	// Verify it can be parsed back
	k2, err := ParseBytes([]byte(output))
	if err != nil {
		t.Fatalf("Failed to parse written KML: %v", err)
	}

	// Verify the parsed document matches original
	doc, ok := k2.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document after round-trip, got %T", k2.Feature)
	}

	if doc.Name != "Test Document" {
		t.Errorf("Expected document name 'Test Document' after round-trip, got '%s'", doc.Name)
	}

	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature after round-trip, got %d", len(doc.Features))
	}

	placemark, ok := doc.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("Expected Placemark after round-trip, got %T", doc.Features[0])
	}

	if placemark.Name != "Test Placemark" {
		t.Errorf("Expected placemark name 'Test Placemark' after round-trip, got '%s'", placemark.Name)
	}

	point, ok := placemark.Geometry.(*Point)
	if !ok {
		t.Fatalf("Expected Point geometry after round-trip, got %T", placemark.Geometry)
	}

	if point.Coordinates.Lon != -122.5 || point.Coordinates.Lat != 37.5 {
		t.Errorf("Expected coordinates (-122.5, 37.5) after round-trip, got (%f, %f)",
			point.Coordinates.Lon, point.Coordinates.Lat)
	}
}

// TestWriteIndent tests the WriteIndent method produces formatted XML
func TestWriteIndent(t *testing.T) {
	k := createTestKML()

	var buf bytes.Buffer
	err := k.WriteIndent(&buf, "", "  ")
	if err != nil {
		t.Fatalf("Failed to write indented KML: %v", err)
	}

	output := buf.String()

	// Verify XML declaration is present
	if !strings.HasPrefix(output, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("Expected XML declaration at start of output")
	}

	// Verify indentation is present (look for newlines and spaces)
	if !strings.Contains(output, "\n") {
		t.Error("Expected newlines in indented output")
	}

	// Count indentation levels - should have nested elements with increasing indentation
	lines := strings.Split(output, "\n")
	hasIndentation := false
	for _, line := range lines {
		if strings.HasPrefix(line, "  ") {
			hasIndentation = true
			break
		}
	}

	if !hasIndentation {
		t.Error("Expected indentation in output")
	}

	// Verify it can still be parsed
	k2, err := ParseBytes([]byte(output))
	if err != nil {
		t.Fatalf("Failed to parse indented KML: %v", err)
	}

	doc, ok := k2.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document after parsing indented KML, got %T", k2.Feature)
	}

	if doc.Name != "Test Document" {
		t.Errorf("Expected document name 'Test Document', got '%s'", doc.Name)
	}
}

// TestNewKML verifies NewKML() creates KML with correct defaults
func TestNewKML(t *testing.T) {
	k := NewKML()

	if k == nil {
		t.Fatal("NewKML() returned nil")
	}

	// Verify default namespace
	if k.Xmlns != DefaultNamespace {
		t.Errorf("Expected xmlns '%s', got '%s'", DefaultNamespace, k.Xmlns)
	}

	// Verify feature is nil (empty document)
	if k.Feature != nil {
		t.Errorf("Expected Feature to be nil, got %T", k.Feature)
	}

	// Verify we can set a feature and write it
	doc := &Document{
		Name: "New Document",
		Features: []Feature{
			&Placemark{
				Name: "Point",
				Geometry: &Point{
					Coordinates: Coord(-120.0, 38.0),
				},
			},
		},
	}
	k.Feature = doc

	var buf bytes.Buffer
	err := k.Write(&buf)
	if err != nil {
		t.Fatalf("Failed to write KML created with NewKML(): %v", err)
	}

	// Verify it contains expected content
	output := buf.String()
	if !strings.Contains(output, "New Document") {
		t.Error("Expected 'New Document' in output")
	}

	if !strings.Contains(output, DefaultNamespace) {
		t.Errorf("Expected default namespace '%s' in output", DefaultNamespace)
	}
}

// TestKMLBytes tests the Bytes() method
func TestKMLBytes(t *testing.T) {
	k := createTestKML()

	data, err := k.Bytes()
	if err != nil {
		t.Fatalf("Failed to get bytes: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty byte slice")
	}

	// Verify XML declaration is present
	output := string(data)
	if !strings.HasPrefix(output, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("Expected XML declaration at start of bytes")
	}

	// Verify it can be parsed
	k2, err := ParseBytes(data)
	if err != nil {
		t.Fatalf("Failed to parse bytes: %v", err)
	}

	doc, ok := k2.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document from bytes, got %T", k2.Feature)
	}

	if doc.Name != "Test Document" {
		t.Errorf("Expected document name 'Test Document' from bytes, got '%s'", doc.Name)
	}
}

// TestWriteKMLWithMultipleGeometries tests writing KML with various geometry types
func TestWriteKMLWithMultipleGeometries(t *testing.T) {
	k := NewKML()
	doc := &Document{
		Name: "Multi-Geometry Test",
		Features: []Feature{
			&Placemark{
				Name: "Line",
				Geometry: &LineString{
					Coordinates: []Coordinate{
						Coord(-122.0, 37.0),
						Coord(-122.1, 37.1),
						Coord(-122.2, 37.2),
					},
				},
			},
			&Placemark{
				Name: "Polygon",
				Geometry: &Polygon{
					OuterBoundary: LinearRing{
						Coordinates: []Coordinate{
							Coord(-122.0, 37.0),
							Coord(-122.1, 37.0),
							Coord(-122.1, 37.1),
							Coord(-122.0, 37.1),
							Coord(-122.0, 37.0),
						},
					},
				},
			},
		},
	}
	k.Feature = doc

	var buf bytes.Buffer
	err := k.Write(&buf)
	if err != nil {
		t.Fatalf("Failed to write multi-geometry KML: %v", err)
	}

	output := buf.String()

	// Verify LineString is present
	if !strings.Contains(output, "<LineString>") {
		t.Error("Expected LineString element in output")
	}

	// Verify Polygon is present
	if !strings.Contains(output, "<Polygon>") {
		t.Error("Expected Polygon element in output")
	}

	// Parse it back and verify
	k2, err := ParseBytes([]byte(output))
	if err != nil {
		t.Fatalf("Failed to parse multi-geometry KML: %v", err)
	}

	doc2, ok := k2.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document, got %T", k2.Feature)
	}

	if len(doc2.Features) != 2 {
		t.Fatalf("Expected 2 features, got %d", len(doc2.Features))
	}

	// Verify LineString
	pm1, ok := doc2.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("Expected first feature to be Placemark, got %T", doc2.Features[0])
	}

	lineString, ok := pm1.Geometry.(*LineString)
	if !ok {
		t.Fatalf("Expected LineString geometry, got %T", pm1.Geometry)
	}

	if len(lineString.Coordinates) != 3 {
		t.Errorf("Expected 3 coordinates in LineString, got %d", len(lineString.Coordinates))
	}

	// Verify Polygon
	pm2, ok := doc2.Features[1].(*Placemark)
	if !ok {
		t.Fatalf("Expected second feature to be Placemark, got %T", doc2.Features[1])
	}

	polygon, ok := pm2.Geometry.(*Polygon)
	if !ok {
		t.Fatalf("Expected Polygon geometry, got %T", pm2.Geometry)
	}

	if len(polygon.OuterBoundary.Coordinates) != 5 {
		t.Errorf("Expected 5 coordinates in Polygon outer boundary, got %d",
			len(polygon.OuterBoundary.Coordinates))
	}
}

// TestParseEmptyDocument tests that parsing KML without features returns an error
func TestParseEmptyDocument(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
</kml>`

	_, err := ParseBytes([]byte(kmlData))
	if err == nil {
		t.Error("Expected error when parsing empty KML, got nil")
	}

	if err != ErrEmptyDocument {
		t.Errorf("Expected ErrEmptyDocument, got %v", err)
	}
}

// TestRoundTripWithXMLDecoder tests that KML can be parsed using xml.Decoder
func TestRoundTripWithXMLDecoder(t *testing.T) {
	k := createTestKML()

	// Marshal to bytes
	data, err := k.Bytes()
	if err != nil {
		t.Fatalf("Failed to marshal KML: %v", err)
	}

	// Unmarshal using xml.Decoder
	var k2 KML
	decoder := xml.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&k2)
	if err != nil {
		t.Fatalf("Failed to decode KML with xml.Decoder: %v", err)
	}

	// Verify namespace
	if k2.Xmlns != DefaultNamespace {
		t.Errorf("Expected xmlns '%s', got '%s'", DefaultNamespace, k2.Xmlns)
	}

	// Verify document
	doc, ok := k2.Feature.(*Document)
	if !ok {
		t.Fatalf("Expected Document, got %T", k2.Feature)
	}

	if doc.Name != "Test Document" {
		t.Errorf("Expected document name 'Test Document', got '%s'", doc.Name)
	}
}

// createTestKML creates a simple test KML document for testing
func createTestKML() *KML {
	k := NewKML()
	doc := &Document{
		Name: "Test Document",
		Features: []Feature{
			&Placemark{
				Name:        "Test Placemark",
				Description: "A test point",
				Geometry: &Point{
					Coordinates: Coord(-122.5, 37.5),
				},
			},
		},
	}
	k.Feature = doc
	return k
}
