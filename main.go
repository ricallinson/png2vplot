package main

import(
    "fmt"
    "strconv"
)

// Five shades of gray.
func pixel(x int, y int, shade int) (pixel string) {
    size := 10
    for xoff := 0; xoff < size; xoff = xoff + shade {
        pixel += "m " + strconv.Itoa(x + xoff) + " " + strconv.Itoa(y) + "\n";
        pixel += "m " + strconv.Itoa(x + xoff) + " " + strconv.Itoa(y + size) + "\n";
    }
    pixel += "m " + strconv.Itoa(x + size) + " " + strconv.Itoa(y) + "\n";
    return pixel
}

func main() {
    fmt.Println(pixel(0, 0, 5))
}
