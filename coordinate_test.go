package kml

import (
	"encoding/xml"
	"errors"
	"testing"
)

// TestCoord tests the Coord() function
func TestCoord(t *testing.T) {
	tests := []struct {
		name string
		lon  float64
		lat  float64
		alt  []float64
		want Coordinate
	}{
		{
			name: "Coord with just lon, lat",
			lon:  -122.084,
			lat:  37.422,
			alt:  nil,
			want: Coordinate{Lon: -122.084, Lat: 37.422, Alt: 0},
		},
		{
			name: "Coord with lon, lat, alt",
			lon:  -122.084,
			lat:  37.422,
			alt:  []float64{100.5},
			want: Coordinate{Lon: -122.084, Lat: 37.422, Alt: 100.5},
		},
		{
			name: "Coord with zero altitude explicitly",
			lon:  10.0,
			lat:  20.0,
			alt:  []float64{0},
			want: Coordinate{Lon: 10.0, Lat: 20.0, Alt: 0},
		},
		{
			name: "Coord with negative values",
			lon:  -180.0,
			lat:  -90.0,
			alt:  []float64{-50.0},
			want: Coordinate{Lon: -180.0, Lat: -90.0, Alt: -50.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Coordinate
			if tt.alt != nil {
				got = Coord(tt.lon, tt.lat, tt.alt...)
			} else {
				got = Coord(tt.lon, tt.lat)
			}

			if got.Lon != tt.want.Lon {
				t.Errorf("Coord().Lon = %v, want %v", got.Lon, tt.want.Lon)
			}
			if got.Lat != tt.want.Lat {
				t.Errorf("Coord().Lat = %v, want %v", got.Lat, tt.want.Lat)
			}
			if got.Alt != tt.want.Alt {
				t.Errorf("Coord().Alt = %v, want %v", got.Alt, tt.want.Alt)
			}
		})
	}
}

// TestParseCoordinates tests the ParseCoordinates() function
func TestParseCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Coordinate
		wantErr bool
		errType error
	}{
		{
			name:  "Single coordinate without altitude",
			input: "1.0,2.0",
			want: []Coordinate{
				{Lon: 1.0, Lat: 2.0, Alt: 0},
			},
			wantErr: false,
		},
		{
			name:  "Single coordinate with altitude",
			input: "1.0,2.0,100",
			want: []Coordinate{
				{Lon: 1.0, Lat: 2.0, Alt: 100},
			},
			wantErr: false,
		},
		{
			name:  "Multiple coordinates",
			input: "1.0,2.0 3.0,4.0 5.0,6.0",
			want: []Coordinate{
				{Lon: 1.0, Lat: 2.0, Alt: 0},
				{Lon: 3.0, Lat: 4.0, Alt: 0},
				{Lon: 5.0, Lat: 6.0, Alt: 0},
			},
			wantErr: false,
		},
		{
			name:  "Multiple coordinates with mixed altitudes",
			input: "1.0,2.0,10 3.0,4.0 5.0,6.0,30",
			want: []Coordinate{
				{Lon: 1.0, Lat: 2.0, Alt: 10},
				{Lon: 3.0, Lat: 4.0, Alt: 0},
				{Lon: 5.0, Lat: 6.0, Alt: 30},
			},
			wantErr: false,
		},
		{
			name:  "With extra whitespace",
			input: "  1.0,2.0   3.0,4.0  ",
			want: []Coordinate{
				{Lon: 1.0, Lat: 2.0, Alt: 0},
				{Lon: 3.0, Lat: 4.0, Alt: 0},
			},
			wantErr: false,
		},
		{
			name:  "With tabs and newlines",
			input: "1.0,2.0\t\n3.0,4.0",
			want: []Coordinate{
				{Lon: 1.0, Lat: 2.0, Alt: 0},
				{Lon: 3.0, Lat: 4.0, Alt: 0},
			},
			wantErr: false,
		},
		{
			name:  "With negative values",
			input: "-122.084,37.422,100.5",
			want: []Coordinate{
				{Lon: -122.084, Lat: 37.422, Alt: 100.5},
			},
			wantErr: false,
		},
		{
			name:  "With scientific notation",
			input: "1.23e2,4.56e1",
			want: []Coordinate{
				{Lon: 123, Lat: 45.6, Alt: 0},
			},
			wantErr: false,
		},
		{
			name:    "Empty string returns error",
			input:   "",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "Only whitespace returns error",
			input:   "   \t\n   ",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "Invalid format - only one value",
			input:   "1.0",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "Invalid format - too many values",
			input:   "1.0,2.0,3.0,4.0",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "Invalid longitude",
			input:   "abc,2.0",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "Invalid latitude",
			input:   "1.0,def",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "Invalid altitude",
			input:   "1.0,2.0,xyz",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "Missing latitude",
			input:   "1.0,",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
		{
			name:    "One valid, one invalid coordinate",
			input:   "1.0,2.0 3.0,abc",
			want:    nil,
			wantErr: true,
			errType: ErrInvalidCoordinate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCoordinates(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCoordinates() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("ParseCoordinates() error = %v, want error type %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseCoordinates() unexpected error = %v", err)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ParseCoordinates() returned %d coordinates, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].Lon != tt.want[i].Lon {
					t.Errorf("ParseCoordinates()[%d].Lon = %v, want %v", i, got[i].Lon, tt.want[i].Lon)
				}
				if got[i].Lat != tt.want[i].Lat {
					t.Errorf("ParseCoordinates()[%d].Lat = %v, want %v", i, got[i].Lat, tt.want[i].Lat)
				}
				if got[i].Alt != tt.want[i].Alt {
					t.Errorf("ParseCoordinates()[%d].Alt = %v, want %v", i, got[i].Alt, tt.want[i].Alt)
				}
			}
		})
	}
}

// TestCoordinateString tests the String() method
func TestCoordinateString(t *testing.T) {
	tests := []struct {
		name  string
		coord Coordinate
		want  string
	}{
		{
			name:  "Coordinate without altitude",
			coord: Coordinate{Lon: 1.0, Lat: 2.0, Alt: 0},
			want:  "1,2",
		},
		{
			name:  "Coordinate with altitude",
			coord: Coordinate{Lon: 1.0, Lat: 2.0, Alt: 100},
			want:  "1,2,100",
		},
		{
			name:  "Coordinate with negative values",
			coord: Coordinate{Lon: -122.084, Lat: -37.422, Alt: 0},
			want:  "-122.084,-37.422",
		},
		{
			name:  "Coordinate with negative altitude",
			coord: Coordinate{Lon: 10.0, Lat: 20.0, Alt: -50.5},
			want:  "10,20,-50.5",
		},
		{
			name:  "Coordinate with floating point values",
			coord: Coordinate{Lon: -122.084, Lat: 37.422, Alt: 100.5},
			want:  "-122.084,37.422,100.5",
		},
		{
			name:  "Coordinate with zero lon/lat",
			coord: Coordinate{Lon: 0, Lat: 0, Alt: 0},
			want:  "0,0",
		},
		{
			name:  "Coordinate with zero lon/lat but with altitude",
			coord: Coordinate{Lon: 0, Lat: 0, Alt: 50},
			want:  "0,0,50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.coord.String()
			if got != tt.want {
				t.Errorf("Coordinate.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCoordinatesXML tests XML marshaling/unmarshaling of Coordinates type
func TestCoordinatesXML(t *testing.T) {
	t.Run("Marshal single coordinate", func(t *testing.T) {
		coords := Coordinates{
			{Lon: 1.0, Lat: 2.0, Alt: 0},
		}

		type Point struct {
			XMLName xml.Name    `xml:"Point"`
			Coords  Coordinates `xml:"coordinates"`
		}

		p := Point{Coords: coords}
		data, err := xml.Marshal(p)
		if err != nil {
			t.Fatalf("xml.Marshal() error = %v", err)
		}

		got := string(data)
		want := "<Point><coordinates>1,2</coordinates></Point>"
		if got != want {
			t.Errorf("xml.Marshal() = %v, want %v", got, want)
		}
	})

	t.Run("Marshal multiple coordinates", func(t *testing.T) {
		coords := Coordinates{
			{Lon: 1.0, Lat: 2.0, Alt: 0},
			{Lon: 3.0, Lat: 4.0, Alt: 0},
			{Lon: 5.0, Lat: 6.0, Alt: 100},
		}

		type LineString struct {
			XMLName xml.Name    `xml:"LineString"`
			Coords  Coordinates `xml:"coordinates"`
		}

		ls := LineString{Coords: coords}
		data, err := xml.Marshal(ls)
		if err != nil {
			t.Fatalf("xml.Marshal() error = %v", err)
		}

		got := string(data)
		want := "<LineString><coordinates>1,2 3,4 5,6,100</coordinates></LineString>"
		if got != want {
			t.Errorf("xml.Marshal() = %v, want %v", got, want)
		}
	})

	t.Run("Marshal coordinates with altitude", func(t *testing.T) {
		coords := Coordinates{
			{Lon: -122.084, Lat: 37.422, Alt: 100.5},
		}

		type Point struct {
			XMLName xml.Name    `xml:"Point"`
			Coords  Coordinates `xml:"coordinates"`
		}

		p := Point{Coords: coords}
		data, err := xml.Marshal(p)
		if err != nil {
			t.Fatalf("xml.Marshal() error = %v", err)
		}

		got := string(data)
		want := "<Point><coordinates>-122.084,37.422,100.5</coordinates></Point>"
		if got != want {
			t.Errorf("xml.Marshal() = %v, want %v", got, want)
		}
	})

	t.Run("Marshal empty coordinates", func(t *testing.T) {
		coords := Coordinates{}

		type Point struct {
			XMLName xml.Name    `xml:"Point"`
			Coords  Coordinates `xml:"coordinates"`
		}

		p := Point{Coords: coords}
		data, err := xml.Marshal(p)
		if err != nil {
			t.Fatalf("xml.Marshal() error = %v", err)
		}

		got := string(data)
		want := "<Point><coordinates></coordinates></Point>"
		if got != want {
			t.Errorf("xml.Marshal() = %v, want %v", got, want)
		}
	})

	t.Run("Unmarshal single coordinate", func(t *testing.T) {
		input := "<Point><coordinates>1.0,2.0</coordinates></Point>"

		type Point struct {
			XMLName xml.Name    `xml:"Point"`
			Coords  Coordinates `xml:"coordinates"`
		}

		var p Point
		err := xml.Unmarshal([]byte(input), &p)
		if err != nil {
			t.Fatalf("xml.Unmarshal() error = %v", err)
		}

		if len(p.Coords) != 1 {
			t.Errorf("Unmarshal returned %d coordinates, want 1", len(p.Coords))
			return
		}

		want := Coordinate{Lon: 1.0, Lat: 2.0, Alt: 0}
		if p.Coords[0] != want {
			t.Errorf("Unmarshal coordinate = %v, want %v", p.Coords[0], want)
		}
	})

	t.Run("Unmarshal multiple coordinates", func(t *testing.T) {
		input := "<LineString><coordinates>1.0,2.0 3.0,4.0 5.0,6.0,100</coordinates></LineString>"

		type LineString struct {
			XMLName xml.Name    `xml:"LineString"`
			Coords  Coordinates `xml:"coordinates"`
		}

		var ls LineString
		err := xml.Unmarshal([]byte(input), &ls)
		if err != nil {
			t.Fatalf("xml.Unmarshal() error = %v", err)
		}

		want := []Coordinate{
			{Lon: 1.0, Lat: 2.0, Alt: 0},
			{Lon: 3.0, Lat: 4.0, Alt: 0},
			{Lon: 5.0, Lat: 6.0, Alt: 100},
		}

		if len(ls.Coords) != len(want) {
			t.Errorf("Unmarshal returned %d coordinates, want %d", len(ls.Coords), len(want))
			return
		}

		for i := range ls.Coords {
			if ls.Coords[i] != want[i] {
				t.Errorf("Unmarshal coordinate[%d] = %v, want %v", i, ls.Coords[i], want[i])
			}
		}
	})

	t.Run("Unmarshal with whitespace", func(t *testing.T) {
		input := "<LineString><coordinates>  1.0,2.0   3.0,4.0  </coordinates></LineString>"

		type LineString struct {
			XMLName xml.Name    `xml:"LineString"`
			Coords  Coordinates `xml:"coordinates"`
		}

		var ls LineString
		err := xml.Unmarshal([]byte(input), &ls)
		if err != nil {
			t.Fatalf("xml.Unmarshal() error = %v", err)
		}

		want := []Coordinate{
			{Lon: 1.0, Lat: 2.0, Alt: 0},
			{Lon: 3.0, Lat: 4.0, Alt: 0},
		}

		if len(ls.Coords) != len(want) {
			t.Errorf("Unmarshal returned %d coordinates, want %d", len(ls.Coords), len(want))
			return
		}

		for i := range ls.Coords {
			if ls.Coords[i] != want[i] {
				t.Errorf("Unmarshal coordinate[%d] = %v, want %v", i, ls.Coords[i], want[i])
			}
		}
	})

	t.Run("Unmarshal invalid format returns error", func(t *testing.T) {
		input := "<Point><coordinates>invalid</coordinates></Point>"

		type Point struct {
			XMLName xml.Name    `xml:"Point"`
			Coords  Coordinates `xml:"coordinates"`
		}

		var p Point
		err := xml.Unmarshal([]byte(input), &p)
		if err == nil {
			t.Error("xml.Unmarshal() error = nil, want error")
			return
		}

		if !errors.Is(err, ErrInvalidCoordinate) {
			t.Errorf("xml.Unmarshal() error = %v, want ErrInvalidCoordinate", err)
		}
	})

	t.Run("Round-trip marshal/unmarshal", func(t *testing.T) {
		original := Coordinates{
			{Lon: -122.084, Lat: 37.422, Alt: 100.5},
			{Lon: -122.085, Lat: 37.423, Alt: 0},
			{Lon: -122.086, Lat: 37.424, Alt: 200},
		}

		type LineString struct {
			XMLName xml.Name    `xml:"LineString"`
			Coords  Coordinates `xml:"coordinates"`
		}

		// Marshal
		ls1 := LineString{Coords: original}
		data, err := xml.Marshal(ls1)
		if err != nil {
			t.Fatalf("xml.Marshal() error = %v", err)
		}

		// Unmarshal
		var ls2 LineString
		err = xml.Unmarshal(data, &ls2)
		if err != nil {
			t.Fatalf("xml.Unmarshal() error = %v", err)
		}

		// Compare
		if len(ls2.Coords) != len(original) {
			t.Errorf("Round-trip returned %d coordinates, want %d", len(ls2.Coords), len(original))
			return
		}

		for i := range ls2.Coords {
			if ls2.Coords[i] != original[i] {
				t.Errorf("Round-trip coordinate[%d] = %v, want %v", i, ls2.Coords[i], original[i])
			}
		}
	})
}
