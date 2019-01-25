package imghash

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"os"
)

const INTERP = resize.Bicubic

func Ahash(img image.Image) uint64 {
	resizedImage:= resize.Resize(8, 8, img, INTERP)

	grayImage := imageToGray(&resizedImage)


	// saveImage(resizedImage, "resizedImage")
	// saveImage(grayImage, "grayImage")

	average := getAverage(grayImage)

	hash := uint64(0)

	// hashImage := image.NewGray(image.Rect(0,0,8,8))

	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			value := uint64(grayImage.GrayAt(i, j).Y)
			if value > average {
				hash += 1
				//hashImage.SetGray(i,j,color.Gray{Y:255})
			} //else {
			//	hashImage.SetGray(i, j, color.Gray{Y: 0})
			//}
			hash = hash << 1
		}
	}
	// saveImage(hashImage, "hashImage")
	return hash
}

func getAverage(img *image.Gray) uint64 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	sum := uint64(0)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			sum += uint64(img.GrayAt(i, j).Y)
		}
	}
	return sum / uint64(width * height)
}

func saveImage(img image.Image, filename string) {
	out, err := os.Create("./" + filename + ".png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
