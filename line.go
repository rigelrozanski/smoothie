package main

import (
	"math/big"
)

// a classic "point"
type Point struct {
	X, Y *big.Float
}

// boring ol' straight line
type Line struct {
	Start, End Point
	M, B       *big.Float //  from y = mx +b
}

func NewLine(start, end Point) Line {

	initFloat(start.X)
	initFloat(start.Y)
	initFloat(end.X)
	initFloat(end.Y)

	m := newFloat().Quo(
		newFloat().Sub(end.Y, start.Y),
		newFloat().Sub(end.X, start.X))

	b := newFloat().Sub(start.Y,
		newFloat().Mul(m, start.X))

	return Line{start, end, m, b}
}

// min of two floats
func Min(f1, f2 *big.Float) *big.Float {
	if f1.Cmp(f2) < 0 {
		return f1
	}
	return f2
}

// min of two floats
func Max(f1, f2 *big.Float) *big.Float {
	if f1.Cmp(f2) > 0 {
		return f1
	}
	return f2
}

//_______________________________________________________________________

// length of a line!
// sqrt((x2 - x1)^2 + (y2 - y1)^2)
func (l Line) Length() *big.Float {
	inter1 := newFloat().Sub(l.End.X, l.Start.X)
	inter2 := newFloat().Sub(l.End.Y, l.Start.Y)
	inter3 := inter1.Mul(inter1, inter1)
	inter4 := inter2.Mul(inter2, inter2)
	inter5 := inter1.Add(inter3, inter4)
	return inter1.Sqrt(inter5)
}

// y-axis end of line l is within end of l2
func (l Line) WithinL2XBound(l2 Line) bool {
	return l.End.X.Cmp(l2.End.X) <= 0
}

var zero = big.NewFloat(0)

// point at which two lines intercept
func (l Line) Intercept(l2 Line) (intercept Point, withinBounds bool) {
	//  y  = (b2 m1 - b1 m2)/(m1 - m2)
	inter1 := newFloat().Mul(l2.B, l.M)

	inter2 := newFloat().Mul(l.B, l2.M)
	inter3 := newFloat().Sub(inter1, inter2)
	inter4 := newFloat().Sub(l.M, l2.M)
	if inter4.Cmp(zero) == 0 {
		return intercept, false
	}
	y := newFloat().Quo(inter3, inter4)
	x := newFloat().Quo(newFloat().Sub(y, l.B), l.M)
	intercept = Point{x, y}

	// check if intercept is in precision error amount
	proximityToZero := newFloat().Abs(newFloat().Sub(l.Start.X, intercept.X))
	if proximityToZero.Cmp(big.NewFloat(precCutoff)) < 0 {
		return intercept, false
	}

	withinBounds = false
	// to deal with precision errors,
	//  if any of the points are equal we know we're not within bounds
	if intercept.X.Cmp(l.Start.X) > 0 &&
		intercept.X.Cmp(l.End.X) < 0 &&
		intercept.X.Cmp(l2.Start.X) > 0 &&
		intercept.X.Cmp(l2.End.X) < 0 &&
		intercept.Y.Cmp(l.Start.Y) < 0 &&
		intercept.Y.Cmp(l.End.Y) > 0 &&
		intercept.Y.Cmp(l2.Start.Y) < 0 &&
		intercept.Y.Cmp(l2.End.Y) > 0 {
		withinBounds = true
	}

	return intercept, withinBounds
}
