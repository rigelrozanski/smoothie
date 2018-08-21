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
type Curve map[int64]Point

// 2D function which constructs the curve
type CurveFn func(x Dec) (y Dec)

func NewRegularCurve(vertices int64, startPoint Point, xBoundMax Dec, fn CurveFn) Curve {

	// create boring polygon
	regularCurve := make(Curve)
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

	case x.GT(pt.X):
		if lookupIndex == int64(len(c))-1 {
			panic("already at the largest point on the curve")
		}
		nextPt := c[lookupIndex+1]
		if nextPt.X.GT(x) {
			interpolateBackwards = false
		} else {
			return c.PointWithX(lookupIndex+1, x)
		}
	case x.LT(pt.X):
		if lookupIndex == 0 {
			panic("already at the largest point on the curve")
		}
		prevPt := c[lookupIndex-1]
		if prevPt.X.LT(x) {
			interpolateBackwards = true
		} else {
			return c.PointWithX(lookupIndex-1, x)
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

	m, b := GetMB(start, end)
	return Point{x, (x.Mul(m)).Add(b)}
}

// from y = mx +b
func GetMB(start, end Point) (m, b Dec) {
	denom := end.X.Sub(start.X)
	if denom.Equal(zero) {
		return ZeroDec(), start.Y
	}
	m = (end.Y.Sub(start.Y)).Quo(end.X.Sub(start.X))
	b = start.Y.Sub(m.Mul(start.X))
	return m, b
}

//_________________________________________________________________________________________

// total length and area for all the lines of the curve
// length = sqrt((x2 - x1)^2 + (y2 - y1)^2)
// area = (x2 - x1) * (y2 + y1)/2
func (c Curve) GetLengthArea() (length, area Dec) {
	length, area = ZeroDec(), ZeroDec()
	for i := int64(1); i < int64(len(c)); i++ {

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
		comma := ","
		if i == 0 {
			comma = ""
		}
		out += fmt.Sprintf("%v{%v, %v}", comma, c[i].X.String(), c[i].Y.String())
	}
	out += "}"
	return out
}

// shift all points uniformly along the x axis, for thexpe starting points use a reflection
// (which is correct for circles) TODO upgrade to use extra function points
//
// CONTRACT - the first and last points of the input curve (c) touch the curve function
// CONTRACT - do not offset more than the first-line-order's width
func (c Curve) OffsetCurve(xAxisForwardShift, xBoundMax Dec, fn CurveFn) Curve {

	// TODO reflection more generic, assumes reflect along the x axis
	// reflect the points along the Y axis
	// by the amount which need to be shifted
	reflected := make(Curve)
	for i := int64(0); i < int64(len(c))-1; i++ { // should not reach the end
		if c[i].X.LT(xAxisForwardShift) {
			reflected[i] = Point{c[i].X.Neg(), c[i].Y}
			continue
		}
		if ((c[i].X.Sub(xAxisForwardShift)).Abs()).LT(precErr) { // equal
			reflected[i] = Point{c[i].X.Neg(), c[i].Y}
			break
		}

		// MUST be the final point - interpolate
		m, b := GetMB(c[i-1], c[i])
		reflected[i] = Point{c[i].X.Neg(), (c[i].X.Mul(m)).Add(b)}
		break
	}

	// add the reflected points to a new curve (in reverse order)
	combined := make(Curve)
	combinedI := int64(0)
	for i := int64(len(reflected)) - 1; i >= 0; i-- { // should not reach the end
		combined[combinedI] = reflected[i]
		combinedI++
	}
	for i := int64(0); i < int64(len(c)); i++ { // now add the rest of the curve
		combined[combinedI] = c[i]
		combinedI++
	}

	// now generate the offset with the provided points
	offsetCurve := make(Curve)
	offsetCurveI := int64(0)
	for i := int64(0); i < int64(len(combined)); i++ {
		pt := combined[i]
		newX := pt.X.Add(xAxisForwardShift)

		// TODO factor out circle specific logic, should somehow be in the function
		// this assumes the circle function and requires it to trim the final point
		negTrim := false
		if newX.GT(xBoundMax) {
			newX = xBoundMax.Sub(newX.Sub(xBoundMax)) // 1.1 -> 0.9
			negTrim = true
		}

		shiftY := fn(pt.X).Sub(fn(newX))
		newY := pt.Y.Sub(shiftY)

		if !negTrim {
			offsetCurve[offsetCurveI] = Point{newX, newY}

		} else { // if negative than must be the final point

			// trim the final point
			untrimmedPt := Point{newX, newY.Neg()}
			m, b := GetMB(offsetCurve[offsetCurveI-1], untrimmedPt)

			finalY := zero
			if !m.Equal(zero) {
				offsetCurve[offsetCurveI] = Point{(finalY.Sub(b)).Quo(m), finalY}
			} else {
				offsetCurve[offsetCurveI] = Point{newX, finalY} // vertical line
			}
			break
		}
		offsetCurveI++
	}

	return offsetCurve
}

//_________________________________________________________________________________________________________________

// get the superset curve of two curves
func SupersetCurve(c1, c2 Curve, fn CurveFn) (superset Curve,
	supersetLength, supersetArea, c1Length, c1Area, c2Length, c2Area Dec, err error) {

	superset = make(Curve)

	// counters for the curves
	supersetI, c1I, c2I := int64(0), int64(0), int64(0)

	for {
		var newPt Point
		c1Pt, c2Pt := c1[c1I], c2[c2I]
		fmt.Printf("debug c1Pt: %v\n", c1Pt)
		fmt.Printf("debug c2Pt: %v\n", c2Pt)

		switch {
		case ((c1Pt.X.Sub(c2Pt.X)).Abs()).LT(precErr): // equal
			fmt.Println("hit1")

			newPt = Point{c1Pt.X, MaxDec(c1Pt.Y, c2Pt.Y)} // TODO don't use MAX (only applies to circle)
			c1I++
			c2I++
		case c1Pt.X.LT(c2Pt.X): // pt1 > pt2
			fmt.Println("hit2")
			c2Interpolated := c2.PointWithX(c2I, c1Pt.X)
			newPt = Point{c1Pt.X, MaxDec(c1Pt.Y, c2Interpolated.Y)}
			c1I++
		case c2Pt.X.LT(c1Pt.X): // pt1 > pt2
			fmt.Println("hit3")
			c1Interpolated := c1.PointWithX(c1I, c2Pt.X)
			newPt = Point{c2Pt.X, MaxDec(c2Pt.Y, c1Interpolated.Y)}
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
