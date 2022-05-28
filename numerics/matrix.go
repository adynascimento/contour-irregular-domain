package numerics

import (
	"heatmap/polygonFunctions"
	"math"
	"sync"

	"github.com/lvisei/go-kriging/ordinarykriging"
	"github.com/paulmach/orb"
	"gonum.org/v1/gonum/mat"
)

type Mesh struct {
	X, Y *mat.Dense
}

type MatrixData struct {
	NGrid        int
	Mesh         Mesh
	KnownPoints  polygonFunctions.KnownPoints
	MultiPolygon orb.MultiPolygon
	Z            *mat.Dense
}

func Meshgrid(bounds orb.Ring, nGrid int) Mesh {
	xMin := bounds.Bound().Min[0]
	yMin := bounds.Bound().Min[1]
	xMax := bounds.Bound().Max[0]
	yMax := bounds.Bound().Max[1]

	return Mesh{
		X: mat.NewDense(1, nGrid, Linspace(xMin, xMax, nGrid)),
		Y: mat.NewDense(1, nGrid, Linspace(yMin, yMax, nGrid)),
	}
}

func PopulateMatrix(multiPolygon orb.MultiPolygon, knownPoints polygonFunctions.KnownPoints, mesh Mesh, nGrid int) *mat.Dense {
	z := mat.NewDense(nGrid, nGrid, nil)

	matrixData := MatrixData{
		NGrid:        nGrid,
		Mesh:         mesh,
		KnownPoints:  knownPoints,
		MultiPolygon: multiPolygon,
		Z:            z,
	}

	var wg sync.WaitGroup

	for i := 0; i < nGrid; i++ {
		wg.Add(1)
		go matrixData.interpol(i, &wg)
	}
	wg.Wait()

	return z
}

func (m MatrixData) interpol(i int, wg *sync.WaitGroup) {
	for j := 0; j < m.NGrid; j++ {
		point := orb.Point{m.Mesh.X.At(0, j), m.Mesh.Y.At(0, m.NGrid-1-i)}
		if polygonFunctions.PointsInMultiPolygon(m.MultiPolygon, point) {
			m.Z.Set(i, j, InterpolateKriging(m.KnownPoints, point))
		} else {
			m.Z.Set(i, j, math.NaN())
		}
	}

	wg.Done()
}

func InterpolateKriging(knownPoints polygonFunctions.KnownPoints, point orb.Point) float64 {
	ordinaryKriging := ordinarykriging.NewOrdinary(knownPoints.Values, knownPoints.Lon, knownPoints.Lat)
	variogram, _ := ordinaryKriging.Train(ordinarykriging.Spherical, 0.0, 100.0)
	return variogram.Predict(point[0], point[1])
}

func FlipVertically(matrix *mat.Dense) *mat.Dense {
	r, c := matrix.Dims()
	flippedMatrix := mat.NewDense(r, c, nil)
	for i := 0; i < r; i++ {
		flippedMatrix.SetRow(r-1-i, matrix.RawRowView(i))
	}
	return flippedMatrix
}
