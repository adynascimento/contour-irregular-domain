package plotlib

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/palette"
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

// struct that defines methods to match the ColorMap interface defined in gonum plot library
// used in colorbar plots
type colorMap struct {
	gradient colorgrad.Gradient // colormap
	n_levels int                // colormap levels
	max, min float64            // max and min of colormap in colorbar
}

// At implements the palette.ColorMap interface.
func (p *colorMap) At(v float64) (color.Color, error) { return p.gradient.At(v), nil }

// SetMax implements the palette.ColorMap interface.
func (p *colorMap) SetMax(v float64) {}

// SetMin implements the palette.ColorMap interface.
func (p *colorMap) SetMin(v float64) {}

// Max implements the palette.ColorMap interface.
func (p *colorMap) Max() float64 { 
	_, p.max = p.gradient.Domain()
	return p.max 
}

// Min implements the palette.ColorMap interface.
func (p *colorMap) Min() float64 { 
	p.min, _ = p.gradient.Domain()
	return p.min
 }

// SetAlpha sets the opacity value of this color map. Zero is transparent
// and one is completely opaque.
// The function will panic is alpha is not between zero and one.
func (p *colorMap) SetAlpha(alpha float64) {}

// Alpha returns the opacity value of this color map.
func (p *colorMap) Alpha() float64 { return 1.0 }

// Palette returns a palette.Palette with the specified number of colors.
func (p *colorMap) Palette(n int) palette.Palette {
	return colorsGradient{ColorList: p.gradient.Colors(uint(p.n_levels))}
}
