package kml

import (
	"encoding/json"
	"testing"
)

func TestPointToGeoJSON(t *testing.T) {
	p := &Point{
		Coordinates: Coord(-122.0, 37.0),
	}

	gj := p.ToGeoJSON()

	if gj.Type != "Point" {
		t.Errorf("expected type Point, got %s", gj.Type)
	}

	coords, ok := gj.Coordinates.([]float64)
	if !ok {
		t.Fatalf("expected []float64 coordinates, got %T", gj.Coordinates)
	}

	if len(coords) != 2 {
		t.Errorf("expected 2 coordinates, got %d", len(coords))
	}

	if coords[0] != -122.0 || coords[1] != 37.0 {
		t.Errorf("expected [-122.0, 37.0], got %v", coords)
	}
}

func TestPointWithAltitudeToGeoJSON(t *testing.T) {
	p := &Point{
		Coordinates: Coord(-122.0, 37.0, 100.0),
	}

	gj := p.ToGeoJSON()

	coords, ok := gj.Coordinates.([]float64)
	if !ok {
		t.Fatalf("expected []float64 coordinates, got %T", gj.Coordinates)
	}

	if len(coords) != 3 {
		t.Errorf("expected 3 coordinates, got %d", len(coords))
	}

	if coords[2] != 100.0 {
		t.Errorf("expected altitude 100.0, got %f", coords[2])
	}
}

func TestLineStringToGeoJSON(t *testing.T) {
	ls := &LineString{
		Coordinates: []Coordinate{
			Coord(-122.0, 37.0),
			Coord(-121.0, 38.0),
			Coord(-120.0, 37.5),
		},
	}

	gj := ls.ToGeoJSON()

	if gj.Type != "LineString" {
		t.Errorf("expected type LineString, got %s", gj.Type)
	}

	coords, ok := gj.Coordinates.([][]float64)
	if !ok {
		t.Fatalf("expected [][]float64 coordinates, got %T", gj.Coordinates)
	}

	if len(coords) != 3 {
		t.Errorf("expected 3 coordinate pairs, got %d", len(coords))
	}

	if coords[0][0] != -122.0 || coords[0][1] != 37.0 {
		t.Errorf("expected first coord [-122.0, 37.0], got %v", coords[0])
	}
}

func TestPolygonToGeoJSON(t *testing.T) {
	p := &Polygon{
		OuterBoundary: LinearRing{
			Coordinates: []Coordinate{
				Coord(-122.0, 37.0),
				Coord(-121.0, 37.0),
				Coord(-121.0, 38.0),
				Coord(-122.0, 38.0),
				Coord(-122.0, 37.0),
			},
		},
	}

	gj := p.ToGeoJSON()

	if gj.Type != "Polygon" {
		t.Errorf("expected type Polygon, got %s", gj.Type)
	}

	rings, ok := gj.Coordinates.([][][]float64)
	if !ok {
		t.Fatalf("expected [][][]float64 coordinates, got %T", gj.Coordinates)
	}

	if len(rings) != 1 {
		t.Errorf("expected 1 ring, got %d", len(rings))
	}

	if len(rings[0]) != 5 {
		t.Errorf("expected 5 coordinates in outer ring, got %d", len(rings[0]))
	}
}

func TestPolygonWithHolesToGeoJSON(t *testing.T) {
	p := &Polygon{
		OuterBoundary: LinearRing{
			Coordinates: []Coordinate{
				Coord(-122.0, 37.0),
				Coord(-121.0, 37.0),
				Coord(-121.0, 38.0),
				Coord(-122.0, 38.0),
				Coord(-122.0, 37.0),
			},
		},
		InnerBoundaries: []LinearRing{
			{
				Coordinates: []Coordinate{
					Coord(-121.8, 37.2),
					Coord(-121.2, 37.2),
					Coord(-121.2, 37.8),
					Coord(-121.8, 37.8),
					Coord(-121.8, 37.2),
				},
			},
		},
	}

	gj := p.ToGeoJSON()

	rings, ok := gj.Coordinates.([][][]float64)
	if !ok {
		t.Fatalf("expected [][][]float64 coordinates, got %T", gj.Coordinates)
	}

	if len(rings) != 2 {
		t.Errorf("expected 2 rings (outer + 1 hole), got %d", len(rings))
	}

	if len(rings[1]) != 5 {
		t.Errorf("expected 5 coordinates in inner ring, got %d", len(rings[1]))
	}
}

func TestMultiGeometryToGeoJSON(t *testing.T) {
	mg := &MultiGeometry{
		Geometries: []Geometry{
			&Point{Coordinates: Coord(-122.0, 37.0)},
			&LineString{
				Coordinates: []Coordinate{
					Coord(-122.0, 37.0),
					Coord(-121.0, 38.0),
				},
			},
		},
	}

	gj := mg.ToGeoJSON()

	if gj.Type != "GeometryCollection" {
		t.Errorf("expected type GeometryCollection, got %s", gj.Type)
	}

	if len(gj.Geometries) != 2 {
		t.Errorf("expected 2 geometries, got %d", len(gj.Geometries))
	}

	if gj.Geometries[0].Type != "Point" {
		t.Errorf("expected first geometry to be Point, got %s", gj.Geometries[0].Type)
	}

	if gj.Geometries[1].Type != "LineString" {
		t.Errorf("expected second geometry to be LineString, got %s", gj.Geometries[1].Type)
	}
}

func TestGeometryInterfaceToGeoJSON(t *testing.T) {
	tests := []struct {
		name     string
		geometry Geometry
		expected string
	}{
		{
			name:     "Point",
			geometry: &Point{Coordinates: Coord(-122.0, 37.0)},
			expected: "Point",
		},
		{
			name: "LineString",
			geometry: &LineString{
				Coordinates: []Coordinate{Coord(-122.0, 37.0), Coord(-121.0, 38.0)},
			},
			expected: "LineString",
		},
		{
			name: "Polygon",
			geometry: &Polygon{
				OuterBoundary: LinearRing{
					Coordinates: []Coordinate{
						Coord(-122.0, 37.0),
						Coord(-121.0, 37.0),
						Coord(-121.0, 38.0),
						Coord(-122.0, 37.0),
					},
				},
			},
			expected: "Polygon",
		},
		{
			name: "MultiGeometry",
			geometry: &MultiGeometry{
				Geometries: []Geometry{&Point{Coordinates: Coord(0, 0)}},
			},
			expected: "GeometryCollection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call ToGeoJSON through the interface
			gj := tt.geometry.ToGeoJSON()
			if gj.Type != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, gj.Type)
			}
		})
	}
}

func TestGeoJSONMarshalJSON(t *testing.T) {
	p := &Point{Coordinates: Coord(-122.0, 37.0)}
	gj := p.ToGeoJSON()

	data, err := json.Marshal(gj)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"type":"Point","coordinates":[-122,37]}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestGeoJSONGeometryCollectionMarshalJSON(t *testing.T) {
	mg := &MultiGeometry{
		Geometries: []Geometry{
			&Point{Coordinates: Coord(0, 0)},
		},
	}
	gj := mg.ToGeoJSON()

	data, err := json.Marshal(gj)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[0,0]}]}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestGeoJSONString(t *testing.T) {
	p := &Point{Coordinates: Coord(-122.5, 37.5)}
	gj := p.ToGeoJSON()

	str := gj.String()
	if str == "" {
		t.Error("expected non-empty string")
	}

	// Verify it's valid JSON
	var parsed map[string]any
	if err := json.Unmarshal([]byte(str), &parsed); err != nil {
		t.Errorf("String() output is not valid JSON: %v", err)
	}

	if parsed["type"] != "Point" {
		t.Errorf("expected type Point in JSON, got %v", parsed["type"])
	}
}

func TestLinearRingToGeoJSON(t *testing.T) {
	lr := &LinearRing{
		Coordinates: []Coordinate{
			Coord(-122.0, 37.0),
			Coord(-121.0, 37.0),
			Coord(-121.0, 38.0),
			Coord(-122.0, 37.0),
		},
	}

	gj := lr.ToGeoJSON()

	// LinearRing becomes a LineString in GeoJSON
	if gj.Type != "LineString" {
		t.Errorf("expected type LineString, got %s", gj.Type)
	}

	coords, ok := gj.Coordinates.([][]float64)
	if !ok {
		t.Fatalf("expected [][]float64 coordinates, got %T", gj.Coordinates)
	}

	if len(coords) != 4 {
		t.Errorf("expected 4 coordinates, got %d", len(coords))
	}
}
