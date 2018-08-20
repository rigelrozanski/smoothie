//
//      /||||\    {this is a good smoothie}
//     |-o-o-~|  /
//    _   ~
//   /        '\
//  |    \ /   |   ___%
//  |     -    \__ \s/
//   \            ' |
//    |)      |
// ___\___      \
///____/ |  | | |
//| | || |  |_| |_
//|   |  |____]___]

package main

import (
	"errors"
	"fmt"
)

// nolint
const Precision = 15

var startDivision = int64(75)

func circleFn(x Dec) (y Dec) {
	inter1 := x.Mul(x)
	inter2 := OneDec().Sub(inter1)
	inter3 := inter2.Sqrt()
	return inter3
}

func main() {

	xBoundMax := OneDec()
	startPt := Point{ZeroDec(), OneDec()} // top of the circle

	// phase 1: construct the unrotated superset
	superset := NewRegularDivisionCurve(startDivision, startPt, xBoundMax, circleFn)
	division := startDivision + 1
	for ; true; division++ { //division < startDivision*2; division++ {

		// get the superset curve of two curves
		subset := NewRegularDivisionCurve(division, startPt, xBoundMax, circleFn)
		newSuperset, supersetLength, supersetArea, subsetLength, subsetArea, oldSupersetLength, oldSupersetArea, err := SupersetCurve(subset, superset, circleFn)

		// make relative to pi, just for printing :)
		supersetLength, supersetArea = two.Mul(supersetLength), four.Mul(supersetArea)
		subsetLength, subsetArea = two.Mul(subsetLength), four.Mul(subsetArea)
		oldSupersetLength, oldSupersetArea = two.Mul(oldSupersetLength), four.Mul(oldSupersetArea)

		output := "---------------------------------------------------------------\n"
		output += fmt.Sprintf("Subset: %v\t\t\tlength %v\tarea %v\nSuperset\t# points %v,\tlength %v\tarea %v\n",
			division, subsetLength.String(), subsetArea.String(),
			len(superset), supersetLength.String(), supersetArea.String())
		fmt.Println(output)

		///////////////////////////////////////////////////////////////////////////////////
		// additional checks
		if err != nil {
			switch {
			case !(newSuperset[int64(len(newSuperset)-1)].End.X).Equal(OneDec()):
				err = errors.New("non-shifted newSuperset doesn't end at {1,0}")
			case (supersetLength).LT(subsetLength):
				err = errors.New("subset > superset length")
			}
		}
		if err != nil {
			insanity := fmt.Sprintf("Error: %v\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				err, subset.String(), superset.String(), newSuperset.String())
			panic(insanity)
		}

		// lastly set the new superset curve and continue
		superset = newSuperset
	}

	// PHASE 2 - shift the superset curve

	//shifted := superset.ShiftAlongX(NewDecWithPrec(1, 2), zero, zero, xBoundMax, division, circleFn)

	//fmt.Printf("superset =Line[\n%v];\nshifted =Line[\n%v];\n",
	//superset.String(), shifted.String())
}

// primes only
//primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83,
//89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193,
//197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311,
//313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433,
//439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541}
//superset := regularDivision(primes[0], xBoundMax)
//for j := 1; j < len(primes); j++ {
//division := primes[j]
