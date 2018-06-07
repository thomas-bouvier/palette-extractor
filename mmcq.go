package main

import (
	"container/heap"
	"fmt"
	"github.com/jinzhu/copier"
	"os"
)

const BITSIG = 5
const RSHIFT = 8 - BITSIG
const ITMAX = 1000
const FRACTPOPULATION = 0.75

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

	doQuantizeIteration(&vboxes, &histogram, float32(count)*FRACTPOPULATION)
}

func doQuantizeIteration(vboxes *VBoxes, histogram *map[int]int, target float32) {
	nbColor := 1
	it := 0

	for it < ITMAX {
		vbox := heap.Pop(vboxes).(*VBox)
		fmt.Print(it)

		if vbox.Count() == 0 {
			heap.Push(vboxes, vbox)

			it++
			continue
		}

		// do the cut

		vbox1, vbox2, count := applyMedianCut(vbox, histogram)

		if count > 0 {
			heap.Push(vboxes, vbox1)

			if count > 1 {
				heap.Push(vboxes, vbox2)
				nbColor++
			}

			if float32(nbColor) >= target {
				return
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

func applyMedianCut(vbox *VBox, histogram *map[int]int) (vbox1 VBox, vbox2 VBox, count int) {
	if vbox.Count() == 0 {
		return vbox1, vbox2, 0
	}

	// only one pixel, no split

	if vbox.Count() == 1 {
		return *vbox, vbox2, 1
	}

	rw := vbox.r2 - vbox.r1
	gw := vbox.g2 - vbox.g1
	bw := vbox.b2 - vbox.b1

	// finding the partial sum arrays along the selected axis

	var partialSum []int
	total := 0

	dim1, dim2 := 0, 0
	br, bg := false, false

	switch max(max(rw, gw), bw) {
	case rw:
		br = true

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
			partialSum = append(partialSum, total)
		}

		dim1 = vbox.r1
		dim2 = vbox.r2

	case gw:
		bg = true

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
			partialSum = append(partialSum, total)
		}

		dim1 = vbox.g1
		dim2 = vbox.g2

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
			partialSum = append(partialSum, total)
		}

		dim1 = vbox.b1
		dim2 = vbox.b2
	}

	lookAheadSum := make([]int, len(partialSum))

	for i := 0; i < len(partialSum); i++ {
		lookAheadSum[i] = total - partialSum[i]
	}

	// determining the cut planes

	for i := dim1; i <= dim2; i++ {
		if partialSum[i] > total/2 {
			vbox1, vbox2 := VBox{}, VBox{}
			copier.Copy(&vbox1, vbox)
			copier.Copy(&vbox2, vbox)

			l := i - dim1
			r := dim2 - i
			new_dim := 0

			if l <= r {
				new_dim = min(dim2-1, i+(r/2))
			} else {
				new_dim = max(dim1, (i-1)-(l/2))
			}

			// avoid 0-count boxes

			for partialSum[new_dim] == 0 {
				new_dim++
			}

			count2 := lookAheadSum[new_dim]

			for !(count2 == 0 && partialSum[new_dim-1] == 0) {
				new_dim--
				count2 = lookAheadSum[new_dim]
			}

			// set dimensions

			if br {
				vbox1.r2 = new_dim
				vbox2.r1 = vbox1.r2 + 1
			} else if bg {
				vbox1.g2 = new_dim
				vbox2.g1 = vbox1.g2 + 1
			} else {
				vbox1.b2 = new_dim
				vbox2.b1 = vbox1.b2 + 1
			}

			return vbox1, vbox2, 2
		}
	}

	return vbox1, vbox2, 0
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
