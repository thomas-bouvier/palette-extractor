package main

import (
	"image"
	"io"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

func getPixels(file io.Reader, quality int) ([]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels []Pixel
	for i := 0; i < width*height; i += quality {
		pixel := rgbaToPixel(img.At(i%width, i/width).RGBA())

		if pixel.A >= 125 {
			if !(pixel.R > 250 && pixel.G > 250 && pixel.B > 250) {
				pixels = append(pixels, pixel)
			}
		}
	}

	return pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
