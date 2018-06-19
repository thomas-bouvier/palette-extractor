package main

import "fmt"

func main() {
	// Creating an Extractor object
	extractor := NewExtractor("image.png", 10)

	// Retrieving the associated color palette
	fmt.Println(extractor.GetPalette(6))

	// Retrieving the dominant color
	fmt.Println(extractor.GetColor())
}
