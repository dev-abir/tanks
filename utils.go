// utils.go
package main

import (
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func RemoveElementFromBulletSlice(slice []Bullet, index int) []Bullet {
	// source : https://stackoverflow.com/a/37335777
	// TODO : How to make it generic, i.e., it can remove an element from a slice of any kind
	slice[index] = slice[len(slice)-1] // No bounds check = panic(on index out of bounds)
	return slice[:len(slice)-1]
}

func RemoveElementFromEnemyTankSlice(slice []EnemyTank, index int) []EnemyTank {
	// source : https://stackoverflow.com/a/37335777
	// TODO : How to make it generic, i.e., it can remove an element from a slice of any kind
	slice[index] = slice[len(slice)-1] // No bounds check = panic(on index out of bounds)
	return slice[:len(slice)-1]
}

func RemoveElementFromExplosionSlice(slice []Explosion, index int) []Explosion {
	// source : https://stackoverflow.com/a/37335777
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

func GetSoundEffect(path string) *mix.Chunk {
	result, err := mix.LoadWAV(path)
	if err != nil {
		HandleError("Cannot load "+EXPLOSION_SOUND_PATH+", you may play without it: ", err)
	}
	return result
}

func PlaySoundEffect(soundEffect *mix.Chunk) {
	if _, err := soundEffect.Play(-1, 0); err != nil {
		HandleError("Error on playing sound effect: ", err)
	}
}

func GetRandomFloat32(min float32, max float32, r *rand.Rand) float32 {
	return min + (rand.Float32() * (max - min))
}

func DrawTexture(renderer *sdl.Renderer, texture *sdl.Texture, boundingBox *sdl.FRect, rotationAngle float32) {
	/*

		TODO : For some reason CopyExF is not working.......

	*/

	renderer.CopyEx(texture, nil, &sdl.Rect{
		int32(boundingBox.X),
		int32(boundingBox.Y),
		int32(boundingBox.W),
		int32(boundingBox.H)}, float64(rotationAngle), nil, sdl.FLIP_NONE)
}

func SetPositionOfEnemyTanks(enemyTanks []EnemyTank, playerTankBoundingBox sdl.FRect, r *rand.Rand) {
	for index, _ := range enemyTanks {
		enemyTanks[index].boundingBox = GetPositionOfOneEnemyTank(enemyTanks[index].boundingBox, enemyTanks[:index], playerTankBoundingBox, r)
	}
}

func GetPositionOfOneEnemyTank(enemyTankBoundingBox sdl.FRect, otherEnemyTanks []EnemyTank, playerTankBoundingBox sdl.FRect, r *rand.Rand) sdl.FRect {
	experimentalTankBoundingBox := sdl.FRect{
		X: r.Float32() * float32(SCREEN_WIDTH),
		Y: r.Float32() * float32(SCREEN_HEIGHT),
		W: enemyTankBoundingBox.W,
		H: enemyTankBoundingBox.H,
	}
	if !ValidPosition(experimentalTankBoundingBox, otherEnemyTanks, playerTankBoundingBox) {
		return GetPositionOfOneEnemyTank(enemyTankBoundingBox, otherEnemyTanks, playerTankBoundingBox, r)
	}
	return experimentalTankBoundingBox
}

func ValidPosition(experimentalTankBoundingBox sdl.FRect, otherEnemyTanks []EnemyTank, playerTankBoundingBox sdl.FRect) bool {
	for idx, _ := range otherEnemyTanks {
		if experimentalTankBoundingBox.HasIntersection(&otherEnemyTanks[idx].boundingBox) {
			return false
		}
	}
	if experimentalTankBoundingBox.HasIntersection(&playerTankBoundingBox) || !IsInsideWindow(experimentalTankBoundingBox) {
		return false
	}
	return true
}

func IsInsideWindow(bounds sdl.FRect) bool {
	return ((bounds.X > 0.0) &&
		(bounds.Y > 0.0) &&
		((bounds.X + bounds.W) < float32(SCREEN_WIDTH)) &&
		((bounds.Y + bounds.H) < float32(SCREEN_HEIGHT)))
}
