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

	// if start from the same intercept they cannot be intercepting going forward
	if l.Start.X.Equal(l2.Start.X) && l.Start.Y.Equal(l2.Start.Y) {
		return intercept, false
	}

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

//primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541}
