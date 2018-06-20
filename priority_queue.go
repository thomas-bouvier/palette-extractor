package extractor

import (
	"container/heap"
	"fmt"
)

type SortingStrategy int

const (
	Count SortingStrategy = iota
	CountTimesVolume
)

type PriorityQueue struct {
	boxes           []*Box
	sortingStrategy SortingStrategy
}

func NewVBoxes(strategy SortingStrategy) *PriorityQueue {
	vboxes := &PriorityQueue{make([]*Box, 0), strategy}
	heap.Init(vboxes)
	return vboxes
}

func (vboxes PriorityQueue) Len() int {
	return len(vboxes.boxes)
}

func (vboxes PriorityQueue) Less(i, j int) bool {
	switch vboxes.sortingStrategy {
	case Count:
		return vboxes.boxes[i].Count() > vboxes.boxes[j].Count()
	case CountTimesVolume:
		return vboxes.boxes[i].Count()*vboxes.boxes[i].volume() > vboxes.boxes[j].Count()*vboxes.boxes[j].volume()
	default:
		return vboxes.boxes[i].Count() > vboxes.boxes[j].Count()
	}
}

func (vboxes PriorityQueue) Swap(i, j int) {
	vboxes.boxes[i], vboxes.boxes[j] = vboxes.boxes[j], vboxes.boxes[i]
	vboxes.boxes[i].index = i
	vboxes.boxes[j].index = j
}

func (vboxes *PriorityQueue) Pop() interface{} {
	old := *vboxes
	n := len(old.boxes)
	item := old.boxes[n-1]
	item.index = -1
	(*vboxes).boxes = old.boxes[0 : n-1]
	return item
}

func (vboxes *PriorityQueue) Push(x interface{}) {
	n := len((*vboxes).boxes)
	item := x.(*Box)
	item.index = n
	(*vboxes).boxes = append((*vboxes).boxes, item)
}

func (vboxes PriorityQueue) Print() {
	fmt.Println(fmt.Sprintf("Len: %d", vboxes.Len()))
	for i := 0; i < vboxes.Len(); i++ {
		vboxes.boxes[i].Print()
	}
}
