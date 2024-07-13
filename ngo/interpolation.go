package ngo

import (
	"math"
	"sync"

	"github.com/lvisei/go-kriging/ordinarykriging"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"gonum.org/v1/gonum/mat"
)

func PolygonInterpolation(mesh Mesh, multiPolygon orb.MultiPolygon, knownPoints Points) *mat.Dense {
	nX := mesh.X.RawMatrix().Cols
	nY := mesh.Y.RawMatrix().Cols

	// kriging interpolation
	ordinaryKriging := ordinarykriging.NewOrdinary(knownPoints.Values,
		knownPoints.Lon, knownPoints.Lat)
	variogram, _ := ordinaryKriging.Train(ordinarykriging.Spherical, 0.0, 100.0)

	var wg sync.WaitGroup
	workers := make(chan struct{}, 10)

	z := mat.NewDense(nX, nY, nil)
	for i, xv := range mesh.X.RawMatrix().Data {
		wg.Add(1)
		workers <- struct{}{}
		go func(i int, xv float64) {
			defer wg.Done()
			defer func() { <-workers }()

			for j, yv := range mesh.Y.RawMatrix().Data {
				point := orb.Point{xv, yv}
				if planar.MultiPolygonContains(multiPolygon, point) {
					z.Set(j, i, variogram.Predict(point[0], point[1]))
				} else {
					z.Set(j, i, math.NaN())
				}
			}
		}(i, xv)
	}
	wg.Wait()

	return z
}
