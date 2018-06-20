package extractor

import (
	"container/heap"
	"fmt"
	"os"
)

const BITSIG = 5
const RSHIFT = 8 - BITSIG
const ITMAX = 1000
const FRACTPOPULATION = 0.75

func quantize(pixels []Pixel, count int) *CMap {
	if pixels == nil || len(pixels) == 0 {
		fmt.Fprintf(os.Stderr, "empty pixels when quantizing")
		os.Exit(1)
	}

	if count < 2 || count > 256 {
		fmt.Fprintf(os.Stderr, "wrong number of max colors when quantizing")
		os.Exit(1)
	}

	histogram := computeHistogram(pixels)
	if len(histogram) <= count {
		fmt.Fprintf(os.Stderr, "insufficient number of levels of quantification")
		os.Exit(1)
	}

	vbox := computeVBox(pixels, histogram)

	vboxes := NewVBoxes(Count)
	vboxes.Push(vbox)

	doQuantizeIteration(vboxes, histogram, float32(count)*FRACTPOPULATION)

	vboxes2 := NewVBoxes(CountTimesVolume)
	for vboxes.Len() > 0 {
		vboxes2.Push(heap.Pop(vboxes))
	}

	doQuantizeIteration(vboxes2, histogram, float32(count-vboxes2.Len()))

	cmap := NewCMap()
	for vboxes2.Len() > 0 {
		cmap.Push(heap.Pop(vboxes2).(*Box))
	}

	return cmap
}

func doQuantizeIteration(vboxes *PriorityQueue, histogram map[int]int, target float32) {
	nbColor := 1
	it := 0

	for it < ITMAX {
		vbox := heap.Pop(vboxes).(*Box)

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
		}

		if float32(nbColor) >= target {
			return
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

func applyMedianCut(vbox *Box, histogram map[int]int) (vbox1 *Box, vbox2 *Box, count int) {
	if vbox.Count() == 0 {
		return vbox1, vbox2, 0
	}

	// only one pixel, no split

	if vbox.Count() == 1 {
		vbox1 := vbox.Copy()
		return vbox1, nil, 1
	}

	rw := vbox.r2 - vbox.r1 + 1
	gw := vbox.g2 - vbox.g1 + 1
	bw := vbox.b2 - vbox.b1 + 1

	// finding the partial sum arrays along the selected axis

	partialSum := make(map[int]int)
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
					if val, ok := histogram[getColorIndex(i, j, k)]; ok {
						sum += val
					}
				}
			}

			total += sum
			partialSum[i] = total
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
			partialSum[i] = total
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
			partialSum[i] = total
		}

		dim1 = vbox.b1
		dim2 = vbox.b2
	}

	lookAheadSum := make(map[int]int, len(partialSum))
	for k, v := range partialSum {
		lookAheadSum[k] = total - v
	}

	// determining the cut planes

	for i := dim1; i <= dim2; i++ {
		if partialSum[i] > total/2 {
			vbox1 := vbox.Copy()
			vbox2 := vbox.Copy()

			l := i - dim1
			r := dim2 - i
			newDim := 0

			if l <= r {
				newDim = min(dim2-1, i+(r/2))
			} else {
				newDim = max(dim1, (i-1)-(l/2))
			}

			// avoid 0-count boxes

			ko := true
			for ko {
				if _, ok := partialSum[newDim]; !ok {
					newDim++
				} else {
					ko = false
				}
			}

			count2 := lookAheadSum[newDim]

			ko = true
			for ko {
				if _, ok := partialSum[newDim-1]; !ok {
					newDim--
					count2 = lookAheadSum[newDim]
				} else {
					if count2 == 0 {
						newDim--
						count2 = lookAheadSum[newDim]
					} else {
						ko = false
					}
				}
			}

			// set dimensions

			if br {
				vbox1.r2 = newDim
				vbox2.r1 = vbox1.r2 + 1
			} else if bg {
				vbox1.g2 = newDim
				vbox2.g1 = vbox1.g2 + 1
			} else {
				vbox1.b2 = newDim
				vbox2.b1 = vbox1.b2 + 1
			}

			return vbox1, vbox2, 2
		}
	}

	return vbox1, vbox2, 0
}

func computeVBox(pixels []Pixel, histogram map[int]int) *Box {
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

	return &Box{r1: rmin, r2: rmax, g1: gmin, g2: gmax, b1: bmin, b2: bmax, histogram: histogram}
}

func getColorIndex(r int, g int, b int) int {
	return r<<(BITSIG*2) + g<<BITSIG + b
}
