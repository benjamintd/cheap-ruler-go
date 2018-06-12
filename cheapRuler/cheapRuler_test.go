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
