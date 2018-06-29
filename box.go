package extractor

import (
	"fmt"
)

type box struct {
	r1, r2    int
	g1, g2    int
	b1, b2    int
	histogram map[int]int
	color     *pixel

	index int
}

func newBox(r1 int, r2 int, g1 int, g2 int, b1 int, b2 int) *box {
	return &box{r1: r1, r2: r2, g1: g1, g2: g2, b1: b1, b2: b2}
}

func (box *box) volume() int {
	subr := box.r2 - box.r1
	subg := box.g2 - box.g1
	subb := box.b2 - box.b1

	return (subr + 1) * (subg + 1) * (subb + 1)
}

func (box *box) count() int {
	n := 0

	for i := box.r1; i <= box.r2; i++ {
		for j := box.g1; j <= box.g2; j++ {
			for k := box.b1; k <= box.b2; k++ {
				if val, ok := box.histogram[getColorIndex(i, j, k)]; ok {
					n += val
				}
			}
		}
	}

	return n
}

func (box *box) copy() *box {
	rbox := newBox(box.r1, box.r2, box.g1, box.g2, box.b1, box.b2)

	rbox.histogram = make(map[int]int, len(box.histogram))
	for k, v := range box.histogram {
		rbox.histogram[k] = v
	}

	return rbox
}

func (box *box) average() *pixel {
	n := 0
	mult := 1 << (8 - bitsig)
	pixel := &pixel{0, 0, 0, 255}

	for i := box.r1; i <= box.r2; i++ {
		for j := box.g1; j <= box.g2; j++ {
			for k := box.b1; k <= box.b2; k++ {
				if val, ok := box.histogram[getColorIndex(i, j, k)]; ok {
					n += val

					pixel.R += val*i*mult + val*mult/2
					pixel.G += val*j*mult + val*mult/2
					pixel.B += val*k*mult + val*mult/2
				}
			}
		}
	}

	if n > 0 {
		pixel.R /= n
		pixel.G /= n
		pixel.B /= n
	} else {
		pixel.R = mult * (box.r1 + box.r2 + 1) / 2
		pixel.G = mult * (box.g1 + box.g2 + 1) / 2
		pixel.B = mult * (box.b1 + box.b2 + 1) / 2
	}

	return pixel
}

func (box *box) contains(pixel *pixel) bool {
	r := pixel.R >> rshift
	g := pixel.G >> rshift
	b := pixel.B >> rshift

	return r >= box.r1 && r <= box.r2 && g >= box.g1 && g <= box.g2 && b >= box.b1 && b <= box.b2
}

func (box *box) print() {
	fmt.Println(box)
	fmt.Println(fmt.Sprintf("\tr1: %d, r2: %d", box.r1, box.r2))
	fmt.Println(fmt.Sprintf("\tg1: %d, g2: %d", box.g1, box.g2))
	fmt.Println(fmt.Sprintf("\tb1: %d, b2: %d", box.b1, box.b2))
	fmt.Println(fmt.Sprintf("Count: %d", box.count()))
}
