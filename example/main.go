package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	_ "embed"

	"github.com/lewisgibson/go-binpack"
	"golang.org/x/image/draw"
)

//go:embed images/*
var imagesFS embed.FS

func main() {
	// Read the images from the embedded filesystem.
	dirents, err := imagesFS.ReadDir("images")
	if err != nil {
		panic(err)
	}

	// Loop through the images and decode them.
	var images = make([]image.Image, 0, len(dirents))
	for _, dirent := range dirents {
		// Read the image from the filesystem.
		b, err := imagesFS.ReadFile("images/" + dirent.Name())
		if err != nil {
			panic(err)
		}

		// Decode the image.
		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			panic(err)
		}

		// Append the image to the images slice.
		images = append(images, img)
	}

	// Create a new Collager.
	c := &Collager{
		Images: images,
		// Locations is a pre-allocated slice of image.Point structs with the same length as images.
		Locations: make([]image.Point, len(images)),
	}

	// Pack the images into a collage.
	width, height := binpack.Pack(c)

	// Render the collage to an image.
	var canvas = image.NewRGBA(image.Rect(0, 0, width, height))
	for n, img := range c.Images {
		location := img.Bounds().Add(c.Locations[n])
		draw.Draw(canvas, location, img, image.Point{}, draw.Over)
	}

	// Create a file to write the image to.
	f, err := os.Create("collage.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write the image to the file.
	if err := png.Encode(f, canvas); err != nil {
		panic(err)
	}
}

// Collager is a struct that implements the binpack.Packer interface.
type Collager struct {
	Images    []image.Image
	Locations []image.Point
}

// Len returns the number of images in the Collager.
func (c *Collager) Len() int {
	return len(c.Images)
}

// Width returns the width of the image at index n.
func (c *Collager) Rectangle(n int) binpack.Rectangle {
	return binpack.Rectangle{
		Width:  c.Images[n].Bounds().Dx(),
		Height: c.Images[n].Bounds().Dy(),
	}
}

// Place sets the location of the image at index n.
func (c *Collager) Place(n, x, y int) {
	c.Locations[n] = image.Point{x, y}
}
