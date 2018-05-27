package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
)

func main() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	reader, err := os.Open("./image.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	defer reader.Close()

	pixels, err := getPixels(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	fmt.Println(pixels)
}

func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

type Pixel struct {
	R int
	G int
	B int
	A int
}
