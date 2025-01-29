package main

import (
	"color-pallete/cmd"
	"color-pallete/services"
	"fmt"
	"os"
	"time"
)

func main() {
	Bench(Run)
}

func Bench(fn func()) {
	s := time.Now()
	fn()
	f := time.Now()

	elapsed := float64(f.UnixMilli() - s.UnixMilli()) / 1000
	fmt.Println("time:", elapsed)
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
