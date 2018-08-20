package main

import (
	"fmt"
)

// nolint
const Precision = 15

var startDivision = int64(3)

func main() {

	circleFn := func(x Dec) (y Dec) {
		inter1 := x.Mul(x)
		inter2 := OneDec().Sub(inter1)
		inter3 := inter2.Sqrt()
		return inter3
	}

	xBoundMax := OneDec()
	startPt := Point{ZeroDec(), OneDec()} // top of the circle

	// phase 1: construct the unrotated superset
	for divisions := startDivision + 1; true; divisions++ {
		supersetPolygon := NewRegularDivisionCurve(startDivision, startPt, xBoundMax, circleFn)

		// primes only
		//primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83,
		//89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193,
		//197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311,
		//313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433,
		//439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541}
		//supersetPolygon := regularDivision(primes[0], xBoundMax)
		//for j := 1; j < len(primes); j++ {
		//divisions := primes[j]

		// polygon to add to the construction of the superset polygon
		newSupersetPolygon := make(Curve)

		subsetPolygon := NewRegularDivisionCurve(divisions, startPt, xBoundMax, circleFn)
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

			doInterceptSwitch := false
			if withinBounds && !sameStartingPt {
				doInterceptSwitch = true
			} else if sameStartingPt && !tracingSubset && !justIntercepted {

				// do the ol' switcharoo
				nextTracing := comparing
				nextComparing := tracing
				tracing = nextTracing
				comparing = nextComparing
				tracingSubset = !tracingSubset
			} // otherwise continue on the subset!

			switch {
			case doInterceptSwitch:

				newSupersetPolygon[newSideN] = NewLine(tracing.Start, interceptPt, tracing.Division)
				newSideN++

				nextTracing := NewLine(interceptPt, comparing.End, comparing.Division)
				nextComparing := NewLine(interceptPt, tracing.End, tracing.Division)
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

		supersetLength, supersetArea := newSupersetPolygon.GetLengthArea()
		supersetLength, supersetArea = two.Mul(supersetLength), four.Mul(supersetArea)

		subsetLength, subsetArea := subsetPolygon.GetLengthArea()
		subsetLength, subsetArea = two.Mul(subsetLength), four.Mul(subsetArea)

		oldSupersetLength, oldSubsetArea := supersetPolygon.GetLengthArea()
		oldSupersetLength, oldSubsetArea = two.Mul(oldSupersetLength), four.Mul(oldSubsetArea)

		output := "---------------------------------------------------------------\n"
		output += fmt.Sprintf("Subset: %v\t\t\tlength %v\tarea %v\nSuperset\t# points %v,\tlength %v\tarea %v\n",
			divisions, subsetLength.String(), subsetArea.String(),
			len(newSupersetPolygon), supersetLength.String(), supersetArea.String())
		fmt.Println(output)

		///////////////////////////////////////////////////////////////////////////////////
		// SANITY
		// NOTE once in a while the oldsubset length > newsubset length - is actually correct
		insanity := ""
		switch {
		case !(newSupersetPolygon[int64(len(newSupersetPolygon)-1)].End.X).Equal(OneDec()):
			insanity = "doesn't end at {1,0}!"
		case (supersetLength).LT(subsetLength):
			insanity = "subset > superset length!"
		case (supersetArea).LT(subsetArea):
			insanity = "subset > superset area!"
		case (supersetArea).LT(oldSubsetArea):
			insanity = "old superset > superset area!"
		}
		if insanity != "" {
			insanity += fmt.Sprintf("\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				subsetPolygon.String(), supersetPolygon.String(), newSupersetPolygon.String())
			panic(insanity)
		}

		// lastly set the new superset polygon and continue
		supersetPolygon = newSupersetPolygon
	}
}
