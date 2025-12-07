package kml

import (
	"bytes"
	"os"
	"testing"
)

// Benchmark parsing

func BenchmarkParseSimple(b *testing.B) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Placemark>
    <name>Test</name>
    <Point><coordinates>-122.0,37.0</coordinates></Point>
  </Placemark>
</kml>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseBytes(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseDocument(b *testing.B) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Test Document</name>
    <Style id="style1">
      <LineStyle><color>ff0000ff</color><width>2</width></LineStyle>
      <PolyStyle><color>7f00ff00</color></PolyStyle>
    </Style>
    <Folder>
      <name>Folder 1</name>
      <Placemark>
        <name>Point 1</name>
        <styleUrl>#style1</styleUrl>
        <Point><coordinates>-122.0,37.0</coordinates></Point>
      </Placemark>
      <Placemark>
        <name>Point 2</name>
        <Point><coordinates>-121.0,38.0</coordinates></Point>
      </Placemark>
    </Folder>
  </Document>
</kml>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseBytes(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLargeFile(b *testing.B) {
	data, err := os.ReadFile("testdata/KML_Samples.kml")
	if err != nil {
		b.Skip("KML_Samples.kml not found")
	}

	b.ResetTimer()
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		_, err := ParseBytes(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseVeryLargeFile(b *testing.B) {
	data, err := os.ReadFile("testdata/catalina-points.kml")
	if err != nil {
		b.Skip("catalina-points.kml not found")
	}

	b.ResetTimer()
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		_, err := ParseBytes(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark writing

func BenchmarkWriteSimple(b *testing.B) {
	docBuilder := NewKMLBuilder().Document("Test")
	docBuilder.Placemark("Point").Point(-122.0, 37.0)
	doc := docBuilder.Build()

	buf := &bytes.Buffer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := doc.Write(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteComplex(b *testing.B) {
	docBuilder := NewKMLBuilder().Document("Test")
	docBuilder.Style("style1").
		LineStyle().Color(Red).Width(2).Done().
		PolyStyle().Color(Green).Fill(true).Done()
	folder := docBuilder.Folder("Folder 1")
	folder.Placemark("P1").Point(-122.0, 37.0).StyleURL("#style1")
	folder.Placemark("P2").Point(-121.0, 38.0)
	folder.Placemark("P3").LineString(
		Coord(-122.0, 37.0),
		Coord(-121.0, 38.0),
		Coord(-120.0, 37.5),
	)
	doc := docBuilder.Build()

	buf := &bytes.Buffer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := doc.Write(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteIndent(b *testing.B) {
	docBuilder := NewKMLBuilder().Document("Test")
	folder := docBuilder.Folder("Folder")
	folder.Placemark("P1").Point(-122.0, 37.0)
	folder.Placemark("P2").Point(-121.0, 38.0)
	doc := docBuilder.Build()

	buf := &bytes.Buffer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := doc.WriteIndent(buf, "", "  ")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark round-trip

func BenchmarkRoundTrip(b *testing.B) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Test</name>
    <Folder>
      <Placemark><name>P1</name><Point><coordinates>-122.0,37.0</coordinates></Point></Placemark>
      <Placemark><name>P2</name><Point><coordinates>-121.0,38.0</coordinates></Point></Placemark>
    </Folder>
  </Document>
</kml>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := ParseBytes(data)
		if err != nil {
			b.Fatal(err)
		}
		_, err = doc.Bytes()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark coordinate parsing

func BenchmarkParseCoordinates(b *testing.B) {
	coordStr := "-122.084075,37.4220033612141,0 -122.085135,37.4220033612141,0 -122.085135,37.4230033612141,0 -122.084075,37.4230033612141,0 -122.084075,37.4220033612141,0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseCoordinates(coordStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseCoordinatesLarge(b *testing.B) {
	// Generate a large coordinate string (100 points)
	var buf bytes.Buffer
	for i := 0; i < 100; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString("-122.0,37.0,0")
	}
	coordStr := buf.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseCoordinates(coordStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark color parsing

func BenchmarkParseColor(b *testing.B) {
	colorStr := "ff0000ff"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseColor(colorStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkColorHex(b *testing.B) {
	color := RGBA(255, 0, 0, 255)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = color.Hex()
	}
}

// Benchmark traversal

func BenchmarkWalk(b *testing.B) {
	data, err := os.ReadFile("testdata/KML_Samples.kml")
	if err != nil {
		b.Skip("KML_Samples.kml not found")
	}

	doc, err := ParseBytes(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		doc.Walk(func(f Feature) error {
			count++
			return nil
		})
	}
}

func BenchmarkPlacemarks(b *testing.B) {
	data, err := os.ReadFile("testdata/KML_Samples.kml")
	if err != nil {
		b.Skip("KML_Samples.kml not found")
	}

	doc, err := ParseBytes(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.Placemarks()
	}
}

func BenchmarkFilter(b *testing.B) {
	data, err := os.ReadFile("testdata/KML_Samples.kml")
	if err != nil {
		b.Skip("KML_Samples.kml not found")
	}

	doc, err := ParseBytes(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.Filter(func(f Feature) bool {
			_, ok := f.(*Placemark)
			return ok
		})
	}
}

func BenchmarkBounds(b *testing.B) {
	data, err := os.ReadFile("testdata/KML_Samples.kml")
	if err != nil {
		b.Skip("KML_Samples.kml not found")
	}

	doc, err := ParseBytes(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = doc.Bounds()
	}
}

func BenchmarkFindByID(b *testing.B) {
	data, err := os.ReadFile("testdata/KML_Samples.kml")
	if err != nil {
		b.Skip("KML_Samples.kml not found")
	}

	doc, err := ParseBytes(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.FindByID("transPurpleLineGreenPoly")
	}
}

// Benchmark builder

func BenchmarkBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		docBuilder := NewKMLBuilder().Document("Test")
		docBuilder.Style("s1").LineStyle().Color(Red).Width(2).Done()
		folder := docBuilder.Folder("F1")
		folder.Placemark("P1").Point(-122.0, 37.0)
		folder.Placemark("P2").Point(-121.0, 38.0)
		_ = docBuilder.Build()
	}
}

// Memory allocation benchmarks

func BenchmarkParseAllocs(b *testing.B) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Test</name>
    <Placemark><name>P1</name><Point><coordinates>-122.0,37.0</coordinates></Point></Placemark>
  </Document>
</kml>`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseBytes(data)
	}
}

func BenchmarkCoordParseAllocs(b *testing.B) {
	coordStr := "-122.0,37.0,0 -121.0,38.0,0 -120.0,37.5,0"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseCoordinates(coordStr)
	}
}

// Parallel benchmarks

func BenchmarkParseParallel(b *testing.B) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Test</name>
    <Placemark><name>P1</name><Point><coordinates>-122.0,37.0</coordinates></Point></Placemark>
  </Document>
</kml>`)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := ParseBytes(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
