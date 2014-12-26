package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"strconv"
)

// This pixel size should be a cli arg and should ajust the shading min/max.
// 50: Good for a "Fine Faber-Castell" as each pixel is 5mm x 5mm.
// 100: Good for chalk as each pixel is 10mm x 10mm.
var pixelSize = 50 
// These should not be here.
var START_X = 1500 // hack to match the cords of my plotter.
var START_Y = 2000 // hack to match the cords of my plotter.

func pixel(x int, y int, c color.Color) (pixel string) {
	r, g, b, _ := c.RGBA()
	shade := int((r>>12 + g>>12 + b>>12) / 3)
	fmt.Println(shade)
	if shade >= 15 {
		// Not too light.
		return "M " + strconv.Itoa(x+pixelSize) + " " + strconv.Itoa(y) + "\n"
	}
	if shade <= 4 {
		// Not too dark.
		shade = 4
	}
	dir := true
	for xoff := x; xoff < x+pixelSize; xoff = xoff + shade {
		pixel += "L " + strconv.Itoa(xoff) + " " + strconv.Itoa(y) + "\n"
		if dir {
			pixel += "L " + strconv.Itoa(xoff) + " " + strconv.Itoa(y+pixelSize) + "\n"
			dir = false
		} else {
			pixel += "L " + strconv.Itoa(xoff) + " " + strconv.Itoa(y) + "\n"
			dir = true
		}
	}
	return pixel
}

func convert(file *os.File) (string, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	plots := ""
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        plots += "M " + strconv.Itoa(START_X) + " " + strconv.Itoa((y*pixelSize) + START_Y) + "\n"
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			plots += pixel((x*pixelSize) + START_X, (y*pixelSize) + START_Y, img.At(x, y))
		}
	}
	return plots, nil
}

func main() {

	flag.Parse()
	source := flag.Arg(0)
	dest := flag.Arg(1)

	if source == "" {
		fmt.Println("A source jpeg file must be provide as the first argument.")
		return
	}

	if dest == "" {
		fmt.Println("A destination vplot file must be provide as the second argument.")
		return
	}

	sFile, err1 := os.Open(source)
	if err1 != nil {
		fmt.Println("Could not open the source jpeg file.")
		return
	}
	defer sFile.Close()

	dFile, err2 := os.Create(dest)
	if err2 != nil {
		fmt.Println("Could not create the destination vplot file.")
		return
	}
	defer sFile.Close()

	w := bufio.NewWriter(dFile)
	plots, err3 := convert(sFile)
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	w.WriteString(plots)
	w.Flush()
}
