package main

import (
	"fmt"
)

// nolint
const Precision = 100

var (
	one = OneDec()

	// nolint
	N = 11

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

func regularDivision(divisions int64) map[int64]Line {

	// create boring polygon
	regularDivision = make(map[int64]Line)

	startPoint := Point{ZeroDec(), OneDec()}  // top of the circle
	width := XBoundMax.Quo(NewDec(divisions)) // width of all these pieces

	for side := int64(0); side < divisions; side++ {
		x2 := startPoint.X.Add(width)
		if x2.GT(XBoundMax) { // precision correction
			x2 = XBoundMax
		}
		y2 := Fn(x2)
		endPoint := Point{x2, y2}
		boringPolygon[side] = NewLine(startPoint, endPoint)
		startPoint = endPoint
	}
	return regularDivision
}

func main() {

	// starting superset
	supersetPolygon := regularDivision(3)

	for divisions := 4; i < len(primes); divisions++ {

		// polygon to add to the construction of the superset polygon
		subsetPolygon := regularDivision(divisions)
		newSupersetPolygon := make(map[int64]Line)

		newSideN, subsetSideN, oldSideN := int64(0), int64(0), int64(0) // side counters of the new and old supersetPolygon
		maxSubsetSides, maxOldSides := int64(len(subsetPolygon)), int64(len(supersetPolygon))

		tracingSubset := true // is the superset tracing the addon polygon or the old superset
		tracing := subsetPolygon[subsetSideN]
		comparing := supersetPolygon[oldSideN]

		for {
			if oldSideN > maxOldSides-1 || subsetSideN > maxSubsetSides-1 {
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

				tracingSubset = !tracingSubset // start tracing the other

			case tracingSubset:
				if tracing.WithinL2XBound(comparing) {
					newSupersetPolygon[newSideN] = tracing
					newSideN++
					subsetSideN++
					tracing = subsetPolygon[subsetSideN]
				} else if comparing.WithinL2XBound(tracing) {
					oldSideN++
					comparing = supersetPolygon[oldSideN]
				}
				//comparing = supersetPolygon[oldSideN]

			case !tracingSubset:
				if tracing.WithinL2XBound(comparing) {
					newSupersetPolygon[newSideN] = tracing
					newSideN++
					oldSideN++
					tracing = supersetPolygon[oldSideN]
				} else if comparing.WithinL2XBound(tracing) {
					subsetSideN++
					comparing = subsetPolygon[subsetSideN]
				}

			default:
				panic("wierd!")
			}

		}

		supersetLength := lenLines(newSupersetPolygon)
		subsetLength := lenLines(subsetPolygon)
		fmt.Printf("Subset: %v, len %v\n Superset: # points %v, length %v\n",
			divisions, subsetLength.String(), len(newSupersetPolygon), supersetLength.String())

		if (lengthSuperSet).LT(subsetLength) {
			msg := fmt.Sprintf("subset > superset length!\n subset")
			panic(msg)
		}

		supersetPolygon = newSupersetPolygon
	}
}
