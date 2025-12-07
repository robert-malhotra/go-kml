package kml

import (
	"testing"
)

// TestBuilderSimple tests basic builder usage
func TestBuilderSimple(t *testing.T) {
	kml := NewKMLBuilder().
		Document("Test").
		Build()

	// Verify KML structure is correct
	if kml == nil {
		t.Fatal("Expected non-nil KML")
	}

	if kml.Xmlns != DefaultNamespace {
		t.Errorf("Expected namespace %q, got %q", DefaultNamespace, kml.Xmlns)
	}

	// Verify document name is set
	doc, ok := kml.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Feature to be *Document")
	}

	if doc.Name != "Test" {
		t.Errorf("Expected document name %q, got %q", "Test", doc.Name)
	}
}

// TestBuilderWithPlacemark tests adding placemarks
func TestBuilderWithPlacemark(t *testing.T) {
	kml := NewKMLBuilder().
		Document("Test").
		Placemark("Test Point").
		Point(-122.0, 37.0).
		Done().(*DocumentBuilder).
		Build()

	// Verify KML structure
	doc, ok := kml.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Feature to be *Document")
	}

	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature, got %d", len(doc.Features))
	}

	// Verify placemark
	placemark, ok := doc.Features[0].(*Placemark)
	if !ok {
		t.Fatal("Expected feature to be *Placemark")
	}

	if placemark.Name != "Test Point" {
		t.Errorf("Expected placemark name %q, got %q", "Test Point", placemark.Name)
	}

	// Verify geometry
	point, ok := placemark.Geometry.(*Point)
	if !ok {
		t.Fatal("Expected geometry to be *Point")
	}

	if point.Coordinates.Lon != -122.0 {
		t.Errorf("Expected longitude -122.0, got %f", point.Coordinates.Lon)
	}

	if point.Coordinates.Lat != 37.0 {
		t.Errorf("Expected latitude 37.0, got %f", point.Coordinates.Lat)
	}

	if point.Coordinates.Alt != 0 {
		t.Errorf("Expected altitude 0, got %f", point.Coordinates.Alt)
	}
}

// TestBuilderWithFolder tests nested folders
func TestBuilderWithFolder(t *testing.T) {
	kml := NewKMLBuilder().
		Document("Test").
		Folder("Folder1").
		Placemark("P1").Point(-122.0, 37.0).Done().(*FolderBuilder).
		Folder("Subfolder").
		Placemark("P2").Point(-121.0, 38.0).Done().(*FolderBuilder).
		Done().(*FolderBuilder).
		Done().(*DocumentBuilder).
		Build()

	// Verify document
	doc, ok := kml.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Feature to be *Document")
	}

	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature in document, got %d", len(doc.Features))
	}

	// Verify Folder1
	folder1, ok := doc.Features[0].(*Folder)
	if !ok {
		t.Fatal("Expected feature to be *Folder")
	}

	if folder1.Name != "Folder1" {
		t.Errorf("Expected folder name %q, got %q", "Folder1", folder1.Name)
	}

	if len(folder1.Features) != 2 {
		t.Fatalf("Expected 2 features in Folder1, got %d", len(folder1.Features))
	}

	// Verify P1 placemark
	p1, ok := folder1.Features[0].(*Placemark)
	if !ok {
		t.Fatal("Expected first feature in Folder1 to be *Placemark")
	}

	if p1.Name != "P1" {
		t.Errorf("Expected placemark name %q, got %q", "P1", p1.Name)
	}

	point1, ok := p1.Geometry.(*Point)
	if !ok {
		t.Fatal("Expected P1 geometry to be *Point")
	}

	if point1.Coordinates.Lon != -122.0 || point1.Coordinates.Lat != 37.0 {
		t.Errorf("Expected P1 coordinates (-122.0, 37.0), got (%f, %f)",
			point1.Coordinates.Lon, point1.Coordinates.Lat)
	}

	// Verify Subfolder
	subfolder, ok := folder1.Features[1].(*Folder)
	if !ok {
		t.Fatal("Expected second feature in Folder1 to be *Folder")
	}

	if subfolder.Name != "Subfolder" {
		t.Errorf("Expected folder name %q, got %q", "Subfolder", subfolder.Name)
	}

	if len(subfolder.Features) != 1 {
		t.Fatalf("Expected 1 feature in Subfolder, got %d", len(subfolder.Features))
	}

	// Verify P2 placemark
	p2, ok := subfolder.Features[0].(*Placemark)
	if !ok {
		t.Fatal("Expected feature in Subfolder to be *Placemark")
	}

	if p2.Name != "P2" {
		t.Errorf("Expected placemark name %q, got %q", "P2", p2.Name)
	}

	point2, ok := p2.Geometry.(*Point)
	if !ok {
		t.Fatal("Expected P2 geometry to be *Point")
	}

	if point2.Coordinates.Lon != -121.0 || point2.Coordinates.Lat != 38.0 {
		t.Errorf("Expected P2 coordinates (-121.0, 38.0), got (%f, %f)",
			point2.Coordinates.Lon, point2.Coordinates.Lat)
	}
}

// TestBuilderWithStyle tests style definitions
func TestBuilderWithStyle(t *testing.T) {
	kml := NewKMLBuilder().
		Document("Test").
		Style("myStyle").
		LineStyle().Color(Red).Width(2).Done().
		PolyStyle().Fill(true).Done().
		Done().
		Build()

	// Verify document
	doc, ok := kml.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Feature to be *Document")
	}

	if len(doc.Styles) != 1 {
		t.Fatalf("Expected 1 style, got %d", len(doc.Styles))
	}

	// Verify style
	style := doc.Styles[0]
	if style.ID != "myStyle" {
		t.Errorf("Expected style ID %q, got %q", "myStyle", style.ID)
	}

	// Verify LineStyle
	if style.LineStyle == nil {
		t.Fatal("Expected LineStyle to be set")
	}

	if style.LineStyle.Color != Red {
		t.Errorf("Expected LineStyle color to be Red (%v), got %v", Red, style.LineStyle.Color)
	}

	if style.LineStyle.Width != 2 {
		t.Errorf("Expected LineStyle width 2, got %f", style.LineStyle.Width)
	}

	// Verify PolyStyle
	if style.PolyStyle == nil {
		t.Fatal("Expected PolyStyle to be set")
	}

	if style.PolyStyle.Fill == nil {
		t.Fatal("Expected PolyStyle.Fill to be set")
	}

	if !*style.PolyStyle.Fill {
		t.Error("Expected PolyStyle.Fill to be true")
	}
}

// TestBuilderLineString tests LineString geometry
func TestBuilderLineString(t *testing.T) {
	coords := []Coordinate{
		Coord(-122.0, 37.0),
		Coord(-122.1, 37.1),
		Coord(-122.2, 37.2),
	}

	kml := NewKMLBuilder().
		Document("Test").
		Placemark("Test Line").
		LineString(coords...).
		Done().(*DocumentBuilder).
		Build()

	// Verify document
	doc, ok := kml.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Feature to be *Document")
	}

	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature, got %d", len(doc.Features))
	}

	// Verify placemark
	placemark, ok := doc.Features[0].(*Placemark)
	if !ok {
		t.Fatal("Expected feature to be *Placemark")
	}

	if placemark.Name != "Test Line" {
		t.Errorf("Expected placemark name %q, got %q", "Test Line", placemark.Name)
	}

	// Verify LineString geometry
	lineString, ok := placemark.Geometry.(*LineString)
	if !ok {
		t.Fatal("Expected geometry to be *LineString")
	}

	if len(lineString.Coordinates) != 3 {
		t.Fatalf("Expected 3 coordinates, got %d", len(lineString.Coordinates))
	}

	// Verify coordinates
	for i, expected := range coords {
		actual := lineString.Coordinates[i]
		if actual.Lon != expected.Lon || actual.Lat != expected.Lat || actual.Alt != expected.Alt {
			t.Errorf("Coordinate %d: expected (%f, %f, %f), got (%f, %f, %f)",
				i, expected.Lon, expected.Lat, expected.Alt,
				actual.Lon, actual.Lat, actual.Alt)
		}
	}
}

// TestBuilderPolygon tests Polygon geometry
func TestBuilderPolygon(t *testing.T) {
	outer := []Coordinate{
		Coord(-122.0, 37.0),
		Coord(-122.0, 38.0),
		Coord(-121.0, 38.0),
		Coord(-121.0, 37.0),
		Coord(-122.0, 37.0), // Close the ring
	}

	inner := []Coordinate{
		Coord(-121.8, 37.2),
		Coord(-121.8, 37.8),
		Coord(-121.2, 37.8),
		Coord(-121.2, 37.2),
		Coord(-121.8, 37.2), // Close the ring
	}

	kml := NewKMLBuilder().
		Document("Test").
		Placemark("Test Polygon").
		Polygon(outer, inner).
		Done().(*DocumentBuilder).
		Build()

	// Verify document
	doc, ok := kml.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Feature to be *Document")
	}

	if len(doc.Features) != 1 {
		t.Fatalf("Expected 1 feature, got %d", len(doc.Features))
	}

	// Verify placemark
	placemark, ok := doc.Features[0].(*Placemark)
	if !ok {
		t.Fatal("Expected feature to be *Placemark")
	}

	if placemark.Name != "Test Polygon" {
		t.Errorf("Expected placemark name %q, got %q", "Test Polygon", placemark.Name)
	}

	// Verify Polygon geometry
	polygon, ok := placemark.Geometry.(*Polygon)
	if !ok {
		t.Fatal("Expected geometry to be *Polygon")
	}

	// Verify outer boundary
	if len(polygon.OuterBoundary.Coordinates) != 5 {
		t.Fatalf("Expected 5 outer coordinates, got %d", len(polygon.OuterBoundary.Coordinates))
	}

	for i, expected := range outer {
		actual := polygon.OuterBoundary.Coordinates[i]
		if actual.Lon != expected.Lon || actual.Lat != expected.Lat {
			t.Errorf("Outer coordinate %d: expected (%f, %f), got (%f, %f)",
				i, expected.Lon, expected.Lat, actual.Lon, actual.Lat)
		}
	}

	// Verify inner boundary
	if len(polygon.InnerBoundaries) != 1 {
		t.Fatalf("Expected 1 inner boundary, got %d", len(polygon.InnerBoundaries))
	}

	if len(polygon.InnerBoundaries[0].Coordinates) != 5 {
		t.Fatalf("Expected 5 inner coordinates, got %d", len(polygon.InnerBoundaries[0].Coordinates))
	}

	for i, expected := range inner {
		actual := polygon.InnerBoundaries[0].Coordinates[i]
		if actual.Lon != expected.Lon || actual.Lat != expected.Lat {
			t.Errorf("Inner coordinate %d: expected (%f, %f), got (%f, %f)",
				i, expected.Lon, expected.Lat, actual.Lon, actual.Lat)
		}
	}
}

// TestBuilderChaining verifies all methods return correct builder types for chaining
func TestBuilderChaining(t *testing.T) {
	// Test KMLBuilder chaining
	var kb *KMLBuilder = NewKMLBuilder()
	if kb == nil {
		t.Fatal("NewKMLBuilder() returned nil")
	}

	// Test DocumentBuilder chaining
	var db *DocumentBuilder = kb.Document("Test")
	if db == nil {
		t.Fatal("Document() returned nil")
	}

	db = db.Name("Updated Name")
	if db == nil {
		t.Fatal("Name() returned nil")
	}

	db = db.Description("Test description")
	if db == nil {
		t.Fatal("Description() returned nil")
	}

	db = db.Open(true)
	if db == nil {
		t.Fatal("Open() returned nil")
	}

	// Test StyleBuilder chaining
	var sb *StyleBuilder = db.Style("testStyle")
	if sb == nil {
		t.Fatal("Style() returned nil")
	}

	// Test LineStyleBuilder chaining
	var lsb *LineStyleBuilder = sb.LineStyle()
	if lsb == nil {
		t.Fatal("LineStyle() returned nil")
	}

	lsb = lsb.Color(Red)
	if lsb == nil {
		t.Fatal("LineStyle.Color() returned nil")
	}

	lsb = lsb.Width(2.5)
	if lsb == nil {
		t.Fatal("LineStyle.Width() returned nil")
	}

	sb = lsb.Done()
	if sb == nil {
		t.Fatal("LineStyleBuilder.Done() returned nil")
	}

	// Test PolyStyleBuilder chaining
	var psb *PolyStyleBuilder = sb.PolyStyle()
	if psb == nil {
		t.Fatal("PolyStyle() returned nil")
	}

	psb = psb.Color(Blue)
	if psb == nil {
		t.Fatal("PolyStyle.Color() returned nil")
	}

	psb = psb.Fill(true)
	if psb == nil {
		t.Fatal("PolyStyle.Fill() returned nil")
	}

	psb = psb.Outline(false)
	if psb == nil {
		t.Fatal("PolyStyle.Outline() returned nil")
	}

	sb = psb.Done()
	if sb == nil {
		t.Fatal("PolyStyleBuilder.Done() returned nil")
	}

	// Test IconStyleBuilder chaining
	var isb *IconStyleBuilder = sb.IconStyle()
	if isb == nil {
		t.Fatal("IconStyle() returned nil")
	}

	isb = isb.Color(Green)
	if isb == nil {
		t.Fatal("IconStyle.Color() returned nil")
	}

	isb = isb.Scale(1.5)
	if isb == nil {
		t.Fatal("IconStyle.Scale() returned nil")
	}

	isb = isb.Heading(45.0)
	if isb == nil {
		t.Fatal("IconStyle.Heading() returned nil")
	}

	isb = isb.Icon("http://example.com/icon.png")
	if isb == nil {
		t.Fatal("IconStyle.Icon() returned nil")
	}

	isb = isb.HotSpot(0.5, 0.5, "fraction", "fraction")
	if isb == nil {
		t.Fatal("IconStyle.HotSpot() returned nil")
	}

	sb = isb.Done()
	if sb == nil {
		t.Fatal("IconStyleBuilder.Done() returned nil")
	}

	// Test LabelStyleBuilder chaining
	var lblsb *LabelStyleBuilder = sb.LabelStyle()
	if lblsb == nil {
		t.Fatal("LabelStyle() returned nil")
	}

	lblsb = lblsb.Color(White)
	if lblsb == nil {
		t.Fatal("LabelStyle.Color() returned nil")
	}

	lblsb = lblsb.Scale(2.0)
	if lblsb == nil {
		t.Fatal("LabelStyle.Scale() returned nil")
	}

	sb = lblsb.Done()
	if sb == nil {
		t.Fatal("LabelStyleBuilder.Done() returned nil")
	}

	db = sb.Done()
	if db == nil {
		t.Fatal("StyleBuilder.Done() returned nil")
	}

	// Test FolderBuilder chaining
	var fb *FolderBuilder = db.Folder("TestFolder")
	if fb == nil {
		t.Fatal("Folder() returned nil")
	}

	fb = fb.Name("Updated Folder")
	if fb == nil {
		t.Fatal("FolderBuilder.Name() returned nil")
	}

	fb = fb.Description("Folder description")
	if fb == nil {
		t.Fatal("FolderBuilder.Description() returned nil")
	}

	fb = fb.Open(true)
	if fb == nil {
		t.Fatal("FolderBuilder.Open() returned nil")
	}

	// Test nested FolderBuilder
	var fb2 *FolderBuilder = fb.Folder("NestedFolder")
	if fb2 == nil {
		t.Fatal("Nested Folder() returned nil")
	}

	// Test PlacemarkBuilder from FolderBuilder
	var pb *PlacemarkBuilder = fb2.Placemark("TestPoint")
	if pb == nil {
		t.Fatal("Placemark() returned nil")
	}

	pb = pb.Name("Updated Point")
	if pb == nil {
		t.Fatal("PlacemarkBuilder.Name() returned nil")
	}

	pb = pb.Description("Point description")
	if pb == nil {
		t.Fatal("PlacemarkBuilder.Description() returned nil")
	}

	pb = pb.StyleURL("#myStyle")
	if pb == nil {
		t.Fatal("PlacemarkBuilder.StyleURL() returned nil")
	}

	pb = pb.Point(-122.0, 37.0, 100.0)
	if pb == nil {
		t.Fatal("PlacemarkBuilder.Point() returned nil")
	}

	// Test LineString chaining
	coords := []Coordinate{Coord(-122.0, 37.0), Coord(-121.0, 38.0)}
	pb = pb.LineString(coords...)
	if pb == nil {
		t.Fatal("PlacemarkBuilder.LineString() returned nil")
	}

	// Test Polygon chaining
	outer := []Coordinate{
		Coord(-122.0, 37.0),
		Coord(-122.0, 38.0),
		Coord(-121.0, 38.0),
		Coord(-122.0, 37.0),
	}
	pb = pb.Polygon(outer)
	if pb == nil {
		t.Fatal("PlacemarkBuilder.Polygon() returned nil")
	}

	// Test Done() returns to correct parent
	parentFb := pb.Done()
	if parentFb == nil {
		t.Fatal("PlacemarkBuilder.Done() returned nil")
	}

	// Verify the type is interface{} but actually FolderBuilder
	if _, ok := parentFb.(*FolderBuilder); !ok {
		t.Fatal("PlacemarkBuilder.Done() did not return *FolderBuilder")
	}

	// Test Build shortcuts
	kmlFromDoc := db.Build()
	if kmlFromDoc == nil {
		t.Fatal("DocumentBuilder.Build() returned nil")
	}

	kmlFromKb := kb.Build()
	if kmlFromKb == nil {
		t.Fatal("KMLBuilder.Build() returned nil")
	}
}

// TestBuilderComplexDocument tests a complex document with multiple features and styles
func TestBuilderComplexDocument(t *testing.T) {
	kml := NewKMLBuilder().
		Document("Complex Test").
		Description("A complex KML document").
		Open(true).
		Style("redLine").
		LineStyle().Color(Red).Width(3).Done().
		Done().
		Style("bluePoly").
		PolyStyle().Color(Blue).Fill(true).Outline(true).Done().
		Done().
		Folder("Points").
		Placemark("Point 1").Point(-122.0, 37.0).StyleURL("#redLine").Done().(*FolderBuilder).
		Placemark("Point 2").Point(-121.0, 38.0, 100.0).Done().(*FolderBuilder).
		Done().(*DocumentBuilder).
		Folder("Lines").
		Placemark("Line 1").
		LineString(Coord(-122.0, 37.0), Coord(-121.0, 38.0)).
		StyleURL("#redLine").
		Done().(*FolderBuilder).
		Done().(*DocumentBuilder).
		Placemark("Polygon 1").
		Polygon([]Coordinate{
			Coord(-122.0, 37.0),
			Coord(-122.0, 38.0),
			Coord(-121.0, 38.0),
			Coord(-122.0, 37.0),
		}).
		StyleURL("#bluePoly").
		Done().(*DocumentBuilder).
		Build()

	// Verify document
	doc, ok := kml.Feature.(*Document)
	if !ok {
		t.Fatal("Expected Feature to be *Document")
	}

	if doc.Name != "Complex Test" {
		t.Errorf("Expected document name %q, got %q", "Complex Test", doc.Name)
	}

	if doc.Description != "A complex KML document" {
		t.Errorf("Expected description %q, got %q", "A complex KML document", doc.Description)
	}

	if !doc.Open {
		t.Error("Expected document to be open")
	}

	// Verify styles
	if len(doc.Styles) != 2 {
		t.Fatalf("Expected 2 styles, got %d", len(doc.Styles))
	}

	// Verify features (2 folders + 1 placemark)
	if len(doc.Features) != 3 {
		t.Fatalf("Expected 3 features, got %d", len(doc.Features))
	}

	// Verify Points folder
	pointsFolder, ok := doc.Features[0].(*Folder)
	if !ok {
		t.Fatal("Expected first feature to be *Folder")
	}

	if pointsFolder.Name != "Points" {
		t.Errorf("Expected folder name %q, got %q", "Points", pointsFolder.Name)
	}

	if len(pointsFolder.Features) != 2 {
		t.Fatalf("Expected 2 placemarks in Points folder, got %d", len(pointsFolder.Features))
	}

	// Verify Lines folder
	linesFolder, ok := doc.Features[1].(*Folder)
	if !ok {
		t.Fatal("Expected second feature to be *Folder")
	}

	if linesFolder.Name != "Lines" {
		t.Errorf("Expected folder name %q, got %q", "Lines", linesFolder.Name)
	}

	if len(linesFolder.Features) != 1 {
		t.Fatalf("Expected 1 placemark in Lines folder, got %d", len(linesFolder.Features))
	}

	// Verify standalone Polygon placemark
	polygonPlacemark, ok := doc.Features[2].(*Placemark)
	if !ok {
		t.Fatal("Expected third feature to be *Placemark")
	}

	if polygonPlacemark.Name != "Polygon 1" {
		t.Errorf("Expected placemark name %q, got %q", "Polygon 1", polygonPlacemark.Name)
	}

	if polygonPlacemark.StyleURL != "#bluePoly" {
		t.Errorf("Expected style URL %q, got %q", "#bluePoly", polygonPlacemark.StyleURL)
	}

	polygon, ok := polygonPlacemark.Geometry.(*Polygon)
	if !ok {
		t.Fatal("Expected geometry to be *Polygon")
	}

	if len(polygon.OuterBoundary.Coordinates) != 4 {
		t.Errorf("Expected 4 outer coordinates, got %d", len(polygon.OuterBoundary.Coordinates))
	}
}
