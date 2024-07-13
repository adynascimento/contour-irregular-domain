package ngo

import (
	"github.com/paulmach/orb"
	"gonum.org/v1/gonum/mat"
)

type Mesh struct {
	X, Y *mat.Dense
}

type Points struct {
	Lon    []float64
	Lat    []float64
	Values []float64
}

// return coordinate matrices from polygon bound
func Meshgrid(polygonBound orb.Bound, nGrid int) Mesh {
	xMin := polygonBound.Min[0]
	yMin := polygonBound.Min[1]
	xMax := polygonBound.Max[0]
	yMax := polygonBound.Max[1]

	return Mesh{
		X: mat.NewDense(1, nGrid, Linspace(xMin, xMax, nGrid)),
		Y: mat.NewDense(1, nGrid, Linspace(yMin, yMax, nGrid)),
	}
}

// generate linearly spaced slice of float64
func Linspace(start, stop float64, num int) []float64 {
	var step float64
	if num == 1 {
		return []float64{start}
	}
	step = (stop - start) / float64(num-1)

	r := make([]float64, num)
	for i := 0; i < num; i++ {
		r[i] = start + float64(i)*step
	}
	return r
}
