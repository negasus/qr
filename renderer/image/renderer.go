package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

var (
	defaultModuleSize      = 5
	defaultForegroundColor = color.Black
	defaultBackgroundColor = color.White
)

type Options struct {
	ModuleSize      int
	ForegroundColor color.Color
	BackgroundColor color.Color
}

func Save(data []byte, filename string, opts *Options) error {
	img, err := Render(data, opts)
	if err != nil {
		return err
	}

	w := bytes.NewBuffer(nil)

	errEncode := png.Encode(w, img)
	if errEncode != nil {
		return errEncode
	}

	return os.WriteFile(filename, w.Bytes(), 0644)
}

func Render(data []byte, opts *Options) (image.Image, error) {
	o := &Options{
		ModuleSize:      defaultModuleSize,
		ForegroundColor: defaultForegroundColor,
		BackgroundColor: defaultBackgroundColor,
	}
	if opts != nil {
		if o.ModuleSize > 0 {
			o.ModuleSize = opts.ModuleSize
		}
		if opts.ForegroundColor != nil {
			o.ForegroundColor = opts.ForegroundColor
		}
		if opts.BackgroundColor != nil {
			o.BackgroundColor = opts.BackgroundColor
		}
	}

	size := int(math.Sqrt(float64(len(data))))
	if size*size != len(data) {
		return nil, fmt.Errorf("data length must be a square number")
	}

	img := image.NewRGBA(image.Rect(0, 0, (size+8)*o.ModuleSize, (size+8)*o.ModuleSize))

	draw.Draw(img, img.Bounds(), &image.Uniform{C: o.BackgroundColor}, image.Point{}, draw.Src)

	borderPadding := 4 * o.ModuleSize

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if data[y*size+x] == 1 {
				for dy := 0; dy < o.ModuleSize; dy++ {
					for dx := 0; dx < o.ModuleSize; dx++ {
						img.Set(x*o.ModuleSize+dx+borderPadding, y*o.ModuleSize+dy+borderPadding, o.ForegroundColor)
					}
				}
			}
		}
	}

	return img, nil
}
