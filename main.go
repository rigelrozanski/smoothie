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

var startOrder = int64(2) // number of divisions in first curve
var numberOfOffsets = int64(2)

func circleFn(x Dec) (y Dec) {
	inter1 := x.Mul(x)
	inter2 := OneDec().Sub(inter1)
	inter3 := inter2.Sqrt()
	return inter3
}

func main() {

	xBoundMax := OneDec()
	startPt := Point{ZeroDec(), OneDec()} // top of the circle

	// PHASE 1: construct the non-offset superset of different order curves
	fmt.Println("---------------------------------------------PHASE-1----------------------------------------------------")
	superset := NewRegularCurve(startOrder, startPt, xBoundMax, circleFn)
	order := startOrder + 1
	for ; order < startOrder*2; order++ {

		// get the superset curve of two curves
		subset := NewRegularCurve(order, startPt, xBoundMax, circleFn)
		newSuperset, supersetLength, supersetArea, subsetLength, subsetArea, oldSupersetLength, oldSupersetArea, err := SupersetCurve(subset, superset, circleFn)

		// make relative to pi, just for printing :)
		supersetLength, supersetArea = two.Mul(supersetLength), four.Mul(supersetArea)
		subsetLength, subsetArea = two.Mul(subsetLength), four.Mul(subsetArea)
		oldSupersetLength, oldSupersetArea = two.Mul(oldSupersetLength), four.Mul(oldSupersetArea)

		output := "---------------------------------------------------------------\n"
		output += fmt.Sprintf("PHASE-1 Subset: %v\t\tlength %v\tarea %v\nSuperset\t# points %v,\tlength %v\tarea %v\n",
			order, subsetLength.String(), subsetArea.String(),
			len(superset), supersetLength.String(), supersetArea.String())
		fmt.Println(output)

		///////////////////////////////////////////////////////////////////////////////////
		// additional checks
		// NOTE the superset length can decrease in this process
		if err != nil {
			switch {
			case !(newSuperset[int64(len(newSuperset)-1)].End.X).Equal(OneDec()):
				err = errors.New("non-offset newSuperset doesn't end at {1,0}")
			case (supersetLength).LT(subsetLength):
				err = errors.New("subset > superset length")
			}
		}
		if err != nil {
			panic(fmt.Sprintf("Error: %v\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				err, subset.String(), superset.String(), newSuperset.String()))
		}

		// lastly set the new superset curve and continue
		superset = newSuperset
	}

	// PHASE 2 - offset the superset curve
	fmt.Println("---------------------------------------------PHASE-2----------------------------------------------------")
	fmt.Printf("debug startOrder: %v\n", startOrder)
	finalOrder := startOrder*2 - 1
	//maxOffset := xBoundMax.Quo(NewDec(startOrder)) //////////////////////////////////////////////////////////////////////////////// TODO this one actually correct but need to mirror
	maxOffset := xBoundMax.Quo(NewDec(finalOrder))
	phase1Superset := superset
	for offsetI := int64(1); offsetI <= numberOfOffsets; offsetI++ {
		output := "---------------------------------------------------------------\n"
		offsetWidth := (maxOffset.Mul(NewDec(offsetI))).Quo(NewDec(numberOfOffsets))
		fmt.Printf("debug maxOffset: %v\n", maxOffset.String())
		fmt.Printf("debug offsetWidth: %v\n", offsetWidth.String())
		offset := phase1Superset.OffsetCurve(offsetWidth, zero, zero, xBoundMax, finalOrder, circleFn)

		newSuperset, supersetLength, supersetArea, offsetLength, offsetArea, oldSupersetLength, oldSupersetArea, err := SupersetCurve(superset, offset, circleFn)

		// make relative to pi, just for printing :)
		supersetLength, supersetArea = two.Mul(supersetLength), four.Mul(supersetArea)
		offsetLength, offsetArea = two.Mul(offsetLength), four.Mul(offsetArea)
		oldSupersetLength, oldSupersetArea = two.Mul(oldSupersetLength), four.Mul(oldSupersetArea)

		output += fmt.Sprintf("PHASE-2 Offset: %v\t\tlength %v\tarea %v\nSuperset\t# points %v,\tlength %v\tarea %v\n",
			offsetI, offsetLength.String(), offsetArea.String(),
			len(superset), supersetLength.String(), supersetArea.String())
		fmt.Println(output)

		//if err != nil {
		panic(fmt.Sprintf("Error: %v\n\noffset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
			err, offset.String(), superset.String(), newSuperset.String()))
		//}

		// lastly set the new superset curve and continue
		superset = newSuperset
	}
}

// primes only
//primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83,
//89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193,
//197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311,
//313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433,
//439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541}
//superset := regularDivision(primes[0], xBoundMax)
//for j := 1; j < len(primes); j++ {
//order := primes[j]
