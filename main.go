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

var startOrder = int64(2) // number of vertex in first curve estimation
var numberOfOffsets = startOrder

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
		//for ; true; order++ {

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
			case !(newSuperset[int64(len(newSuperset)-1)].X).Equal(OneDec()):
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
	//fmt.Printf("debug startOrder: %v\n", startOrder)
	finalOrder := startOrder*2 - 1
	//maxOffset := xBoundMax.Quo(NewDec(startOrder)) //////////////////////////////////////////////////////////////////////////////// TODO this one actually correct but need to mirror
	maxOffset := xBoundMax.Quo(NewDec(finalOrder))
	phase1Superset := superset
	for offsetI := int64(1); offsetI <= numberOfOffsets; offsetI++ {
		output := "---------------------------------------------------------------\n"
		offsetWidth := (maxOffset.Mul(NewDec(offsetI))).Quo(NewDec(numberOfOffsets))
		//fmt.Printf("debug maxOffset: %v\n", maxOffset.String())
		//fmt.Printf("debug offsetWidth: %v\n", offsetWidth.String())
		offset := phase1Superset.OffsetCurve(offsetWidth, xBoundMax, circleFn)

		//panic(fmt.Sprintf("\n\noffset =Line[\n%v];\noldsuperset =Line[\n%v];\n",
		//offset.String(), superset.String()))

		newSuperset, supersetLength, supersetArea, offsetLength, offsetArea, oldSupersetLength, oldSupersetArea, err := SupersetCurve(superset, offset, circleFn)

		// make relative to pi, just for printing :)
		supersetLength, supersetArea = two.Mul(supersetLength), four.Mul(supersetArea)
		offsetLength, offsetArea = two.Mul(offsetLength), four.Mul(offsetArea)
		oldSupersetLength, oldSupersetArea = two.Mul(oldSupersetLength), four.Mul(oldSupersetArea)

		output += fmt.Sprintf("PHASE-2 Offset: %v\t\tlength %v\tarea %v\nSuperset\t# points %v,\tlength %v\tarea %v\n",
			offsetI, offsetLength.String(), offsetArea.String(),
			len(superset), supersetLength.String(), supersetArea.String())
		fmt.Println(output)

		//fmt.Printf("\noffset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
		//offset.String(), superset.String(), newSuperset.String())

		if err != nil {
			panic(fmt.Sprintf("\nError: %v\n\noffset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				err, offset.String(), superset.String(), newSuperset.String()))
		}

		// lastly set the new superset curve and continue
		superset = newSuperset
	}
}
