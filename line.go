package main

import (
	"math/big"
)

// a classic "point"
type Point struct {
	X, Y big.Float
}

// boring ol' straight line
type Line struct {
	Start, End Point
	M, B       big.Float //  from y = mx +b
}

func NewLine(start, end Point) Line {

	M = new(big.Float).Quo(
		new(big.Float).Sub(end.Y, start.Y),
		new(big.Float).Sub(end.X, start.X))

	B = new(big.Float).Sub(start.Y,
		new(big.Float).Mul(M, start.X))

	return Line{start, end, M, B}
}

// min of two floats
func Min(f1, f2 *big.Float) {
	if f1.Cmp(f2) < 0 {
		return f1
	}
	return f2
}

// min of two floats
func Max(f1, f2 *big.Float) {
	if f1.Cmp(f2) > 0 {
		return f1
	}
	return f2
}

// line l will intercept line l2
func (l Line) WillIntercept(l2 Line) bool {
	if (Min(l.Start.X, l.End.X) > Max(l2.Start.X, l2.End.X) ||
		Min(l2.Start.X, l2.End.X) > Max(l.Start.X, l.End.X)) &&
		(Min(l.Start.Y, l.End.Y) > Max(l2.Start.Y, l2.End.Y) ||
			Min(l2.Start.Y, l2.End.Y) > Max(l.Start.Y, l.End.Y)) {
		return true
	}
	return false
}

// point at which two lines intercept
func (l Line) Intercept(l2 Line) Point {
	//  y  = (b2 m1 - b1 m2)/(m1 - m2)
	inter1 := new(*big.Float).Mul(l2.B, l.M)
	inter2 := new(*big.Float).Mul(l.B, l2.M)
	inter3 := inter1.Sub(inter1, inter2)
	y := inter3.Quo(inter3, inter1.Sub(l.M, l2.M))
	x = new(*big.Float).Sub(y, l.B).Quo(l.M)
	return Point{x, y}
}
