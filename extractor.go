package extractor

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Extractor struct {
	filename string
	quality  int
	pixels   []pixel
}

func NewExtractor(filename string, quality int) *Extractor {
	extractor := &Extractor{}

	reader, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer reader.Close()

	pixels, err := getPixels(reader, quality)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	extractor.filename = filename
	extractor.quality = quality
	extractor.pixels = pixels

	return extractor
}

func (extractor *Extractor) GetPalette(count int) [][]int {
	ret := make([][]int, count)
	for i := range ret {
		ret[i] = make([]int, 3)
	}

	pixels := quantize(extractor.pixels, count).getPalette()
	for i := 0; i < count; i++ {
		ret[i][0] = pixels[i].R
		ret[i][1] = pixels[i].G
		ret[i][2] = pixels[i].B
	}

	return ret
}

func (extractor *Extractor) GetColor() []int {
	return extractor.GetPalette(5)[0]
}
