package kml

// NewKMLBuilder creates a new KMLBuilder for fluent construction of KML documents.
func NewKMLBuilder() *KMLBuilder {
	return &KMLBuilder{
		kml: NewKML(),
	}
}

// KMLBuilder provides a fluent API for building KML documents.
type KMLBuilder struct {
	kml *KML
}

// Document creates and adds a Document to the KML and returns a DocumentBuilder.
func (kb *KMLBuilder) Document(name string) *DocumentBuilder {
	doc := &Document{
		Name: name,
	}
	kb.kml.Feature = doc
	return &DocumentBuilder{
		kml:      kb.kml,
		document: doc,
		parent:   kb,
	}
}

// Build returns the constructed KML document.
func (kb *KMLBuilder) Build() *KML {
	return kb.kml
}

// DocumentBuilder provides a fluent API for building Document elements.
type DocumentBuilder struct {
	kml      *KML
	document *Document
	parent   *KMLBuilder
}

// Name sets the name of the document.
func (db *DocumentBuilder) Name(name string) *DocumentBuilder {
	db.document.Name = name
	return db
}

// Description sets the description of the document.
func (db *DocumentBuilder) Description(desc string) *DocumentBuilder {
	db.document.Description = desc
	return db
}

// Open sets the open state of the document.
func (db *DocumentBuilder) Open(open bool) *DocumentBuilder {
	db.document.Open = open
	return db
}

// Style creates a style with the given ID and returns a StyleBuilder.
func (db *DocumentBuilder) Style(id string) *StyleBuilder {
	style := Style{
		ID: id,
	}
	db.document.Styles = append(db.document.Styles, style)
	// Get pointer to the newly added style
	stylePtr := &db.document.Styles[len(db.document.Styles)-1]
	return &StyleBuilder{
		style:  stylePtr,
		parent: db,
	}
}

// Folder creates a folder with the given name and returns a FolderBuilder.
func (db *DocumentBuilder) Folder(name string) *FolderBuilder {
	folder := &Folder{
		Name: name,
	}
	db.document.Features = append(db.document.Features, folder)
	return &FolderBuilder{
		folder: folder,
		finisher: func() interface{} {
			return db
		},
	}
}

// Placemark creates a placemark with the given name and returns a PlacemarkBuilder.
func (db *DocumentBuilder) Placemark(name string) *PlacemarkBuilder {
	placemark := &Placemark{
		Name: name,
	}
	db.document.Features = append(db.document.Features, placemark)
	return &PlacemarkBuilder{
		placemark: placemark,
		finisher: func() interface{} {
			return db
		},
	}
}

// Done returns to the parent KMLBuilder.
func (db *DocumentBuilder) Done() *KMLBuilder {
	return db.parent
}

// Build is a shortcut to build the KML directly from the DocumentBuilder.
func (db *DocumentBuilder) Build() *KML {
	return db.kml
}

// FolderBuilder provides a fluent API for building Folder elements.
type FolderBuilder struct {
	folder   *Folder
	finisher func() interface{} // Returns to parent (DocumentBuilder or FolderBuilder)
}

// Name sets the name of the folder.
func (fb *FolderBuilder) Name(name string) *FolderBuilder {
	fb.folder.Name = name
	return fb
}

// Description sets the description of the folder.
func (fb *FolderBuilder) Description(desc string) *FolderBuilder {
	fb.folder.Description = desc
	return fb
}

// Open sets the open state of the folder.
func (fb *FolderBuilder) Open(open bool) *FolderBuilder {
	fb.folder.Open = open
	return fb
}

// Folder creates a nested folder with the given name and returns a new FolderBuilder.
func (fb *FolderBuilder) Folder(name string) *FolderBuilder {
	nested := &Folder{
		Name: name,
	}
	fb.folder.Features = append(fb.folder.Features, nested)
	return &FolderBuilder{
		folder: nested,
		finisher: func() interface{} {
			return fb
		},
	}
}

// Placemark creates a placemark with the given name and returns a PlacemarkBuilder.
func (fb *FolderBuilder) Placemark(name string) *PlacemarkBuilder {
	placemark := &Placemark{
		Name: name,
	}
	fb.folder.Features = append(fb.folder.Features, placemark)
	return &PlacemarkBuilder{
		placemark: placemark,
		finisher: func() interface{} {
			return fb
		},
	}
}

// Done returns to the parent builder (DocumentBuilder or FolderBuilder).
func (fb *FolderBuilder) Done() interface{} {
	return fb.finisher()
}

// PlacemarkBuilder provides a fluent API for building Placemark elements.
type PlacemarkBuilder struct {
	placemark *Placemark
	finisher  func() interface{} // Returns to parent (DocumentBuilder or FolderBuilder)
}

// Name sets the name of the placemark.
func (pb *PlacemarkBuilder) Name(name string) *PlacemarkBuilder {
	pb.placemark.Name = name
	return pb
}

// Description sets the description of the placemark.
func (pb *PlacemarkBuilder) Description(desc string) *PlacemarkBuilder {
	pb.placemark.Description = desc
	return pb
}

// StyleURL sets the style URL reference for the placemark.
func (pb *PlacemarkBuilder) StyleURL(url string) *PlacemarkBuilder {
	pb.placemark.StyleURL = url
	return pb
}

// Point creates a Point geometry at the specified coordinates.
// Takes longitude, latitude, and optional altitude values.
func (pb *PlacemarkBuilder) Point(lon, lat float64, alt ...float64) *PlacemarkBuilder {
	pb.placemark.Geometry = &Point{
		Coordinates: Coord(lon, lat, alt...),
	}
	return pb
}

// LineString creates a LineString geometry from the provided coordinates.
func (pb *PlacemarkBuilder) LineString(coords ...Coordinate) *PlacemarkBuilder {
	pb.placemark.Geometry = &LineString{
		Coordinates: coords,
	}
	return pb
}

// Polygon creates a Polygon geometry with the specified outer boundary and optional inner boundaries (holes).
func (pb *PlacemarkBuilder) Polygon(outer []Coordinate, inner ...[]Coordinate) *PlacemarkBuilder {
	polygon := &Polygon{
		OuterBoundary: LinearRing{
			Coordinates: outer,
		},
	}

	for _, innerRing := range inner {
		polygon.InnerBoundaries = append(polygon.InnerBoundaries, LinearRing{
			Coordinates: innerRing,
		})
	}

	pb.placemark.Geometry = polygon
	return pb
}

// Done returns to the parent builder (DocumentBuilder or FolderBuilder).
func (pb *PlacemarkBuilder) Done() interface{} {
	return pb.finisher()
}

// StyleBuilder provides a fluent API for building Style elements.
type StyleBuilder struct {
	style  *Style
	parent *DocumentBuilder
}

// IconStyle starts building an IconStyle and returns an IconStyleBuilder.
func (sb *StyleBuilder) IconStyle() *IconStyleBuilder {
	if sb.style.IconStyle == nil {
		sb.style.IconStyle = &IconStyle{}
	}
	return &IconStyleBuilder{
		iconStyle: sb.style.IconStyle,
		parent:    sb,
	}
}

// LineStyle starts building a LineStyle and returns a LineStyleBuilder.
func (sb *StyleBuilder) LineStyle() *LineStyleBuilder {
	if sb.style.LineStyle == nil {
		sb.style.LineStyle = &LineStyle{}
	}
	return &LineStyleBuilder{
		lineStyle: sb.style.LineStyle,
		parent:    sb,
	}
}

// PolyStyle starts building a PolyStyle and returns a PolyStyleBuilder.
func (sb *StyleBuilder) PolyStyle() *PolyStyleBuilder {
	if sb.style.PolyStyle == nil {
		sb.style.PolyStyle = &PolyStyle{}
	}
	return &PolyStyleBuilder{
		polyStyle: sb.style.PolyStyle,
		parent:    sb,
	}
}

// LabelStyle starts building a LabelStyle and returns a LabelStyleBuilder.
func (sb *StyleBuilder) LabelStyle() *LabelStyleBuilder {
	if sb.style.LabelStyle == nil {
		sb.style.LabelStyle = &LabelStyle{}
	}
	return &LabelStyleBuilder{
		labelStyle: sb.style.LabelStyle,
		parent:     sb,
	}
}

// Done returns to the parent DocumentBuilder.
func (sb *StyleBuilder) Done() *DocumentBuilder {
	return sb.parent
}

// IconStyleBuilder provides a fluent API for building IconStyle elements.
type IconStyleBuilder struct {
	iconStyle *IconStyle
	parent    *StyleBuilder
}

// Color sets the color of the icon.
func (isb *IconStyleBuilder) Color(color Color) *IconStyleBuilder {
	isb.iconStyle.Color = color
	return isb
}

// Scale sets the scale of the icon.
func (isb *IconStyleBuilder) Scale(scale float64) *IconStyleBuilder {
	isb.iconStyle.Scale = scale
	return isb
}

// Heading sets the heading/rotation of the icon.
func (isb *IconStyleBuilder) Heading(heading float64) *IconStyleBuilder {
	isb.iconStyle.Heading = heading
	return isb
}

// Icon sets the icon URL.
func (isb *IconStyleBuilder) Icon(href string) *IconStyleBuilder {
	isb.iconStyle.Icon = &Icon{
		Href: href,
	}
	return isb
}

// HotSpot sets the anchor point for the icon.
func (isb *IconStyleBuilder) HotSpot(x, y float64, xUnits, yUnits string) *IconStyleBuilder {
	isb.iconStyle.HotSpot = &HotSpot{
		X:      x,
		Y:      y,
		XUnits: xUnits,
		YUnits: yUnits,
	}
	return isb
}

// Done returns to the parent StyleBuilder.
func (isb *IconStyleBuilder) Done() *StyleBuilder {
	return isb.parent
}

// LineStyleBuilder provides a fluent API for building LineStyle elements.
type LineStyleBuilder struct {
	lineStyle *LineStyle
	parent    *StyleBuilder
}

// Color sets the color of the line.
func (lsb *LineStyleBuilder) Color(color Color) *LineStyleBuilder {
	lsb.lineStyle.Color = color
	return lsb
}

// Width sets the width of the line.
func (lsb *LineStyleBuilder) Width(width float64) *LineStyleBuilder {
	lsb.lineStyle.Width = width
	return lsb
}

// Done returns to the parent StyleBuilder.
func (lsb *LineStyleBuilder) Done() *StyleBuilder {
	return lsb.parent
}

// PolyStyleBuilder provides a fluent API for building PolyStyle elements.
type PolyStyleBuilder struct {
	polyStyle *PolyStyle
	parent    *StyleBuilder
}

// Color sets the color of the polygon.
func (psb *PolyStyleBuilder) Color(color Color) *PolyStyleBuilder {
	psb.polyStyle.Color = color
	return psb
}

// Fill sets whether the polygon is filled.
func (psb *PolyStyleBuilder) Fill(fill bool) *PolyStyleBuilder {
	psb.polyStyle.Fill = &fill
	return psb
}

// Outline sets whether the polygon has an outline.
func (psb *PolyStyleBuilder) Outline(outline bool) *PolyStyleBuilder {
	psb.polyStyle.Outline = &outline
	return psb
}

// Done returns to the parent StyleBuilder.
func (psb *PolyStyleBuilder) Done() *StyleBuilder {
	return psb.parent
}

// LabelStyleBuilder provides a fluent API for building LabelStyle elements.
type LabelStyleBuilder struct {
	labelStyle *LabelStyle
	parent     *StyleBuilder
}

// Color sets the color of the label.
func (lsb *LabelStyleBuilder) Color(color Color) *LabelStyleBuilder {
	lsb.labelStyle.Color = color
	return lsb
}

// Scale sets the scale of the label.
func (lsb *LabelStyleBuilder) Scale(scale float64) *LabelStyleBuilder {
	lsb.labelStyle.Scale = scale
	return lsb
}

// Done returns to the parent StyleBuilder.
func (lsb *LabelStyleBuilder) Done() *StyleBuilder {
	return lsb.parent
}
