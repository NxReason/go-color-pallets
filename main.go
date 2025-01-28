package main

import (
	"color-pallete/cmd"
	"color-pallete/services"
	"fmt"
	"os"
)

func main() {
	Run()
}

func Run() {
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
	
	errs = services.ProcessFiles(config)
	for _, err := range errs {
		fmt.Printf("image processing error: %s\n", err.Error())
	}
}
