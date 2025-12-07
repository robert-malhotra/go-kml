package kml

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRealWorldKMLSamples tests parsing real KML files from googlearchive/kml-samples
func TestRealWorldKMLSamples(t *testing.T) {
	testdataDir := "testdata"

	// Check if testdata directory exists
	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found, skipping real-world tests")
	}

	// List of KML files to test
	kmlFiles := []string{
		"KML_Samples.kml",
		"placemark.kml",
		"basic.kml",
		"linestring-tessellate.kml",
		"linestring-styled.kml",
		"polygon-inner.kml",
		"multi-linestrings.kml",
		"polygon-point.kml",
		"balloon-styles.kml",
		"extendeddata.kml",
		"doc-with-id.kml",
		"catalina-points.kml",
	}

	for _, filename := range kmlFiles {
		t.Run(filename, func(t *testing.T) {
			path := filepath.Join(testdataDir, filename)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Skipf("file %s not found", filename)
				return
			}

			// Test parsing
			doc, err := ParseFile(path)
			if err != nil {
				t.Fatalf("failed to parse %s: %v", filename, err)
			}

			if doc == nil {
				t.Fatalf("parsed document is nil for %s", filename)
			}

			// Test round-trip
			bytes, err := doc.Bytes()
			if err != nil {
				t.Fatalf("failed to write %s: %v", filename, err)
			}

			// Parse the round-tripped content
			doc2, err := ParseBytes(bytes)
			if err != nil {
				t.Fatalf("failed to re-parse %s: %v", filename, err)
			}

			if doc2 == nil {
				t.Fatalf("re-parsed document is nil for %s", filename)
			}
		})
	}
}

// TestSimplePlacemark tests parsing the simple placemark.kml
func TestSimplePlacemark(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Placemark>
    <name>My office</name>
    <description>This is the location of my office.</description>
    <Point>
      <coordinates>-122.087461,37.422069</coordinates>
    </Point>
  </Placemark>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	pm, ok := doc.Feature.(*Placemark)
	if !ok {
		t.Fatalf("expected Placemark, got %T", doc.Feature)
	}

	if pm.Name != "My office" {
		t.Errorf("expected name 'My office', got '%s'", pm.Name)
	}

	if pm.Description != "This is the location of my office." {
		t.Errorf("unexpected description: %s", pm.Description)
	}

	pt, ok := pm.Geometry.(*Point)
	if !ok {
		t.Fatalf("expected Point geometry, got %T", pm.Geometry)
	}

	if !floatNear(pt.Coordinates.Lon, -122.087461, 0.000001) {
		t.Errorf("expected lon -122.087461, got %f", pt.Coordinates.Lon)
	}

	if !floatNear(pt.Coordinates.Lat, 37.422069, 0.000001) {
		t.Errorf("expected lat 37.422069, got %f", pt.Coordinates.Lat)
	}
}

// TestLineStringTessellate tests parsing LineString with tessellate flag
func TestLineStringTessellate(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Placemark>
    <name>Tessellated paths</name>
    <visibility>0</visibility>
    <LineString>
      <tessellate>1</tessellate>
      <coordinates>
        -112.0814237830345,36.10677870477137,0
        -112.0870267752693,36.0905099328766,0
      </coordinates>
    </LineString>
  </Placemark>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	pm, ok := doc.Feature.(*Placemark)
	if !ok {
		t.Fatalf("expected Placemark, got %T", doc.Feature)
	}

	if pm.Name != "Tessellated paths" {
		t.Errorf("expected name 'Tessellated paths', got '%s'", pm.Name)
	}

	// Check visibility is false (0)
	if pm.Visibility == nil || *pm.Visibility != false {
		t.Errorf("expected visibility false")
	}

	ls, ok := pm.Geometry.(*LineString)
	if !ok {
		t.Fatalf("expected LineString geometry, got %T", pm.Geometry)
	}

	if !ls.Tessellate {
		t.Error("expected tessellate to be true")
	}

	if len(ls.Coordinates) != 2 {
		t.Fatalf("expected 2 coordinates, got %d", len(ls.Coordinates))
	}

	// Check first coordinate
	if !floatNear(ls.Coordinates[0].Lon, -112.0814237830345, 0.0000001) {
		t.Errorf("unexpected first lon: %f", ls.Coordinates[0].Lon)
	}
}

// TestPolygonWithInnerBoundaries tests parsing polygon with holes
func TestPolygonWithInnerBoundaries(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Placemark>
    <name>innerBoundaries</name>
    <Polygon>
      <outerBoundaryIs>
        <LinearRing>
          <coordinates>
            -122.0,37.0
            -121.9,37.0
            -121.9,37.1
            -122.0,37.1
            -122.0,37.0
          </coordinates>
        </LinearRing>
      </outerBoundaryIs>
      <innerBoundaryIs>
        <LinearRing>
          <coordinates>
            -121.99,37.01
            -121.96,37.01
            -121.96,37.04
            -121.99,37.04
            -121.99,37.01
          </coordinates>
        </LinearRing>
      </innerBoundaryIs>
      <innerBoundaryIs>
        <LinearRing>
          <coordinates>
            -121.94,37.01
            -121.91,37.01
            -121.91,37.04
            -121.94,37.04
            -121.94,37.01
          </coordinates>
        </LinearRing>
      </innerBoundaryIs>
    </Polygon>
  </Placemark>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	pm, ok := doc.Feature.(*Placemark)
	if !ok {
		t.Fatalf("expected Placemark, got %T", doc.Feature)
	}

	poly, ok := pm.Geometry.(*Polygon)
	if !ok {
		t.Fatalf("expected Polygon geometry, got %T", pm.Geometry)
	}

	// Check outer boundary has 5 points (closed ring)
	if len(poly.OuterBoundary.Coordinates) != 5 {
		t.Errorf("expected 5 outer boundary coords, got %d", len(poly.OuterBoundary.Coordinates))
	}

	// Check we have 2 inner boundaries (holes)
	if len(poly.InnerBoundaries) != 2 {
		t.Errorf("expected 2 inner boundaries, got %d", len(poly.InnerBoundaries))
	}

	// Each inner boundary should have 5 points
	for i, inner := range poly.InnerBoundaries {
		if len(inner.Coordinates) != 5 {
			t.Errorf("inner boundary %d: expected 5 coords, got %d", i, len(inner.Coordinates))
		}
	}
}

// TestMultiGeometryLineStrings tests parsing MultiGeometry with multiple LineStrings
func TestMultiGeometryLineStrings(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>MultiGeometry - linestrings</name>
    <Style id="bangormarina">
      <LineStyle>
        <color>ff0000ff</color>
        <width>3</width>
      </LineStyle>
    </Style>
    <Folder>
      <Placemark>
        <name>MultiGeometry</name>
        <styleUrl>#bangormarina</styleUrl>
        <MultiGeometry>
          <LineString>
            <tessellate>1</tessellate>
            <coordinates>
              -5.670104418698614,54.66484515395317,0
              -5.669619432714838,54.66364764916655,0
            </coordinates>
          </LineString>
          <LineString>
            <tessellate>1</tessellate>
            <coordinates>
              -5.671178963580644,54.66531451472103,0
              -5.670934252664708,54.66473294883198,0
            </coordinates>
          </LineString>
        </MultiGeometry>
      </Placemark>
    </Folder>
  </Document>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	document, ok := doc.Feature.(*Document)
	if !ok {
		t.Fatalf("expected Document, got %T", doc.Feature)
	}

	if document.Name != "MultiGeometry - linestrings" {
		t.Errorf("unexpected document name: %s", document.Name)
	}

	// Check style
	if len(document.Styles) != 1 {
		t.Fatalf("expected 1 style, got %d", len(document.Styles))
	}

	style := document.Styles[0]
	if style.ID != "bangormarina" {
		t.Errorf("expected style ID 'bangormarina', got '%s'", style.ID)
	}

	if style.LineStyle == nil {
		t.Fatal("expected LineStyle")
	}

	if style.LineStyle.Width != 3 {
		t.Errorf("expected line width 3, got %f", style.LineStyle.Width)
	}

	// Navigate to folder and placemark
	if len(document.Features) != 1 {
		t.Fatalf("expected 1 feature (folder), got %d", len(document.Features))
	}

	folder, ok := document.Features[0].(*Folder)
	if !ok {
		t.Fatalf("expected Folder, got %T", document.Features[0])
	}

	if len(folder.Features) != 1 {
		t.Fatalf("expected 1 placemark in folder, got %d", len(folder.Features))
	}

	pm, ok := folder.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("expected Placemark, got %T", folder.Features[0])
	}

	if pm.StyleURL != "#bangormarina" {
		t.Errorf("expected styleUrl '#bangormarina', got '%s'", pm.StyleURL)
	}

	mg, ok := pm.Geometry.(*MultiGeometry)
	if !ok {
		t.Fatalf("expected MultiGeometry, got %T", pm.Geometry)
	}

	if len(mg.Geometries) != 2 {
		t.Errorf("expected 2 geometries in MultiGeometry, got %d", len(mg.Geometries))
	}

	// Verify both are LineStrings
	for i, g := range mg.Geometries {
		ls, ok := g.(*LineString)
		if !ok {
			t.Errorf("geometry %d: expected LineString, got %T", i, g)
			continue
		}
		if len(ls.Coordinates) != 2 {
			t.Errorf("geometry %d: expected 2 coords, got %d", i, len(ls.Coordinates))
		}
	}
}

// TestExtendedData tests parsing placemarks with ExtendedData
func TestExtendedData(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Data+BalloonStyle</name>
    <Placemark>
      <name>Club house</name>
      <ExtendedData>
        <Data name="holeNumber">
          <value>1</value>
        </Data>
        <Data name="holePar">
          <value>4</value>
        </Data>
        <Data name="holeYardage">
          <value>234</value>
        </Data>
      </ExtendedData>
    </Placemark>
  </Document>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	document, ok := doc.Feature.(*Document)
	if !ok {
		t.Fatalf("expected Document, got %T", doc.Feature)
	}

	if len(document.Features) != 1 {
		t.Fatalf("expected 1 placemark, got %d", len(document.Features))
	}

	pm, ok := document.Features[0].(*Placemark)
	if !ok {
		t.Fatalf("expected Placemark, got %T", document.Features[0])
	}

	if pm.ExtendedData == nil {
		t.Fatal("expected ExtendedData")
	}

	if len(pm.ExtendedData.Data) != 3 {
		t.Fatalf("expected 3 Data items, got %d", len(pm.ExtendedData.Data))
	}

	// Check data values
	dataMap := make(map[string]string)
	for _, d := range pm.ExtendedData.Data {
		dataMap[d.Name] = d.Value
	}

	if dataMap["holeNumber"] != "1" {
		t.Errorf("expected holeNumber=1, got %s", dataMap["holeNumber"])
	}
	if dataMap["holePar"] != "4" {
		t.Errorf("expected holePar=4, got %s", dataMap["holePar"])
	}
	if dataMap["holeYardage"] != "234" {
		t.Errorf("expected holeYardage=234, got %s", dataMap["holeYardage"])
	}
}

// TestKMLSamplesMainFile tests parsing the main KML_Samples.kml file
func TestKMLSamplesMainFile(t *testing.T) {
	path := "testdata/KML_Samples.kml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("KML_Samples.kml not found")
	}

	doc, err := ParseFile(path)
	if err != nil {
		t.Fatalf("failed to parse KML_Samples.kml: %v", err)
	}

	document, ok := doc.Feature.(*Document)
	if !ok {
		t.Fatalf("expected Document, got %T", doc.Feature)
	}

	if document.Name != "KML Samples" {
		t.Errorf("expected name 'KML Samples', got '%s'", document.Name)
	}

	// Should have multiple styles
	if len(document.Styles) < 5 {
		t.Errorf("expected at least 5 styles, got %d", len(document.Styles))
	}

	// Use Walk to count all features
	var featureCount int
	var placemarkCount int
	var folderCount int

	err = doc.Walk(func(f Feature) error {
		featureCount++
		switch f.(type) {
		case *Placemark:
			placemarkCount++
		case *Folder:
			folderCount++
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	t.Logf("KML_Samples.kml contains: %d features, %d placemarks, %d folders",
		featureCount, placemarkCount, folderCount)

	// Should have placemarks
	if placemarkCount == 0 {
		t.Error("expected at least some placemarks")
	}

	// Should have folders
	if folderCount == 0 {
		t.Error("expected at least some folders")
	}
}

// TestLargeFileCatalinaPoints tests parsing a larger file with many points
func TestLargeFileCatalinaPoints(t *testing.T) {
	path := "testdata/catalina-points.kml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("catalina-points.kml not found")
	}

	doc, err := ParseFile(path)
	if err != nil {
		t.Fatalf("failed to parse catalina-points.kml: %v", err)
	}

	placemarks := doc.Placemarks()
	t.Logf("catalina-points.kml contains %d placemarks", len(placemarks))

	if len(placemarks) == 0 {
		t.Error("expected placemarks in catalina-points.kml")
	}

	// Check that all placemarks have geometry
	for i, pm := range placemarks {
		if pm.Geometry == nil {
			t.Errorf("placemark %d (%s) has no geometry", i, pm.Name)
		}
	}
}

// TestBoundsCalculation tests calculating bounds on real-world data
func TestBoundsCalculation(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <Placemark>
      <name>P1</name>
      <Point><coordinates>-122.0,37.0</coordinates></Point>
    </Placemark>
    <Placemark>
      <name>P2</name>
      <Point><coordinates>-121.0,38.0</coordinates></Point>
    </Placemark>
    <Placemark>
      <name>P3</name>
      <Point><coordinates>-123.0,36.0</coordinates></Point>
    </Placemark>
  </Document>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	sw, ne := doc.Bounds()

	// SW should be min lon, min lat
	if !floatNear(sw.Lon, -123.0, 0.0001) {
		t.Errorf("expected SW lon -123.0, got %f", sw.Lon)
	}
	if !floatNear(sw.Lat, 36.0, 0.0001) {
		t.Errorf("expected SW lat 36.0, got %f", sw.Lat)
	}

	// NE should be max lon, max lat
	if !floatNear(ne.Lon, -121.0, 0.0001) {
		t.Errorf("expected NE lon -121.0, got %f", ne.Lon)
	}
	if !floatNear(ne.Lat, 38.0, 0.0001) {
		t.Errorf("expected NE lat 38.0, got %f", ne.Lat)
	}
}

// TestFindByID tests finding features by ID
func TestFindByID(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document id="doc1">
    <name>Test Document</name>
    <Folder id="folder1">
      <name>Test Folder</name>
      <Placemark id="pm1">
        <name>Test Placemark</name>
        <Point><coordinates>0,0</coordinates></Point>
      </Placemark>
    </Folder>
  </Document>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Find document
	f := doc.FindByID("doc1")
	if f == nil {
		t.Error("failed to find doc1")
	} else if d, ok := f.(*Document); !ok || d.Name != "Test Document" {
		t.Error("found wrong feature for doc1")
	}

	// Find folder
	f = doc.FindByID("folder1")
	if f == nil {
		t.Error("failed to find folder1")
	} else if folder, ok := f.(*Folder); !ok || folder.Name != "Test Folder" {
		t.Error("found wrong feature for folder1")
	}

	// Find placemark
	f = doc.FindByID("pm1")
	if f == nil {
		t.Error("failed to find pm1")
	} else if pm, ok := f.(*Placemark); !ok || pm.Name != "Test Placemark" {
		t.Error("found wrong feature for pm1")
	}

	// Search for non-existent ID
	f = doc.FindByID("nonexistent")
	if f != nil {
		t.Error("expected nil for non-existent ID")
	}
}

// TestStyleReferences tests that style references work correctly
func TestStyleReferences(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <Style id="myStyle">
      <LineStyle>
        <color>ff0000ff</color>
        <width>4</width>
      </LineStyle>
      <PolyStyle>
        <color>7f00ff00</color>
      </PolyStyle>
    </Style>
    <Placemark>
      <name>Styled Placemark</name>
      <styleUrl>#myStyle</styleUrl>
      <LineString>
        <coordinates>0,0 1,1 2,2</coordinates>
      </LineString>
    </Placemark>
  </Document>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	document := doc.Feature.(*Document)

	// Check style was parsed
	if len(document.Styles) != 1 {
		t.Fatalf("expected 1 style, got %d", len(document.Styles))
	}

	style := document.Styles[0]
	if style.ID != "myStyle" {
		t.Errorf("expected style ID 'myStyle', got '%s'", style.ID)
	}

	// Verify LineStyle color (KML AABBGGRR format: ff0000ff = opaque red)
	if style.LineStyle == nil {
		t.Fatal("expected LineStyle")
	}

	// The color should parse as A=ff, B=00, G=00, R=ff (red)
	expectedColor := RGBA(255, 0, 0, 255) // Creates Color with correct KML ordering
	if style.LineStyle.Color != expectedColor {
		t.Errorf("unexpected line color: got A=%d B=%d G=%d R=%d, expected A=%d B=%d G=%d R=%d",
			style.LineStyle.Color.A, style.LineStyle.Color.B, style.LineStyle.Color.G, style.LineStyle.Color.R,
			expectedColor.A, expectedColor.B, expectedColor.G, expectedColor.R)
	}

	// Check placemark references the style
	pm := document.Features[0].(*Placemark)
	if pm.StyleURL != "#myStyle" {
		t.Errorf("expected styleUrl '#myStyle', got '%s'", pm.StyleURL)
	}
}

// TestRoundTripRealWorldFile tests that real-world files round-trip correctly
func TestRoundTripRealWorldFile(t *testing.T) {
	path := "testdata/placemark.kml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("placemark.kml not found")
	}

	// Parse original
	doc1, err := ParseFile(path)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Write to bytes
	bytes1, err := doc1.Bytes()
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	// Parse again
	doc2, err := ParseBytes(bytes1)
	if err != nil {
		t.Fatalf("failed to re-parse: %v", err)
	}

	// Get placemarks from both
	pms1 := doc1.Placemarks()
	pms2 := doc2.Placemarks()

	if len(pms1) != len(pms2) {
		t.Fatalf("placemark count mismatch: %d vs %d", len(pms1), len(pms2))
	}

	// Compare first placemark
	if len(pms1) > 0 {
		if pms1[0].Name != pms2[0].Name {
			t.Errorf("name mismatch: %s vs %s", pms1[0].Name, pms2[0].Name)
		}

		pt1, ok1 := pms1[0].Geometry.(*Point)
		pt2, ok2 := pms2[0].Geometry.(*Point)

		if ok1 && ok2 {
			if !floatNear(pt1.Coordinates.Lon, pt2.Coordinates.Lon, 0.000001) {
				t.Errorf("lon mismatch: %f vs %f", pt1.Coordinates.Lon, pt2.Coordinates.Lon)
			}
			if !floatNear(pt1.Coordinates.Lat, pt2.Coordinates.Lat, 0.000001) {
				t.Errorf("lat mismatch: %f vs %f", pt1.Coordinates.Lat, pt2.Coordinates.Lat)
			}
		}
	}
}

// TestCDATAHandling tests that CDATA sections in descriptions are preserved
func TestCDATAHandling(t *testing.T) {
	kmlData := `<?xml version="1.0" encoding="utf-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Placemark>
    <name>CDATA Test</name>
    <description><![CDATA[
      <b>Bold text</b><br/>
      <i>Italic text</i>
    ]]></description>
    <Point><coordinates>0,0</coordinates></Point>
  </Placemark>
</kml>`

	doc, err := ParseBytes([]byte(kmlData))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	pm := doc.Feature.(*Placemark)

	// Description should contain the HTML tags
	if !strings.Contains(pm.Description, "<b>Bold text</b>") {
		t.Errorf("expected HTML in description, got: %s", pm.Description)
	}
}

// Helper function
func floatNear(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
