# go-binpack

[![Build Workflow](https://github.com/lewisgibson/go-binpack/actions/workflows/build.yaml/badge.svg)](https://github.com/lewisgibson/go-binpack/actions/workflows/build.yaml)
[![Pkg Go Dev](https://pkg.go.dev/badge/github.com/lewisgibson/go-binpack)](https://pkg.go.dev/github.com/lewisgibson/go-binpack)

A Go package for packing rectangles into a bin using the MaxRects algorithm.

## Resources

-   [Discussions](https://github.com/lewisgibson/go-binpack/discussions)
-   [Reference](https://pkg.go.dev/github.com/lewisgibson/go-binpack)
-   [Examples](https://pkg.go.dev/github.com/lewisgibson/go-binpack#pkg-examples)

## Installation

```sh
go get github.com/lewisgibson/go-binpack
```

## Quickstart

```go
package main

import (
	"github.com/lewisgibson/go-binpack"
)

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

// Create a new Collager.
c := &Collager{
    Images: images,
    // Locations is a pre-allocated slice of image.Point structs with the same length as images.
    Locations: make([]image.Point, len(images)),
}

// Pack the images into a collage.
width, height := binpack.Pack(c)
```
