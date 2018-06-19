package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type Extractor struct {
	filename string
	quality  int
	pixels   []Pixel
}

func NewExtractor(filename string, quality int) *Extractor {
	extractor := &Extractor{}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

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

func (extractor *Extractor) GetPalette(count int) []Pixel {
	return quantize(extractor.pixels, count).GetPalette()
}

func (extractor *Extractor) GetColor() Pixel {
	return extractor.GetPalette(5)[0]
}
