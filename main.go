package main

import (
	"countour-irregular-domain/ngo"

	"github.com/adynascimento/plot/plotter"
	"github.com/mazznoer/colorgrad"
	"github.com/paulmach/orb"
)

func main() {
	nGrid := 501

	multiPolygon := orb.MultiPolygon{
		orb.Polygon{{{0, 0}, {1, 0}, {0.5, 1}, {0, 0}}},
	}
	knownPoints := ngo.Points{
		Lon:    []float64{0.2, 0.4, 0.6, 0.8, 0.4, 0.6, 0.5},
		Lat:    []float64{0.2, 0.2, 0.2, 0.2, 0.4, 0.4, 0.8},
		Values: []float64{5.0, 16.0, 8.0, 20.0, 55.0, 27.0, 40.0},
	}

	mesh := ngo.Meshgrid(multiPolygon.Bound(), nGrid)
	matrix := ngo.PolygonInterpolation(mesh, multiPolygon, knownPoints)

	plt := plotter.NewPlot()
	plt.FigSize(12, 12)

	plt.ContourF(mesh.X, mesh.Y, matrix,
		plotter.WithLevels(15),
		plotter.WithGradient(colorgrad.Turbo()),
		plotter.WithColorbar(plotter.Vertical))
	plt.Save("triangle.png")
	plt.Show()
}
