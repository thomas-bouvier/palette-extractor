package main

import (
	"container/heap"
	"fmt"
	"os"
)

const BITSIG = 5
const RSHIFT = 8 - BITSIG
const ITMAX = 1000

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

	vbox := computeVBox(pixels, histogram)

	vboxes := make(VBoxes, 1)
	vboxes[0] = &vbox

	heap.Init(&vboxes)

	doQuantizeIteration(&vboxes)
}

func doQuantizeIteration(vboxes *VBoxes) {
	it := 0

	for it < ITMAX {
		if vboxes.Len() > 0 {
			vbox := heap.Pop(vboxes).(*VBox)

			if vbox.Count() == 0 {
				fmt.Print("test")
				heap.Push(vboxes, vbox)
			}
		}

		it++
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

func computeVBox(pixels []Pixel, histogram map[int]int) VBox {
	rmin, rmax := int(^uint(0)>>1), 0
	gmin, gmax := int(^uint(0)>>1), 0
	bmin, bmax := int(^uint(0)>>1), 0

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
