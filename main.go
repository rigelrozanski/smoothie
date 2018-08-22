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
	"math"
)

// nolint
const Precision = 15

var startOrder = int64(2) //0number of vertex in first curve estimation
var numberOfOffsets = int64(2000)

func circleFn(x Dec) (y Dec) {
	inter1 := x.Mul(x)
	inter2 := OneDec().Sub(inter1)
	inter3 := inter2.Sqrt()
	return inter3
}

// return a group of primes
func SieveOfEratosthenes(value int) []int64 {

	res := make([]int64, value)
	resI := 0
	f := make([]bool, value)
	for i := 2; i <= int(math.Sqrt(float64(value))); i++ {
		if f[i] == false {
			for j := i * i; j < value; j += i {
				f[j] = true
			}
		}
	}
	for i := 2; i < value; i++ {
		if f[i] == false {
			res[resI] = int64(i)
			resI++
		}
	}
	return res[:resI]
}

func main() {

	xBoundMax := OneDec()
	startPt := Point{ZeroDec(), OneDec()} // top of the circle

	// PHASE 1: construct the non-offset superset of different order curves
	fmt.Println("---------------------------------------------PHASE-1----------------------------------------------------")
	superset := NewRegularCurve(startOrder, startPt, xBoundMax, circleFn)
	order := startOrder + 1
	//for ; order < startOrder*2; order++ {
	for ; true; order++ {
		//break // no supersets only rotations

		//primes only
		//primes := SieveOfEratosthenes(300)
		//superset := NewRegularCurve(primes[0], startPt, xBoundMax, circleFn)
		//for j := 1; j < len(primes); j++ {
		//order := primes[j]
		//fmt.Printf("debug order: %v\n", order)

		// get the superset curve of two curves
		subset := NewRegularCurve(order, startPt, xBoundMax, circleFn)
		newSuperset, supersetLength, supersetArea, subsetLength, subsetArea, oldSupersetLength, oldSupersetArea, err := SupersetCurve(subset, superset, circleFn)

		// make relative to pi, just for printing :)
		supersetLength, supersetArea = two.Mul(supersetLength), four.Mul(supersetArea)
		subsetLength, subsetArea = two.Mul(subsetLength), four.Mul(subsetArea)
		oldSupersetLength, oldSupersetArea = two.Mul(oldSupersetLength), four.Mul(oldSupersetArea)

		output := "--------------------PHASE-1-------------------------------------------\n"
		output += fmt.Sprintf("Subset: # points %v\t\tlength %v\tarea %v\nSuperset\t# points %v,\tlength %v\tarea %v\n",
			len(subset), subsetLength.String(), subsetArea.String(),
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

		//fmt.Printf("\noffset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
		//subset.String(), superset.String(), newSuperset.String())

		if err != nil {
			panic(fmt.Sprintf("Error: %v\n\nsubset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				err, subset.String(), superset.String(), newSuperset.String()))
		}

		// lastly set the new superset curve and continue
		superset = newSuperset

	}

	return

	// PHASE 2 - offset the superset curve
	fmt.Println("---------------------------------------------PHASE-2----------------------------------------------------")
	//fmt.Printf("debug startOrder: %v\n", startOrder)
	maxOffset := xBoundMax.Quo(NewDec(startOrder))
	//finalOrder := startOrder * 2 * 2
	//maxOffset := xBoundMax.Quo(NewDec(finalOrder))
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
			return
			panic(fmt.Sprintf("\nError: %v\n\noffset =Line[\n%v];\noldsuperset =Line[\n%v];\nsuperset =Line[\n%v];\n",
				err, offset.String(), superset.String(), newSuperset.String()))
		}

		// lastly set the new superset curve and continue
		superset = newSuperset
	}
}
