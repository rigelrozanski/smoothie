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

	m := new(big.Float).Quo(
		new(big.Float).Sub(end.Y, start.Y),
		new(big.Float).Sub(end.X, start.X))

	b := new(big.Float).Sub(start.Y,
		new(big.Float).Mul(m, start.X))

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
	inter1 := new(big.Float).Sub(l.End.X, l.Start.X)
	inter2 := new(big.Float).Sub(l.End.Y, l.Start.Y)
	inter3 := inter1.Mul(inter1, inter1)
	inter4 := inter2.Mul(inter2, inter2)
	inter5 := inter1.Add(inter3, inter4)
	return inter1.Sqrt(inter5)
}

// line l will intercept line l2
func (l Line) WillIntercept(l2 Line) bool {
	if (Min(l.Start.X, l.End.X).Cmp(Max(l2.Start.X, l2.End.X)) > 0 ||
		Min(l2.Start.X, l2.End.X).Cmp(Max(l.Start.X, l.End.X)) > 0) &&
		(Min(l.Start.Y, l.End.Y).Cmp(Max(l2.Start.Y, l2.End.Y)) > 0 ||
			Min(l2.Start.Y, l2.End.Y).Cmp(Max(l.Start.Y, l.End.Y)) > 0) {
		return true
	}
	return false
}

// y-axis end of line l is within end of l2
func (l Line) WithinL2YBound(l2 Line) bool {
	return l.End.Y.Cmp(l2.End.Y) <= 0
}

// point at which two lines intercept
func (l Line) Intercept(l2 Line) Point {
	//  y  = (b2 m1 - b1 m2)/(m1 - m2)
	inter1 := new(big.Float).Mul(l2.B, l.M)
	inter2 := new(big.Float).Mul(l.B, l2.M)
	inter3 := inter1.Sub(inter1, inter2)
	y := inter3.Quo(inter3, inter1.Sub(l.M, l2.M))
	inter4 := new(big.Float).Sub(y, l.B)
	x := inter4.Quo(inter4, l.M)
	return Point{x, y}
}
