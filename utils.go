// utils.go
package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
)

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func HandleFatalError(err error) {
	if err != nil {
		//fmt.Fprintf(os.Stderr, err.Error())
		panic(err)
	}
}

func GetBoundingBox(position pixel.Vec, sprite *pixel.Sprite) pixel.Rect {
	return pixel.Rect{
		position,
		sprite.Frame().Max.Add(position),
	}
}

func RemoveElementFromBulletSlice(slice []Bullet, index int) []Bullet {
	// TODO : How to make it generic, i.e., it can remove an element from a slice of any kind
	slice[index] = slice[len(slice)-1] // No bounds check = panic(on index out of bounds)
	return slice[:len(slice)-1]
}
