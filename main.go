package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"strconv"
)

var pixelSize = 10

func pixel(x int, y int, shade int) (pixel string) {
	if shade < 1 {
		shade = 1
	}
	dir := true
	for xoff := x; xoff < x+pixelSize; xoff = xoff + shade {
		pixel += "M " + strconv.Itoa(xoff) + " " + strconv.Itoa(y) + "\n"
		if dir {
			pixel += "M " + strconv.Itoa(xoff) + " " + strconv.Itoa(y+pixelSize) + "\n"
			dir = false
		} else {
			pixel += "M " + strconv.Itoa(xoff) + " " + strconv.Itoa(y) + "\n"
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
        plots += "L 0 " + strconv.Itoa(y*pixelSize) + "\n"
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			shade := int((r>>13 + g>>13 + b>>13) / 3)
			// fmt.Println(shade)
			plots += pixel(x*pixelSize, y*pixelSize, shade)
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
