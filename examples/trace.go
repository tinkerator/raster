// Program trace creates a trace-like rendering of something you might
// see on a PCB.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"zappem.net/pub/graphics/raster"
)

var (
	dest   = flag.String("dest", "", "destination png file")
	width  = flag.Int("width", 200, "width of image")
	height = flag.Int("height", 100, "height of image")
)

func main() {
	flag.Parse()

	if *dest == "" {
		log.Fatal("please specifiy --dest=image.png")
	}
	f, err := os.Create(*dest)
	if err != nil {
		log.Fatalf("unable to create %q: %v", *dest, err)
	}
	defer f.Close()

	w, h := *width, *height
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(im, im.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Over)

	d := float64(h) * .2
	cX, cY := float64(w)*0.5, float64(h)*0.5

	conv := func(a, b float64) (col, row float64) {
		col = a*d + cX
		row = cY - b*d
		return
	}

	pad1X, pad1Y := conv(-4, -1)
	pt1X, pt1Y := conv(-1, -1)
	pt2X, pt2Y := conv(1, 1)
	pad2X, pad2Y := conv(4, 1)

	r := raster.NewRasterizer(w, h)
	raster.SquareAt(r, pad1X, pad1Y, d)
	raster.LineTo(r, true, pad1X, pad1Y, pt1X, pt1Y, d/3)
	raster.LineTo(r, true, pt1X, pt1Y, pt2X, pt2Y, d/3)
	raster.LineTo(r, true, pt2X, pt2Y, pad2X, pad2Y, d/3)
	raster.PointAt(r, pad2X, pad2Y, d)
	raster.DrawAt(im, r.R, 0, 0, color.Black)
	r.Reset(w, h)

	raster.PointAt(r, pad1X, pad1Y, d*0.6)
	raster.PointAt(r, pad2X, pad2Y, d*0.6)
	raster.DrawAt(im, r.R, 0, 0, color.White)

	png.Encode(f, im)
}
