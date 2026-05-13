package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	invertImages := flag.Bool("invert-images", false, "also invert raster images embedded in the PDF")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pdf-inverter [--invert-images] <input.pdf> <output.pdf>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	cfg := Config{InvertImages: *invertImages}
	if err := Process(args[0], args[1], cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
