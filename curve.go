package main

import (
	"errors"
	"fmt"
)

// a classic "point"
type Point struct {
	X, Y Dec
}

func (p Point) String() string {
	return fmt.Sprintf("{%v, %v}", p.X.String(), p.Y.String())
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
		x := xAxisForwardShift
		m, b := GetMB(c[i-1], c[i])
		reflected[i] = Point{x.Neg(), (x.Mul(m)).Add(b)}
		break
	}

	// add the reflected points to a new curve (in reverse order)
	combined := make(Curve)
	combinedI := int64(0)
	for i := int64(len(reflected)) - 1; i > 0; i-- { // skip the final reflection point
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

// add all the intercepting points to each curve
// TODO optimize this, I bet we can do everything in one step
func AddIntercepts(c1, c2 Curve) (c1Out, c2Out Curve) {

	c1Out, c2Out = make(Curve), make(Curve)
	c1OutI, c2OutI := int64(0), int64(0)

	// gather all the intercepts points to a curve
	intercepts := make(Curve)
	interceptI := int64(0)

	startC1, startC2 := c1[0], c2[0]
	c1I, c2I := int64(1), int64(1)
	endC1, endC2 := c1[c1I], c2[c2I]

	c1Out[c1OutI] = startC1
	c1OutI++
	c2Out[c2OutI] = startC2
	c2OutI++
	//fmt.Printf("debug c1: %v\n", c1.String())
	//fmt.Printf("debug c2: %v\n", c2.String())
	//fmt.Printf("debug len(c1): %v\n", len(c1))
	//fmt.Printf("debug len(c2): %v\n", len(c2))
	for {
		if c1I >= int64(len(c1)) || c2I >= int64(len(c2)) {
			break
		}
		//fmt.Printf("debug c2I: %v\n", c2I)
		//fmt.Printf("debug c1I: %v\n", c1I)

		//fmt.Printf("debug startC1: %v\n", startC1)
		//fmt.Printf("debug endC1: %v\n", endC1)
		//fmt.Printf("debug startC2: %v\n", startC2)
		//fmt.Printf("debug endC2: %v\n", endC2)
		if startC1.X.GTE(endC2.X) {
			c2I++
			startC2, endC2 = c2[c2I-1], c2[c2I]
			c2Out[c2OutI] = startC2
			c2OutI++
			continue
		}
		if startC2.X.GTE(endC1.X) {
			c1I++
			startC1, endC1 = c1[c1I-1], c1[c1I]
			c1Out[c1OutI] = startC1
			c1OutI++
			continue
		}

		m1, b1 := GetMB(startC1, endC1)
		m2, b2 := GetMB(startC2, endC2)

		//  y  = (b2 m1 - b1 m2)/(m1 - m2)
		num := (b2.Mul(m1)).Sub(b1.Mul(m2))
		denom := m1.Sub(m2)
		if !denom.Equal(zero) {
			y := num.Quo(denom)
			var intercept Point
			if !m1.Equal(zero) {
				intercept = Point{(y.Sub(b1)).Quo(m1), y}
			} else {
				intercept = Point{(y.Sub(b2)).Quo(m2), y}
			}
			//fmt.Printf("debug early intercept: %v\n", intercept)

			lXLine, l2XLine, lYLine, l2YLine := false, false, false, false
			if ((startC1.X.Sub(endC1.X)).Abs()).LT(precErr) {
				lXLine = true
			}
			if ((startC2.X.Sub(endC2.X)).Abs()).LT(precErr) {
				l2XLine = true
			}
			if ((startC1.Y.Sub(endC1.Y)).Abs()).LT(precErr) {
				lYLine = true
			}
			if ((startC2.Y.Sub(endC2.Y)).Abs()).LT(precErr) {
				l2YLine = true
			}

			// valid range if X is bigger than the maximum startX
			//   and less than the minimum endX
			//if intercept.X.GT(MaxDec(startC1.X, startC2.X)) &&
			//intercept.X.LT(MinDec(endC1.X, endC2.X)) &&
			//!(((intercept.X.Sub(startC1.X)).Abs()).LT(precErr) ||
			//((intercept.X.Sub(endC1.X)).Abs()).LT(precErr) ||
			//((intercept.X.Sub(startC2.X)).Abs()).LT(precErr) ||
			//((intercept.X.Sub(endC2.X)).Abs()).LT(precErr)) {
			if (intercept.X.GT(MinDec(startC1.X, endC1.X)) || lXLine) &&
				(intercept.X.LT(MaxDec(startC1.X, endC1.X)) || lXLine) &&
				(intercept.X.GT(MinDec(startC2.X, endC2.X)) || l2XLine) &&
				(intercept.X.LT(MaxDec(startC2.X, endC2.X)) || l2XLine) &&
				(intercept.Y.LT(MaxDec(startC1.Y, endC1.Y)) || lYLine) &&
				(intercept.Y.GT(MinDec(startC1.Y, endC1.Y)) || lYLine) &&
				(intercept.Y.LT(MaxDec(startC2.Y, endC2.Y)) || l2YLine) &&
				(intercept.Y.GT(MinDec(startC2.Y, endC2.Y)) || l2YLine) {

				intercepts[interceptI] = intercept
				interceptI++
				c1Out[c1OutI] = intercept
				c1OutI++
				c2Out[c2OutI] = intercept
				c2OutI++
				//fmt.Printf("debug intercept: %v\n", intercept)
			}
		}

		if endC1.X.LT(endC2.X) {
			c1I++
			startC1, endC1 = c1[c1I-1], c1[c1I]
			c1Out[c1OutI] = startC1
			c1OutI++
		} else {
			c2I++
			startC2, endC2 = c2[c2I-1], c2[c2I]
			c2Out[c2OutI] = startC2
			c2OutI++
		}
	}

	// tack on any leftover end points
	// BUG - check for duplicates
	if c1I == int64(len(c1))-1 {
		c1Out[c1OutI] = c1[int64(len(c1))-1]
	}
	if c2I == int64(len(c2))-1 {
		c2Out[c2OutI] = c2[int64(len(c2))-1]
	}

	return c1Out, c2Out

	//// add all those intercepts to each curve
	//maxInterceptI := int64(len(intercepts))
	////fmt.Printf("debug maxInterceptI: %v\n", maxInterceptI)

	//// curve 1
	//interceptI = int64(0)
	//c1I = int64(0)
	////fmt.Printf("debug c1: %v\n", c1.String())
	//for {
	//if interceptI >= maxInterceptI {
	//break
	//}

	//intercept := intercepts[interceptI]
	//cPt := c1[c1I]
	//if intercept.X.LT(cPt.X) {
	//c1Out[c1OutI] = intercept
	//interceptI++
	//} else {
	//c1Out[c1OutI] = cPt
	//c1I++
	//}
	//c1OutI++
	//}
	//for ; c1I < int64(len(c1)); c1I++ { // any remaining points
	//c1Out[c1OutI] = c1[c1I]
	//c1OutI++
	//}
	////fmt.Printf("debug len(c1Out): %v\n", len(c1Out))
	////fmt.Printf("debug c1Out: %v\n", c1Out)

	//// curve 2
	//interceptI = int64(0)
	//c2I = int64(0)
	////fmt.Printf("debug c2: %v\n", c2.String())
	//for {
	//if interceptI >= maxInterceptI {
	//break
	//}

	//intercept := intercepts[interceptI]
	//cPt := c2[c2I]
	//if intercept.X.LT(cPt.X) {
	//c2Out[c2OutI] = intercept
	//interceptI++
	//} else {
	//c2Out[c2OutI] = cPt
	//c2I++
	//}
	//c2OutI++
	//}
	//for ; c2I < int64(len(c2)); c2I++ { // any remaining points
	//c2Out[c2OutI] = c2[c2I]
	//c2OutI++
	//}
	////fmt.Printf("debug len(c2Out): %v\n", len(c2Out))
	////fmt.Printf("debug c2Out: %v\n", c2Out)

	//return c1Out, c2Out
}

// get the superset curve of two curves
func SupersetCurve(c1, c2 Curve, fn CurveFn) (superset Curve,
	supersetLength, supersetArea, c1Length, c1Area, c2Length, c2Area Dec, err error) {

	superset = make(Curve)
	c1Inter, c2Inter := AddIntercepts(c1, c2)
	//fmt.Printf("\ndebug c2: %v\n", c2.String())
	//fmt.Printf("\ndebug c1: %v\n", c1.String())
	//fmt.Printf("\ndebug c1Inter: %v\n", c1Inter.String())
	//fmt.Printf("\ndebug c2Inter: %v\n", c2Inter.String())
	if len(c1Inter) < len(c1) {
		panic("why 1")
	}
	if len(c2Inter) < len(c2) {
		panic("why 2")
	}

	// counters for the curves
	supersetI, c1I, c2I := int64(0), int64(0), int64(0)

	for {
		var newPt Point
		c1Pt, c2Pt := c1Inter[c1I], c2Inter[c2I]
		//fmt.Printf("debug c1Pt: %v\n", c1Pt.String())
		//fmt.Printf("debug c2Pt: %v\n", c2Pt.String())

		switch {
		case ((c1Pt.X.Sub(c2Pt.X)).Abs()).LT(precErr): // equal
			//fmt.Println("hit1")

			newPt = Point{c1Pt.X, MaxDec(c1Pt.Y, c2Pt.Y)} // TODO don't use MAX (only applies to circle)
			c1I++
			c2I++
		case c1Pt.X.LT(c2Pt.X): // pt1 > pt2
			//fmt.Println("hit2")

			if c1I == int64(len(c1Inter))-1 { // if the final point just give it to c2
				newPt = c2Pt
			} else {
				c2Interpolated := c2Inter.PointWithX(c2I, c1Pt.X)
				newPt = Point{c1Pt.X, MaxDec(c1Pt.Y, c2Interpolated.Y)}
			}
			c1I++
		case c2Pt.X.LT(c1Pt.X): // pt1 > pt2
			//fmt.Println("hit3")
			if c2I == int64(len(c2Inter))-1 { // if the final point just give it to c1
				//fmt.Println("hit3.1")
				newPt = c1Pt
			} else {
				c1Interpolated := c1Inter.PointWithX(c1I, c2Pt.X)
				newPt = Point{c2Pt.X, MaxDec(c2Pt.Y, c1Interpolated.Y)}
			}
			c2I++
		default:
			panic("why")
		}

		superset[supersetI] = newPt
		supersetI++

		if c1I >= int64(len(c1Inter)) || c2I >= int64(len(c2Inter)) {
			break
		}
	}

	supersetLength, supersetArea = superset.GetLengthArea()
	c1Length, c1Area = c1Inter.GetLengthArea()
	c2Length, c2Area = c2Inter.GetLengthArea()

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
