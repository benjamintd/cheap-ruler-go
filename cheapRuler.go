package cheapRuler

import (
	"errors"
	"math"
)

// Ruler is the interface implemented by ruler objects.
type CheapRuler interface {
	along(l Line, dist float64) Point
	area(p Polygon) float64
	bearing(a Point, b Point) float64
	destination(p Point, d float64, b float64) Point
	distance(a Point, b Point) float64
	LineDistance(l Line) float64
	LineSlice(start Point, end Point, l Line) Line
	offset(p Point, dx float64, dy float64) float64
	PointOnLine(l Line, p Point) PointOnLine
}

type Ruler struct {
	kx, ky float64
}

type Point [2]float64

type Bbox [4]float64

type Line []Point

type Polygon []Line

type PointOnLine struct {
	point Point
	index int
	t     float64
}

var Units = map[string]float64{
	"kilometers":    1,
	"miles":         1000 / 1609.344,
	"nauticalmiles": 1000 / 1852,
	"meters":        1000,
	"metres":        1000,
	"yards":         1000 / 0.9144,
	"feet":          1000 / 0.3048,
	"inches":        1000 / 0.0254,
}

// CheapRuler instantiates a new ruler from a latitude and a unit. The default unit is kilometers.
func NewRuler(lat float64, unit string) (Ruler, error) {
	var m float64
	var e error = nil
	if scale, ok := Units[unit]; ok {
		m = scale
	} else {
		// falling back to the default kilometers
		m = 1
		e = errors.New(unit + " is not a valid unit")
	}

	cos := math.Cos(lat * math.Pi / 180)
	cos2 := 2*cos*cos - 1
	cos3 := 2*cos*cos2 - cos
	cos4 := 2*cos*cos3 - cos2
	cos5 := 2*cos*cos4 - cos3

	// multipliers for converting longitude and latitude degrees into distance (http://1.usa.gov/1Wb1bv7)
	kx := m * (111.41513*cos - 0.09455*cos3 + 0.00012*cos5)
	ky := m * (111.13209 - 0.56605*cos2 + 0.0012*cos4)

	return Ruler{kx: kx, ky: ky}, e
}

func (r Ruler) distance(a Point, b Point) float64 {
	dx := (a[0] - b[0]) * r.kx
	dy := (a[1] - b[1]) * r.ky
	return math.Sqrt(dx*dx + dy*dy)
}

func (r Ruler) bearing(a Point, b Point) float64 {
	dx := (a[0] - b[0]) * r.kx
	dy := (a[1] - b[1]) * r.ky
	if dx == 0 && dy == 0 {
		return 0
	}

	bearing := math.Atan2(dx, dy) * 180 / math.Pi
	if bearing > 180 {
		bearing -= 360
	}
	return bearing
}

func (r Ruler) offset(p Point, dx float64, dy float64) Point {
	return Point{p[0] + dx/r.kx, p[1] + dy/r.ky}
}

func (r Ruler) LineDistance(l Line) float64 {
	var distance float64 = 0

	for i := 0; i < len(l)-1; i++ {
		distance += r.distance(l[i], l[i+1])
	}
	return distance
}

func (r Ruler) area(p Polygon) float64 {
	var sum float64 = 0

	for i := 0; i < len(p); i++ {
		var ring Line = p[i]
		for j, len, k := 0, len(ring), len(ring)-1; j < len; k, j = j+1, j+1 {
			var isNotHole float64 = 1
			if i > 0 {
				isNotHole = -1
			}
			sum += (ring[j][0] - ring[k][0]) * (ring[j][1] + ring[k][1]) * isNotHole
		}
	}

	return (math.Abs(sum) / 2) * r.kx * r.ky
}

func (r Ruler) along(l Line, dist float64) Point {
	var sum float64 = 0

	if dist <= 0 {
		return l[0]
	}

	for i := 0; i < len(l)-1; i++ {
		p0 := l[i]
		p1 := l[i+1]
		d := r.distance(p0, p1)
		sum += d
		if sum > dist {
			return interpolate(p0, p1, (dist-(sum-d))/d)
		}
	}

	return l[len(l)-1]
}

func (r Ruler) PointOnLine(l Line, p Point) PointOnLine {
	var minDist float64 = math.Inf(1)
	var minX, minY, minT, x, y, dx, dy, t float64
	var minI int

	for i := 0; i < len(l)-1; i++ {

		x = l[i][0]
		y = l[i][1]
		dx = (l[i+1][0] - x) * r.kx
		dy = (l[i+1][1] - y) * r.ky

		if dx != 0 || dy != 0 {

			t = ((p[0]-x)*r.kx*dx + (p[1]-y)*r.ky*dy) / (dx*dx + dy*dy)

			if t > 1 {
				x = l[i+1][0]
				y = l[i+1][1]

			} else if t > 0 {
				x += (dx / r.kx) * t
				y += (dy / r.ky) * t
			}
		}

		dx = (p[0] - x) * r.kx
		dy = (p[1] - y) * r.ky

		var sqDist = dx*dx + dy*dy
		if sqDist < minDist {
			minDist = sqDist
			minX = x
			minY = y
			minI = i
			minT = t
		}
	}

	return PointOnLine{
		point: Point{minX, minY},
		index: minI,
		t:     math.Max(0, math.Min(1, minT)),
	}
}

func (r Ruler) LineSlice(start Point, end Point, l Line) Line {
	p1 := r.PointOnLine(l, start)
	p2 := r.PointOnLine(l, end)

	if p1.index > p2.index || (p1.index == p2.index && p1.t < p2.t) {
		p1, p2 = p2, p1
	}

	var slice Line = []Point{p1.point}

	left := p1.index + 1
	right := p2.index

	if l[left] != slice[0] && left <= right {
		slice = append(slice, l[left])
	}

	for i := left + 1; i <= right; i++ {
		slice = append(slice, l[i])
	}

	if l[right] != p2.point {
		slice = append(slice, p2.point)
	}

	return slice
}

func (r Ruler) LineSliceAlong(start float64, stop float64, l Line) Line {
	var sum float64 = 0
	var slice []Point

	for i := 0; i < len(l)-1; i++ {
		p0 := l[i]
		p1 := l[i+1]
		d := r.distance(p0, p1)

		sum += d

		if sum > start && len(slice) == 0 {
			slice = append(slice, interpolate(p0, p1, (start-(sum-d))/d))
		}

		if sum >= stop {
			slice = append(slice, interpolate(p0, p1, (stop-(sum-d))/d))
		}

		if sum > start {
			slice = append(slice, p1)
		}
	}

	return slice
}

func (r Ruler) bufferPoint(p Point, buffer float64) Bbox {
	v := buffer / r.kx
	h := buffer / r.ky

	return Bbox{
		p[0] - h,
		p[1] - v,
		p[0] + h,
		p[1] + v,
	}
}

func (r Ruler) bufferBbox(b Bbox, buffer float64) Bbox {
	v := buffer / r.kx
	h := buffer / r.ky

	return Bbox{
		b[0] - h,
		b[1] - v,
		b[2] + h,
		b[3] + v,
	}
}

func (r Ruler) insideBbox(p Point, b Bbox) bool {
	return p[0] >= b[0] &&
		p[0] <= b[2] &&
		p[1] >= b[1] &&
		p[1] <= b[3]
}

func interpolate(a Point, b Point, t float64) Point {
	dx := b[0] - a[0]
	dy := b[1] - a[1]
	return Point{a[0] + dx*t, a[1] + dy*t}
}
