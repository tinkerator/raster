// Package raster contains some convenience functions for creating
// "golang.org/x/image/vector" paths and rendering them into images.
package raster

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/image/vector"
)

// Scriber is the interface for drawing vector graphics.
type Scriber interface {
	MoveTo(x, y float64)
	LineTo(x, y float64)
	CubeTo(a, b, c, d, e, f float64)
	ClosePath()
}

// Rasterizer is a wrapper for the
// golang.org/x/image/vector.Rasterizer type which maps float64
// arguments to float32 bit calls.
type Rasterizer struct {
	R *vector.Rasterizer
}

// NewRasterizer allocates a new rasterizer.
func NewRasterizer(w, h int) *Rasterizer {
	r := vector.NewRasterizer(w, h)
	return &Rasterizer{R: r}
}

// Reset resets the memory of the rasterizer and sets the size of its
// clipping rectangle.
func (r *Rasterizer) Reset(w, h int) {
	r.R.Reset(w, h)
}

// MoveTo sets the rasterizer pen to the coordinate (x,y).
func (r *Rasterizer) MoveTo(x, y float64) {
	r.R.MoveTo(float32(x), float32(y))
}

// LineTo constructs a straight line from the pen to the target (x,y)
// coordinate, and updates the pen to this location.
func (r *Rasterizer) LineTo(x, y float64) {
	r.R.LineTo(float32(x), float32(y))
}

// CubeTo constructs a cubic Bezier curve using the supplied
// parameters, from the pen location to point (e,f), which becomes the
// updated pen location.
func (r *Rasterizer) CubeTo(a, b, c, d, e, f float64) {
	r.R.CubeTo(float32(a), float32(b), float32(c), float32(d), float32(e), float32(f))
}

// ClosePath forms a loop back line from the pen to the start of the
// path.
func (r *Rasterizer) ClosePath() {
	r.R.ClosePath()
}

// The circular approximation, that uses this constant, with Bezier
// curves is nicely explained in:
//
//	https://stackoverflow.com/questions/1734745/how-to-create-circle-with-b%C3%A9zier-curves
const partial = 0.552284749831

// PointAt renders an approximate "circle" via 4 cubic Bezier curves
// describing the arc of the 4 quadrants.
func PointAt(r Scriber, x, y, width float64) {
	d := 0.5 * width
	p := partial * d
	r.MoveTo(x+d, y)
	r.CubeTo(x+d, y+p, x+p, y+d, x, y+d)
	r.CubeTo(x-p, y+d, x-d, y+p, x-d, y)
	r.CubeTo(x-d, y-p, x-p, y-d, x, y-d)
	r.CubeTo(x+p, y-d, x+d, y-p, x+d, y)
	r.ClosePath()
}

// LineTo renders a line segment from (oX,oY) to (nX,nY) with the
// specified perpendicular width. The capped value adds rounded
// end-caps to the line of radius half of the width (as approximated
// with Bezier curves).
func LineTo(r Scriber, capped bool, oX, oY, nX, nY, width float64) {
	if oX == nX && oY == nY {
		return // nothing to draw
	}
	dX, dY := nX-oX, nY-oY
	d := .5 * width / math.Sqrt(dX*dX+dY*dY)
	dX, dY = dX*d, dY*d
	r.MoveTo(oX-dY, oY+dX)
	if capped {
		r.CubeTo(oX-dY-partial*dX, oY+dX-partial*dY,
			oX-dX-partial*dY, oY-dY+partial*dX,
			oX-dX, oY-dY)
		r.CubeTo(oX-dX+partial*dY, oY-dY-partial*dX,
			oX+dY-partial*dX, oY-dX-partial*dY,
			oX+dY, oY-dX)
	} else {
		r.LineTo(oX+dY, oY-dX)
	}
	r.LineTo(nX+dY, nY-dX)
	if capped {
		r.CubeTo(nX+dY+partial*dX, nY-dX+partial*dY,
			nX+dX+partial*dY, nY+dY-partial*dX,
			nX+dX, nY+dY)
		r.CubeTo(nX+dX-partial*dY, nY+dY+partial*dX,
			nX-dY+partial*dX, nY+dX+partial*dY,
			nX-dY, nY+dX)
	} else {
		r.LineTo(nX-dY, nY+dX)
	}
	r.ClosePath()
}

// SquareAt renders a width by width square centered at (x,y).
func SquareAt(r Scriber, x, y, width float64) {
	if width <= 0 {
		return // nothing to draw
	}
	d := 0.5 * width
	r.MoveTo(x-d, y-d)
	r.LineTo(x+d, y-d)
	r.LineTo(x+d, y+d)
	r.LineTo(x-d, y+d)
	r.ClosePath()
}

// DrawAt places the r into an image aligning (x,y) of r with the
// (0,0) coordinate of the image.
func DrawAt(im draw.Image, r *vector.Rasterizer, x, y float64, col color.Color) {
	r.Draw(im, im.Bounds(), image.NewUniform(col), image.Point{X: int(x), Y: int(y)})
}
