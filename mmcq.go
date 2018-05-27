package main

import (
	"fmt"
	"os"
)

const BITSIG = 5
const RSHIFT = 8 - BITSIG

func quantize(pixels []Pixel, count int) {
	if count < 2 || count > 256 {
		fmt.Fprintf(os.Stderr, "wrong number of max colors when quantize")
		os.Exit(1)
	}

	histogram := computeHistogram(pixels)

	if len(histogram) <= count {
		fmt.Fprintf(os.Stderr, "insufficient number of levels of quantification")
		os.Exit(1)
	}
}

func computeHistogram(pixels []Pixel) map[int]int {
	var index int
	histogram := make(map[int]int)

	for _, pixel := range pixels {
		index = getColorIndex(pixel.R>>RSHIFT, pixel.G>>RSHIFT, pixel.B>>RSHIFT)
		if val, ok := histogram[index]; ok {
			histogram[index] = val + 1
		} else {
			histogram[index] = 1
		}

	}

	return histogram
}

func getColorIndex(r int, g int, b int) int {
	return r<<(BITSIG*2) + g<<BITSIG + b
}
