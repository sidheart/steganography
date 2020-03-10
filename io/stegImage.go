package io

import (
	"image"
	"image/color"
)

type SteganographicImage struct {
	i     image.Image
	edits map[image.Point]color.Color
}

type ImageIter func() (image.Point, error)

func NewSteganographicImage(i image.Image) SteganographicImage {
	return SteganographicImage{i, make(map[image.Point]color.Color)}
}

func (img SteganographicImage) ColorModel() color.Model {
	return img.i.ColorModel()
}

func (img SteganographicImage) Bounds() image.Rectangle {
	return img.i.Bounds()
}

func (img SteganographicImage) At(x, y int) color.Color {
	if pixelColor, present := img.edits[image.Point{x, y}]; present {
		return pixelColor
	}
	return img.i.At(x, y)
}

func (img SteganographicImage) WritePixel(point image.Point, color color.Color) {
	img.edits[point] = color
}
