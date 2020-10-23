package main

import (
	"fmt"
	"flag"
	"image"
	"io"
	"os"
	"github.com/nfnt/resize"
)
import _ "image/png"
import _ "image/jpeg"
import _ "image/gif"

var (
	inputFilename  = flag.String("input", "", "Input the input-file name")
	outputFilename = flag.String("output", "", "Input the output-file name")
	width  = flag.Uint("width", 0, "Input the width")
	height = flag.Uint("height", 0, "Input the height")
)

func imageToRGB565(in io.Reader, out io.Writer, width uint, height uint) error {
	var img, _, decodeErr = image.Decode(in)

	if width != 0 {
		img = resize.Resize(width, height, img, resize.Lanczos3)
	}

	if decodeErr != nil {
		fmt.Println(decodeErr)
		fmt.Printf("Could not decode IMAGE\n\n")
		os.Exit(1)
	}

	var size = (img.Bounds().Max.X) * (img.Bounds().Max.Y)
	rgb565s := make([]byte, size*2)
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for y := img.Bounds().Min.Y; y < h; y++ {
		for x := img.Bounds().Min.X; x < w; x++ {

			c := img.At(x, y)
			r, g, b, _ := c.RGBA()

			r5 := ((r >> 3) & 0x1f) << 11
			g6 := ((g >> 2) & 0x3f) << 5
			b5 := (b >> 3) & 0x1f
			rgb565 := r5 | g6 | b5

			rgb565s[2*(y*w+x)] = byte(rgb565 >> 8)
			rgb565s[2*(y*w+x)+1] = byte(rgb565)
		}
	}
	out.Write(rgb565s)
	return nil
}

func main() {
	flag.Parse()
	in, err := os.Open(*inputFilename)
	if err != nil {
		panic("Please setup -input")
	}
	out, err := os.Create(*outputFilename)
	if err != nil {
		panic("Please setup -ouput")
	}
	if err := imageToRGB565(in, out, *width, *height); err != nil {
		panic(err)
	}
}
