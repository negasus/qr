package simple

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

func WithDotSize(size int) Option {
	return func(r *Renderer) {
		r.dotSize = size
	}
}

func WithForegroundColor(c color.Color) Option {
	return func(r *Renderer) {
		r.foregroundColor = c
	}
}

func WithBackgroundColor(c color.Color) Option {
	return func(r *Renderer) {
		r.backgroundColor = c
	}
}

const (
	defaultDotSize = 5
)

var (
	defaultForegroundColor = color.Black
	defaultBackgroundColor = color.White
)

type Option func(*Renderer)

type Renderer struct {
	dotSize         int
	foregroundColor color.Color
	backgroundColor color.Color
}

func New(opts ...Option) *Renderer {
	r := &Renderer{
		dotSize:         defaultDotSize,
		foregroundColor: defaultForegroundColor,
		backgroundColor: defaultBackgroundColor,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Renderer) Render(data []byte) (image.Image, error) {
	if r.dotSize <= 0 {
		return nil, fmt.Errorf("dotSize must be greater than 0")
	}

	size := int(math.Sqrt(float64(len(data))))
	if size*size != len(data) {
		return nil, fmt.Errorf("data length must be a square number")
	}

	img := image.NewRGBA(image.Rect(0, 0, (size+8)*r.dotSize, (size+8)*r.dotSize))

	draw.Draw(img, img.Bounds(), &image.Uniform{C: r.backgroundColor}, image.Point{}, draw.Src)

	borderPadding := 4 * r.dotSize

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if data[y*size+x] == 1 {
				for dy := 0; dy < r.dotSize; dy++ {
					for dx := 0; dx < r.dotSize; dx++ {
						img.Set(x*r.dotSize+dx+borderPadding, y*r.dotSize+dy+borderPadding, r.foregroundColor)
					}
				}
			}
		}
	}

	return img, nil
}
