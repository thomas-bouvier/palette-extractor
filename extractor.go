package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	reader, err := os.Open("./image.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	defer reader.Close()

	pixels, err := getPixels(reader, 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	fmt.Print(getPalette(pixels, 5))
}

func getPalette(pixels []Pixel, count int) []Pixel {
	return quantize(pixels, count).GetPalette()
}
