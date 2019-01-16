package main

import (
	"fmt"
	"golang.org/x/tour/pic"
)

func Pic(dx, dy int) [][]uint8 {
	fmt.Println("Size", dx, dy)

	var result [][]uint8
	result = make([][]uint8, dy, dy)
	for y := 0; y < dy; y++ {
		result[y] = make([]uint8, dx, dx)
	}

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			result[x][y] = uint8(x ^ y)
		}
	}

	return result
}

func main() {
	fmt.Println("Paste this into a base64 image decoder")
	pic.Show(Pic)
}
