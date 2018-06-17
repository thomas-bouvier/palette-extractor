package main

import (
	"container/heap"
	"fmt"
	"math"
)

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
	color     *Pixel

	index int
}

type VBoxes struct {
	boxes           []*VBox
	sortingStrategy SortingStrategy
}

type CMap struct {
	vboxes VBoxes
}

func (vbox *VBox) volume() int {
	subr := vbox.r2 - vbox.r1
	subg := vbox.g2 - vbox.g1
	subb := vbox.b2 - vbox.b1

	return (subr + 1) * (subg + 1) * (subb + 1)
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

func (vbox *VBox) Copy() *VBox {
	rvbox := &VBox{r1: vbox.r1, r2: vbox.r2, g1: vbox.g1, g2: vbox.g2, b1: vbox.b1, b2: vbox.b2}

	rvbox.histogram = make(map[int]int, len(vbox.histogram))
	for k, v := range vbox.histogram {
		rvbox.histogram[k] = v
	}

	return rvbox
}

func (vbox *VBox) Print() {
	fmt.Println(vbox)
	fmt.Println(fmt.Sprintf("\tr1: %d, r2: %d", vbox.r1, vbox.r2))
	fmt.Println(fmt.Sprintf("\tg1: %d, g2: %d", vbox.g1, vbox.g2))
	fmt.Println(fmt.Sprintf("\tb1: %d, b2: %d", vbox.b1, vbox.b2))
	fmt.Println(fmt.Sprintf("Count: %d", vbox.Count()))
}

func (vbox *VBox) average() *Pixel {
	n := 0
	mult := 1 << (8 - BITSIG)
	pixel := &Pixel{0, 0, 0, 255}

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

func NewVBoxes(strategy SortingStrategy) *VBoxes {
	vboxes := &VBoxes{make([]*VBox, 0), strategy}
	heap.Init(vboxes)
	return vboxes
}

func (vboxes VBoxes) Len() int {
	return len(vboxes.boxes)
}

func (vboxes VBoxes) Less(i, j int) bool {
	switch vboxes.sortingStrategy {
	case Count:
		return vboxes.boxes[i].Count() > vboxes.boxes[j].Count()
	case CountTimesVolume:
		return vboxes.boxes[i].Count()*vboxes.boxes[i].volume() < vboxes.boxes[j].Count()*vboxes.boxes[j].volume()
	default:
		return vboxes.boxes[i].Count() > vboxes.boxes[j].Count()
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

func (vboxes VBoxes) Print() {
	fmt.Println(fmt.Sprintf("Len: %d", vboxes.Len()))
	for i := 0; i < vboxes.Len(); i++ {
		vboxes.boxes[i].Print()
	}
}

func NewCMap() *CMap {
	cmap := &CMap{}
	cmap.vboxes = VBoxes{make([]*VBox, 0), CountTimesVolume}
	heap.Init(&cmap.vboxes)
	return cmap
}

func (cmap *CMap) Push(vbox *VBox) {
	vbox.color = vbox.average()
	cmap.vboxes.Push(vbox)
}

func (cmap *CMap) Map(color *Pixel) *Pixel {
	for i := 0; i < cmap.vboxes.Len(); i++ {
		vbox := cmap.vboxes.boxes[i]

		if vbox.Contains(color) {
			return vbox.color
		}
	}

	return cmap.nearest(color)
}

func (cmap *CMap) nearest(color *Pixel) *Pixel {
	d1 := -10000.0
	var ret *Pixel

	for i := 0; i < cmap.vboxes.Len(); i++ {
		d2 := math.Sqrt(math.Pow(float64(color.R-cmap.vboxes.boxes[i].color.R), 2) + math.Pow(float64(color.G-cmap.vboxes.boxes[i].color.G), 2) + math.Pow(float64(color.B-cmap.vboxes.boxes[i].color.B), 2))

		if d2 < d1 {
			d1 = d2
			ret = cmap.vboxes.boxes[i].color
		}
	}

	return ret
}

func (cmap *CMap) Len() int {
	return cmap.vboxes.Len()
}

func (cmap *CMap) GetPalette() []Pixel {
	var pixels []Pixel

	for i := 0; i < cmap.Len(); i++ {
		pixels = append(pixels, *cmap.vboxes.boxes[i].color)
	}

	return pixels
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
