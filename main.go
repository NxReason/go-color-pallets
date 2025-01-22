package main

import (
	"color-pallete/cmd"
	"color-pallete/services"
	"fmt"
	"os"
)

func main() {

	errs := cmd.ParseArgs()
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	fmt.Println(cmd.Ifa.GetFiles())

	for _, file := range cmd.Ifa.GetFiles() {
		img, _, err := services.ReadImage(file)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		colors := services.GetColors(img)
		tiles := services.MakeTiles(len(colors[0]), len(colors), 16, 16)
		fmt.Println(tiles)
		copy := services.DrawPallete(img, tiles)
		services.SaveImage(copy, "copy-" + file)
	}
}