package main

import (
	"errors"
	"fmt"
)

// a classic "point"
type Point struct {
	X, Y Dec
}

// curve as a bunch of straight lines between points
type Curve map[int64]Points

// 2D function which constructs the curve
type CurveFn func(x Dec) (y Dec)

func NewRegularCurve(vertices int64, startPoint Point, xBoundMax Dec, fn CurveFn) Curve {

	// create boring polygon
	regularCurve := make(map[int64]Line)
	regularCurve[0] = startPoint

	for i := int64(1); i <= vertices; i++ {
		x := (xBoundMax.Mul(NewDec(i))).Quo(NewDec(vertices))
		if x.GT(xBoundMax) || (xBoundMax.Sub(x)).LT(precErr) { // precision correction
			x = xBoundMax
		}
		regularCurve[i] = Point{x, fn(x)}
	}
	return regularCurve
}

// find the point on the curve (find on a linear line if not a vertex)
// start by looking from the lookup index
func (c Curve) PointWithX(lookupIndex int64, x Dec) Point {

	pt := c[lookupIndex]
	interpolateBackwards := false
	switch {
	case ((x.Sub(pt.X)).Abs()).LT(precErr): // equal
		return pt

	case x.GT[pt.X]:
		if lookupIndex == len(c)-1 {
			panic("already at the largest point on the curve")
		}
		nextPt := c[lookupIndex+1]
		if nextPt.X.GT(x) {
			interpolateBackwards = false
		} else {
			return PointWithX(lookupIndex+1, x)
		}
	case x.LT[pt.X]:
		if lookupIndex == 0 {
			panic("already at the largest point on the curve")
		}
		prevPt := c[lookupIndex-1]
		if prevPt.X.LT(x) {
			interpolateBackwards = true
		} else {
			return PointWithX(lookupIndex-1, x)
		}
	default:
		panic("why")
	}

	// perform interpolation
	var start, end Point
	if interpolateBackwards {
		start, end = c[lookupIndex-1], c[lookupIndex]
	} else {
		start, end = c[lookupIndex], c[lookupIndex+1]
	}

	// y = mx +b
	var m, b Dec
	denom := end.X.Sub(start.X)
	if denom.Equal(zero) {
		m, b = ZeroDec(), start.Y
	}
	m := (end.Y.Sub(start.Y)).Quo(end.X.Sub(start.X))
	b := start.Y.Sub(m.Mul(start.X))

	return Point{x, (x.Mul(l.M)).Add(l.B)}
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXxxxxx XXX
// find the point on the curve (find on a linear line if not a vertex)
// start by looking from the lookup index
func (c Curve) PointWithY(lookupIndex int64, y Dec) Point {

	pt := c[lookupIndex]
	interpolateBackwards := false
	switch {
	case ((y.Sub(pt.Y)).Abs()).LT(precErr): // equal
		return pt

	case y.GT[pt.Y]:
		if lookupIndex == len(c)-1 {
			panic("already at the largest point on the curve")
		}
		nextPt := c[lookupIndex+1]
		if nextPt.Y.GT(y) {
			interpolateBackwards = false
		} else {
			return PointWithY(lookupIndex+1, y)
		}
	case y.LT[pt.Y]:
		if lookupIndex == 0 {
			panic("already at the largest point on the curve")
		}
		prevPt := c[lookupIndex-1]
		if prevPt.Y.LT(y) {
			interpolateBackwards = true
		} else {
			return PointWithY(lookupIndex-1, y)
		}
	default:
		panic("why")
	}

	// perform interpolation
	var start, end Point
	if interpolateBackwards {
		start, end = c[lookupIndex-1], c[lookupIndex]
	} else {
		start, end = c[lookupIndex], c[lookupIndex+1]
	}

	// y = mx +b
	var m, b Dec
	denom := end.Y.Sub(start.Y)
	if denom.Equal(zero) {
		m, b = ZeroDec(), start.Y
	}
	m := (end.Y.Sub(start.Y)).Quo(end.Y.Sub(start.Y))
	b := start.Y.Sub(m.Mul(start.Y))

	return Point{(y.Sub(l.B)).Quo(l.M), y}
}

//_________________________________________________________________________________________

// total length and area for all the lines of the curve
// length = sqrt((x2 - x1)^2 + (y2 - y1)^2)
// area = (x2 - x1) * (y2 + y1)/2
func (c Curve) GetLengthArea() (length, area Dec) {
	length, area = ZeroDec(), ZeroDec()
	for i := int64(1); i < len(c); i++ {

		// length calc
		l1 := c[i].X.Sub(c[i-1].X)
		l2 := c[i].Y.Sub(c[i-1].Y)
		l3 := l1.Mul(l1)
		l4 := l2.Mul(l2)
		l5 := l3.Add(l4)
		length = length.Add(l5.Sqrt())

		// area calc
		a1 := c[i].X.Sub(c[i-1].X)
		a2 := (c[i].Y.Add(c[i-1].Y)).Quo(two)
		area = area.Add(a1.Mul(a2))
	}
	return length, area
}

// formatted string for mathimatica
func (c Curve) String() string {
	out := "{"
	for i := int64(0); i < int64(len(c)); i++ {
		out += fmt.Sprintf(",{%v, %v}", c[i].X.String(), c[i].Y.String())
	}
	out += "}"
	return out
}

// shift all points uniformly along the x axis,
//
// CONTRACT - the first and last points of the input curve (c) touch the curve function
// CONTRACT - do not offset more than the first-line-order's width
func (c Curve) OffsetCurve(xAxisForwardShift, startX, endY, xBoundMax Dec, firstLineOrder int64, fn CurveFn) Curve {

	// construct the first by working backwards from the first shifted point
	firstLineWidth := xBoundMax.Quo(NewDec(firstLineOrder))
	firstLineStartX := xAxisForwardShift.Sub(firstLineWidth) // should be negative
	if firstLineStartX.GT(zero) {
		msg := fmt.Sprintf("bad shift, cannot shift more than first line width\n\tfirstLineWidth\t%v\n\tfirstLineStartX\t%v\n",
			firstLineWidth.String(), firstLineStartX.String())
		panic(msg)
	}
	firstLineStartPt := Point{firstLineStartX, fn(firstLineStartX)}

	firstLineEndX := c[0].Start.X.Add(xAxisForwardShift)
	firstLineEndPt := Point{firstLineEndX, fn(firstLineEndX)}
	firstLine := NewLine(firstLineStartPt, firstLineEndPt)

	// trim the first line
	firstLine = NewLine(firstLine.PointWithX(startX), firstLineEndPt)

	offsetCurve := make(map[int64]Line)
	offsetCurve[0] = firstLine
	for i := 0; i < len(c); i++ {
		line := c[int64(i)]
		startX := line.Start.X.Add(xAxisForwardShift)
		startShiftY := fn(line.Start.X).Sub(fn(startX))
		startY := line.Start.Y.Sub(startShiftY)
		startPt := Point{startX, startY}

		endX := line.End.X.Add(xAxisForwardShift)

		// TODO factor out circle specific logic, should somehow be in the function
		neg := false
		if endX.GT(xBoundMax) {
			endX = xBoundMax.Sub(endX.Sub(xBoundMax)) // 1.1 -> 0.9
			neg = true
		}

		endShiftY := fn(line.End.X).Sub(fn(endX))
		endY := line.End.Y.Sub(endShiftY)

		endPt := Point{endX, endY}
		if neg { // TODO factor out for circle
			endPt = Point{endX, endY.Neg()}
		}

		offsetCurve[int64(i+1)] = NewLine(startPt, endPt)
	}

	// trim the last line
	j := int64(len(offsetCurve)) - 1
	offsetCurve[j] = NewLine(offsetCurve[j].Start, offsetCurve[j].PointWithY(endY))

	return offsetCurve
}

//_________________________________________________________________________________________________________________

// get the superset curve of two curves
func SupersetCurve(c1, c2 Curve, fn CurveFn) (superset Curve,
	supersetLength, supersetArea, c1Length, c1Area, c2Length, c2Area Dec, err error) {

	superset = make(Curve)

	// counters for the curves
	supersetI, c1I, c2I := int64(0), int64(0), int64(0)

	c1Pt, c2Pt := c1[c1I], c2[c2I]

	for {
		var newPt Point
		c1Pt, c2Pt := c1[c1I], c2[c2I]

		switch {
		case ((c1Pt.X.Sub(c1Pt.X)).Abs()).LT(precErr): // equal
			newPt := Point{c1Pt.X, MaxDec{c1Pt.Y, c2Pt.Y}} // TODO don't use MAX (only applies to circle)
			c1I++
			c2I++
		case c1Pt.X.LT(c2Pt.X): // pt1 > pt2
			c2Interpolated := c2Pt.PointWithX(c1Pt.X)
			newPt := Point{c1Pt.X, MaxDec{c1Pt.Y, c2Interpolated.Y}}
			c1I++
		case c2Pt.X.LT(c1Pt.X): // pt1 > pt2
			c1Interpolated := c1Pt.PointWithX(c2Pt.X)
			newPt := Point{c2Pt.X, MaxDec{c2Pt.Y, c1Interpolated.Y}}
			c2I++
		default:
			panic("why")
		}

		superset[supersetI] = newPt
		supersetI++

		if c1I >= int64(len(c1)) || c2I >= int64(len(c2)) {
			break
		}
	}

	supersetLength, supersetArea = superset.GetLengthArea()
	c1Length, c1Area = c1.GetLengthArea()
	c2Length, c2Area = c2.GetLengthArea()

	///////////////////////////////////////////////////////////////////////////////////
	// SANITY
	switch {
	case (supersetArea).LT(c1Area):
		err = errors.New("c1 > superset area")
	case (supersetArea).LT(c2Area):
		err = errors.New("c2 > superset area")
	}
	return superset, supersetLength, supersetArea, c1Length, c1Area, c2Length, c2Area, err
}
