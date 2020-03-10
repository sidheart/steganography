package io

import (
	"github.com/stretchr/testify/assert"
	"image"
	"testing"
)

func TestContains(t *testing.T) {
	bounds := image.Rectangle{image.Point{0, 0}, image.Point{10, 10}}
	testPoint := image.Point{0, -1}
	assert.False(t, Contains(bounds, testPoint))
	testPoint = image.Point{-1, 0}
	assert.False(t, Contains(bounds, testPoint))
	testPoint = image.Point{10, 11}
	assert.False(t, Contains(bounds, testPoint))
	testPoint = image.Point{11, 10}
	assert.False(t, Contains(bounds, testPoint))
	testPoint = image.Point{0, 0}
	assert.True(t, Contains(bounds, testPoint))
	testPoint = image.Point{10, 10}
	assert.True(t, Contains(bounds, testPoint))
	testPoint = image.Point{1, 9}
	assert.True(t, Contains(bounds, testPoint))
}

func TestNext(t *testing.T) {
	bounds := image.Rectangle{image.Point{0, 0}, image.Point{10, 10}}
	start := image.Point{0, 0}
	var err error
	var matches int
	for x, y := 0, 0; err == nil; start, err = Next(bounds, start) {
		assert.Equal(t, image.Point{x, y}, start)
		if x == 10 {
			x = 0
			y++
		} else {
			x++
		}
		matches++
	}
	assert.Equal(t, 121, matches)
	assert.Equal(t, OOB, err)
}
