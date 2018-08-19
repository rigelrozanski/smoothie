package main

import (
	"fmt"
)

// nolint
const Precision = 20

// evaluation function for a circle
func Fn(x Dec) (y Dec) {
	inter1 := x.Mul(x)
	inter2 := NewDec(1).Sub(inter1)
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

func regularDivision(divisions int64, XBoundMax Dec) map[int64]Line {

	// create boring polygon
	regularDivision := make(map[int64]Line)

	startPoint := Point{NewDec(0), NewDec(1)} // top of the circle
	width := XBoundMax.Quo(NewDec(divisions)) // width of all these pieces

	for side := int64(0); side < divisions; side++ {
		x2 := startPoint.X.Add(width)
		if x2.GT(XBoundMax) { // precision correction
			x2 = XBoundMax
		}
		y2 := Fn(x2)
		endPoint := Point{x2, y2}
		regularDivision[side] = NewLine(startPoint, endPoint)
		startPoint = endPoint
	}
	return regularDivision
}

func main() {

	// nolint- bounds for a quarter of the circle
	XBoundMax := NewDec(1)

	// starting superset
	supersetPolygon := regularDivision(3, XBoundMax)

	for divisions := int64(4); true; divisions++ {

		// polygon to add to the construction of the superset polygon
		subsetPolygon := regularDivision(divisions, XBoundMax)
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

		oldSupersetLength := lenLines(newSupersetPolygon)
		supersetLength := lenLines(newSupersetPolygon)
		subsetLength := lenLines(subsetPolygon)
		output := "---------------------------------------------------------------\n"
		output += fmt.Sprintf("Subset: %v\t\t\tlength %v\nSuperset\t# points %v,\tlength %v\n",
			divisions, subsetLength.String(), len(newSupersetPolygon), supersetLength.String())
		fmt.Println(output)

		// sanity
		if (supersetLength).LT(subsetLength) {
			msg := fmt.Sprintf("subset > superset length!\n subset:\n%v\nold superset:\n%v\nsuperset:\n%v\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}
		if (supersetLength).LT(oldSupersetLength) {
			msg := fmt.Sprintf("old superset > superset length!\n subset:\n%v\nold superset:\n%v\nsuperset:\n%v\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}

		supersetPolygon = newSupersetPolygon
	}
}
