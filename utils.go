package main

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
