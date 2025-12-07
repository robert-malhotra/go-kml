package kml

import (
	"errors"
	"fmt"
)

// ParseError represents an error that occurred during KML parsing.
// It includes location information (line and column) to help identify
// where in the KML document the error occurred.
type ParseError struct {
	Line    int    // Line number where the error occurred (1-based)
	Column  int    // Column number where the error occurred (1-based)
	Message string // Human-readable error message
	Cause   error  // Underlying error that caused this parse error
}

// Error returns a formatted error message including location information.
func (e *ParseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("kml: parse error at line %d, column %d: %s: %v", e.Line, e.Column, e.Message, e.Cause)
	}
	return fmt.Sprintf("kml: parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// Unwrap returns the underlying cause of the parse error.
// This enables error unwrapping with errors.Is and errors.As.
func (e *ParseError) Unwrap() error {
	return e.Cause
}

// ValidationError represents an error that occurred during KML validation.
// It identifies the specific element and field that failed validation.
type ValidationError struct {
	Element string // The KML element that failed validation (e.g., "Placemark", "Point")
	Field   string // The specific field that failed validation (e.g., "coordinates", "name")
	Message string // Human-readable error message describing the validation failure
}

// Error returns a formatted error message identifying the element and field.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("kml: validation error in %s.%s: %s", e.Element, e.Field, e.Message)
	}
	return fmt.Sprintf("kml: validation error in %s: %s", e.Element, e.Message)
}

// Sentinel errors for common error conditions.
// These can be used with errors.Is for error checking.
var (
	// ErrInvalidCoordinate indicates that a coordinate string could not be parsed.
	// Coordinates must be in the format "longitude,latitude[,altitude]".
	ErrInvalidCoordinate = errors.New("kml: invalid coordinate format")

	// ErrInvalidColor indicates that a color string could not be parsed.
	// Colors must be in the format "aabbggrr" (hex format with alpha).
	ErrInvalidColor = errors.New("kml: invalid color format")

	// ErrEmptyDocument indicates that a KML document contains no features.
	// A valid KML document should contain at least one Placemark or other feature.
	ErrEmptyDocument = errors.New("kml: document contains no features")

	// ErrMissingGeometry indicates that a Placemark has no associated geometry.
	// Every Placemark should have at least one geometry element (Point, LineString, etc.).
	ErrMissingGeometry = errors.New("kml: placemark has no geometry")
)
