package main

import (
	"os"
	"log"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/gographics/imagick.v2/imagick"
	"strconv"
	"fmt"
)

func showImage(c *gin.Context, filename string) {

	imagick.Initialize()
	defer imagick.Terminate()

	img, err := openImage(filename)
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	defer img.Close()

	width := c.DefaultQuery("width", "0")
	height := c.DefaultQuery("height", "0")
	int_width, _ := strconv.Atoi(width)
	int_height, _ := strconv.Atoi(height)
	img.resize(int_width, int_height)


	c.Writer.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	c.Writer.Write(img.render())
}

type Image struct {
	path  string
	image *imagick.MagickWand
	file  *os.File
}

func openImage(path string) (Image, error) {
	img := Image{path:path, image:imagick.NewMagickWand()}
	return img, img.load()
}

func (img *Image) Close() {
	img.image.Destroy()
	if img.file != nil {
		filename := img.file.Name()
		img.file.Close()
		fmt.Println(filename)
		os.Remove(filename)
	}
}

func (img *Image) render() []byte {
	return img.image.GetImageBlob()
}

func (img *Image) load() error {
	img.image = imagick.NewMagickWand()
	return img.image.ReadImage(img.path)
}

func (img *Image) resize(width int, height int) {
	o_height := int(img.image.GetImageHeight())
	o_width := int(img.image.GetImageWidth())

	if height == 0 && width == 0 {
		return
	}

	if height == 0 {
		height = o_height*width/o_width
	}

	if width == 0 {
		width = o_width*height/o_height
	}

	img.image.ResizeImage(uint(width), uint(height), imagick.FILTER_BOX, 1)
}


