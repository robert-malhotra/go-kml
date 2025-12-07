package kml

import (
	"encoding/xml"
)

// Placemark represents a geographic feature with geometry.
// It implements the Feature interface.
type Placemark struct {
	ID           string        `xml:"id,attr,omitempty"`
	Name         string        `xml:"name,omitempty"`
	Description  string        `xml:"description,omitempty"`
	Visibility   *bool         `xml:"visibility,omitempty"`
	StyleURL     string        `xml:"styleUrl,omitempty"`
	Style        *Style        `xml:"Style,omitempty"`
	Geometry     Geometry      `xml:"-"` // Point, LineString, Polygon, etc. - needs custom XML
	ExtendedData *ExtendedData `xml:"ExtendedData,omitempty"`
}

// featureType implements the Feature interface.
func (p *Placemark) featureType() string {
	return "Placemark"
}

// MarshalXML implements custom XML marshaling for Placemark.
// This handles the polymorphic Geometry field by type-asserting and
// writing the appropriate element (Point, LineString, Polygon, MultiGeometry).
func (p *Placemark) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Placemark"

	// Handle id attribute
	if p.ID != "" {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "id"},
			Value: p.ID,
		})
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Encode simple fields
	if p.Name != "" {
		if err := e.EncodeElement(p.Name, xml.StartElement{Name: xml.Name{Local: "name"}}); err != nil {
			return err
		}
	}

	if p.Description != "" {
		if err := e.EncodeElement(p.Description, xml.StartElement{Name: xml.Name{Local: "description"}}); err != nil {
			return err
		}
	}

	if p.Visibility != nil {
		v := 0
		if *p.Visibility {
			v = 1
		}
		if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "visibility"}}); err != nil {
			return err
		}
	}

	if p.StyleURL != "" {
		if err := e.EncodeElement(p.StyleURL, xml.StartElement{Name: xml.Name{Local: "styleUrl"}}); err != nil {
			return err
		}
	}

	if p.Style != nil {
		if err := e.Encode(p.Style); err != nil {
			return err
		}
	}

	// Encode Geometry - type assert to determine the concrete type
	if p.Geometry != nil {
		if err := e.Encode(p.Geometry); err != nil {
			return err
		}
	}

	if p.ExtendedData != nil {
		if err := e.Encode(p.ExtendedData); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML implements custom XML unmarshaling for Placemark.
// This handles reading polymorphic geometry elements and assigning them
// to the Geometry field.
func (p *Placemark) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Handle id attribute
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			p.ID = attr.Value
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "name":
				if err := d.DecodeElement(&p.Name, &el); err != nil {
					return err
				}
			case "description":
				if err := d.DecodeElement(&p.Description, &el); err != nil {
					return err
				}
			case "visibility":
				var v int
				if err := d.DecodeElement(&v, &el); err != nil {
					return err
				}
				vis := v != 0
				p.Visibility = &vis
			case "styleUrl":
				if err := d.DecodeElement(&p.StyleURL, &el); err != nil {
					return err
				}
			case "Style":
				var style Style
				if err := d.DecodeElement(&style, &el); err != nil {
					return err
				}
				p.Style = &style
			case "Point":
				var point Point
				if err := d.DecodeElement(&point, &el); err != nil {
					return err
				}
				p.Geometry = &point
			case "LineString":
				var lineString LineString
				if err := d.DecodeElement(&lineString, &el); err != nil {
					return err
				}
				p.Geometry = &lineString
			case "LinearRing":
				var linearRing LinearRing
				if err := d.DecodeElement(&linearRing, &el); err != nil {
					return err
				}
				p.Geometry = &linearRing
			case "Polygon":
				var polygon Polygon
				if err := d.DecodeElement(&polygon, &el); err != nil {
					return err
				}
				p.Geometry = &polygon
			case "MultiGeometry":
				var multiGeometry MultiGeometry
				if err := d.DecodeElement(&multiGeometry, &el); err != nil {
					return err
				}
				p.Geometry = &multiGeometry
			case "ExtendedData":
				var extendedData ExtendedData
				if err := d.DecodeElement(&extendedData, &el); err != nil {
					return err
				}
				p.ExtendedData = &extendedData
			default:
				// Skip unknown elements
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// ExtendedData allows you to add custom data to a KML feature.
// It supports two ways of adding data: Data elements and SchemaData elements.
type ExtendedData struct {
	Data       []Data       `xml:"Data,omitempty"`
	SchemaData []SchemaData `xml:"SchemaData,omitempty"`
}

// Data represents a single name-value pair with an optional display name.
type Data struct {
	Name        string `xml:"name,attr"`
	DisplayName string `xml:"displayName,omitempty"`
	Value       string `xml:"value"`
}

// SchemaData associates structured data with a schema definition.
type SchemaData struct {
	SchemaURL  string       `xml:"schemaUrl,attr,omitempty"`
	SimpleData []SimpleData `xml:"SimpleData,omitempty"`
}

// SimpleData represents a single name-value pair within a SchemaData element.
type SimpleData struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
