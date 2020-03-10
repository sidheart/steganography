package main

import (
	stegio "github.com/sidheart/steganography/io"
	"image"
	"image/png"
	"io"
	"math"
	"os"
)

const WORD = 64
var HEADER = []byte{ // 64 bytes of DEADBEEF, yum!
	0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
	0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
	0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
	0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
}


// encode takes a dataFile and writes its data into
func encode(dataFile *os.File, img image.Image) image.Image {
	fileInfo, err := dataFile.Stat()
	if err != nil {
		panic(err.Error())
	}
	requiredPx := math.Ceil((float64(fileInfo.Size() + int64(len(HEADER))) * 8) / 3.0)
	availablePx := img.Bounds().Dx() * img.Bounds().Dy()
	if int(requiredPx) > availablePx {
		panic("the provided image is too small to hold the data")
	}
	buffer := make([]byte, WORD)
	stegImg := stegio.NewSteganographicImage(img)
	writer := stegio.NewSteganographyWriter(stegImg)
	n, err := writer.Write(HEADER)
	if err != nil {
		panic(err.Error())
	}
	if n != 64 {
		panic("couldn't write header information into image")
	}
	for n, err = dataFile.Read(buffer); n > 0 && err == nil; n, err = dataFile.Read(buffer) {
		_, werr := writer.Write(buffer[:n])
		if werr != nil {
			panic(werr.Error())
		}
	}
	return stegImg
}

func decode(img image.Image, outFile string) *os.File {
	newFile, err := os.Create(outFile)
	if err != nil {
		panic(err.Error())
	}
	reader := stegio.NewSteganographyReader(img)
	buffer := make([]byte, WORD)
	for {
		n, err := reader.Read(buffer)
		if n == 0 && err == io.EOF {
			break
		} else if err != io.EOF && err != nil {
			panic(err.Error())
		} else {
			_, werr := newFile.Write(buffer)
			if werr != nil {
				panic(werr.Error())
			}
		}
	}
	return newFile
}

// TODO I didn't write the whole file I only wrote the image portion
func main() {
	imgFile, err := os.Open("assets/fractal.png")
	if err != nil {
		panic(err.Error())
	}
	img, err := png.Decode(imgFile)
	if err != nil {
		panic(err.Error())
	}
	msgFile, err := os.Open("assets/dickbutt.png")
	newImage := encode(msgFile, img)
	newFile, err := os.Create("encoded_output.png")
	err = png.Encode(newFile, newImage)
	if err != nil {
		panic(err.Error())
	}
	stegImage := newImage.(stegio.SteganographicImage)
	decode(stegImage, "decoded_output.png")
}
