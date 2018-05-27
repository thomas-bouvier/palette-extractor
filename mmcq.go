package main

import (
	"fmt"
	"os"
)

const BITSIG = 5
const RSHIFT = 8 - BITSIG

type VBox struct {
	r1, r2    int
	g1, g2    int
	b1, b2    int
	histogram map[int]int
}

func (vbox *VBox) volume() int {
	sub_r := vbox.r2 - vbox.r1
	sub_g := vbox.g2 - vbox.g1
	sub_b := vbox.b2 - vbox.b1

	return (sub_r + 1) * (sub_g + 1) * (sub_b + 1)
}

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

	fmt.Print(computeVBox(pixels, histogram))
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

func computeVBox(pixels []Pixel, histogram map[int]int) VBox {
	rmin := 1000000
	rmax := 0
	gmin := 1000000
	gmax := 0
	bmin := 1000000
	bmax := 0

	for _, pixel := range pixels {
		r := pixel.R >> RSHIFT
		g := pixel.G >> RSHIFT
		b := pixel.B >> RSHIFT

		rmin = min(r, rmin)
		rmax = max(r, rmax)
		gmin = min(g, gmin)
		gmax = max(g, gmax)
		bmin = min(b, bmin)
		bmax = max(b, bmax)
	}

	return VBox{r1: rmin, r2: rmax, g1: gmin, g2: gmax, b1: bmin, b2: bmax, histogram: histogram}
}

func getColorIndex(r int, g int, b int) int {
	return r<<(BITSIG*2) + g<<BITSIG + b
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
