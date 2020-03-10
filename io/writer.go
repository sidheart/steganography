package io

import (
	"errors"
	"image"
	"image/color"
)

const (
	R = 0
	G
	B
)

var MASK = byte(1)

type SteganographyWriter struct {
	img        SteganographicImage
	curPoint   image.Point
	curChannel byte
}

func NewSteganographyWriter(img SteganographicImage) SteganographyWriter {
	if img.ColorModel() != color.RGBAModel {
		panic("SteganographyWriter only works with RGBA images")
	}
	rect := img.Bounds()
	return SteganographyWriter{img, image.Point{X: rect.Min.X, Y: rect.Min.Y}, R}
}

func (w SteganographyWriter) RemainingBits() int {
	rect := w.img.Bounds()
	if !Contains(rect, w.curPoint)  {
		return 0
	}
	bits := int(3 - w.curChannel) // if curChannel is R, there are 3 usable bits left in curPoint, 2 for G, and 1 for B
	for next, err := Next(rect, w.curPoint); err == nil; next, err = Next(rect, next) {
		bits += 3
	}
	return bits
}

func writeLSB(b byte, lsb byte) byte {
	if lsb == 0 {
		return b & 0xFE
	} else if lsb == 1 {
		return b | 0x01
	} else {
		panic("writeLSB called with invalid lsb value")
	}
}

func (w SteganographyWriter) Write(p []byte) (n int, err error) {
	k := len(p) * 8
	if k == 0 {
		return 0, nil
	}
	if k > w.RemainingBits() {
		return 0, errors.New("not enough bits left to write data")
	}
	curColor := w.img.At(w.curPoint.X, w.curPoint.Y).(color.RGBA)
	for _, b := range p {
		for i := uint(0); i < 8; i++ { // TODO test that this overflow works
			bit := (b & (MASK << i)) >> i
			if w.curChannel == R {
				curColor.R = writeLSB(curColor.R, bit)
				w.curChannel = G
			} else if w.curChannel == G {
				curColor.G = writeLSB(curColor.G, bit)
				w.curChannel = B
			} else {
				curColor.B = writeLSB(curColor.B, bit)
				w.img.WritePixel(w.curPoint, curColor)
				w.curChannel = R
				w.curPoint, err = Next(w.img.Bounds(), w.curPoint)
				if err != nil {
					panic("Ran out of pixels while attempting to write byte")
				}
				curColor = w.img.At(w.curPoint.X, w.curPoint.Y).(color.RGBA)
			}
		}
		n++
	}
	w.img.WritePixel(w.curPoint, curColor) // this write ensures that partial data is also written
	return
}
