package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/color/palette"
	"io"
	"math/cmplx"
	"os"
	"fmt"
)

var colors []color.Color
var mMin float64 = -2.0
var mMax float64 =  2.0

func init() {
	colors = palette.Plan9
}

func main() {
	mandelbrotAnimation(os.Stdout)
}

// Translates x and y coordinates from one coordinate domain to another
func coordinateTranslate(x, srcMin, srcMax, destMin, destMax float64) (ret float64) {
	return (((x - srcMin) / (srcMax - srcMin)) * (destMax - destMin)) + destMin
}

func mandelbrotFrame(dStartX, dEndX, dStartY, dEndY, nPow float64, iterMax, imgW, imgH int) *image.Paletted {
	// Number of fractal formula iterations
	var iter int
	// z[n+1] = z[n]^2 + c
	// where z[0] = 0, 0
	// where c = tx, ty
	var c, z complex128
	// Translated x and y
	var tx, ty float64

	rect  := image.Rect(0, 0, imgW, imgH)
	frame := image.NewPaletted(rect, colors)

	for y := 0; y < imgH; y++ {
		ty = coordinateTranslate(float64(y), 0.0, float64(imgH), dStartY, dEndY)
		for x := 0; x < imgW; x++ {
			tx = coordinateTranslate(float64(x), 0.0, float64(imgW), dStartX, dEndX)

			c = complex(tx, ty)
			z = complex(0, 0)
			for iter = 0; cmplx.Abs(z) < dEndX && iter < iterMax; iter++ {
				z = cmplx.Pow(z, complex(nPow, 0)) + c
			}

			frame.SetColorIndex(x, y, uint8(iter % len(colors)))
		}
	}

	return frame
}

// Mandelbrot set:
// f[c] = z^2 + c
func mandelbrotAnimation(out io.Writer) {
	imageH := 2000
	imageW := 2000

	frameTotal := 1000
	frameDelay := 8

	iterMax := 100

	animation := gif.GIF{LoopCount: frameTotal}

	var nPow float64
	for i := 0; i < frameTotal; i++ {
		nPow = 2.0 + 10.0 * (float64(i) / float64(frameTotal))

		fmt.Fprintf(os.Stderr, "frame: %d\n", i+1)
		frame := mandelbrotFrame(mMin, mMax, mMin, mMax, nPow, int(iterMax), int(imageW), int(imageH))

		animation.Delay = append(animation.Delay, frameDelay)
		animation.Image = append(animation.Image, frame)
	}

	gif.EncodeAll(out, &animation)
}








