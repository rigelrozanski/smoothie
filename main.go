package main

import (
	"fmt"
)

// nolint
const Precision = 15

// evaluation function for a circle
func Fn(x Dec) (y Dec) {
	inter1 := x.Mul(x)
	inter2 := OneDec().Sub(inter1)
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

	startPoint := Point{ZeroDec(), OneDec()}  // top of the circle
	width := XBoundMax.Quo(NewDec(divisions)) // width of all these pieces

	for side := int64(0); side < divisions; side++ {
		x2 := startPoint.X.Add(width)
		if x2.GT(XBoundMax) || (XBoundMax.Sub(x2)).LT(precErr) { // precision correction
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
	XBoundMax := OneDec()

	// starting superset
	startDivision := int64(3)
	supersetPolygon := regularDivision(startDivision, XBoundMax)

	for divisions := startDivision + 1; true; divisions++ {

		// polygon to add to the construction of the superset polygon
		subsetPolygon := regularDivision(divisions, XBoundMax)
		newSupersetPolygon := make(map[int64]Line)

		newSideN, subsetSideN, oldSideN := int64(0), int64(0), int64(0) // side counters of the new and old supersetPolygon
		maxSubsetSides, maxOldSides := int64(len(subsetPolygon)), int64(len(supersetPolygon))

		tracingSubset := true // is the superset tracing the addon polygon or the old superset
		tracing := subsetPolygon[subsetSideN]
		comparing := supersetPolygon[oldSideN]

		justIntercepted := false

		for {
			if oldSideN > maxOldSides-1 || subsetSideN > maxSubsetSides-1 {
				break
			}

			interceptPt, withinBounds, sameStartingPt := tracing.Intercept(comparing)
			//fmt.Printf("debug interceptPt: %v withinBounds %v, samestart %v\n", interceptPt, withinBounds, sameStartingPt)

			doInterceptSwitch := false
			if withinBounds && !sameStartingPt {
				doInterceptSwitch = true
			} else if sameStartingPt && !tracingSubset && !justIntercepted { /////////////////////////////////////////// XXX but not nessisarily!

				// pull the ol' switcharoo
				nextTracing := comparing
				nextComparing := tracing
				tracing = nextTracing
				comparing = nextComparing
				tracingSubset = !tracingSubset
			} // otherwise continue on the subset!

			switch {
			case doInterceptSwitch:

				newSupersetPolygon[newSideN] = NewLine(tracing.Start, interceptPt)
				newSideN++

				nextTracing := NewLine(interceptPt, comparing.End)
				nextComparing := NewLine(interceptPt, tracing.End)
				tracing = nextTracing
				comparing = nextComparing

				tracingSubset = !tracingSubset // start tracing the other

				justIntercepted = true

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
				justIntercepted = false

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
				justIntercepted = false

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
			msg := fmt.Sprintf("subset > superset length!\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}
		if (supersetLength).LT(oldSupersetLength) {
			msg := fmt.Sprintf("old superset > superset length!\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}

		supersetPolygon = newSupersetPolygon
	}
}
