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
	"runtime"
	"bufio"
)

var colors []color.Color
var mMin float64 = -2.0
var mMax float64 =  2.0
var sema chan struct{}

type Frame struct {
	Imgseq int
	Img *image.Paletted
}

func init() {
	colors = palette.Plan9
	sema = make(chan struct{}, runtime.NumCPU())
}

func main() {
	mandelbrotAnimation(os.Stdout)
}

// Translates x and y coordinates from one coordinate domain to another
func coordinateTranslate(x, srcMin, srcMax, destMin, destMax float64) (ret float64) {
	return (((x - srcMin) / (srcMax - srcMin)) * (destMax - destMin)) + destMin
}

func mandelbrotFrame(dStartX, dEndX, dStartY, dEndY, nPow float64, iterMax, imgW, imgH, imgseq int, ch chan *Frame) {
	// Number of fractal formula iterations
	var iter int
	// z[n+1] = z[n]^2 + c
	// where z[0] = 0, 0
	// where c = tx, ty
	var c, z complex128
	// Translated x and y
	var tx, ty float64

	var ret Frame
	ret.Imgseq = imgseq

	sema <- struct{}{}

	rect  := image.Rect(0, 0, imgW, imgH)
	frame := image.NewPaletted(rect, colors)

	for y := 0; y < imgH; y++ {
		ty = coordinateTranslate(float64(y), 0.0, float64(imgH) - 1, dStartY, dEndY)
		for x := 0; x < imgW; x++ {
			tx = coordinateTranslate(float64(x), 0.0, float64(imgW) - 1, dStartX, dEndX)

			c = complex(tx, ty)
			z = complex(0, 0)
			for iter = 0; cmplx.Abs(z) < dEndX && iter < iterMax; iter++ {
				z = cmplx.Pow(z, complex(nPow, 0)) + c
			}

			frame.SetColorIndex(x, y, uint8(iter % len(colors)))
		}
	}
	ret.Img = frame

	<-sema
	ch <- &ret
}

// Mandelbrot set:
// f[c] = z^2 + c
func mandelbrotAnimation(out io.Writer) {
	const (
		// Image
		imageH = 2000
		imageW = 2000

		// gif
		frameTotal = 1000
		frameDelay = 20

		// Mandelbrot
		iterMax = 100
	)

	var nPow float64

	var frames []*Frame

	animation := gif.GIF{LoopCount: frameTotal}

	ch := make(chan *Frame)

	for i := 0; i < frameTotal; i++ {
		nPow = 2.0 + 10.0 * (float64(i) / float64(frameTotal))
		go mandelbrotFrame(mMin, mMax, mMin, mMax, nPow, int(iterMax), int(imageW), int(imageH), i, ch)
	}

	w := bufio.NewWriter(os.Stderr)
	for ;len(frames) < frameTotal; {
		f := <-ch
		frames = append(frames, f)
		fmt.Fprintf(w, "\rFrame: %4d / %d", len(frames), frameTotal)
		w.Flush()
	}
	fmt.Fprintln(w)

	for i := 0; i < len(frames); i++ {
		for j := i+1; j < len(frames); j++ {
			if frames[i].Imgseq > frames[j].Imgseq {
				tmp := frames[i]
				frames[i] = frames[j]
				frames[j] = tmp
			}
		}
	}


	for i := 0; i < frameTotal; i++ {
		animation.Delay = append(animation.Delay, frameDelay)
		animation.Image = append(animation.Image, frames[i].Img)
	}

	gif.EncodeAll(out, &animation)
}








