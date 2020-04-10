// utils.go
package main

import (
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func RemoveElementFromBulletSlice(slice []Bullet, index int) []Bullet {
	// TODO : How to make it generic, i.e., it can remove an element from a slice of any kind
	slice[index] = slice[len(slice)-1] // No bounds check = panic(on index out of bounds)
	return slice[:len(slice)-1]
}

func RemoveElementFromEnemyTankSlice(slice []EnemyTank, index int) []EnemyTank {
	// TODO : How to make it generic, i.e., it can remove an element from a slice of any kind
	slice[index] = slice[len(slice)-1] // No bounds check = panic(on index out of bounds)
	return slice[:len(slice)-1]
}

func DegreeToRadian(angleInDegree float64) float64 {
	return angleInDegree * (math.Pi / 180.0)
}

func GetTexture(texturePath string, renderer *sdl.Renderer) (*sdl.Surface, *sdl.Texture, int) {
	image, err := img.Load(texturePath)
	if err != nil {
		HandleError("Failed to load image "+texturePath, err)
		return nil, nil, ERROR_FAILED_TO_LOAD_IMAGE
	}

	texture, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		HandleError("Failed to create texture: ", err)
		return image, nil, ERROR_FAILED_TO_CREATE_TEXTURE_FROM_IMAGE
	}
	return image, texture, 0
}

func GetRandomFloat32(min float32, max float32, r *rand.Rand) float32 {
	return min + (rand.Float32() * (max - min))
}

func IsInsideWindow(bounds sdl.FRect) bool {
	return ((bounds.X > 0.0) &&
		(bounds.Y > 0.0) &&
		((bounds.X + bounds.W) < float32(SCREEN_WIDTH)) &&
		((bounds.Y + bounds.H) < float32(SCREEN_HEIGHT)))
}
