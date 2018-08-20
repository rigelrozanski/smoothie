package main

var zero, precErr, two, four Dec

func init() {
	precErr = NewDecWithPrec(2, Precision) // XXX NEED A BETTER WAY OF DEALING WITH PRECISION LOSSES - maybe switch to big rational
	zero = ZeroDec()
	two = NewDec(2)
	four = NewDec(4)
}

// a classic "point"
type Point struct {
	X, Y Dec
}

// boring ol' straight line
type Line struct {
	Start, End Point
	M, B       Dec   //  from y = mx +b
	Order      int64 // source order of this line
}

func NewLine(start, end Point, order int64) Line {
	denom := end.X.Sub(start.X)
	if denom.Equal(zero) {
		return Line{start, end, ZeroDec(), start.Y, order}
	}
	m := (end.Y.Sub(start.Y)).Quo(end.X.Sub(start.X))
	b := start.Y.Sub(m.Mul(start.X))
	return Line{start, end, m, b, order}
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

// y = mx +b
func (l Line) PointWithX(x Dec) Point {
	return Point{x, (x.Mul(l.M)).Add(l.B)}
}

// x = (y - b)/m
func (l Line) PointWithY(y Dec) Point {
	if l.M.Equal(zero) {
		return Point{l.Start.X, y}
	}
	return Point{(y.Sub(l.B)).Quo(l.M), y}
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
	//return MinDec(l.End.X, l.Start.X).LTE(MaxDec(l2.End.X, l2.Start.X)) // iXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX THIS GUY IS BAD
	if l2.End.X.LT(l2.Start.X) {
		return l.End.X.LTE(l2.Start.X) // towards the end of the polygon this can happen
	}
	//fmt.Printf("debug l.End.X: %v\n", l.End.X)
	//fmt.Printf("debug l2.End.X: %v\n", l2.End.X)
	//fmt.Printf("debug l.End.X.LTE(l2.End.X): %v\n", l.End.X.LTE(l2.End.X))

	return l.End.X.LTE(l2.End.X)
}

// point at which two lines intercept,
// ... is it within bounds, do the two points start from the same origin?
func (l Line) Intercept(l2 Line) (intercept Point, withinBounds, sameStartingPt bool) {

	lXLine, l2XLine, lYLine, l2YLine := false, false, false, false
	if l.Start.X.Equal(l.End.X) {
		lXLine = true
	} else if ((l.Start.X.Sub(l.End.X)).Abs()).LT(precErr) {
		lXLine = true
	}
	if l2.Start.X.Equal(l2.End.X) {
		l2XLine = true
	} else if ((l2.Start.X.Sub(l2.End.X)).Abs()).LT(precErr) {
		l2XLine = true
	}
	if l.Start.Y.Equal(l.End.Y) {
		lYLine = true
	} else if ((l.Start.Y.Sub(l.End.Y)).Abs()).LT(precErr) {
		lYLine = true
	}
	if l2.Start.Y.Equal(l2.End.Y) {
		l2YLine = true
	} else if ((l2.Start.Y.Sub(l2.End.Y)).Abs()).LT(precErr) {
		l2YLine = true
	}

	// if start from the same intercept they cannot be intercepting going forward
	if l.Start.X.Equal(l2.Start.X) && l.Start.Y.Equal(l2.Start.Y) {
		return l.Start, false, true
	}
	if ((l.Start.X.Sub(l2.Start.X)).Abs()).LT(precErr) &&
		((l.Start.Y.Sub(l2.Start.Y)).Abs()).LT(precErr) {
		return l.Start, false, true
	}

	if lXLine {
		intercept = l2.PointWithX(l.Start.X)
	} else if l2XLine {
		intercept = l.PointWithX(l2.Start.X)
	} else if lYLine {
		intercept = l2.PointWithY(l.Start.Y)
	} else if l2YLine {
		intercept = l.PointWithY(l2.Start.Y)
	} else {

		//  y  = (b2 m1 - b1 m2)/(m1 - m2)
		num := (l2.B.Mul(l.M)).Sub(l.B.Mul(l2.M))
		denom := l.M.Sub(l2.M)
		if denom.Equal(zero) {
			return intercept, false, false
		}
		y := num.Quo(denom)
		intercept = l.PointWithY(y)
	}

	//fmt.Printf("debug intercept.X: %v\n", intercept.X)
	//fmt.Printf("debug intercept.Y: %v\n", intercept.Y)
	//fmt.Printf("debug l.Start.X: %v\n", l.Start.X)
	//fmt.Printf("debug l.End.X: %v\n", l.End.X)
	//fmt.Printf("debug l.Start.Y: %v\n", l.Start.Y)
	//fmt.Printf("debug l.End.Y: %v\n", l.End.Y)
	//fmt.Printf("debug l2.Start.X: %v\n", l2.Start.X)
	//fmt.Printf("debug l2.End.X: %v\n", l2.End.X)
	//fmt.Printf("debug l2.Start.Y: %v\n", l2.Start.Y)
	//fmt.Printf("debug l2.End.Y: %v\n", l2.End.Y)
	//fmt.Println(intercept.X.GT(MinDec(l.Start.X, l.End.X)))
	//fmt.Println(intercept.X.LT(MaxDec(l.Start.X, l.End.X)))
	//fmt.Println(intercept.X.GT(MinDec(l2.Start.X, l2.End.X)))
	//fmt.Println(intercept.X.LT(MaxDec(l2.Start.X, l2.End.X)))
	//fmt.Println(intercept.Y.LT(MaxDec(l.Start.Y, l.End.Y)))
	//fmt.Println(intercept.Y.GT(MinDec(l.Start.Y, l.End.Y)))
	//fmt.Println(intercept.Y.LT(MaxDec(l2.Start.Y, l2.End.Y)))
	//fmt.Println(intercept.Y.GT(MinDec(l2.Start.Y, l2.End.Y)))

	withinBounds = false
	if (intercept.X.GT(MinDec(l.Start.X, l.End.X)) || lXLine) &&
		(intercept.X.LT(MaxDec(l.Start.X, l.End.X)) || lXLine) &&
		(intercept.X.GT(MinDec(l2.Start.X, l2.End.X)) || l2XLine) &&
		(intercept.X.LT(MaxDec(l2.Start.X, l2.End.X)) || l2XLine) &&
		(intercept.Y.LT(MaxDec(l.Start.Y, l.End.Y)) || lYLine) &&
		(intercept.Y.GT(MinDec(l.Start.Y, l.End.Y)) || lYLine) &&
		(intercept.Y.LT(MaxDec(l2.Start.Y, l2.End.Y)) || l2YLine) && // problematic if straight line
		(intercept.Y.GT(MinDec(l2.Start.Y, l2.End.Y)) || l2YLine) {
		withinBounds = true
	}

	//withinBounds = false
	//if intercept.X.GT(MinDec(MinDec(l.Start.X, l.End.X), MinDec(l2.Start.X, l2.End.X))) &&
	//intercept.X.LT(MaxDec(MaxDec(l.Start.X, l.End.X), MaxDec(l2.Start.X, l2.End.X))) &&
	//intercept.Y.LT(MaxDec(MaxDec(l.Start.Y, l.End.Y), MaxDec(l2.Start.Y, l2.End.Y))) &&
	//intercept.Y.GT(MinDec(MinDec(l.Start.Y, l.End.Y), MinDec(l2.Start.Y, l2.End.Y))) {
	//withinBounds = true
	//}

	return intercept, withinBounds, false
}
