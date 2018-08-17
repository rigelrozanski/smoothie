package main

import (
	"fmt"
	"math/big"
)

// bounds for the calculation
var (
	one = big.NewFloat(1)

	// bounds for a quarter of the circle
	XBoundMin = big.NewFloat(0)
	XBoundMax = one
)

// evaluation function for a circle
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

	for i := int64(0); i < int64(len(lines)); i++ {
		line := lines[i]
		out += fmt.Sprintf(",{%v, %v}", line.End.X, line.End.Y)
	}
	out += "}"
	return out
}

func main() {

	var n int64 = 3  // starting number of divisions
	maxN := int64(4) //2*n - 1 // maximum number of sides for boring polygons

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
		newSupersetPolygon := make(map[int64]Line)

		newSideN, addonSideN, oldSideN := int64(0), int64(0), int64(0) // side counters of the new and old supersetPolygon
		maxAddonSides, maxOldSides := int64(len(addonPolygon)), int64(len(supersetPolygon))

		tracingAddon := true // is the superset tracing the addon polygon or the old superset

		for {
			if (addonSideN == maxAddonSides-1 && oldSideN > maxOldSides-1) ||
				(oldSideN == maxOldSides-1 && addonSideN > maxAddonSides-1) {
				break
			}

			addonSide := addonPolygon[addonSideN]
			oldSide := supersetPolygon[oldSideN]

			fmt.Printf("debug newSideN %v, addonSideN %v, oldSideN %v\n", newSideN, addonSideN, oldSideN)
			fmt.Printf("debug addonSide: %v\n", addonSide)
			fmt.Printf("debug oldSide: %v\n", oldSide)
			//fmt.Printf("debug WithinL2YBound: %v\n", addonSide.WithinL2XBound(oldSide))

			interceptPt, withinBounds := addonSide.Intercept(oldSide)
			fmt.Printf("debug interceptPt: %v\n", interceptPt)
			switch {
			case withinBounds:
				fmt.Printf("debug withinBounds: %v\n", withinBounds)

				var newLine1, newLine2 Line
				if tracingAddon {
					newLine1 = NewLine(addonSide.Start, interceptPt)
					newLine2 = NewLine(interceptPt, oldSide.End)
					oldSideN++
				} else {
					newLine1 = NewLine(oldSide.Start, interceptPt)
					newLine2 = NewLine(interceptPt, addonSide.End)
					addonSideN++
				}
				fmt.Printf("debug tracingAddon: %v\n", tracingAddon)
				fmt.Printf("debug newLine1: %v\n", newLine1)
				newSupersetPolygon[newSideN] = newLine1
				newSideN++
				newSupersetPolygon[newSideN] = newLine2
				newSideN++

				tracingAddon = !tracingAddon // start tracing the other

			case tracingAddon:
				if addonSide.WithinL2XBound(oldSide) {
					newSupersetPolygon[newSideN] = addonSide
					newSideN++
				}
				addonSideN++
			case !tracingAddon:
				if oldSide.WithinL2XBound(addonSide) {
					newSupersetPolygon[newSideN] = oldSide
					newSideN++
				}
				oldSideN++
			default:
				panic("wierd!")
			}

		}

		supersetPolygon = newSupersetPolygon
	}

	fmt.Printf("superset polygon, num points %v, length %v\nformatted: %v\n", len(supersetPolygon),
		lenLines(supersetPolygon), formattedLines(supersetPolygon))
}
