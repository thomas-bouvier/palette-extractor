package extractor

import (
	"container/heap"
	"fmt"
	"os"
)

const bitsig = 5
const rshift = 8 - bitsig
const itmax = 1000
const fractpopulation = 0.75

func quantize(pixels []pixel, maxColor int) *colorMap {
	if pixels == nil || len(pixels) == 0 {
		fmt.Fprintf(os.Stderr, "empty pixels when quantizing")
		os.Exit(1)
	}

	if maxColor < 2 || maxColor > 256 {
		fmt.Fprintf(os.Stderr, "wrong number of max colors when quantizing")
		os.Exit(1)
	}

	histogram := computeHistogram(pixels)
	if len(histogram) <= maxColor {
		fmt.Fprintf(os.Stderr, "insufficient number of levels of quantification")
		os.Exit(1)
	}

	boxes := newBoxes(count)
	boxes.Push(computeBox(pixels, histogram))

	doQuantizeIteration(boxes, histogram, float32(maxColor)*fractpopulation)

	boxes2 := newBoxes(countTimesVolume)
	for boxes.Len() > 0 {
		boxes2.Push(heap.Pop(boxes))
	}

	doQuantizeIteration(boxes2, histogram, float32(maxColor-boxes2.Len()))

	colorMap := newColorMap()
	for boxes2.Len() > 0 {
		colorMap.push(heap.Pop(boxes2).(*box))
	}

	return colorMap
}

func doQuantizeIteration(boxes *priorityQueue, histogram map[int]int, target float32) {
	nbColor := 1
	it := 0

	for it < itmax {
		box := heap.Pop(boxes).(*box)

		if box.count() == 0 {
			heap.Push(boxes, box)

			it++
			continue
		}

		// do the cut

		box1, box2, count := applyMedianCut(box, histogram)

		if count > 0 {
			heap.Push(boxes, box1)

			if count > 1 {
				heap.Push(boxes, box2)
				nbColor++
			}
		}

		if float32(nbColor) >= target {
			return
		}

		it++
	}
}

func computeHistogram(pixels []pixel) map[int]int {
	var index int
	histogram := make(map[int]int)

	for _, pixel := range pixels {
		index = getColorIndex(pixel.R>>rshift, pixel.G>>rshift, pixel.B>>rshift)
		if val, ok := histogram[index]; ok {
			histogram[index] = val + 1
		} else {
			histogram[index] = 1
		}
	}

	return histogram
}

func applyMedianCut(box *box, histogram map[int]int) (box1 *box, box2 *box, count int) {
	if box.count() == 0 {
		return box1, box2, 0
	}

	// only one pixel, no split

	if box.count() == 1 {
		box1 := box.copy()
		return box1, nil, 1
	}

	rw := box.r2 - box.r1 + 1
	gw := box.g2 - box.g1 + 1
	bw := box.b2 - box.b1 + 1

	// finding the partial sum arrays along the selected axis

	partialSum := make(map[int]int)
	total := 0

	dim1, dim2 := 0, 0
	br, bg := false, false

	switch max(max(rw, gw), bw) {
	case rw:
		br = true

		for i := box.r1; i <= box.r2; i++ {
			sum := 0

			for j := box.g1; j <= box.g2; j++ {
				for k := box.b1; k <= box.b2; k++ {
					if val, ok := histogram[getColorIndex(i, j, k)]; ok {
						sum += val
					}
				}
			}

			total += sum
			partialSum[i] = total
		}

		dim1 = box.r1
		dim2 = box.r2

	case gw:
		bg = true

		for i := box.g1; i <= box.g2; i++ {
			sum := 0

			for j := box.r1; j <= box.r2; j++ {
				for k := box.b1; k <= box.b2; k++ {
					if val, ok := box.histogram[getColorIndex(j, i, k)]; ok {
						sum += val
					}
				}
			}

			total += sum
			partialSum[i] = total
		}

		dim1 = box.g1
		dim2 = box.g2

	default:
		for i := box.b1; i <= box.b2; i++ {
			sum := 0

			for j := box.r1; j <= box.r2; j++ {
				for k := box.g1; k <= box.g2; k++ {
					if val, ok := box.histogram[getColorIndex(j, k, i)]; ok {
						sum += val
					}
				}
			}

			total += sum
			partialSum[i] = total
		}

		dim1 = box.b1
		dim2 = box.b2
	}

	lookAheadSum := make(map[int]int, len(partialSum))
	for k, v := range partialSum {
		lookAheadSum[k] = total - v
	}

	// determining the cut planes

	for i := dim1; i <= dim2; i++ {
		if partialSum[i] > total/2 {
			box1 := box.copy()
			box2 := box.copy()

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
				box1.r2 = newDim
				box2.r1 = box1.r2 + 1
			} else if bg {
				box1.g2 = newDim
				box2.g1 = box1.g2 + 1
			} else {
				box1.b2 = newDim
				box2.b1 = box1.b2 + 1
			}

			return box1, box2, 2
		}
	}

	return box1, box2, 0
}

func computeBox(pixels []pixel, histogram map[int]int) *box {
	rmin, rmax := int(^uint(0)>>1), 0
	gmin, gmax := int(^uint(0)>>1), 0
	bmin, bmax := int(^uint(0)>>1), 0

	for _, pixel := range pixels {
		r := pixel.R >> rshift
		g := pixel.G >> rshift
		b := pixel.B >> rshift

		rmin = min(r, rmin)
		rmax = max(r, rmax)
		gmin = min(g, gmin)
		gmax = max(g, gmax)
		bmin = min(b, bmin)
		bmax = max(b, bmax)
	}

	return &box{r1: rmin, r2: rmax, g1: gmin, g2: gmax, b1: bmin, b2: bmax, histogram: histogram}
}

func getColorIndex(r int, g int, b int) int {
	return r<<(bitsig*2) + g<<bitsig + b
}
