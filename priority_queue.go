package extractor

import (
	"container/heap"
	"fmt"
)

type sortingStrategy int

const (
	count sortingStrategy = iota
	countTimesVolume
)

type priorityQueue struct {
	boxes           []*box
	sortingStrategy sortingStrategy
}

func newBoxes(strategy sortingStrategy) *priorityQueue {
	boxes := &priorityQueue{make([]*box, 0), strategy}
	heap.Init(boxes)
	return boxes
}

func (boxes priorityQueue) Len() int {
	return len(boxes.boxes)
}

func (boxes priorityQueue) Less(i, j int) bool {
	switch boxes.sortingStrategy {
	case count:
		return boxes.boxes[i].count() > boxes.boxes[j].count()
	case countTimesVolume:
		return boxes.boxes[i].count()*boxes.boxes[i].volume() > boxes.boxes[j].count()*boxes.boxes[j].volume()
	default:
		return boxes.boxes[i].count() > boxes.boxes[j].count()
	}
}

func (boxes priorityQueue) Swap(i, j int) {
	boxes.boxes[i], boxes.boxes[j] = boxes.boxes[j], boxes.boxes[i]
	boxes.boxes[i].index = i
	boxes.boxes[j].index = j
}

func (boxes *priorityQueue) Pop() interface{} {
	old := *boxes
	n := len(old.boxes)
	item := old.boxes[n-1]
	item.index = -1
	(*boxes).boxes = old.boxes[0 : n-1]
	return item
}

func (boxes *priorityQueue) Push(x interface{}) {
	n := len((*boxes).boxes)
	item := x.(*box)
	item.index = n
	(*boxes).boxes = append((*boxes).boxes, item)
}

func (boxes priorityQueue) print() {
	fmt.Println(fmt.Sprintf("Len: %d", boxes.Len()))
	for i := 0; i < boxes.Len(); i++ {
		boxes.boxes[i].print()
	}
}
