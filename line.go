package main

// a classic "point"
type Point struct {
	X, Y Dec
}

// boring ol' straight line
type Line struct {
	Start, End Point
	M, B       Dec   //  from y = mx +b
	Division   int64 // source division order of this line
}

func NewLine(start, end Point, division int64) Line {
	m := (end.Y.Sub(start.Y)).Quo(end.X.Sub(start.X))
	b := start.Y.Sub(m.Mul(start.X))
	return Line{start, end, m, b, division}
}

//_______________________________________________________________________

// length of a line!
// sqrt((x2 - x1)^2 + (y2 - y1)^2)
func (l Line) Length() Dec {
	inter1 := l.End.X.Sub(l.Start.X)
	inter2 := l.End.Y.Sub(l.Start.Y)
	inter3 := inter1.Mul(inter1)
	inter4 := inter2.Mul(inter2)
	inter5 := inter3.Add(inter4)
	return inter5.Sqrt()
}

// Area under the line
// (x2 - x1) * (y2 + y1)/2
func (l Line) Area() Dec {
	inter1 := l.End.X.Sub(l.Start.X)
	inter2 := (l.End.Y.Add(l.Start.Y)).Quo(two)
	return inter1.Mul(inter2)
}

// y-axis end of line l is within end of l2
func (l Line) WithinL2XBound(l2 Line) bool {
	return l.End.X.LTE(l2.End.X)
}

var zero, precErr, two, four Dec

func init() {
	precErr = NewDecWithPrec(2, Precision) // XXX NEED A BETTER WAY OF DEALING WITH PRECISION LOSSES - maybe switch to big rational
	zero = ZeroDec()
	two = NewDec(2)
	four = NewDec(4)
}

// point at which two lines intercept,
// ... is it within bounds, do the two points start from the same origin?
func (l Line) Intercept(l2 Line) (intercept Point, withinBounds, sameStartingPt bool) {

	// if start from the same intercept they cannot be intercepting going forward
	if l.Start.X.Equal(l2.Start.X) && l.Start.Y.Equal(l2.Start.Y) {
		return l.Start, false, true
	}
	if ((l.Start.X.Sub(l2.Start.X)).Abs()).LT(precErr) &&
		((l.Start.Y.Sub(l2.Start.Y)).Abs()).LT(precErr) {
		return l.Start, false, true
	}

	//  y  = (b2 m1 - b1 m2)/(m1 - m2)
	num := (l2.B.Mul(l.M)).Sub(l.B.Mul(l2.M))
	denom := l.M.Sub(l2.M)
	if denom.Equal(zero) {
		return intercept, false, false
	}
	y := num.Quo(denom)
	x := (y.Sub(l.B)).Quo(l.M)
	intercept = Point{x, y}

	withinBounds = false
	if intercept.X.GT(l.Start.X) &&
		intercept.X.LT(l.End.X) &&
		intercept.X.GT(l2.Start.X) &&
		intercept.X.LT(l2.End.X) &&
		intercept.Y.LT(l.Start.Y) &&
		intercept.Y.GT(l.End.Y) &&
		intercept.Y.LT(l2.Start.Y) &&
		intercept.Y.GT(l2.End.Y) {
		withinBounds = true
	}

	return intercept, withinBounds, false
}
