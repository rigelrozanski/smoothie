package main

import (
	"errors"
	"fmt"
)

// curve as a bunch of lines
type Curve map[int64]Line

// 2D function which constructs the curve
type CurveFn func(x Dec) (y Dec)

func NewRegularCurve(order int64, startPoint Point, xBoundMax Dec, fn CurveFn) Curve {

	// create boring polygon
	regularCurve := make(map[int64]Line)

	for side := int64(0); side < order; side++ {
		x2 := (xBoundMax.Mul(NewDec(side + 1))).Quo(NewDec(order))
		if x2.GT(xBoundMax) || (xBoundMax.Sub(x2)).LT(precErr) { // precision correction
			x2 = xBoundMax
		}
		y2 := fn(x2)
		endPoint := Point{x2, y2}
		regularCurve[side] = NewLine(startPoint, endPoint, order)
		startPoint = endPoint
	}
	return regularCurve
}

//_________________________________________________________________________________________

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
	firstLine := NewLine(firstLineStartPt, firstLineEndPt, firstLineOrder)

	// trim the first line
	firstLine = NewLine(firstLine.PointWithX(startX), firstLineEndPt, firstLineOrder)

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

		offsetCurve[int64(i+1)] = NewLine(startPt, endPt, line.Order)
	}

	// trim the last line
	j := int64(len(offsetCurve)) - 1
	offsetCurve[j] = NewLine(offsetCurve[j].Start, offsetCurve[j].PointWithY(endY), offsetCurve[j].Order)

	return offsetCurve
}

//_________________________________________________________________________________________________________________

// get the superset curve of two curves
func SupersetCurve(c1, c2 Curve, fn CurveFn) (superset Curve,
	supersetLength, supersetArea, c1Length, c1Area, c2Length, c2Area Dec, err error) {

	superset = make(Curve)

	newSideN, c1SideN, c2SideN := int64(0), int64(0), int64(0) // side counters of the curves
	maxC1Sides, maxC2Sides := int64(len(c1)), int64(len(c2))

	tracingC1 := true // is the superset tracing curve c1 or c2
	tracing := c1[c1SideN]
	comparing := c2[c2SideN]

	// This term is to avoid auto-switching to the highest order term
	// when the two lines are only starting at the same point because they were just intercepted!
	justIntercepted := false

	for {
		if c2SideN > maxC2Sides-1 || c1SideN > maxC1Sides-1 {
			break
		}

		interceptPt, withinBounds, sameStartingPt := tracing.Intercept(comparing)
		//fmt.Printf("debug justIntercepted: %v\n", justIntercepted)
		//fmt.Printf("debug sameStartingPt: %v\n", sameStartingPt)
		//fmt.Printf("debug withinBounds: %v\n", withinBounds)
		//fmt.Printf("debug interceptPt: %v\n", interceptPt)
		//fmt.Printf("debug comparing: %v\n", comparing)
		//fmt.Printf("debug tracing: %v\n", tracing)

		doInterceptSwitch := false
		if withinBounds && !sameStartingPt {
			//fmt.Println("Hit1")
			doInterceptSwitch = true
		} else if sameStartingPt && !justIntercepted {

			//fmt.Println("Hit2")
			// get min X
			switcharoo := false
			if tracing.End.X.LTE(comparing.End.X) {
				comparePt := comparing.PointWithX(tracing.End.X)
				if comparePt.Y.GT(tracing.End.Y) {
					switcharoo = true
				}
			} else {
				tracingPt := tracing.PointWithX(comparing.End.X)
				if tracingPt.Y.LT(comparing.End.Y) {
					switcharoo = true
				}
			}
			if switcharoo {
				nextTracing := comparing
				nextComparing := tracing
				tracing = nextTracing
				comparing = nextComparing
				tracingC1 = !tracingC1
			}

			// if the trace and compare have intersecting
			// vertices always switch to the greatest number
			// of order as it will be closer the curve
			//if comparing.Order >= tracing.Order { // PROBLEM - for rotation compating order can be bigger BUT shouldn't move to higher order!
			//fmt.Println("Hit2.1")

			//// the ol' switcharoo
			//nextTracing := comparing
			//nextComparing := tracing
			//tracing = nextTracing
			//comparing = nextComparing
			//tracingC1 = !tracingC1
			//}
		}
		// else - biz as usual!

		switch {
		case doInterceptSwitch:
			//fmt.Println("Hit3")

			superset[newSideN] = NewLine(tracing.Start, interceptPt, tracing.Order)
			newSideN++

			nextTracing := NewLine(interceptPt, comparing.End, comparing.Order)
			nextComparing := NewLine(interceptPt, tracing.End, tracing.Order)
			tracing = nextTracing
			comparing = nextComparing

			tracingC1 = !tracingC1 // start tracing the other

			justIntercepted = true

		case tracingC1:
			//fmt.Println("Hit4")
			if tracing.WithinL2XBound(comparing) {
				//fmt.Println("Hit4.1")
				superset[newSideN] = tracing
				newSideN++
				c1SideN++
				tracing = c1[c1SideN]
			} else if comparing.WithinL2XBound(tracing) {
				//fmt.Println("Hit4.2")
				c2SideN++
				if c2SideN == maxC2Sides {
					superset[newSideN] = tracing
				} else {
					comparing = c2[c2SideN]
				}
			}
			justIntercepted = false

		case !tracingC1:
			//fmt.Println("Hit5")
			if tracing.WithinL2XBound(comparing) {
				//fmt.Println("Hit5.1")
				superset[newSideN] = tracing
				newSideN++
				c2SideN++
				tracing = c2[c2SideN]
			} else if comparing.WithinL2XBound(tracing) {
				//fmt.Println("Hit5.2")
				c1SideN++
				if c1SideN == maxC2Sides {
					superset[newSideN] = tracing
				} else {
					comparing = c1[c1SideN]
				}
			}
			justIntercepted = false

		default:
			panic("weird!")
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
