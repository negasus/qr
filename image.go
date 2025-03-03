package qr

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

func renderImage(data []byte) (image.Image, error) {
	moduleSize := 5
	backgroundColor := color.White
	foregroundColor := color.Black

	size := int(math.Sqrt(float64(len(data))))
	if size*size != len(data) {
		return nil, fmt.Errorf("data length must be a square number")
	}

	img := image.NewRGBA(image.Rect(0, 0, (size+8)*moduleSize, (size+8)*moduleSize))

	draw.Draw(img, img.Bounds(), &image.Uniform{C: backgroundColor}, image.Point{}, draw.Src)

	borderPadding := 4 * moduleSize

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if data[y*size+x] == 1 {
				for dy := 0; dy < moduleSize; dy++ {
					for dx := 0; dx < moduleSize; dx++ {
						img.Set(x*moduleSize+dx+borderPadding, y*moduleSize+dy+borderPadding, foregroundColor)
					}
				}
			}
		}
	}

	return img, nil
}
