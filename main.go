package main

import (
	"fmt"
	"math/big"
)

var one = big.NewFloat(1)

// bounds for the calculation
var (
	XBoundMin = big.NewFloat(0)
	XBoundMax = one
)

// evaluation function
func Fn(x *big.Float) (y *big.Float) {
	inter1 := new(big.Float).Mul(x, x)
	inter2 := inter1.Sub(one, inter1)
	inter3 := inter2.Sqrt(inter2)
	return inter3
}

func lenLines(lines map[int64]Line) *big.Float {
	totalLen := big.NewFloat(0)
	for _, line := range lines {
		fmt.Printf("debug line: %v\n", line)
		fmt.Printf("debug line.Length: %v\n", line.Length())
		totalLen = totalLen.Add(totalLen, line.Length())
	}
	return totalLen
}

func main() {

	var n int64 = 6 // starting number of divisions
	maxN := n       //2*n - 1 // maximum number of sides for boring polygons
	//coolPolygon := make(map[int64]Line) // cool polygon to construct out of boring ones

	boringPolygons := make(map[int64]map[int64]Line) // index 1: number of lines in element, index 2: element line no.

	// init boring polygons
	for i := int64(n); i <= maxN; i++ {
		boringPolygons[i] = make(map[int64]Line)

		startPoint := Point{big.NewFloat(0), big.NewFloat(1)}            // top of the circle
		width := new(big.Float).Quo(XBoundMax, big.NewFloat(float64(i))) // width of all these pieces

		for side := int64(0); side < i; side++ {
			x2 := new(big.Float).Add(startPoint.X, width)
			y2 := Fn(x2)
			endPoint := Point{x2, y2}
			boringPolygons[i][side] = NewLine(startPoint, endPoint)
			startPoint = endPoint
		}

		fmt.Printf("line %v, length %v \n", i, lenLines(boringPolygons[i]))
	}

	fmt.Println(boringPolygons)
}
