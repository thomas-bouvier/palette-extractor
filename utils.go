package main

import "math"

type Pixel struct {
	R int
	G int
	B int
	A int
}

type SortingStrategy int

const (
	Count SortingStrategy = iota
	CountTimesVolume
)

type VBox struct {
	r1, r2    int
	g1, g2    int
	b1, b2    int
	histogram map[int]int

	index int
}

type VBoxes struct {
	boxes           []*VBox
	sortingStrategy SortingStrategy
}

type CMap struct {
	vboxes *VBoxes
	colors []Pixel
}

func (vbox *VBox) volume() int {
	sub_r := vbox.r2 - vbox.r1
	sub_g := vbox.g2 - vbox.g1
	sub_b := vbox.b2 - vbox.b1

	return (sub_r + 1) * (sub_g + 1) * (sub_b + 1)
}

func (vbox *VBox) Count() int {
	n := 0

	for i := vbox.r1; i <= vbox.r2; i++ {
		for j := vbox.g1; j <= vbox.g2; j++ {
			for k := vbox.b1; k <= vbox.b2; k++ {
				if val, ok := vbox.histogram[getColorIndex(i, j, k)]; ok {
					n += val
				}
			}
		}
	}

	return n
}

func (vbox *VBox) average() Pixel {
	n := 0
	mult := 1 << (8 - BITSIG)
	pixel := Pixel{0, 0, 0, 255}

	for i := vbox.r1; i <= vbox.r2; i++ {
		for j := vbox.g1; j <= vbox.g2; j++ {
			for k := vbox.b1; k <= vbox.b2; k++ {
				if val, ok := vbox.histogram[getColorIndex(i, j, k)]; ok {
					n += val

					pixel.R += val*i*mult + val*i*mult/2
					pixel.G += val*j*mult + val*j*mult/2
					pixel.B += val*k*mult + val*k*mult/2
				}
			}
		}
	}

	if n > 0 {
		pixel.R /= n
		pixel.G /= n
		pixel.B /= n
	} else {
		pixel.R = mult * (vbox.r1 + vbox.r2 + 1) / 2
		pixel.G = mult * (vbox.g1 + vbox.g2 + 1) / 2
		pixel.B = mult * (vbox.b1 + vbox.b2 + 1) / 2
	}

	return pixel
}

func (vbox *VBox) Contains(pixel *Pixel) bool {
	r := pixel.R >> RSHIFT
	g := pixel.G >> RSHIFT
	b := pixel.B >> RSHIFT

	return r >= vbox.r1 && r <= vbox.r2 && g >= vbox.g1 && g <= vbox.g2 && b >= vbox.b1 && b <= vbox.b2
}

func (vboxes VBoxes) Len() int {
	return len(vboxes.boxes)
}

func (vboxes VBoxes) Less(i, j int) bool {
	switch vboxes.sortingStrategy {
	case Count:
		return vboxes.boxes[i].Count() < vboxes.boxes[j].Count()
	case CountTimesVolume:
		return vboxes.boxes[i].Count()*vboxes.boxes[i].volume() < vboxes.boxes[j].Count()*vboxes.boxes[j].volume()
	default:
		return vboxes.boxes[i].Count() < vboxes.boxes[j].Count()
	}
}

func (vboxes VBoxes) Swap(i, j int) {
	vboxes.boxes[i], vboxes.boxes[j] = vboxes.boxes[j], vboxes.boxes[i]
	vboxes.boxes[i].index = i
	vboxes.boxes[j].index = j
}

func (vboxes *VBoxes) Pop() interface{} {
	old := *vboxes
	n := len(old.boxes)
	item := old.boxes[n-1]
	item.index = -1
	(*vboxes).boxes = old.boxes[0 : n-1]
	return item
}

func (vboxes *VBoxes) Push(x interface{}) {
	n := len((*vboxes).boxes)
	item := x.(*VBox)
	item.index = n
	(*vboxes).boxes = append((*vboxes).boxes, item)
}

func (cmap *CMap) Push(x interface{}) {
	item := x.(*VBox)
	cmap.vboxes.Push(item)
	cmap.colors = append(cmap.colors, item.average())
}

func (cmap *CMap) Map(color Pixel) *Pixel {
	for i := 0; i < cmap.vboxes.Len(); i++ {
		vbox := cmap.vboxes.boxes[i]

		if vbox.Contains(&color) {
			return &cmap.colors[i]
		}
	}

	return cmap.nearest(color)
}

func (cmap *CMap) nearest(color Pixel) *Pixel {
	d1 := -10000.0
	var ret *Pixel

	for i := 0; i < cmap.vboxes.Len(); i++ {
		d2 := math.Sqrt(
			math.Pow(float64(color.R-cmap.colors[i].R), 2) +
				math.Pow(float64(color.G-cmap.colors[i].G), 2) +
				math.Pow(float64(color.B-cmap.colors[i].B), 2))

		if d2 < d1 {
			d1 = d2
			ret = &cmap.colors[i]
		}
	}

	return ret
}

func (cmap *CMap) Len() int {
	return cmap.vboxes.Len()
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
