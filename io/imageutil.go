package io

import (
	"errors"
	"image"
)

var OOB = errors.New("next will be out of bounds")

func Contains(rect image.Rectangle, point image.Point) bool {
	if rect.Min.X <= point.X && rect.Max.X >= point.X && rect.Min.Y <= point.Y && rect.Max.Y >= point.Y {
		return true
	}
	return false
}

func Next(rect image.Rectangle, point image.Point) (p image.Point, err error) {
	if !Contains(rect, point) {
		return p, OOB
	}
	if point.X >= rect.Max.X {
		if point.Y >= rect.Max.Y {
			return p, OOB
		}
		p = image.Point{X: rect.Min.X, Y: point.Y + 1}
	} else {
		p = image.Point{X: point.X + 1, Y: point.Y}
	}
	return
}
