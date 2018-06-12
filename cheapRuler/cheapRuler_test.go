package cheapRuler

import (
	"math"
	"testing"
)

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

	if math.Abs(distance - 585.71) > 1e-2 {
		t.Fatalf("%f != %f", distance, expected)
	}

	t.Log("OK", distance)
}
