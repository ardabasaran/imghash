package imghash

import (
	"image"
	"image/color"
)

func imageToGray(img *image.Image) *image.Gray {
	bounds := (*img).Bounds()
	gray := image.NewGray(bounds)

	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			rgbaPixel := (*img).At(i, j)
			red, green, blue, _ := rgbaPixel.RGBA()

			// See: https://en.wikipedia.org/wiki/Grayscale#Luma_coding_in_video_systems
			grayValue := 0.299*float64(red) + 0.587*float64(green) + 0.114*float64(blue)
			grayPixel := color.Gray{Y: uint8(grayValue)}

			gray.Set(i, j, grayPixel)
		}
	}
	return gray
}
