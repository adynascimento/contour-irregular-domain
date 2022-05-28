package main

import (
	"fmt"
	"heatmap/numerics"
	"heatmap/plotlib"
	"heatmap/polygonFunctions"
	"time"

	"github.com/mazznoer/colorgrad"
	"github.com/paulmach/orb"
)

func main() {

	nGrid := 500

	multiPolygon := orb.MultiPolygon{
		orb.Polygon{{{0, 0}, {1, 0}, {0.5, 1}, {0, 0}}},
	}
	knownPoints := polygonFunctions.KnownPoints{
		Lon:    []float64{0.2, 0.4, 0.6, 0.8, 0.4, 0.6, 0.5},
		Lat:    []float64{0.2, 0.2, 0.2, 0.2, 0.4, 0.4, 0.8},
		Values: []float64{5.0, 16.0, 8.0, 20.0, 55.0, 27.0, 40.0},
	}

	bounds := polygonFunctions.GetBoxBounds(multiPolygon)
	mesh := numerics.Meshgrid(bounds, nGrid)

	start := time.Now()
	matrix := numerics.PopulateMatrix(multiPolygon, knownPoints, mesh, nGrid)
	elapsed := time.Since(start)
	fmt.Println("time: " + elapsed.String())

	plt := plotlib.NewPlot()
	plt.FigSize(8, 8)
	plt.Contour(mesh, numerics.FlipVertically(matrix), 15, colorgrad.Turbo())
	plt.Save("triangle.png")

}
