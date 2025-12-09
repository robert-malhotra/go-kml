package kml

import "encoding/json"

// GeoJSONGeometry represents a GeoJSON geometry object.
type GeoJSONGeometry struct {
	Type        string            `json:"type"`
	Coordinates any               `json:"coordinates,omitempty"`
	Geometries  []GeoJSONGeometry `json:"geometries,omitempty"`
}

// ToGeoJSON converts a Point to a GeoJSON geometry.
func (p *Point) ToGeoJSON() GeoJSONGeometry {
	return GeoJSONGeometry{
		Type:        "Point",
		Coordinates: coordToGeoJSON(p.Coordinates),
	}
}

// ToGeoJSON converts a LineString to a GeoJSON geometry.
func (ls *LineString) ToGeoJSON() GeoJSONGeometry {
	return GeoJSONGeometry{
		Type:        "LineString",
		Coordinates: coordsToGeoJSON(ls.Coordinates),
	}
}

// ToGeoJSON converts a LinearRing to a GeoJSON geometry (as a closed LineString).
func (lr *LinearRing) ToGeoJSON() GeoJSONGeometry {
	return GeoJSONGeometry{
		Type:        "LineString",
		Coordinates: coordsToGeoJSON(lr.Coordinates),
	}
}

// ToGeoJSON converts a Polygon to a GeoJSON geometry.
func (p *Polygon) ToGeoJSON() GeoJSONGeometry {
	rings := make([][][]float64, 0, 1+len(p.InnerBoundaries))

	// Outer boundary
	rings = append(rings, coordsToGeoJSON(p.OuterBoundary.Coordinates))

	// Inner boundaries (holes)
	for _, inner := range p.InnerBoundaries {
		rings = append(rings, coordsToGeoJSON(inner.Coordinates))
	}

	return GeoJSONGeometry{
		Type:        "Polygon",
		Coordinates: rings,
	}
}

// ToGeoJSON converts a MultiGeometry to a GeoJSON GeometryCollection.
func (mg *MultiGeometry) ToGeoJSON() GeoJSONGeometry {
	geometries := make([]GeoJSONGeometry, 0, len(mg.Geometries))

	for _, geom := range mg.Geometries {
		if g := geometryToGeoJSON(geom); g != nil {
			geometries = append(geometries, *g)
		}
	}

	return GeoJSONGeometry{
		Type:       "GeometryCollection",
		Geometries: geometries,
	}
}

// ToGeoJSON converts any KML Geometry to a GeoJSON geometry.
// Returns nil if the geometry type is not supported.
func ToGeoJSON(g Geometry) *GeoJSONGeometry {
	return geometryToGeoJSON(g)
}

// geometryToGeoJSON converts a Geometry interface to a GeoJSONGeometry.
func geometryToGeoJSON(g Geometry) *GeoJSONGeometry {
	if g == nil {
		return nil
	}

	var result GeoJSONGeometry
	switch geom := g.(type) {
	case *Point:
		result = geom.ToGeoJSON()
	case *LineString:
		result = geom.ToGeoJSON()
	case *LinearRing:
		result = geom.ToGeoJSON()
	case *Polygon:
		result = geom.ToGeoJSON()
	case *MultiGeometry:
		result = geom.ToGeoJSON()
	default:
		return nil
	}
	return &result
}

// coordToGeoJSON converts a single Coordinate to GeoJSON format [lon, lat] or [lon, lat, alt].
func coordToGeoJSON(c Coordinate) []float64 {
	if c.Alt != 0 {
		return []float64{c.Lon, c.Lat, c.Alt}
	}
	return []float64{c.Lon, c.Lat}
}

// coordsToGeoJSON converts a slice of Coordinates to GeoJSON format.
func coordsToGeoJSON(coords []Coordinate) [][]float64 {
	result := make([][]float64, len(coords))
	for i, c := range coords {
		result[i] = coordToGeoJSON(c)
	}
	return result
}

// MarshalJSON implements json.Marshaler for GeoJSONGeometry.
func (g GeoJSONGeometry) MarshalJSON() ([]byte, error) {
	type Alias GeoJSONGeometry

	// Use Geometries for GeometryCollection, Coordinates for others
	if g.Type == "GeometryCollection" {
		return json.Marshal(struct {
			Type       string            `json:"type"`
			Geometries []GeoJSONGeometry `json:"geometries"`
		}{
			Type:       g.Type,
			Geometries: g.Geometries,
		})
	}

	return json.Marshal(struct {
		Type        string `json:"type"`
		Coordinates any    `json:"coordinates"`
	}{
		Type:        g.Type,
		Coordinates: g.Coordinates,
	})
}

// String returns the GeoJSON representation as a JSON string.
func (g GeoJSONGeometry) String() string {
	data, err := json.Marshal(g)
	if err != nil {
		return ""
	}
	return string(data)
}
