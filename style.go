package kml

// Style defines the appearance of features.
type Style struct {
	ID           string        `xml:"id,attr,omitempty"`
	IconStyle    *IconStyle    `xml:"IconStyle,omitempty"`
	LabelStyle   *LabelStyle   `xml:"LabelStyle,omitempty"`
	LineStyle    *LineStyle    `xml:"LineStyle,omitempty"`
	PolyStyle    *PolyStyle    `xml:"PolyStyle,omitempty"`
	BalloonStyle *BalloonStyle `xml:"BalloonStyle,omitempty"`
}

// IconStyle specifies how icons are drawn.
type IconStyle struct {
	Color   Color    `xml:"color,omitempty"`
	Scale   float64  `xml:"scale,omitempty"`
	Heading float64  `xml:"heading,omitempty"`
	Icon    *Icon    `xml:"Icon,omitempty"`
	HotSpot *HotSpot `xml:"hotSpot,omitempty"`
}

// Icon specifies the icon image resource.
type Icon struct {
	Href string `xml:"href,omitempty"`
}

// HotSpot specifies the position within the icon that is anchored to the point.
type HotSpot struct {
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	XUnits string  `xml:"xunits,attr,omitempty"`
	YUnits string  `xml:"yunits,attr,omitempty"`
}

// LabelStyle specifies how labels are drawn.
type LabelStyle struct {
	Color Color   `xml:"color,omitempty"`
	Scale float64 `xml:"scale,omitempty"`
}

// LineStyle specifies how lines are drawn.
type LineStyle struct {
	Color Color   `xml:"color,omitempty"`
	Width float64 `xml:"width,omitempty"`
}

// PolyStyle specifies how polygons are drawn.
type PolyStyle struct {
	Color   Color `xml:"color,omitempty"`
	Fill    *bool `xml:"fill,omitempty"`
	Outline *bool `xml:"outline,omitempty"`
}

// BalloonStyle specifies how description balloons are displayed.
type BalloonStyle struct {
	BgColor   Color  `xml:"bgColor,omitempty"`
	TextColor Color  `xml:"textColor,omitempty"`
	Text      string `xml:"text,omitempty"`
}

// StyleMap maps between normal and highlight styles.
type StyleMap struct {
	ID    string `xml:"id,attr,omitempty"`
	Pairs []Pair `xml:"Pair"`
}

// Pair associates a style state (normal or highlight) with a style URL.
type Pair struct {
	Key      string `xml:"key"`      // "normal" or "highlight"
	StyleURL string `xml:"styleUrl"`
}
