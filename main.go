package main

import (
	"fmt"
)

// nolint
const Precision = 5000

var (
	one = OneDec()

	// nolint
	N          = 11
	precCutoff = 60
	primes     = []int64{3, 4, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541}

	// nolint- bounds for a quarter of the circle
	XBoundMin = ZeroDec()
	XBoundMax = one
)

// evaluation function for a circle
func Fn(x Dec) (y Dec) {
	inter1 := x.Mul(x)
	inter2 := one.Sub(inter1)
	inter3 := inter2.Sqrt()
	return inter3
}

//____________________________________

func lenLines(lines map[int64]Line) Dec {
	totalLen := ZeroDec()
	for _, line := range lines {
		totalLen = totalLen.Add(line.Length())
	}
	return totalLen
}

func formattedLines(lines map[int64]Line) string {
	out := "{"
	out += fmt.Sprintf("{%v, %v}", lines[0].Start.X, lines[0].Start.Y)

	for i := int64(0); i < int64(len(lines)); i++ {
		line := lines[i]
		out += fmt.Sprintf(",{%v, %v}", line.End.X.String(), line.End.Y.String())
	}
	out += "}"
	return out
}

func main() {

	n := int64(N)   // starting number of divisions
	maxN := 2*n - 1 // maximum number of sides for boring polygons

	boringPolygons := make(map[int64]map[int64]Line) // index 1: number of lines in element, index 2: element line no.

	// init boring polygons
	for i := n; i <= maxN; i++ {
		val := primes[i]
		boringPolygons[i] = make(map[int64]Line)

		startPoint := Point{ZeroDec(), OneDec()} // top of the circle
		width := XBoundMax.Quo(NewDec(i))        // width of all these pieces

		for side := int64(0); side < i; side++ {
			x2 := startPoint.X.Add(width)
			if x2.GT(XBoundMax) { // precision correction
				x2 = XBoundMax
			}
			y2 := Fn(x2)
			endPoint := Point{x2, y2}
			boringPolygons[i][side] = NewLine(startPoint, endPoint)
			startPoint = endPoint
		}

		fmt.Printf("line %v, length %v\n", i, lenLines(boringPolygons[i]).String())
		//fmt.Printf("formatted: %v\n", formattedLines(boringPolygons[i]))
	}

	// construct the superset polygon
	supersetPolygon := boringPolygons[n] // start with the smallest
	for i := n + 1; i <= maxN; i++ {

		// polygon to add to the construction of the superset polygon
		addonPolygon := boringPolygons[i]
		newSupersetPolygon := make(map[int64]Line)

		newSideN, addonSideN, oldSideN := int64(0), int64(0), int64(0) // side counters of the new and old supersetPolygon
		maxAddonSides, maxOldSides := int64(len(addonPolygon)), int64(len(supersetPolygon))

		tracingAddon := true // is the superset tracing the addon polygon or the old superset
		tracing := addonPolygon[addonSideN]
		comparing := supersetPolygon[oldSideN]

		for {
			if oldSideN > maxOldSides-1 || addonSideN > maxAddonSides-1 {
				break
			}

			var withinBounds bool
			var interceptPt Point
			interceptPt, withinBounds = tracing.Intercept(comparing)

			switch {
			case withinBounds:

				newSupersetPolygon[newSideN] = NewLine(tracing.Start, interceptPt)
				newSideN++

				nextTracing := NewLine(interceptPt, comparing.End)
				nextComparing := NewLine(interceptPt, tracing.End)
				tracing = nextTracing
				comparing = nextComparing

				tracingAddon = !tracingAddon // start tracing the other

			case tracingAddon:
				if tracing.WithinL2XBound(comparing) {
					newSupersetPolygon[newSideN] = tracing
					newSideN++
					addonSideN++
					tracing = addonPolygon[addonSideN]
				} else if comparing.WithinL2XBound(tracing) {
					oldSideN++
					comparing = supersetPolygon[oldSideN]
				}
				//comparing = supersetPolygon[oldSideN]

			case !tracingAddon:
				if tracing.WithinL2XBound(comparing) {
					newSupersetPolygon[newSideN] = tracing
					newSideN++
					oldSideN++
					tracing = supersetPolygon[oldSideN]
				} else if comparing.WithinL2XBound(tracing) {
					addonSideN++
					comparing = addonPolygon[addonSideN]
				}

			default:
				panic("wierd!")
			}

		}

		supersetPolygon = newSupersetPolygon
	}

	fmt.Printf("superset polygon, num points %v, length %v\n", len(supersetPolygon),
		lenLines(supersetPolygon).String())
	//fmt.Printf("formatted: %v\n", formattedLines(supersetPolygon))
}
