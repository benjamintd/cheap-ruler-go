package cheapRuler

import (
	"math"
	"testing"
)

var testLine Line = Line{
	Point{2.3503875, 48.863598},
	Point{2.3501086, 48.8627334},
	Point{2.3485958, 48.862747},
	Point{2.3482418, 48.86240},
	Point{2.3477053, 48.86240},
	Point{2.3469865, 48.862147},
}

func TestNewRuler(t *testing.T) {
	t.Log("NewRuler returns expected coefficients")

	ruler, err := NewRuler(42.0, "kilometers")

	if err != nil {
		t.Error(err)
	}

	expected := Ruler{kx: 82.853048511947, ky: 111.07174788624647}

	if math.Abs(ruler.kx-expected.kx) > 1e-5 || math.Abs(ruler.ky-expected.ky) > 1e-5 {
		t.Fatalf("%+v != %+v", ruler, expected)
	}

	t.Log("OK", ruler)
}

func TestDistance(t *testing.T) {
	t.Log("ruler distance is correct")

	ruler, _ := NewRuler(48.8629, "meters")
	a := [2]float64{2.344808, 48.862851}
	b := [2]float64{2.352790, 48.862907}
	distance := ruler.Distance(a, b)
	expected := 585.71

	if math.Abs(distance-expected) > 1e-2 {
		t.Fatalf("%f != %f", distance, expected)
	}

	t.Log("OK", distance)
}

func TestLineDistance(t *testing.T) {
	t.Log("ruler line distance is correct")

	ruler, _ := NewRuler(48.8629, "miles")
	distance := ruler.LineDistance(testLine)
	expected := 0.220571

	if math.Abs(distance-expected) > 1e-2 {
		t.Fatalf("%f != %f", distance, expected)
	}

	t.Log("OK", distance)
}

func TestBearing(t *testing.T) {
	t.Log("ruler bearing is correct")

	ruler, _ := NewRuler(48.8629, "miles")
	a := [2]float64{2.344808, 48.862851}
	b := [2]float64{2.352790, 48.862907}
	bearing := ruler.Bearing(a, b)
	expected := 89.39

	if math.Abs(bearing-expected) > 1e-2 {
		t.Fatalf("%f != %f", bearing, expected)
	}

	t.Log("OK", bearing)
}

func TestOffset(t *testing.T) {
	t.Log("ruler offset is correct")

	ruler, _ := NewRuler(48.8629, "miles")
	a := [2]float64{2.344808, 48.862851}
	offset := ruler.Offset(a, 1., -2.)
	expected := [2]float64{2.366741, 48.833907}

	if math.Abs(offset[0]-expected[0]) > 1e-5 || math.Abs(offset[1]-expected[1]) > 1e-5 {
		t.Fatalf("%+v != %+v", offset, expected)
	}

	t.Log("OK", offset)
}

func TestDestination(t *testing.T) {
	t.Log("ruler destination is correct")

	ruler, _ := NewRuler(48.8629, "meters")
	a := [2]float64{2.344808, 48.862851}
	destination := ruler.Destination(a, 1., 30.)
	expected := [2]float64{2.344814, 48.862858}

	if math.Abs(destination[0]-expected[0]) > 1e-5 || math.Abs(destination[1]-expected[1]) > 1e-5 {
		t.Fatalf("%+v != %+v", destination, expected)
	}

	t.Log("OK", destination)
}

func TestAlong(t *testing.T) {
	t.Log("ruler along is correct")

	ruler, _ := NewRuler(48.8629, "meters")
	along := ruler.Along(testLine, 150.)
	expected := [2]float64{2.349404, 48.862739}

	if math.Abs(along[0]-expected[0]) > 1e-5 || math.Abs(along[1]-expected[1]) > 1e-5 {
		t.Fatalf("%+v != %+v", along, expected)
	}

	t.Log("OK", along)
}

func TestPointOnLine(t *testing.T) {
	t.Log("ruler pointOnLine is correct")

	ruler, _ := NewRuler(48.8629, "meters")
	pol := ruler.PointOnLine(testLine, [2]float64{2.350, 48.861})
	var expected PointOnLine = PointOnLine{
		point: [2]float64{2.3500358, 48.862734},
		index: 1,
		t:     0.048116,
	}

	if math.Abs(pol.point[0]-expected.point[0]) > 1e-5 ||
		math.Abs(pol.point[1]-expected.point[1]) > 1e-5 ||
		pol.index != expected.index ||
		math.Abs(pol.t-expected.t) > 1e-5 {
		t.Fatalf("%+v != %+v", pol, expected)
	}

	t.Log("OK", pol)
}

func TestLineSlice(t *testing.T) {
	t.Log("ruler line slice is correct")

	ruler, _ := NewRuler(48.8629, "miles")
	a := Point{2.350054, 48.863154}
	b := Point{2.347555, 48.862145}
	slice := ruler.LineSlice(a, b, testLine)
	sliceDistance := ruler.LineDistance(slice)
	expected := 0.164593

	if math.Abs(sliceDistance-expected) > 1e-5 {
		t.Fatalf("%f != %f", sliceDistance, expected)
	}

	t.Log("OK", slice)
}

func TestLineSliceAlong(t *testing.T) {
	t.Log("ruler line slice along is correct")

	ruler, _ := NewRuler(48.8629, "miles")
	slice := ruler.LineSliceAlong(0.05, 0.15, testLine)
	sliceDistance := ruler.LineDistance(slice)
	expected := 0.1

	if math.Abs(sliceDistance-expected) > 1e-5 {
		t.Fatalf("%f != %f", sliceDistance, expected)
	}

	t.Log("OK", slice)
}

func TestBufferPoint(t *testing.T) {
	t.Log("ruler buffer point is correct")

	ruler, _ := NewRuler(48.8629, "meters")
	a := Point{2.350054, 48.863154}
	bbox := ruler.BufferPoint(a, 12)
	expected := Bbox{2.349946, 48.862990, 2.350162, 48.863318}

	if math.Abs(bbox[0]-expected[0]) > 1e-5 ||
		math.Abs(bbox[1]-expected[1]) > 1e-5 ||
		math.Abs(bbox[2]-expected[2]) > 1e-5 ||
		math.Abs(bbox[3]-expected[3]) > 1e-5 {
		t.Fatalf("%f != %f", bbox, expected)
	}

	t.Log("OK", bbox)
}

func TestBufferBbox(t *testing.T) {
	t.Log("ruler buffer bbox is correct")

	ruler, _ := NewRuler(48.8629, "meters")
	a := Bbox{2.349946, 48.862990, 2.350162, 48.863318}
	bbox := ruler.BufferBbox(a, 12)
	expected := Bbox{2.349838, 48.862826, 2.350270, 48.863482}

	if math.Abs(bbox[0]-expected[0]) > 1e-5 ||
		math.Abs(bbox[1]-expected[1]) > 1e-5 ||
		math.Abs(bbox[2]-expected[2]) > 1e-5 ||
		math.Abs(bbox[3]-expected[3]) > 1e-5 {
		t.Fatalf("%f != %f", bbox, expected)
	}

	t.Log("OK", bbox)
}

func TestInsideBbox(t *testing.T) {
	t.Log("ruler inside bbox is correct")

	ruler, _ := NewRuler(48.8629, "meters")
	bbox := Bbox{2.349946, 48.862990, 2.350162, 48.863318}
	pointIn := Point{2.35, 48.863}
	pointOut := Point{2.349, 48.863}

	if !(ruler.InsideBbox(pointIn, bbox)) || ruler.InsideBbox(pointOut, bbox) {
		t.Fatalf("Inside bbox gives false result")
	}

	t.Log("OK", bbox)
}
