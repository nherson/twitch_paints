package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

// Run loops over a channel, receiving PaintCommands and issuing them
// onto the image
func (ip *ImageProcessor) Run() {
	for {
		pc := <-ip.Chan
		// Set R, G, and B
		ip.Image.Pix[(uint32(pc.Y)*800*4)+(uint32(pc.X)*4)] = pc.Red
		ip.Image.Pix[(uint32(pc.Y)*800*4)+(uint32(pc.X)*4)+1] = pc.Green
		ip.Image.Pix[(uint32(pc.Y)*800*4)+(uint32(pc.X)*4)+2] = pc.Blue
		writeImage(ip.OutFile, &ip.Image)
	}
}

func freshImage() image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	// Iterate over the image to make it all white and turn up the alpha
	for i := 0; i < 800*600*4; i++ {
		img.Pix[i] = 255
	}
	return *img
}

func writeImage(f *os.File, i *image.RGBA) {
	f.Seek(0, 0)
	png.Encode(f, i)
}

// Sample method showing how to manipulate an image, it doesn't get called
// Images are stored in memory with the following structure
/*
{ [r,g,b,a](0,0), [r,g,b,a](1,0) ...
  [r,g,b,a](0,1), [r,g,b,a](1,1) ...
  ...
}
*/
// ... all in one contiguous uint8 array
func paintDemo() {
	// Create an image in memory
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	// Draw a read horizontal line through the middle from the left?
	for i := 0; i < 800*600; i++ {
		img.Pix[i*4+3] = 255
	}
	for i := 0; i <= 400; i++ {
		img.Pix[3200*300+i*4] = 128
	}
	// Write it out
	f, err := os.Create("out.png")
	if err != nil {
		fmt.Println("Problem creating demo PNG")
	}
	defer f.Close()

	png.Encode(f, img)

}
