![Programming Language](https://img.shields.io/badge/Go-Programming%20Language-brightgreen)
![Zero Clause BSD License](https://img.shields.io/badge/License-BSD%20Zero%20Clause-green)

# GoMandelbrot

# Example Output

![Rendered Output](example.gif)

# Description

Gif image generator written in Go.

# Usage

## Example

```
go run gomandelbrot.go -dimension 500 -frames 50 -delay 2 -iterations 1000 -powerStart 1.0 -powerEnd 10.0
```

## Parameters explained

- dimension
	- The height and width of the square image rendered to stdout.
- frames
	- Number of divisions between starting and ending power.
- delay
	- Number of 100th's of a second between frames.
- iterations
	- Maximum number of times to compute the Mandelbrot fractal equation "z\[n+1\] = z\[n\]^m + c".
- powerStart
	- First power to start the series of frames.
- powerEnd
	- Last power to end the series of frames.

# Building

## Requires

```
golang-1.18
```

# License

All code and files in this repository are licensed under the 0-BSD License

