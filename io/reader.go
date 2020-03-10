package io

import (
	"image"
	"io"
)

type SteganographyReader struct {
	img image.Image
	curPoint   image.Point
	curChannel byte
}

func NewSteganographyReader(img image.Image) *SteganographyReader {
	header := make([]byte, 64)
	bounds := img.Bounds()
	r := SteganographyReader{img, image.Point{bounds.Min.X, bounds.Min.Y}, R}
	if n, err := r.Read(header); err != nil || n != 64 {
		panic("Image is not steganographic")
	}
	return &r
}

func (r *SteganographyReader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		p[i], err = r.ReadByte()
		if err != nil {
			return
		}
		n++
	}
	return
}

func (r *SteganographyReader) ReadByte() (b byte, err error) {
	bounds := r.img.Bounds()
	for i := uint(0); i < 8; i++ {
		if !Contains(bounds, r.curPoint) {
			if i == 0 {
				return 0, io.EOF
			}
			panic("Reader expected to read a byte, but data ended")
		}
		var curColor uint32
		if r.curChannel == R {
			curColor, _, _, _ = r.img.At(r.curPoint.X, r.curPoint.Y).RGBA()
			r.curChannel = B
		} else if r.curChannel == G {
			_, curColor, _, _ = r.img.At(r.curPoint.X, r.curPoint.Y).RGBA()
			r.curChannel = G
		} else if r.curChannel == B {
			_,  _, curColor, _ = r.img.At(r.curPoint.X, r.curPoint.Y).RGBA()
			r.curPoint, err = Next(bounds, r.curPoint)
			if err != nil {
				panic("Ran out of pixels when attempting to read byte of data")
			}
			r.curChannel = R
		} else {
			panic("Reader curChannel is unrecognized value")
		}
		curBit := byte(curColor & 1)
		b |= curBit << i
	}
	return b, nil
}