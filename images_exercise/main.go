package main

import (
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

type Image struct {
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	return m.Bounds()
}

func (i Image) At(x, y int) color.Color {
	return color.RGBA{uint8((x ^ y) * 2), uint8((x ^ y) * 2), 255, 255}
}

func main() {
	m := Image{}
	pic.ShowImage(m)
}
