package main

import (
	"color-pallete/cmd"
	"color-pallete/services"
	"fmt"
	"os"
)

func main() {
	config, errs := cmd.ParseArgs()

	for _, err := range errs {
		fmt.Printf("args parsing error: %s\n", err.Error())
	}
	if len(errs) > 0 {
		os.Exit(1)
	}

	config.SetDefaults()
	errs = config.Validate()

	for _, err := range errs {
		fmt.Printf("configuration error: %s\n", err.Error())
	}
	if len(errs) > 0 {
		os.Exit(1)
	}
}

func RunExample() {
	file := "path"
	img, _, err := services.ReadImage(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	colors := services.GetColors(img)
	tiles := services.MakeTiles(len(colors[0]), len(colors), 16, 16)
	copy := services.DrawPallete(img, tiles)
	services.SaveImage(copy, "copy-" + file)

}