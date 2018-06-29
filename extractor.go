// Package extractor extracts the dominant color or a representative color palette
// from an image.
package extractor

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// Extractor represents the color palette extractor instance associated with a single image.
type Extractor struct {
	filename string
	quality  int
	pixels   []pixel
}

// NewExtractor returns a new instance of an Extractor.
// The provider filename parameter can contain a path, and must include the file extension.
// The provided quality parameter behaves in the following manner:
// the bigger the number, the faster a color will be returned but the greater the likelihood
// that it will not be the visually most dominant color. 1 then refers the highest quality.
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

// GetPalette builds a color palette.
// We are using the median cut algorithm to cluster similar colors.
// The provided count parameter defines how many colors should be extracted.
// Each color is returned as [r g b].
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

// GetColor selects the dominant color.
// It corresponds to the first color in the palette.
// The color is returned as [r g b].
func (extractor *Extractor) GetColor() []int {
	return extractor.GetPalette(5)[0]
}
