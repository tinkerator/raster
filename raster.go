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

// The circular approximation, that uses this constant, with Bezier
// curves is nicely explained in:
//
//	https://stackoverflow.com/questions/1734745/how-to-create-circle-with-b%C3%A9zier-curves
const partial = 0.552284749831

// PointAt renders an approximate "circle" via 4 cubic Bezier curves
// describing the arc of the 4 quadrants.
func PointAt(r *vector.Rasterizer, x, y, width float64) {
	d := 0.5 * width
	p := partial * d
	r.MoveTo(float32(x+d), float32(y))
	r.CubeTo(float32(x+d), float32(y+p), float32(x+p), float32(y+d), float32(x), float32(y+d))
	r.CubeTo(float32(x-p), float32(y+d), float32(x-d), float32(y+p), float32(x-d), float32(y))
	r.CubeTo(float32(x-d), float32(y-p), float32(x-p), float32(y-d), float32(x), float32(y-d))
	r.CubeTo(float32(x+p), float32(y-d), float32(x+d), float32(y-p), float32(x+d), float32(y))
	r.ClosePath()
}

// LineTo renders a line segment from (oX,oY) to (nX,nY) with the
// specified perpendicular width. The capped value adds rounded
// end-caps to the line of radius half of the width (as approximated
// with Bezier curves).
func LineTo(r *vector.Rasterizer, capped bool, oX, oY, nX, nY, width float64) {
	if oX == nX && oY == nY {
		return // nothing to draw
	}
	dX, dY := nX-oX, nY-oY
	d := .5 * width / math.Sqrt(dX*dX+dY*dY)
	dX, dY = dX*d, dY*d
	r.MoveTo(float32(oX-dY), float32(oY+dX))
	if capped {
		r.CubeTo(float32(oX-dY-partial*dX), float32(oY+dX-partial*dY),
			float32(oX-dX-partial*dY), float32(oY-dY+partial*dX),
			float32(oX-dX), float32(oY-dY))
		r.CubeTo(float32(oX-dX+partial*dY), float32(oY-dY-partial*dX),
			float32(oX+dY-partial*dX), float32(oY-dX-partial*dY),
			float32(oX+dY), float32(oY-dX))
	} else {
		r.LineTo(float32(oX+dY), float32(oY-dX))
	}
	r.LineTo(float32(nX+dY), float32(nY-dX))
	if capped {
		r.CubeTo(float32(nX+dY+partial*dX), float32(nY-dX+partial*dY),
			float32(nX+dX+partial*dY), float32(nY+dY-partial*dX),
			float32(nX+dX), float32(nY+dY))
		r.CubeTo(float32(nX+dX-partial*dY), float32(nY+dY+partial*dX),
			float32(nX-dY+partial*dX), float32(nY+dX+partial*dY),
			float32(nX-dY), float32(nY+dX))
	} else {
		r.LineTo(float32(nX-dY), float32(nY+dX))
	}
	r.ClosePath()
}

// SquareAt renders a width by width square centered at (x,y).
func SquareAt(r *vector.Rasterizer, x, y, width float64) {
	if width <= 0 {
		return // nothing to draw
	}
	d := 0.5 * width
	r.MoveTo(float32(x-d), float32(y-d))
	r.LineTo(float32(x+d), float32(y-d))
	r.LineTo(float32(x+d), float32(y+d))
	r.LineTo(float32(x-d), float32(y+d))
	r.ClosePath()
}

// DrawAt places the r into an image aligning (x,y) of r with the
// (0,0) coordinate of the image.
func DrawAt(im draw.Image, r *vector.Rasterizer, x, y float64, col color.Color) {
	r.Draw(im, im.Bounds(), image.NewUniform(col), image.Point{X: int(x), Y: int(y)})
}
