package main

import (
	"fmt"
)

// nolint
const Precision = 15

var startDivision = int64(2)

// evaluation function for a circle
func Fn(x Dec) (y Dec) {
	inter1 := x.Mul(x)
	inter2 := OneDec().Sub(inter1)
	inter3 := inter2.Sqrt()
	return inter3
}

//____________________________________

// total length and area for all the lines
func LengthAreaLines(lines map[int64]Line) (length, area Dec) {
	length, area = ZeroDec(), ZeroDec()
	for _, line := range lines {
		length = length.Add(line.Length())
		area = area.Add(line.Area())
	}
	return length, area
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

	startPoint := Point{ZeroDec(), OneDec()} // top of the circle
	//width := XBoundMax.Quo(NewDec(divisions)) // width of all these pieces

	for side := int64(0); side < divisions; side++ {
		//x2 := startPoint.X.Add(width)
		x2 := (XBoundMax.Mul(NewDec(side + 1))).Quo(NewDec(divisions))
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
	supersetPolygon := regularDivision(startDivision, XBoundMax)
	for divisions := startDivision + 1; true; divisions++ {

		// primes only
		//primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83,
		//89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193,
		//197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311,
		//313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433,
		//439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541}
		//supersetPolygon := regularDivision(primes[0], XBoundMax)
		//for j := 1; j < len(primes); j++ {
		//divisions := primes[j]

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

		supersetLength, supersetArea := LengthAreaLines(newSupersetPolygon)
		supersetLength, supersetArea = two.Mul(supersetLength), four.Mul(supersetArea)

		subsetLength, subsetArea := LengthAreaLines(subsetPolygon)
		subsetLength, subsetArea = two.Mul(subsetLength), four.Mul(subsetArea)

		oldSupersetLength, oldSubsetArea := LengthAreaLines(supersetPolygon)
		oldSupersetLength, oldSubsetArea = two.Mul(oldSupersetLength), four.Mul(oldSubsetArea)

		output := "---------------------------------------------------------------\n"
		output += fmt.Sprintf("Subset: %v\t\t\tlength %v\tarea %v\nSuperset\t# points %v,\tlength %v\tarea %v\n",
			divisions, subsetLength.String(), subsetArea.String(),
			len(newSupersetPolygon), supersetLength.String(), supersetArea.String())
		fmt.Println(output)

		///////////////////////////////////////////////////////////////////////////////////
		// sanity
		//NOTE once in a while the oldsubset length > newsubset length - is actually correct
		if !(newSupersetPolygon[int64(len(newSupersetPolygon)-1)].End.X).Equal(OneDec()) {
			msg := fmt.Sprintf("doesn't end at {1,0} !\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}
		if (supersetLength).LT(subsetLength) {
			msg := fmt.Sprintf("subset > superset length!\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}
		if (supersetArea).LT(subsetArea) {
			msg := fmt.Sprintf("subset > superset Area!\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}
		if (supersetArea).LT(oldSubsetArea) {
			msg := fmt.Sprintf("old superset > superset Area!\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				formattedLines(subsetPolygon), formattedLines(supersetPolygon), formattedLines(newSupersetPolygon))
			panic(msg)
		}

		// lastly set the new superset polygon and continue
		supersetPolygon = newSupersetPolygon
	}
}
