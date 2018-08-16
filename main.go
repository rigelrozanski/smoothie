package main

import (
	"fmt"
	"math/big"
)

// nolint
type (
	NPolygon int64 // sides of the polygon
	NSide    int64 // which side on a polyon
)

var one = big.NewFloat(1)

// evaluation function
func Fn(x *big.Float) (y *big.Float) {
	return y.Sqrt(new(*big.Float).Sub(one, new(*big.Float).Pow(x, 2)))
}

func main() {

	var n int64 = 8                     // starting number of sides
	maxN := 2*n - 1                     // maximum number of sides for boring polygons
	coolPolygon := make(map[NSide]Line) // cool polygon to construct out of boring ones

	// init boring polygons
	for i := n; i <= maxN; i++ {
		boringPolygons := make(map[NPolygon]map[NSide]Line)
		boringPolygons[i] = make(map[NSide]Line)
	}

	// first create a group of all the boring polygons

	lines[1][2] = "hoot"
	fmt.Println(lines[1][2])
}
