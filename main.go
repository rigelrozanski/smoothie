package main

import (
	"fmt"
	"math/big"
)

// bounds for the calculation
var (
	one       = big.NewFloat(1)
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
		totalLen = totalLen.Add(totalLen, line.Length())
	}
	return totalLen
}

func formattedLines(lines map[int64]Line) string {
	out := "{"
	out += fmt.Sprintf("{%v, %v}", lines[0].Start.X, lines[0].Start.Y)
	for _, line := range lines {
		out += fmt.Sprintf(",{%v, %v}", line.End.X, line.End.Y)
	}
	out += "}"
	return out
}

func main() {

	var n int64 = 3 // starting number of divisions
	maxN := 2*n - 1 // maximum number of sides for boring polygons

	boringPolygons := make(map[int64]map[int64]Line) // index 1: number of lines in element, index 2: element line no.

	// init boring polygons
	for i := n; i <= maxN; i++ {
		boringPolygons[i] = make(map[int64]Line)

		startPoint := Point{big.NewFloat(0), big.NewFloat(1)}            // top of the circle
		width := new(big.Float).Quo(XBoundMax, big.NewFloat(float64(i))) // width of all these pieces

		for side := int64(0); side < i; side++ {
			x2 := new(big.Float).Add(startPoint.X, width)
			if x2.Cmp(XBoundMax) > 0 { // precision correction
				x2 = XBoundMax
			}
			y2 := Fn(x2)
			endPoint := Point{x2, y2}
			boringPolygons[i][side] = NewLine(startPoint, endPoint)
			startPoint = endPoint
		}

		fmt.Printf("line %v, length %v\nformatted: %v\n", i,
			lenLines(boringPolygons[i]), formattedLines(boringPolygons[i]))
	}

	// construct the superset polygon
	supersetPolygon := boringPolygons[n] // start with the smallest
	for i := n + 1; i <= maxN; i++ {

		// polygon to add to the construction of the superset polygon
		addonPolygon := boringPolygons[i]
		newSupersetPolygon = make(map[int64]Line)

		newSideN, addonSideN, oldSideN := 0, 0, 0 // side counters of the new and old supersetPolygon
		maxAddonSides, maxOldSides := len(addonPolygon), len(maxOldSides)

		tracingAddon := true // is the superset tracing the addon polygon or the old superset

		for {
			addonSide := addonPolygon[addonSideN]
			oldSide := supersetPolygon[oldSideN]

			if addonSide.WillIntercept(oldSide) {
				interceptPt := addonSide.Intercept(oldSide)
				newSupersetPolygon[newSideN] = getLineSplit(tracingAddon, interceptPt, addonSide, oldSide)
				tracingAddon := !tracingAddon // start tracing the other
			} else if side.WithinL2YBound {

				newSideN++
				continue
			}
		}

		supersetPolygon = newSupersetPolygon
	}
}

func getLineSplit(tracingAddon bool, intercept Point, addonSide, oldSide Line) Line {
	if tracingAddon {
		return NewLine(addonSide.Start, intecept)
	}
	return NewLine(oldSide, intercept)
}
