package main

import (
	"fmt"
)

// curve as a bunch of lines
type Curve map[int64]Line

func NewRegularDivisionCurve(divisions int64, startPoint Point, xBoundMax Dec, fn func(x Dec) (y Dec)) Curve {

	// create boring polygon
	regularDivision := make(map[int64]Line)

	for side := int64(0); side < divisions; side++ {
		x2 := (xBoundMax.Mul(NewDec(side + 1))).Quo(NewDec(divisions))
		if x2.GT(xBoundMax) || (xBoundMax.Sub(x2)).LT(precErr) { // precision correction
			x2 = xBoundMax
		}
		y2 := fn(x2)
		endPoint := Point{x2, y2}
		regularDivision[side] = NewLine(startPoint, endPoint, divisions)
		startPoint = endPoint
	}
	return regularDivision
}

//____________________________________

// total length and area for all the lines
func (c Curve) GetLengthArea() (length, area Dec) {
	length, area = ZeroDec(), ZeroDec()
	for _, line := range c {
		length = length.Add(line.Length())
		area = area.Add(line.Area())
	}
	return length, area
}

func (c Curve) String() string {
	out := "{"
	out += fmt.Sprintf("{%v, %v}", c[0].Start.X, c[0].Start.Y)

	for i := int64(0); i < int64(len(c)); i++ {
		line := c[i]
		out += fmt.Sprintf(",{%v, %v}", line.End.X.String(), line.End.Y.String())
	}
	out += "}"
	return out
}
