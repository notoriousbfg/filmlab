package main

import (
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

type Profile interface {
	Adjust(image *image.NRGBA, preset string) *image.NRGBA
}

type Portra160 struct {
}

func (p Portra160) Adjust(image *image.NRGBA, preset string) *image.NRGBA {
	switch preset {
	case "light":
		image = adjustColours(image, 67, 10, -25)
		image = imaging.AdjustSigmoid(image, 0.9, 8.0)
		image = imaging.AdjustContrast(image, 25)
		image = imaging.AdjustBrightness(image, 20)
		image = imaging.AdjustGamma(image, 0.6)
	case "mid", "dark":
		image = adjustColours(image, 70, 12, -30)
		image = imaging.AdjustSigmoid(image, 0.9, 8.0)
		image = imaging.AdjustContrast(image, 25)
		image = imaging.AdjustBrightness(image, 20)
		image = imaging.AdjustGamma(image, 0.7)
	}
	return image
}

type HP5Plus struct {
}

func (h HP5Plus) Adjust(image *image.NRGBA, preset string) *image.NRGBA {
	switch preset {
	case "light":
		image = adjustColours(image, 25, 5, 5)
		image = imaging.AdjustBrightness(image, 10)
		image = imaging.AdjustGamma(image, 0.55)
		image = imaging.AdjustSigmoid(image, 0.9, 8.0)
	case "mid":
		image = adjustColours(image, 30, 5, 5)
		image = imaging.AdjustBrightness(image, 35)
		image = imaging.AdjustGamma(image, 0.6)
		image = imaging.AdjustSigmoid(image, 0.9, 8.0)
	case "dark":
		image = adjustColours(image, 33, 0, 0)
		image = imaging.AdjustBrightness(image, 35)
		image = imaging.AdjustGamma(image, 0.7)
		image = imaging.AdjustSigmoid(image, 0.7, 8.0)
		image = imaging.AdjustContrast(image, 20)
	}
	return image
}

type ColorPlus struct {
}

func (c ColorPlus) Adjust(image *image.NRGBA, preset string) *image.NRGBA {
	image = adjustColours(image, 62, 5, -28)
	image = imaging.AdjustSigmoid(image, 0.9, 9.0)
	image = imaging.AdjustContrast(image, 25)
	image = imaging.AdjustBrightness(image, 25)
	image = imaging.AdjustGamma(image, 0.6)
	return image
}

func adjustColourValue(colour int, amount int) int {
	colour = colour + amount

	if colour > 255 {
		return 255
	} else if colour < 0 {
		return 0
	}

	return colour
}

func adjustColours(image *image.NRGBA, r int, g int, b int) *image.NRGBA {
	image = imaging.AdjustFunc(
		image,
		func(c color.NRGBA) color.NRGBA {
			r := adjustColourValue(int(c.R), r)
			g := adjustColourValue(int(c.G), g)
			b := adjustColourValue(int(c.B), b)

			return color.NRGBA{uint8(r), uint8(g), uint8(b), c.A}
		},
	)
	return image
}
