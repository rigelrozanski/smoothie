package main

// a classic "point"
type Point struct {
	X, Y Dec
}

// boring ol' straight line
type Line struct {
	Start, End Point
	M, B       Dec //  from y = mx +b
}

func NewLine(start, end Point) Line {
	m := (end.Y.Sub(start.Y)).Quo(end.X.Sub(start.X))
	b := start.Y.Sub(m.Mul(start.X))
	return Line{start, end, m, b}
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

// y-axis end of line l is within end of l2
func (l Line) WithinL2XBound(l2 Line) bool {
	return l.End.X.LTE(l2.End.X)
}

var zero = ZeroDec()

// point at which two lines intercept
func (l Line) Intercept(l2 Line) (intercept Point, withinBounds bool) {
	//  y  = (b2 m1 - b1 m2)/(m1 - m2)
	inter1 := l2.B.Mul(l.M)

	inter2 := l.B.Mul(l2.M)
	inter3 := inter1.Sub(inter2)
	inter4 := l.M.Sub(l2.M)
	if inter4.Equal(zero) {
		return intercept, false
	}
	y := inter3.Quo(inter4)
	x := (y.Sub(l.B)).Quo(l.M)
	intercept = Point{x, y}

	// check if intercept is in precision error amount
	proximityToZero := l.Start.X.Sub(intercept.X).Abs()
	if proximityToZero.LT(NewDecWithPrec(1, int64(precCutoff))) {
		return intercept, false
	}

	withinBounds = false
	// to deal with precision errors,
	//  if any of the points are equal we know we're not within bounds
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

	return intercept, withinBounds
}
