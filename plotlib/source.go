package plotlib

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
)

// struct that defines methods to match the GridXYZ interface defined in gonum plot library
// used in heatmap and contour plots
type unitGrid struct {
	x, y, Data *mat.Dense
}

// methods to match the GridXYZ interface defined in gonum plot library
func (g unitGrid) Dims() (c, r int)   { r, c = g.Data.Dims(); return c, r }
func (g unitGrid) Z(c, r int) float64 { return g.Data.At(r, c) }
func (g unitGrid) X(c int) float64    { return g.x.At(0, c) }
func (g unitGrid) Y(r int) float64    { return g.y.At(0, r) }

// struct that defines methods to match the Palette interface defined in gonum plot library
// used in heatmap and contour plots
type colorsGradient struct {
	ColorList []color.Color
}

// methods to match the Palette interface defined in gonum plot library
func (g colorsGradient) Colors() []color.Color { return g.ColorList }

type plotParameters struct {
	contourData contourData        // used in heatmap and contour
	gradient    colorgrad.Gradient // colormap
	n_levels    int                // colormap levels
	figSize     figSize            // xwidth and ywidth of the saved figure
}

type contourData struct{ x, y, z *mat.Dense }
type figSize struct{ xwidth, ywidth int }
