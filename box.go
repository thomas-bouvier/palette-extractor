package main

import "fmt"

type Box struct {
	r1, r2    int
	g1, g2    int
	b1, b2    int
	histogram map[int]int
	color     *Pixel

	index int
}

func (box *Box) volume() int {
	subr := box.r2 - box.r1
	subg := box.g2 - box.g1
	subb := box.b2 - box.b1

	return (subr + 1) * (subg + 1) * (subb + 1)
}

func (box *Box) Count() int {
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

func (box *Box) Copy() *Box {
	rvbox := &Box{r1: box.r1, r2: box.r2, g1: box.g1, g2: box.g2, b1: box.b1, b2: box.b2}

	rvbox.histogram = make(map[int]int, len(box.histogram))
	for k, v := range box.histogram {
		rvbox.histogram[k] = v
	}

	return rvbox
}

func (box *Box) average() *Pixel {
	n := 0
	mult := 1 << (8 - BITSIG)
	pixel := &Pixel{0, 0, 0, 255}

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

func (box *Box) Contains(pixel *Pixel) bool {
	r := pixel.R >> RSHIFT
	g := pixel.G >> RSHIFT
	b := pixel.B >> RSHIFT

	return r >= box.r1 && r <= box.r2 && g >= box.g1 && g <= box.g2 && b >= box.b1 && b <= box.b2
}

func (box *Box) Print() {
	fmt.Println(box)
	fmt.Println(fmt.Sprintf("\tr1: %d, r2: %d", box.r1, box.r2))
	fmt.Println(fmt.Sprintf("\tg1: %d, g2: %d", box.g1, box.g2))
	fmt.Println(fmt.Sprintf("\tb1: %d, b2: %d", box.b1, box.b2))
	fmt.Println(fmt.Sprintf("Count: %d", box.Count()))
}
