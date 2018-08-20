package main

import (
	"errors"
	"fmt"
)

// curve as a bunch of lines
type Curve map[int64]Line

// 2D function which constructs the curve
type CurveFn func(x Dec) (y Dec)

func NewRegularDivisionCurve(divisions int64, startPoint Point, xBoundMax Dec, fn CurveFn) Curve {

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
// CONTRACT - do not rotate more than the first line division width
func (c Curve) ShiftAlongX(xAxisForwardShift, startX, endY, xBoundMax Dec, firstLineDivision int64, fn CurveFn) Curve {

	// construct the first by working backwards from the first shifted point
	firstLineWidth := xBoundMax.Quo(NewDec(firstLineDivision))
	firstLineStartX := xAxisForwardShift.Sub(firstLineWidth) // should be negative
	if firstLineStartX.GT(zero) {
		msg := fmt.Sprintf("bad shift, cannot shift more than first line width\n\tfirstLineWidth\t%v\n\tfirstLineStartX\t%v\n",
			firstLineWidth, firstLineStartX)
		panic(msg)
	}
	firstLineStartPt := Point{firstLineStartX, fn(firstLineStartX)}

	firstLineEndX := c[0].Start.X.Add(xAxisForwardShift)
	firstLineEndPt := Point{firstLineEndX, fn(firstLineEndX)}
	firstLine := NewLine(firstLineStartPt, firstLineEndPt, firstLineDivision)

	// trim the first line
	firstLine = NewLine(firstLine.PointWithX(startX), firstLineEndPt, firstLineDivision)

	shiftedCurve := make(map[int64]Line)
	shiftedCurve[0] = firstLine
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

		shiftedCurve[int64(i+1)] = NewLine(startPt, endPt, line.Division)
	}

	// trim the last line
	j := int64(len(shiftedCurve)) - 1
	shiftedCurve[j] = NewLine(shiftedCurve[j].Start, shiftedCurve[j].PointWithY(endY), shiftedCurve[j].Division)

	return shiftedCurve
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

	// This term is to avoid auto-switching to the highest division term
	// when the two lines are only starting at the same point because they were just intercepted!
	justIntercepted := false

	for {
		if c2SideN > maxC2Sides-1 || c1SideN > maxC1Sides-1 {
			break
		}

		interceptPt, withinBounds, sameStartingPt := tracing.Intercept(comparing)

		doInterceptSwitch := false
		if withinBounds && !sameStartingPt {
			doInterceptSwitch = true
		} else if sameStartingPt && !justIntercepted {

			// if the trace and compare have intersecting
			// vertices always switch to the greatest number
			// of divisions as it will be closer the curve
			if comparing.Division > tracing.Division {

				// the ol' switcharoo
				nextTracing := comparing
				nextComparing := tracing
				tracing = nextTracing
				comparing = nextComparing
				tracingC1 = !tracingC1
			}
		}
		// else - biz as usual!

		switch {
		case doInterceptSwitch:

			superset[newSideN] = NewLine(tracing.Start, interceptPt, tracing.Division)
			newSideN++

			nextTracing := NewLine(interceptPt, comparing.End, comparing.Division)
			nextComparing := NewLine(interceptPt, tracing.End, tracing.Division)
			tracing = nextTracing
			comparing = nextComparing

			tracingC1 = !tracingC1 // start tracing the other

			justIntercepted = true

		case tracingC1:
			if tracing.WithinL2XBound(comparing) {
				superset[newSideN] = tracing
				newSideN++
				c1SideN++
				tracing = c1[c1SideN]
			} else if comparing.WithinL2XBound(tracing) {
				c2SideN++
				comparing = c2[c2SideN]
			}
			justIntercepted = false

		case !tracingC1:
			if tracing.WithinL2XBound(comparing) {
				superset[newSideN] = tracing
				newSideN++
				c2SideN++
				tracing = c2[c2SideN]
			} else if comparing.WithinL2XBound(tracing) {
				c1SideN++
				comparing = c1[c1SideN]
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
