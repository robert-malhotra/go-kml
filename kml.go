package kml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

const (
	// DefaultNamespace is the standard KML 2.2 namespace.
	DefaultNamespace = "http://www.opengis.net/kml/2.2"

	// XMLHeader is the standard XML declaration for KML files.
	XMLHeader = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// KML represents the root element of a KML document.
// The Feature field can contain a Document, Folder, or Placemark.
type KML struct {
	XMLName xml.Name `xml:"kml"`
	Xmlns   string   `xml:"xmlns,attr"`
	Feature Feature  `xml:"-"` // Document, Folder, or Placemark - custom marshaling
}

// NewKML creates a new empty KML document with default namespace.
func NewKML() *KML {
	return &KML{
		Xmlns: DefaultNamespace,
	}
}

// MarshalXML implements custom XML marshaling for KML.
// It writes the kml element with xmlns attribute and the Feature child.
func (k *KML) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "kml"

	// Add xmlns attribute
	xmlns := k.Xmlns
	if xmlns == "" {
		xmlns = DefaultNamespace
	}
	start.Attr = append(start.Attr, xml.Attr{
		Name:  xml.Name{Local: "xmlns"},
		Value: xmlns,
	})

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Encode the Feature based on its concrete type
	if k.Feature != nil {
		switch f := k.Feature.(type) {
		case *Document:
			if err := e.EncodeElement(f, xml.StartElement{Name: xml.Name{Local: "Document"}}); err != nil {
				return err
			}
		case *Folder:
			if err := e.EncodeElement(f, xml.StartElement{Name: xml.Name{Local: "Folder"}}); err != nil {
				return err
			}
		case *Placemark:
			if err := e.EncodeElement(f, xml.StartElement{Name: xml.Name{Local: "Placemark"}}); err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalXML implements custom XML unmarshaling for KML.
// It reads the feature child (Document, Folder, or Placemark).
func (k *KML) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Process attributes
	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" {
			k.Xmlns = attr.Value
		}
	}

	// Set default namespace if not specified
	if k.Xmlns == "" {
		k.Xmlns = DefaultNamespace
	}

	// Process child elements
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return &ParseError{
				Message: "error reading KML document",
				Cause:   err,
			}
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "Document":
				var doc Document
				if err := d.DecodeElement(&doc, &tok); err != nil {
					return &ParseError{
						Message: "error parsing Document element",
						Cause:   err,
					}
				}
				k.Feature = &doc
			case "Folder":
				var folder Folder
				if err := d.DecodeElement(&folder, &tok); err != nil {
					return &ParseError{
						Message: "error parsing Folder element",
						Cause:   err,
					}
				}
				k.Feature = &folder
			case "Placemark":
				var placemark Placemark
				if err := d.DecodeElement(&placemark, &tok); err != nil {
					return &ParseError{
						Message: "error parsing Placemark element",
						Cause:   err,
					}
				}
				k.Feature = &placemark
			default:
				// Skip unknown elements
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tok.Name.Local == "kml" {
				return nil
			}
		}
	}

	return nil
}

// Parse reads a KML document from an io.Reader.
func Parse(r io.Reader) (*KML, error) {
	decoder := xml.NewDecoder(r)

	var k KML
	if err := decoder.Decode(&k); err != nil {
		if parseErr, ok := err.(*ParseError); ok {
			return nil, parseErr
		}
		return nil, &ParseError{
			Message: "error decoding KML document",
			Cause:   err,
		}
	}

	// Validate that the document contains a feature
	if k.Feature == nil {
		return nil, ErrEmptyDocument
	}

	return &k, nil
}

// ParseFile reads a KML document from a file path.
func ParseFile(path string) (*KML, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("kml: error opening file: %w", err)
	}
	defer f.Close()

	k, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("kml: error parsing file %s: %w", path, err)
	}

	return k, nil
}

// ParseBytes reads a KML document from a byte slice.
func ParseBytes(data []byte) (*KML, error) {
	return Parse(bytes.NewReader(data))
}

// Write writes a KML document to an io.Writer.
// It outputs the XML declaration before the KML content.
func (k *KML) Write(w io.Writer) error {
	// Write XML declaration
	if _, err := io.WriteString(w, XMLHeader); err != nil {
		return fmt.Errorf("kml: error writing XML header: %w", err)
	}

	encoder := xml.NewEncoder(w)
	if err := encoder.Encode(k); err != nil {
		return fmt.Errorf("kml: error encoding KML document: %w", err)
	}

	// Write final newline
	if _, err := io.WriteString(w, "\n"); err != nil {
		return fmt.Errorf("kml: error writing final newline: %w", err)
	}

	return nil
}

// WriteIndent writes a KML document with indentation.
// The prefix is written at the beginning of each line, and indent
// specifies the indentation string for each level.
func (k *KML) WriteIndent(w io.Writer, prefix, indent string) error {
	// Write XML declaration
	if _, err := io.WriteString(w, XMLHeader); err != nil {
		return fmt.Errorf("kml: error writing XML header: %w", err)
	}

	encoder := xml.NewEncoder(w)
	encoder.Indent(prefix, indent)

	if err := encoder.Encode(k); err != nil {
		return fmt.Errorf("kml: error encoding KML document: %w", err)
	}

	// Write final newline
	if _, err := io.WriteString(w, "\n"); err != nil {
		return fmt.Errorf("kml: error writing final newline: %w", err)
	}

	return nil
}

// WriteFile writes a KML document to a file.
// The file is created with permissions 0644.
func (k *KML) WriteFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("kml: error creating file: %w", err)
	}
	defer f.Close()

	if err := k.WriteIndent(f, "", "  "); err != nil {
		return fmt.Errorf("kml: error writing to file %s: %w", path, err)
	}

	return nil
}

// Bytes returns the KML document as a byte slice.
func (k *KML) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := k.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
