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

func applyMedianCut(vbox *VBox, histogram map[int]int) (*VBox, *VBox) {
	if vbox.Count() == 0 {
		return nil, nil
	}

	// only one pixel, no split

	if vbox.Count() == 1 {
		return vbox, nil
	}

	rw := vbox.r2 - vbox.r1
	gw := vbox.g2 - vbox.g1
	bw := vbox.b2 - vbox.b1

	// finding the partial sum arrays along the selected axis

	var partialSum []int
	total := 0

	switch max(max(rw, gw), bw) {
	case rw:
		for i := vbox.r1; i <= vbox.r2; i++ {
			sum := 0

			for j := vbox.g1; j <= vbox.g2; j++ {
				for k := vbox.b1; k <= vbox.b2; k++ {
					if val, ok := vbox.histogram[getColorIndex(i, j, k)]; ok {
						sum += val
					}
				}
			}

			total += sum
			partialSum[i] = total
		}

	case gw:
		for i := vbox.g1; i <= vbox.g2; i++ {
			sum := 0

			for j := vbox.r1; j <= vbox.r2; j++ {
				for k := vbox.b1; k <= vbox.b2; k++ {
					if val, ok := vbox.histogram[getColorIndex(j, i, k)]; ok {
						sum += val
					}
				}
			}

			total += sum
			partialSum[i] = total
		}

	default:
		for i := vbox.b1; i <= vbox.b2; i++ {
			sum := 0

			for j := vbox.r1; j <= vbox.r2; j++ {
				for k := vbox.g1; k <= vbox.g2; k++ {
					if val, ok := vbox.histogram[getColorIndex(j, k, i)]; ok {
						sum += val
					}
				}
			}

			total += sum
			partialSum[i] = total
		}
	}

	lookAheadSum := make([]int, len(partialSum))

	for i := 0; i < len(partialSum); i++ {
		lookAheadSum[i] = total - partialSum[i]
	}

	// determining the cut planes

	return nil, nil
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
