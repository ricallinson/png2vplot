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
	"strings"
)

var SHADE_MAX = 16
var SHADE_MIN = 1

func getShade(c color.Color, pixelsize int) int {
	r, g, b, _ := c.RGBA()
	shade := int((r+g+b)/3) / 4000 // Max is 65535
	if shade > SHADE_MAX {
		return SHADE_MAX
	}
	if shade < SHADE_MIN {
		return SHADE_MIN
	}
	return shade
}

func pixelQuake(x int, y int, shade int, pixelsize int, dir bool) (cmds string) {
	pixelhalf := pixelsize / 2
	offset := shade + 4%pixelsize
	down := true
	if dir {
		cmds += "M " + strconv.Itoa(x) + " " + strconv.Itoa(y+pixelhalf) + "\n"
		if shade >= SHADE_MAX {
			cmds += "L " + strconv.Itoa(x+pixelsize) + " " + strconv.Itoa(y+pixelhalf) + "\n"
			return cmds
		}
		for xoff := x; xoff < x+pixelsize; xoff = xoff + offset {
			if down {
				cmds += "L " + strconv.Itoa(xoff) + " " + strconv.Itoa(y+pixelsize) + "\n"
				down = false
			} else {
				cmds += "L " + strconv.Itoa(xoff) + " " + strconv.Itoa(y) + "\n"
				down = true
			}
		}
		cmds += "L " + strconv.Itoa(x+pixelsize) + " " + strconv.Itoa(y+pixelhalf) + "\n"
		dir = false
	} else {
		cmds += "M " + strconv.Itoa(x+pixelsize) + " " + strconv.Itoa(y+pixelhalf) + "\n"
		if shade >= SHADE_MAX {
			cmds += "L " + strconv.Itoa(x) + " " + strconv.Itoa(y+pixelhalf) + "\n"
			return cmds
		}
		for xoff := x + pixelsize; xoff > x; xoff = xoff - offset {
			if down {
				cmds += "L " + strconv.Itoa(xoff) + " " + strconv.Itoa(y+pixelsize) + "\n"
				down = false
			} else {
				cmds += "L " + strconv.Itoa(xoff) + " " + strconv.Itoa(y) + "\n"
				down = true
			}
		}
		cmds += "L " + strconv.Itoa(x) + " " + strconv.Itoa(y+pixelhalf) + "\n"
		dir = true
	}
	return cmds
}

func convert(file *os.File, xoffset int, yoffset int, pixelsize int, pause int) (string, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	p := pause
	dir := true
	plots := ""
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if dir {
			plots += "M " + strconv.Itoa(xoffset) + " " + strconv.Itoa((y*pixelsize)+yoffset) + "\n"
			for x := bounds.Min.X; x < bounds.Max.X; x = x + 1 {
				shade := getShade(img.At(x, y), pixelsize)
				plots += pixelQuake((x*pixelsize)+xoffset, (y*pixelsize)+yoffset, shade, pixelsize, dir)
			}
			dir = false
		} else {
			plots += "M " + strconv.Itoa((bounds.Max.X*pixelsize)+xoffset) + " " + strconv.Itoa((y*pixelsize)+yoffset) + "\n"
			for x := bounds.Max.X - 1; x >= bounds.Min.X; x = x - 1 {
				shade := getShade(img.At(x, y), pixelsize)
				plots += pixelQuake((x*pixelsize)+xoffset, (y*pixelsize)+yoffset, shade, pixelsize, dir)
			}
			dir = true
		}
		p--
		if p == 0 {
			plots += "P\n"
			p = pause
		}
	}
	plots += "h\n"
	return plots, nil
}

func main() {

	var xoffset = flag.Int("x", 0, "the offset to use for the X dimension")
	var yoffset = flag.Int("y", 0, "the offset to use for the Y dimension")
	var pixelsize = flag.Int("p", 50, "the size of pixel to plot (10 = 1mm)")
	var pause = flag.Int("w", 1000, "wait for input before continuing after the given number of rows")

	flag.Parse()
	source := flag.Arg(0)
	dest := flag.Arg(1)

	if source == "" {
		fmt.Println("A source jpeg file must be provide as the first argument.")
		return
	}

	if dest == "" {
		dest = strings.Replace(source, ".png", ".vplot", 1)
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
	plots, err3 := convert(sFile, *xoffset, *yoffset, *pixelsize, *pause)
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	w.WriteString(plots)
	w.Flush()

	fmt.Println(dest)
}
