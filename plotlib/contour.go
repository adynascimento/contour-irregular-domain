package plotlib

import (
	"heatmap/numerics"
	"image/color"
	"log"
	"os"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func NewPlot() plotParameters {
	return plotParameters{}
}

// size of the saved figure
func (plt *plotParameters) FigSize(xwidth, ywidth int) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
}

// parameters to heatmap plot
func (plt *plotParameters) Contour(mesh numerics.Mesh, z *mat.Dense, n_levels int, gradient colorgrad.Gradient) {
	plt.contourData.x = mesh.X
	plt.contourData.y = mesh.Y
	plt.contourData.z = z
	plt.n_levels = n_levels
	plt.gradient = gradient
}

// generate plot and save it to file
func (plt *plotParameters) Save(name string) {
	// create a new plot
	p := plot.New()

	// make a heatmap plotter
	plt.contourPlot(p)

	// save the plot to a PNG file.
	xwdith := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywdith := font.Length(plt.figSize.ywidth) * vg.Centimeter

	c := vgimg.PngCanvas{
		Canvas: vgimg.NewWith(
			vgimg.UseWH(xwdith, ywdith),
			vgimg.UseBackgroundColor(color.Transparent),
		),
	}
	p.Draw(draw.New(c))

	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
}

func (plt *plotParameters) Colorbar() {
	// create a new plot
	p := plot.New()
	p.BackgroundColor = color.Transparent
	p.HideY()
	p.X.Padding = 0
	p.Y.Padding = 0

	colors := plt.gradient.Colors(uint(plt.n_levels))
	min := mat.Min(plt.contourData.z)
	max := mat.Max(plt.contourData.z)
	plt.gradient, _ = colorgrad.NewGradient().Colors(colors...).Domain(min, max).Build()
	colorMap := &colorMap{gradient: plt.gradient, n_levels: plt.n_levels}
	cb := &plotter.ColorBar{ColorMap: colorMap}
	p.Add(cb)

	if err := p.Save(300, 35, "colorbar.png"); err != nil {
		log.Panic(err)
	}
}

func (plt *plotParameters) contourPlot(p *plot.Plot) {
	// prepare data to plot
	m := unitGrid{x: plt.contourData.x, y: plt.contourData.y, Data: plt.contourData.z}

	// add colormap and make a heatmap plotter
	palette := colorsGradient{ColorList: plt.gradient.Colors(uint(plt.n_levels))}
	raster := plotter.NewHeatMap(m, palette)
	raster.Rasterized = true
	raster.NaN = color.Transparent

	p.BackgroundColor = color.Transparent
	p.HideAxes()
	p.X.Padding = 0
	p.Y.Padding = 0
	p.Add(raster)
}
