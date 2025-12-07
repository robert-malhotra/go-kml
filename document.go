package kml

import (
	"encoding/xml"
	"io"
)

// Feature interface represents anything that can be contained in a Document/Folder
type Feature interface {
	featureType() string
}

// Document represents a KML Document element
type Document struct {
	ID          string     `xml:"id,attr,omitempty"`
	Name        string     `xml:"name,omitempty"`
	Description string     `xml:"description,omitempty"`
	Open        bool       `xml:"open,omitempty"`
	Visibility  *bool      `xml:"visibility,omitempty"`
	Styles      []Style    `xml:"Style,omitempty"`
	StyleMaps   []StyleMap `xml:"StyleMap,omitempty"`
	Features    []Feature  `xml:"-"` // Custom unmarshaling required
}

// featureType implements the Feature interface
func (d *Document) featureType() string {
	return "Document"
}

// MarshalXML implements custom XML marshaling for Document
func (d *Document) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Document"

	if d.ID != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "id"}, Value: d.ID})
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Encode basic fields
	if d.Name != "" {
		if err := e.EncodeElement(d.Name, xml.StartElement{Name: xml.Name{Local: "name"}}); err != nil {
			return err
		}
	}

	if d.Description != "" {
		if err := e.EncodeElement(d.Description, xml.StartElement{Name: xml.Name{Local: "description"}}); err != nil {
			return err
		}
	}

	if d.Open {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "open"}}); err != nil {
			return err
		}
	}

	if d.Visibility != nil {
		vis := 0
		if *d.Visibility {
			vis = 1
		}
		if err := e.EncodeElement(vis, xml.StartElement{Name: xml.Name{Local: "visibility"}}); err != nil {
			return err
		}
	}

	// Encode Styles
	for _, style := range d.Styles {
		if err := e.EncodeElement(&style, xml.StartElement{Name: xml.Name{Local: "Style"}}); err != nil {
			return err
		}
	}

	// Encode StyleMaps
	for _, styleMap := range d.StyleMaps {
		if err := e.EncodeElement(&styleMap, xml.StartElement{Name: xml.Name{Local: "StyleMap"}}); err != nil {
			return err
		}
	}

	// Encode Features based on their concrete type
	for _, feature := range d.Features {
		switch f := feature.(type) {
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

// UnmarshalXML implements custom XML unmarshaling for Document
func (d *Document) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// Process attributes
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			d.ID = attr.Value
		}
	}

	// Process child elements
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "name":
				if err := decoder.DecodeElement(&d.Name, &tok); err != nil {
					return err
				}
			case "description":
				if err := decoder.DecodeElement(&d.Description, &tok); err != nil {
					return err
				}
			case "open":
				var open int
				if err := decoder.DecodeElement(&open, &tok); err != nil {
					return err
				}
				d.Open = open != 0
			case "visibility":
				var vis int
				if err := decoder.DecodeElement(&vis, &tok); err != nil {
					return err
				}
				visibility := vis != 0
				d.Visibility = &visibility
			case "Style":
				var style Style
				if err := decoder.DecodeElement(&style, &tok); err != nil {
					return err
				}
				d.Styles = append(d.Styles, style)
			case "StyleMap":
				var styleMap StyleMap
				if err := decoder.DecodeElement(&styleMap, &tok); err != nil {
					return err
				}
				d.StyleMaps = append(d.StyleMaps, styleMap)
			case "Document":
				var doc Document
				if err := decoder.DecodeElement(&doc, &tok); err != nil {
					return err
				}
				d.Features = append(d.Features, &doc)
			case "Folder":
				var folder Folder
				if err := decoder.DecodeElement(&folder, &tok); err != nil {
					return err
				}
				d.Features = append(d.Features, &folder)
			case "Placemark":
				var placemark Placemark
				if err := decoder.DecodeElement(&placemark, &tok); err != nil {
					return err
				}
				d.Features = append(d.Features, &placemark)
			default:
				// Skip unknown elements
				if err := decoder.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tok.Name.Local == "Document" {
				return nil
			}
		}
	}

	return nil
}

// Folder represents a KML Folder element
type Folder struct {
	ID          string    `xml:"id,attr,omitempty"`
	Name        string    `xml:"name,omitempty"`
	Description string    `xml:"description,omitempty"`
	Open        bool      `xml:"open,omitempty"`
	Visibility  *bool     `xml:"visibility,omitempty"`
	Features    []Feature `xml:"-"`
}

// featureType implements the Feature interface
func (f *Folder) featureType() string {
	return "Folder"
}

// MarshalXML implements custom XML marshaling for Folder
func (f *Folder) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Folder"

	if f.ID != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "id"}, Value: f.ID})
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Encode basic fields
	if f.Name != "" {
		if err := e.EncodeElement(f.Name, xml.StartElement{Name: xml.Name{Local: "name"}}); err != nil {
			return err
		}
	}

	if f.Description != "" {
		if err := e.EncodeElement(f.Description, xml.StartElement{Name: xml.Name{Local: "description"}}); err != nil {
			return err
		}
	}

	if f.Open {
		if err := e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "open"}}); err != nil {
			return err
		}
	}

	if f.Visibility != nil {
		vis := 0
		if *f.Visibility {
			vis = 1
		}
		if err := e.EncodeElement(vis, xml.StartElement{Name: xml.Name{Local: "visibility"}}); err != nil {
			return err
		}
	}

	// Encode Features based on their concrete type
	for _, feature := range f.Features {
		switch feat := feature.(type) {
		case *Document:
			if err := e.EncodeElement(feat, xml.StartElement{Name: xml.Name{Local: "Document"}}); err != nil {
				return err
			}
		case *Folder:
			if err := e.EncodeElement(feat, xml.StartElement{Name: xml.Name{Local: "Folder"}}); err != nil {
				return err
			}
		case *Placemark:
			if err := e.EncodeElement(feat, xml.StartElement{Name: xml.Name{Local: "Placemark"}}); err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalXML implements custom XML unmarshaling for Folder
func (f *Folder) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// Process attributes
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			f.ID = attr.Value
		}
	}

	// Process child elements
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "name":
				if err := decoder.DecodeElement(&f.Name, &tok); err != nil {
					return err
				}
			case "description":
				if err := decoder.DecodeElement(&f.Description, &tok); err != nil {
					return err
				}
			case "open":
				var open int
				if err := decoder.DecodeElement(&open, &tok); err != nil {
					return err
				}
				f.Open = open != 0
			case "visibility":
				var vis int
				if err := decoder.DecodeElement(&vis, &tok); err != nil {
					return err
				}
				visibility := vis != 0
				f.Visibility = &visibility
			case "Document":
				var doc Document
				if err := decoder.DecodeElement(&doc, &tok); err != nil {
					return err
				}
				f.Features = append(f.Features, &doc)
			case "Folder":
				var folder Folder
				if err := decoder.DecodeElement(&folder, &tok); err != nil {
					return err
				}
				f.Features = append(f.Features, &folder)
			case "Placemark":
				var placemark Placemark
				if err := decoder.DecodeElement(&placemark, &tok); err != nil {
					return err
				}
				f.Features = append(f.Features, &placemark)
			default:
				// Skip unknown elements
				if err := decoder.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tok.Name.Local == "Folder" {
				return nil
			}
		}
	}

	return nil
}
