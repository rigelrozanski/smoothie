package main

import (
	"fmt"
	"math/big"
)

var (
	one = big.NewFloat(1)

	// nolint- bounds for a quarter of the circle
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

//____________________________________

func initFloat(f *big.Float) {
	f.SetPrec(200)
	f.SetMode(big.ToNearestEven)
}

func newFloat() *big.Float {
	f := new(big.Float)
	initFloat(f)
	return f
}

//____________________________________

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
	fmt.Println("wackydebugoutput main 0")

	var n int64 = 3  // starting number of divisions
	maxN := int64(4) //2*n - 1 // maximum number of sides for boring polygons

	boringPolygons := make(map[int64]map[int64]Line) // index 1: number of lines in element, index 2: element line no.

	// init boring polygons
	for i := n; i <= maxN; i++ {
		boringPolygons[i] = make(map[int64]Line)

		startPoint := Point{big.NewFloat(0), big.NewFloat(1)} // top of the circle
		initFloat(startPoint.X)
		initFloat(startPoint.Y)
		width := newFloat().Quo(XBoundMax, big.NewFloat(float64(i))) // width of all these pieces
		initFloat(width)

		for side := int64(0); side < i; side++ {
			x2 := newFloat().Add(startPoint.X, width)
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
		tracing := addonPolygon[addonSideN]
		comparing := supersetPolygon[oldSideN]

		fmt.Printf("superset polygon, num points %v, length %v\nformatted: %v\n", len(supersetPolygon),
			lenLines(supersetPolygon), formattedLines(supersetPolygon))

		for {
			fmt.Printf("----------------------------\n")
			fmt.Printf("debug oldSideN: %v\n", oldSideN)
			fmt.Printf("debug addonSideN: %v\n", addonSideN)
			if oldSideN > maxOldSides-1 || addonSideN > maxAddonSides-1 {
				break
			}

			var withinBounds bool
			var interceptPt Point
			interceptPt, withinBounds = tracing.Intercept(comparing)

			fmt.Printf("debug interceptPt: %v\n", interceptPt)
			fmt.Printf("debug withinBounds: %v\n", withinBounds)
			fmt.Printf("debug tracingAddon: %v\n", tracingAddon)
			fmt.Printf("debug tracing: %v\n", tracing)
			fmt.Printf("debug comparing: %v\n", comparing)
			switch {
			case withinBounds:

				var newLine1 Line
				if tracingAddon {
					newLine1 = NewLine(tracing.Start, interceptPt)
					nextTracing := NewLine(interceptPt, comparing.End)
					nextComparing := NewLine(interceptPt, tracing.End)
					tracing = nextTracing
					comparing = nextComparing
					//comparing = addonPolygon[addonSideN]

				} else {
					newLine1 = NewLine(tracing.Start, interceptPt)
					nextTracing := NewLine(interceptPt, comparing.End)
					nextComparing := NewLine(interceptPt, tracing.End)
					tracing = nextTracing
					comparing = nextComparing
					//comparing = supersetPolygon[oldSideN]
				}

				newSupersetPolygon[newSideN] = newLine1
				newSideN++

				tracingAddon = !tracingAddon // start tracing the other

			case tracingAddon:
				if tracing.WithinL2XBound(comparing) {
					newSupersetPolygon[newSideN] = tracing
					newSideN++
				}
				addonSideN++
				tracing = addonPolygon[addonSideN]
				comparing = supersetPolygon[oldSideN]

			case !tracingAddon:
				if tracing.WithinL2XBound(comparing) {
					newSupersetPolygon[newSideN] = tracing
					newSideN++
				}
				oldSideN++
				tracing = supersetPolygon[oldSideN]
				comparing = addonPolygon[addonSideN]

			default:
				panic("wierd!")
			}

		}

		supersetPolygon = newSupersetPolygon
	}

	fmt.Printf("superset polygon, num points %v, length %v\nformatted: %v\n", len(supersetPolygon),
		lenLines(supersetPolygon), formattedLines(supersetPolygon))
}
