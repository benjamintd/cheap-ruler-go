[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/benjamintd/cheap-ruler-go/cheapRuler) [![Build Status](https://travis-ci.org/benjamintd/cheap-ruler-go.svg?branch=master)](https://travis-ci.org/benjamintd/cheap-ruler-go)

# cheapRuler

A Golang port of @mourner's cheap-ruler:

> "A collection of very fast approximations to common geodesic measurements.
Useful for performance-sensitive code that measures things on a city scale.

> The approximations are based on an [FCC-approved formula of ellipsoidal Earth projection](https://www.gpo.gov/fdsys/pkg/CFR-2005-title47-vol4/pdf/CFR-2005-title47-vol4-sec73-208.pdf).
For distances under 500 kilometers and not on the poles,
the results are very precise — within [0.1% margin of error](#precision)
compared to [Vincenti formulas](https://en.wikipedia.org/wiki/Vincenty%27s_formulae),
and usually much less for shorter distances."

## Usage

```go
package main

import (
  "github.com/benjamintd/cheap-ruler-go/cheapRuler"
  "fmt"
)

func main() {
  // latitude of Paris
  ruler, _ := cheapRuler.NewRuler(48.8629, "meters")
  a := [2]float64{2.344808, 48.862851}
  b := [2]float64{2.352790, 48.862907}
  distance := ruler.Distance(a, b)
  fmt.Println(distance)
}
// 585.71 meters
```

## License

MIT

## See also

https://github.com/mapbox/cheap-ruler
https://github.com/JamesMilnerUK/cheap-ruler-go
