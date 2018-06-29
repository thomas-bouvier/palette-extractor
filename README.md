palette-extractor
=================

This program extracts the dominant color or a representative color palette from an image.

## Usage

Here's a simple example, where we build a 5 color palette:

```go
package main

import (
	"fmt"
	"github.com/thomas-bouvier/palette-extractor"
)

func main() {
	// Creating the extractor object
	extractor := extractor.NewExtractor("image.png", 10)

    // Displaying the top 5 dominant colors of the image
	fmt.Println(extractor.GetPalette(5))
}
```

## Example

The following image has been used for this example:

![Example](image.png)

The program will give the following output when used with the image above:

```
[[234 231 230] [208 24 44] [59 41 37] [158 149 145] [145 126 114]]
```

## Thanks

Many thanks to [Lokesh Dhakar](https://github.com/lokesh) for [his original work](https://github.com/lokesh/color-thief/) and [Shipeng Feng](https://github.com/fengsp) for [his implementation](https://github.com/fengsp/color-thief-py).