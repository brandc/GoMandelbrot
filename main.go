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

func init() {
	colors = palette.Plan9
}

func main() {
	mandelbrotAnimation(os.Stdout)
}

func coordinateTranslate(coordinate, domainStart, domainEnd, scaleFactor float64) int {
	return int(((coordinate - domainStart) / (domainEnd - domainStart)) * scaleFactor)
}

func mandelbrotFrame(domainStartX, domainStartY, domainEndX, domainEndY float64, iterationMax, imageW, imageH int) *image.Paletted {
	var iteration int
	var c, z complex128

	iterX := (domainEndX - domainStartX) / float64(imageW)
	iterY := (domainEndY - domainStartY) / float64(imageH)

	rect := image.Rect(0, 0, int(imageW), int(imageH))
	frame := image.NewPaletted(rect, colors)

	for y := domainStartY; y < domainEndY; y += iterY {
		for x := domainStartX; x < domainEndX; x += iterX {
			iteration = 0
			c = complex(x, y)
			z = complex(0, 0)
			for ;cmplx.Abs(z) < domainEndX && iteration < iterationMax; iteration++ {
				z = z*z + c
			}

			frameX := coordinateTranslate(x, domainStartX, domainEndX, float64(imageW))
			frameY := coordinateTranslate(y, domainStartY, domainEndY, float64(imageH))
			frame.SetColorIndex(frameX, frameY, uint8(iteration  % len(colors)))
		}
	}

	return frame
}

// Mandelbrot set:
// f[c] = z^2 + c
func mandelbrotAnimation(out io.Writer) {
	imageH := 1000
	imageW := 1000

	frameTotal := 1
	frameDelay := 8

	iterationMax := 255

	domainEndX := 2.0
	domainEndY := 2.0
	domainStartX := -2.0
	domainStartY := -2.0

	animation := gif.GIF{LoopCount: frameTotal}

	for i := 0; i < frameTotal; i++ {
		fmt.Fprintf(os.Stderr, "frame: %d\n", i+1)
		frame := mandelbrotFrame(float64(domainStartX), float64(domainStartY), float64(domainEndX), float64(domainEndY), int(iterationMax), imageW, imageH)

		animation.Delay = append(animation.Delay, frameDelay)
		animation.Image = append(animation.Image, frame)
	}

	gif.EncodeAll(out, &animation)
}








