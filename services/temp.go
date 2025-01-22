package services

import (
	"image/jpeg"
	"log"
	"os"
)

func GetImageSize(path string) (width, height int) {
	file, err := os.Open("jp.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	width, height = bounds.Max.X, bounds.Max.Y
	return
}