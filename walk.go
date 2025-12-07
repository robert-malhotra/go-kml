package kml

import "math"

// Walk traverses all features in a KML document depth-first.
// The callback is called for each feature (Document, Folder, Placemark).
// If the callback returns an error, traversal stops and the error is returned.
func (k *KML) Walk(fn func(Feature) error) error {
	if k.Feature == nil {
		return nil
	}
	return walkFeature(k.Feature, fn)
}

// walkFeature is a recursive helper that walks through features.
func walkFeature(f Feature, fn func(Feature) error) error {
	// Call the callback for this feature
	if err := fn(f); err != nil {
		return err
	}

	// Recursively walk children based on feature type
	switch feature := f.(type) {
	case *Document:
		for _, child := range feature.Features {
			if err := walkFeature(child, fn); err != nil {
				return err
			}
		}
	case *Folder:
		for _, child := range feature.Features {
			if err := walkFeature(child, fn); err != nil {
				return err
			}
		}
	case *Placemark:
		// Placemarks have no child features
	}

	return nil
}

// Placemarks returns all placemarks in the document (recursively searches all folders).
func (k *KML) Placemarks() []*Placemark {
	var placemarks []*Placemark

	k.Walk(func(f Feature) error {
		if p, ok := f.(*Placemark); ok {
			placemarks = append(placemarks, p)
		}
		return nil
	})

	return placemarks
}

// FindByID finds a feature by its ID attribute.
// Returns nil if not found.
func (k *KML) FindByID(id string) Feature {
	var result Feature

	k.Walk(func(f Feature) error {
		// Check ID based on feature type
		switch feature := f.(type) {
		case *Document:
			if feature.ID == id {
				result = feature
				return errStopWalk // Stop walking once found
			}
		case *Folder:
			if feature.ID == id {
				result = feature
				return errStopWalk
			}
		case *Placemark:
			if feature.ID == id {
				result = feature
				return errStopWalk
			}
		}
		return nil
	})

	return result
}

// errStopWalk is a sentinel error used to stop walking.
var errStopWalk = &struct{ error }{error: nil}

// Filter returns all features matching the predicate function.
// Searches recursively through Documents and Folders.
func (k *KML) Filter(fn func(Feature) bool) []Feature {
	var results []Feature

	k.Walk(func(f Feature) error {
		if fn(f) {
			results = append(results, f)
		}
		return nil
	})

	return results
}

// Bounds returns the southwest and northeast corners of a bounding box
// containing all coordinates in the document.
// Returns zero coordinates if the document contains no geometry.
func (k *KML) Bounds() (sw, ne Coordinate) {
	// Initialize with extreme values
	minLon := math.MaxFloat64
	maxLon := -math.MaxFloat64
	minLat := math.MaxFloat64
	maxLat := -math.MaxFloat64

	hasCoords := false

	// Walk all features and collect coordinates
	k.Walk(func(f Feature) error {
		coords := collectCoordinates(f)
		for _, c := range coords {
			hasCoords = true
			if c.Lon < minLon {
				minLon = c.Lon
			}
			if c.Lon > maxLon {
				maxLon = c.Lon
			}
			if c.Lat < minLat {
				minLat = c.Lat
			}
			if c.Lat > maxLat {
				maxLat = c.Lat
			}
		}
		return nil
	})

	// Return zero coordinates if no geometry found
	if !hasCoords {
		return Coordinate{}, Coordinate{}
	}

	// SW corner = (minLon, minLat), NE corner = (maxLon, maxLat)
	sw = Coordinate{Lon: minLon, Lat: minLat}
	ne = Coordinate{Lon: maxLon, Lat: maxLat}

	return sw, ne
}

// collectCoordinates extracts all coordinates from a feature.
func collectCoordinates(f Feature) []Coordinate {
	// Only placemarks have geometry
	p, ok := f.(*Placemark)
	if !ok || p.Geometry == nil {
		return nil
	}

	return getGeometryCoordinates(p.Geometry)
}

// getGeometryCoordinates extracts coordinates from a geometry.
// Handles all geometry types including MultiGeometry recursively.
func getGeometryCoordinates(g Geometry) []Coordinate {
	switch geom := g.(type) {
	case *Point:
		return []Coordinate{geom.Coordinates}

	case *LineString:
		return geom.Coordinates

	case *LinearRing:
		return geom.Coordinates

	case *Polygon:
		// Collect coordinates from outer boundary
		coords := geom.OuterBoundary.Coordinates

		// Add coordinates from inner boundaries (holes)
		for _, inner := range geom.InnerBoundaries {
			coords = append(coords, inner.Coordinates...)
		}

		return coords

	case *MultiGeometry:
		var coords []Coordinate
		// Recursively collect from all child geometries
		for _, child := range geom.Geometries {
			coords = append(coords, getGeometryCoordinates(child)...)
		}
		return coords

	default:
		return nil
	}
}
