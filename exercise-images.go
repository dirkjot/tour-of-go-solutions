package main

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

func notused() {
	// BEWARE: using any print will prevent the image from inlining.
	fmt.Println("")
}

type Image struct {
	width        int
	height       int
	fillFunction func(uint8, uint8) uint8
	canvas       [][]uint8
}

// ColorModel returns the Image's color model.
func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (i Image) Bounds() image.Rectangle {
	// return image.Rect(0, 0, 100,100)
	return image.Rect(0, 0, i.width, i.height)
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (i Image) At(x, y int) color.Color {
	// return image.Black
	v := i.canvas[x][y]
	return color.RGBA{v, v, 255, 255}
}

func NewImage(dx, dy int, fillFunction func(uint8, uint8) uint8) (result Image) {

	i := Image{width: dy, height: dx}
	i.fillFunction = fillFunction
	i.canvas = make([][]uint8, dy, dy)

	for y := 0; y < dy; y++ {
		i.canvas[y] = make([]uint8, dx, dx)
	}

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			i.canvas[x][y] = i.fillFunction(uint8(x), uint8(y))
		}
	}
	// fmt.Printf("NewImage returns type %T\n", i)
	return i
}

func main() {
	// m := Image{}
	// m := Image{width: 100, height: 100}
	var m Image = NewImage(100, 100, func(x, y uint8) uint8 { return x ^ y })

	// BEWARE: using any print will prevent the image from inlining.
	//fmt.Printf("In main, we see type %T\n", m)
	pic.ShowImage(m)
}
