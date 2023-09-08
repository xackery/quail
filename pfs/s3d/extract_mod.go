package s3d

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/quail/helper"
)

func (e *S3D) ExtractMod(path string) (string, error) {

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("creating directory %s/\n", path)
			err = os.MkdirAll(path, 0766)
			if err != nil {
				return "", fmt.Errorf("mkdirall: %w", err)
			}
		}
		fi, err = os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("stat after mkdirall: %w", err)
		}
	}
	if !fi.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}

	extractStdout := ""
	for i, file := range e.files {
		//fmt.Println(fmt.Sprintf("%s/%s", path, file.Name()), len(file.Data()))
		data := file.Data()

		fmt.Println(file.Name())
		if strings.HasSuffix(strings.ToLower(file.Name()), ".bmp") {
			data = modifyBMP(data)
		}

		err = os.WriteFile(fmt.Sprintf("%s/%s", path, file.Name()), data, 0644)
		if err != nil {
			return "", fmt.Errorf("index %d: %w", i, err)
		}
		extractStdout += file.Name() + ", "
	}
	if len(e.files) == 0 {
		return "", fmt.Errorf("no files found to extract")
	}
	extractStdout = extractStdout[0 : len(extractStdout)-2]
	return fmt.Sprintf("extracted %d file%s to %s: %s", len(e.files), helper.Pluralize(len(e.files)), path, extractStdout), nil
}

func modifyBMP(data []byte) []byte {

	if data[0] != 'B' {
		return data
	}
	if data[1] != 'M' {
		return data
	}
	colorTableOffset := binary.LittleEndian.Uint32(data[0x0A:0x0E])

	if colorTableOffset == 0 {
		return data
	}

	firstIndexRGB := color.RGBA{R: data[0x36+2], G: data[0x36+1], B: data[0x36+0], A: 0xFF}

	img, err := bmp.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println("decode failed:", err)
		return data
	}
	fmt.Println("first index rgb", firstIndexRGB)

	// create a new image, same size as the original
	// iterate through each pixel, if it matches the first pixel, replace it with transparent
	bounds := img.Bounds()
	newImg := image.NewAlpha(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixel := img.At(x, y) //bounds.Max.Y-y-1)
			r, g, b, a := pixel.RGBA()
			if r == uint32(firstIndexRGB.R) && g == uint32(firstIndexRGB.G) && b == uint32(firstIndexRGB.B) {
				r = 0
				g = 0
				b = 0
				a = 0
				fmt.Println("replaced pixel at xy", x, y)
			}
			newImg.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	// encode the new image as png
	buf := new(bytes.Buffer)
	err = png.Encode(buf, newImg)
	if err != nil {
		fmt.Println("encode failed:", err)
		return data
	}

	return buf.Bytes()
}
