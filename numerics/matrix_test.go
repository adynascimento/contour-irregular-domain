package numerics

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"heatmap/polygonFunctions"

	"github.com/paulmach/orb"
	"gonum.org/v1/gonum/mat"
)

func TestMeshgrid(t *testing.T) {
	cases := []struct {
		Description  string
		MultiPolygon orb.MultiPolygon
		NGrid        int
		Want         Mesh
	}{
		{
			Description: "Testing a 3x3 square",
			MultiPolygon: orb.MultiPolygon{
				orb.Polygon{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}},
			},
			NGrid: 3,
			Want: Mesh{
				X: mat.NewDense(1, 3, Linspace(0, 2, 3)),
				Y: mat.NewDense(1, 3, Linspace(0, 2, 3)),
			},
		},
		{
			Description: "Testing a 3x3 multi square",
			MultiPolygon: orb.MultiPolygon{
				orb.Polygon{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}},
				orb.Polygon{{{5, 5}, {5, 8}, {9, 9}, {9, 2}, {5, 5}}},
			},
			NGrid: 3,
			Want: Mesh{
				X: mat.NewDense(1, 3, Linspace(0, 9, 3)),
				Y: mat.NewDense(1, 3, Linspace(0, 9, 3)),
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			bounds := polygonFunctions.GetBoxBounds(test.MultiPolygon)
			got := Meshgrid(bounds, test.NGrid)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("got %v, want %v", got, test.Want)
			}
		})
	}
}

func TestPopulateMatrix(t *testing.T) {
	cases := []struct {
		Description  string
		MultiPolygon orb.MultiPolygon
		KnownPoints  polygonFunctions.KnownPoints
		NGrid        int
		Want         *mat.Dense
		Condition    int
	}{
		{
			Description: "Testing matrix with NaN in a triangle case",
			MultiPolygon: orb.MultiPolygon{
				orb.Polygon{{{0, 0}, {2, 0}, {1, 2}, {0, 0}}},
			},
			KnownPoints: polygonFunctions.KnownPoints{
				Lon:    []float64{0.5, 0.75, 1.5},
				Lat:    []float64{0.5, 0.75, 0.5},
				Values: []float64{5, 15, 8},
			},
			NGrid: 3,
			Want: mat.NewDense(
				3,
				3,
				[]float64{math.NaN(), 0, math.NaN(), math.NaN(), 0, math.NaN(), 0, 0, 0},
			),
			Condition: 0,
		},
		{
			Description: "Testing kriging in matrix with NaN",
			MultiPolygon: orb.MultiPolygon{
				orb.Polygon{{{0, 0}, {2, 0}, {1, 2}, {0, 0}}},
			},
			KnownPoints: polygonFunctions.KnownPoints{
				Lon:    []float64{0.5, 0.75, 1.5},
				Lat:    []float64{0.5, 0.75, 0.5},
				Values: []float64{5, 15, 8},
			},
			NGrid: 3,
			Want: mat.NewDense(
				3,
				3,
				[]float64{math.NaN(), 0, math.NaN(), math.NaN(), 0, math.NaN(), 0, 0, 0},
			),
			Condition: 1,
		},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			bounds := polygonFunctions.GetBoxBounds(test.MultiPolygon)
			mesh := Meshgrid(bounds, test.NGrid)
			got := PopulateMatrix(test.MultiPolygon, test.KnownPoints, mesh, test.NGrid)
			var condition bool
			switch test.Condition {
			case 0:
				condition = fmt.Sprint(got.RawMatrix().Data[0]) != fmt.Sprint(test.Want.RawMatrix().Data[0])
			case 1:
				condition = !(got.RawMatrix().Data[1] > test.Want.RawMatrix().Data[1])
			default:
				condition = false
			}
			if condition {
				t.Errorf("got %v, want %v", got, test.Want)
			}

		})
	}
}

func TestFlipVertically(t *testing.T) {
	t.Run("Testing a 2x2 matrix", func(t *testing.T) {
		matrix := mat.NewDense(2, 2, []float64{1, 2, 3, 4})
		got := FlipVertically(matrix)
		want := mat.NewDense(2, 2, []float64{3, 4, 1, 2})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Testing a 3x3 matrix", func(t *testing.T) {
		matrix := mat.NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
		got := FlipVertically(matrix)
		want := mat.NewDense(3, 3, []float64{7, 8, 9, 4, 5, 6, 1, 2, 3})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}
