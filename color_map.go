package main

import (
	"container/heap"
	"math"
)

type CMap struct {
	boxes PriorityQueue
}

func NewCMap() *CMap {
	cmap := &CMap{}
	cmap.boxes = PriorityQueue{make([]*Box, 0), CountTimesVolume}
	heap.Init(&cmap.boxes)
	return cmap
}

func (cmap *CMap) Push(vbox *Box) {
	vbox.color = vbox.average()
	cmap.boxes.Push(vbox)
}

func (cmap *CMap) Map(color *Pixel) *Pixel {
	for i := 0; i < cmap.boxes.Len(); i++ {
		vbox := cmap.boxes.boxes[i]

		if vbox.Contains(color) {
			return vbox.color
		}
	}

	return cmap.nearest(color)
}

func (cmap *CMap) nearest(color *Pixel) *Pixel {
	d1 := -10000.0
	var ret *Pixel

	for i := 0; i < cmap.boxes.Len(); i++ {
		d2 := math.Sqrt(math.Pow(float64(color.R-cmap.boxes.boxes[i].color.R), 2) + math.Pow(float64(color.G-cmap.boxes.boxes[i].color.G), 2) + math.Pow(float64(color.B-cmap.boxes.boxes[i].color.B), 2))

		if d2 < d1 {
			d1 = d2
			ret = cmap.boxes.boxes[i].color
		}
	}

	return ret
}

func (cmap *CMap) Len() int {
	return cmap.boxes.Len()
}

func (cmap *CMap) GetPalette() []Pixel {
	var pixels []Pixel

	for i := 0; i < cmap.Len(); i++ {
		pixels = append(pixels, *cmap.boxes.boxes[i].color)
	}

	return pixels
}
