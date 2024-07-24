package main

import (
	"image"
	"image/color"
	"image/gif"
	p "image/color/palette"
	"io"
	"math/cmplx"
	"os"
	"fmt"
)

var palette []color.Color
var colorsMax int

func init() {
	palette = p.Plan9
	colorsMax = len(palette)
}

func main() {
	mandelbrotAnimation(os.Stdout)
}

func coordinateTranslate(coordinate, domainStart, domainEnd, scaleFactor float64) int {
	return int(((coordinate - domainStart) / (domainEnd - domainStart)) * scaleFactor)
}

func mandelbrotFrame(domainStartX, domainStartY, domainEndX, domainEndY, epsilon float64, imageDimension, iterationMax int) *image.Paletted {
	var iteration int
	var c, z complex128

	rect := image.Rect(0, 0, int(imageDimension), int(imageDimension))
	frame := image.NewPaletted(rect, palette)

	for y := domainStartY; y < domainEndY; y += epsilon {
		for x := domainStartX; x < domainEndX; x += epsilon {
			iteration = 0
			c = complex(x, y)
			z = complex(0, 0)
			for ;cmplx.Abs(z) < domainEndX && iteration < iterationMax; iteration++ {
				z = z*z + c
			}

			frameX := coordinateTranslate(float64(x), float64(domainStartX), float64(domainEndX), float64(imageDimension))
			frameY := coordinateTranslate(float64(y), float64(domainStartY), float64(domainEndY), float64(imageDimension))
			frame.SetColorIndex(frameX, frameY, uint8(iteration  % colorsMax))
		}
	}

	return frame
}

// Mandelbrot set:
// f[c] = z^2 + c
func mandelbrotAnimation(out io.Writer) {
	const (
		epsilon = 0.001

		imageDimension = 2.0 * (1.0 / epsilon)

		frameTotal = 1
		frameDelay = 8

		iterationMax = 255

		domainEndX = 2.0
		domainEndY = 2.0
		domainStartX = -2.0
		domainStartY = -2.0
	)

	animation := gif.GIF{LoopCount: frameTotal}

	for i := 0; i < frameTotal; i++ {
		fmt.Fprintf(os.Stderr, "frame: i: %d\n", i)
		frame := mandelbrotFrame(float64(domainStartX), float64(domainStartY), float64(domainEndX), float64(domainEndY), float64(epsilon), int(imageDimension), int(iterationMax))

		animation.Delay = append(animation.Delay, frameDelay)
		animation.Image = append(animation.Image, frame)
	}

	gif.EncodeAll(out, &animation)
}








