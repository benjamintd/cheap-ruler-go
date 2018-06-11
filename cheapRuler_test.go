package cheapRuler

import "testing"

func TestNewRuler(t *testing.T) {
	t.Log("NewRuler returns expected coefficients")

	ruler, err := NewRuler(42.0, "kilometers")

	if err != nil {
		t.Error(err)
	}

  expected := Ruler{kx: 82.853048511947, ky: 111.07174788624647}

  if ruler != expected {
		t.Fatalf("%s != %s", ruler, expected)
	}

	t.Log("OK", ruler)
}
