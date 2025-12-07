package kml

import (
	"encoding/xml"
	"testing"
)

func TestRGBA(t *testing.T) {
	tests := []struct {
		name string
		r, g, b, a uint8
		want Color
	}{
		{
			name: "red",
			r: 255, g: 0, b: 0, a: 255,
			want: Color{A: 255, B: 0, G: 0, R: 255},
		},
		{
			name: "green",
			r: 0, g: 255, b: 0, a: 255,
			want: Color{A: 255, B: 0, G: 255, R: 0},
		},
		{
			name: "blue",
			r: 0, g: 0, b: 255, a: 255,
			want: Color{A: 255, B: 255, G: 0, R: 0},
		},
		{
			name: "semi-transparent red",
			r: 255, g: 0, b: 0, a: 128,
			want: Color{A: 128, B: 0, G: 0, R: 255},
		},
		{
			name: "fully transparent white",
			r: 255, g: 255, b: 255, a: 0,
			want: Color{A: 0, B: 255, G: 255, R: 255},
		},
		{
			name: "semi-transparent blue",
			r: 0, g: 0, b: 255, a: 64,
			want: Color{A: 64, B: 255, G: 0, R: 0},
		},
		{
			name: "opaque yellow (red + green)",
			r: 255, g: 255, b: 0, a: 255,
			want: Color{A: 255, B: 0, G: 255, R: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RGBA(tt.r, tt.g, tt.b, tt.a)
			if got != tt.want {
				t.Errorf("RGBA(%d, %d, %d, %d) = %+v, want %+v",
					tt.r, tt.g, tt.b, tt.a, got, tt.want)
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Color
		wantErr bool
	}{
		{
			name:  "red (lowercase)",
			input: "ff0000ff",
			want:  Color{A: 255, B: 0, G: 0, R: 255},
		},
		{
			name:  "red (uppercase)",
			input: "FF0000FF",
			want:  Color{A: 255, B: 0, G: 0, R: 255},
		},
		{
			name:  "green",
			input: "ff00ff00",
			want:  Color{A: 255, B: 0, G: 255, R: 0},
		},
		{
			name:  "blue",
			input: "ffff0000",
			want:  Color{A: 255, B: 255, G: 0, R: 0},
		},
		{
			name:  "semi-transparent red",
			input: "800000ff",
			want:  Color{A: 128, B: 0, G: 0, R: 255},
		},
		{
			name:  "mixed case",
			input: "Ff00FF00",
			want:  Color{A: 255, B: 0, G: 255, R: 0},
		},
		{
			name:    "too short",
			input:   "ff00",
			wantErr: true,
		},
		{
			name:    "too long",
			input:   "ff0000ff00",
			wantErr: true,
		},
		{
			name:    "invalid hex characters",
			input:   "zz0000ff",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only 7 characters",
			input:   "ff0000f",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor(%q) error = %v, wantErr %v",
					tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseColor(%q) = %+v, want %+v",
					tt.input, got, tt.want)
			}
		})
	}
}

func TestColorHex(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{
			name:  "red",
			color: Color{A: 255, B: 0, G: 0, R: 255},
			want:  "ff0000ff",
		},
		{
			name:  "green",
			color: Color{A: 255, B: 0, G: 255, R: 0},
			want:  "ff00ff00",
		},
		{
			name:  "blue",
			color: Color{A: 255, B: 255, G: 0, R: 0},
			want:  "ffff0000",
		},
		{
			name:  "semi-transparent red",
			color: Color{A: 128, B: 0, G: 0, R: 255},
			want:  "800000ff",
		},
		{
			name:  "fully transparent",
			color: Color{A: 0, B: 0, G: 0, R: 0},
			want:  "00000000",
		},
		{
			name:  "white",
			color: Color{A: 255, B: 255, G: 255, R: 255},
			want:  "ffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.Hex()
			if got != tt.want {
				t.Errorf("Color.Hex() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColorRoundTrip(t *testing.T) {
	// Test that colors round-trip through Hex() and ParseColor()
	tests := []Color{
		{A: 255, B: 0, G: 0, R: 255},     // red
		{A: 255, B: 0, G: 255, R: 0},     // green
		{A: 255, B: 255, G: 0, R: 0},     // blue
		{A: 128, B: 64, G: 192, R: 32},   // arbitrary color
		{A: 0, B: 0, G: 0, R: 0},         // fully transparent black
		{A: 255, B: 255, G: 255, R: 255}, // white
	}

	for _, original := range tests {
		t.Run(original.Hex(), func(t *testing.T) {
			hex := original.Hex()
			parsed, err := ParseColor(hex)
			if err != nil {
				t.Errorf("ParseColor(%q) unexpected error: %v", hex, err)
				return
			}
			if parsed != original {
				t.Errorf("Round-trip failed: original %+v -> hex %q -> parsed %+v",
					original, hex, parsed)
			}
		})
	}
}

func TestPredefinedColors(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  Color
	}{
		{
			name:  "White",
			color: White,
			want:  Color{A: 255, B: 255, G: 255, R: 255},
		},
		{
			name:  "Black",
			color: Black,
			want:  Color{A: 255, B: 0, G: 0, R: 0},
		},
		{
			name:  "Red",
			color: Red,
			want:  Color{A: 255, B: 0, G: 0, R: 255},
		},
		{
			name:  "Green",
			color: Green,
			want:  Color{A: 255, B: 0, G: 255, R: 0},
		},
		{
			name:  "Blue",
			color: Blue,
			want:  Color{A: 255, B: 255, G: 0, R: 0},
		},
		{
			name:  "Transparent",
			color: Transparent,
			want:  Color{A: 0, B: 255, G: 255, R: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.want {
				t.Errorf("%s = %+v, want %+v", tt.name, tt.color, tt.want)
			}
		})
	}
}

func TestPredefinedColorsHex(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{
			name:  "White",
			color: White,
			want:  "ffffffff",
		},
		{
			name:  "Black",
			color: Black,
			want:  "ff000000",
		},
		{
			name:  "Red",
			color: Red,
			want:  "ff0000ff",
		},
		{
			name:  "Green",
			color: Green,
			want:  "ff00ff00",
		},
		{
			name:  "Blue",
			color: Blue,
			want:  "ffff0000",
		},
		{
			name:  "Transparent",
			color: Transparent,
			want:  "00ffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.Hex()
			if got != tt.want {
				t.Errorf("%s.Hex() = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestColorMarshalXML(t *testing.T) {
	tests := []struct {
		name    string
		color   Color
		want    string
		wantErr bool
	}{
		{
			name:  "red",
			color: Color{A: 255, B: 0, G: 0, R: 255},
			want:  "<wrapper><color>ff0000ff</color></wrapper>",
		},
		{
			name:  "green",
			color: Color{A: 255, B: 0, G: 255, R: 0},
			want:  "<wrapper><color>ff00ff00</color></wrapper>",
		},
		{
			name:  "blue",
			color: Color{A: 255, B: 255, G: 0, R: 0},
			want:  "<wrapper><color>ffff0000</color></wrapper>",
		},
		{
			name:  "semi-transparent",
			color: Color{A: 128, B: 64, G: 192, R: 32},
			want:  "<wrapper><color>8040c020</color></wrapper>",
		},
		{
			name:  "predefined White",
			color: White,
			want:  "<wrapper><color>ffffffff</color></wrapper>",
		},
		{
			name:  "predefined Black",
			color: Black,
			want:  "<wrapper><color>ff000000</color></wrapper>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple struct to wrap the color for marshaling
			type wrapper struct {
				Color Color `xml:"color"`
			}
			w := wrapper{Color: tt.color}

			got, err := xml.Marshal(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("MarshalXML() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestColorUnmarshalXML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Color
		wantErr bool
	}{
		{
			name:  "red",
			input: "<wrapper><color>ff0000ff</color></wrapper>",
			want:  Color{A: 255, B: 0, G: 0, R: 255},
		},
		{
			name:  "green",
			input: "<wrapper><color>ff00ff00</color></wrapper>",
			want:  Color{A: 255, B: 0, G: 255, R: 0},
		},
		{
			name:  "blue",
			input: "<wrapper><color>ffff0000</color></wrapper>",
			want:  Color{A: 255, B: 255, G: 0, R: 0},
		},
		{
			name:  "uppercase hex",
			input: "<wrapper><color>FF0000FF</color></wrapper>",
			want:  Color{A: 255, B: 0, G: 0, R: 255},
		},
		{
			name:  "semi-transparent",
			input: "<wrapper><color>8040c020</color></wrapper>",
			want:  Color{A: 128, B: 64, G: 192, R: 32},
		},
		{
			name:    "invalid hex",
			input:   "<wrapper><color>zz0000ff</color></wrapper>",
			wantErr: true,
		},
		{
			name:    "too short",
			input:   "<wrapper><color>ff00ff</color></wrapper>",
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "<wrapper><color></color></wrapper>",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type wrapper struct {
				Color Color `xml:"color"`
			}
			var w wrapper

			err := xml.Unmarshal([]byte(tt.input), &w)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && w.Color != tt.want {
				t.Errorf("UnmarshalXML() = %+v, want %+v", w.Color, tt.want)
			}
		})
	}
}

func TestColorMarshalUnmarshalRoundTrip(t *testing.T) {
	// Test that colors round-trip through XML marshaling and unmarshaling
	tests := []Color{
		{A: 255, B: 0, G: 0, R: 255},     // red
		{A: 255, B: 0, G: 255, R: 0},     // green
		{A: 255, B: 255, G: 0, R: 0},     // blue
		{A: 128, B: 64, G: 192, R: 32},   // arbitrary color
		{A: 0, B: 0, G: 0, R: 0},         // fully transparent black
		{A: 255, B: 255, G: 255, R: 255}, // white
		White,
		Black,
		Red,
		Green,
		Blue,
		Transparent,
	}

	type wrapper struct {
		Color Color `xml:"color"`
	}

	for _, original := range tests {
		t.Run(original.Hex(), func(t *testing.T) {
			// Marshal
			w := wrapper{Color: original}
			data, err := xml.Marshal(w)
			if err != nil {
				t.Errorf("Marshal() unexpected error: %v", err)
				return
			}

			// Unmarshal
			var w2 wrapper
			err = xml.Unmarshal(data, &w2)
			if err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
				return
			}

			if w2.Color != original {
				t.Errorf("Round-trip failed: original %+v -> xml %q -> parsed %+v",
					original, string(data), w2.Color)
			}
		})
	}
}
