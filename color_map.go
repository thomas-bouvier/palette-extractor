package extractor

import (
	"container/heap"
	"math"
)

type colorMap struct {
	boxes priorityQueue
}

func newColorMap() *colorMap {
	colorMap := &colorMap{}
	colorMap.boxes = priorityQueue{make([]*box, 0), countTimesVolume}
	heap.Init(&colorMap.boxes)
	return colorMap
}

func (colorMap *colorMap) push(vbox *box) {
	vbox.color = vbox.average()
	colorMap.boxes.Push(vbox)
}

func (colorMap *colorMap) mapColor(color *pixel) *pixel {
	for i := 0; i < colorMap.boxes.Len(); i++ {
		box := colorMap.boxes.boxes[i]

		if box.contains(color) {
			return box.color
		}
	}

	return colorMap.getNearestColor(color)
}

func (colorMap *colorMap) getNearestColor(color *pixel) *pixel {
	d1 := -10000.0
	var ret *pixel

	for i := 0; i < colorMap.boxes.Len(); i++ {
		d2 := math.Sqrt(math.Pow(float64(color.R-colorMap.boxes.boxes[i].color.R), 2) + math.Pow(float64(color.G-colorMap.boxes.boxes[i].color.G), 2) + math.Pow(float64(color.B-colorMap.boxes.boxes[i].color.B), 2))

		if d2 < d1 {
			d1 = d2
			ret = colorMap.boxes.boxes[i].color
		}
	}

	return ret
}

func (colorMap *colorMap) len() int {
	return colorMap.boxes.Len()
}

func (colorMap *colorMap) getPalette() []pixel {
	var pixels []pixel

	for i := 0; i < colorMap.len(); i++ {
		pixels = append(pixels, *colorMap.boxes.boxes[i].color)
	}

	return pixels
}
