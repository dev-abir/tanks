// tank.go
package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Bullet struct {
	bulletSprite  *pixel.Sprite
	velocity      float64
	position      pixel.Vec
	rotationAngle float64
}

func (bullet *Bullet) Update(delta float64) {
	bullet.position = bullet.position.Add(pixel.V(bullet.velocity*delta*math.Cos(bullet.rotationAngle), bullet.velocity*delta*math.Sin(bullet.rotationAngle)))
}

func (bullet Bullet) Draw(window *pixelgl.Window) {
	matrix := pixel.IM
	matrix = matrix.Moved(bullet.position)
	// matrix = matrix.Scaled(bullet.position, 1.0) // no need to scale, when scale is 1
	matrix = matrix.Rotated(bullet.position, bullet.rotationAngle)
	bullet.bulletSprite.Draw(window, matrix)
}

func (bullet Bullet) OutOfWindow(window *pixelgl.Window) bool {
	return !window.Bounds().Contains(bullet.position)
}

/*
Almost same as the Bullet struct, it is also inheriting this interface...
Even if I define it, I will not use it most probably

type Tank interface {
	// Init()
	Update()
	Draw(window *pixelgl.Window)
}*/

type EnemyTank struct {
	tankSprite    *pixel.Sprite
	rotationAngle float64
	position      pixel.Vec
	alive         bool
	noUpdateTime  float64
	timer         time.Time
}

func NewEnemyTank(tankSprite *pixel.Sprite, initialRotationAngle float64, alive bool, noUpdateTime float64) EnemyTank {
	return EnemyTank{tankSprite: tankSprite, rotationAngle: initialRotationAngle, position: pixel.ZV, alive: alive, noUpdateTime: noUpdateTime, timer: time.Now()}
}

func SetPositions(enemyTanks []EnemyTank, playerTankBoundingBox pixel.Rect, r *rand.Rand) {
	for index, _ := range enemyTanks {
		if enemyTanks[index].alive {
			enemyTanks[index].position = GetPositionOfOneEnemyTank(enemyTanks[index], enemyTanks[:index], playerTankBoundingBox, r)
		}
	}
}

func GetPositionOfOneEnemyTank(currentEnemyTank EnemyTank, otherEnemyTanks []EnemyTank, playerTankBoundingBox pixel.Rect, r *rand.Rand) pixel.Vec {
	experimentalPosition := pixel.V(r.Float64()*SCREEN_WIDTH, r.Float64()*SCREEN_HEIGHT)
	currentEnemyTank.position = experimentalPosition
	for idx, _ := range otherEnemyTanks {
		otherEnemyTankBoundingBox := GetBoundingBox(otherEnemyTanks[idx].position, otherEnemyTanks[idx].tankSprite)
		currentEnemyTankBoundingBox := GetBoundingBox(currentEnemyTank.position, currentEnemyTank.tankSprite)
		if (otherEnemyTankBoundingBox.Intersects(currentEnemyTankBoundingBox) || currentEnemyTankBoundingBox.Intersects(playerTankBoundingBox)) && otherEnemyTanks[idx].alive {
			return GetPositionOfOneEnemyTank(currentEnemyTank, otherEnemyTanks, playerTankBoundingBox, r)
		}
	}
	return experimentalPosition
}

/*func (tank EnemyTank) Update(delta float64, r *rand.Rand, playerTankPosition pixel.Vec) (EnemyTank, Bullet) {
	var bullet Bullet

	return tank, bullet
}*/

func (tank EnemyTank) MoveInRandomDir(delta float64, r *rand.Rand) EnemyTank {
	switch r.Intn(4) {
	case 0:
		tank.position = tank.position.Add(pixel.V(0, (LEVEL_0_ENEMY_TANK_VELOCITY * delta))) // UP
	case 1:
		tank.position = tank.position.Add(pixel.V(0, -(LEVEL_0_ENEMY_TANK_VELOCITY * delta))) // DOWN
	case 2:
		tank.position = tank.position.Add(pixel.V(-(LEVEL_0_ENEMY_TANK_VELOCITY * delta), 0)) // LEFT
	case 3:
		tank.position = tank.position.Add(pixel.V((LEVEL_0_ENEMY_TANK_VELOCITY * delta), 0)) // RIGHT
	}
	return tank
}

func (tank EnemyTank) SpinAndShoot(delta float64, r *rand.Rand, playerTankPosition pixel.Vec) (EnemyTank, Bullet) {
	switch r.Intn(2) {
	case 0:
		displacementVector := playerTankPosition.Sub(tank.position) // SHOOT THE PLAYER
		tank.rotationAngle = displacementVector.Angle()
	case 1:
		tank.rotationAngle = r.Float64() * (2.0 * math.Pi) // SHOOT ANYWHERE RANDOMLY
	}
	picture, err := LoadPicture(BULLET_TEXTURE_PATH)
	HandleFatalError(err)
	bulletSprite := pixel.NewSprite(picture, picture.Bounds())
	return tank, Bullet{
		bulletSprite:  bulletSprite,
		velocity:      BULLET_VELOCITY,
		position:      tank.position,
		rotationAngle: tank.rotationAngle,
	}
}

func (tank EnemyTank) WillUpdate() bool {
	return time.Since(tank.timer).Seconds() >= tank.noUpdateTime
}

func (tank EnemyTank) Die(delta float64) EnemyTank {
	tank.alive = false
	tank.tankDieAnimation(delta)
	return tank
}

func (tank EnemyTank) tankDieAnimation(delta float64) {
	// TODO
}

func (tank EnemyTank) Draw(window *pixelgl.Window) {
	if tank.alive {
		matrix := pixel.IM
		matrix = matrix.Moved(tank.position)
		// matrix = matrix.Scaled(tank.position, 1.0) // no need to scale, when scale is 1
		matrix = matrix.Rotated(tank.position, tank.rotationAngle)
		tank.tankSprite.Draw(window, matrix)
	}
}

type PlayerTank struct {
	tankSprite    *pixel.Sprite
	rotationAngle float64
	position      pixel.Vec
}

func NewPlayerTank(tankSprite *pixel.Sprite, initialPosition pixel.Vec) *PlayerTank {
	return &PlayerTank{tankSprite: tankSprite, position: initialPosition, rotationAngle: 0.0}
}

/*func (tank PlayerTank) Update() PlayerTank {
	return nil
}*/

func (tank PlayerTank) Draw(window *pixelgl.Window) {
	matrix := pixel.IM
	matrix = matrix.Moved(tank.position)
	// matrix = matrix.Scaled(tank.position, 1.0) // no need to scale, when scale is 1
	matrix = matrix.Rotated(tank.position, tank.rotationAngle)
	tank.tankSprite.Draw(window, matrix)
}

func (tank PlayerTank) Shoot() Bullet {
	picture, err := LoadPicture(BULLET_TEXTURE_PATH)
	HandleFatalError(err)
	bulletSprite := pixel.NewSprite(picture, picture.Bounds())
	return Bullet{
		bulletSprite:  bulletSprite,
		velocity:      BULLET_VELOCITY,
		position:      tank.position,
		rotationAngle: tank.rotationAngle,
	}
}

func (tank *PlayerTank) RotateClockWise(delta float64) *PlayerTank {
	result := *tank
	result.rotationAngle -= TANK_ROTATION_ANGLE * delta
	return &result
}

func (tank *PlayerTank) RotateAntiClockWise(delta float64) *PlayerTank {
	result := *tank
	result.rotationAngle += TANK_ROTATION_ANGLE * delta
	return &result
}

func (tank *PlayerTank) MoveUp(delta float64) *PlayerTank {
	result := *tank
	result.position = result.position.Add(pixel.V(0, (PLAYER_TANK_VELOCITY * delta)))
	return &result
}

func (tank *PlayerTank) MoveDown(delta float64) *PlayerTank {
	result := *tank
	result.position = result.position.Add(pixel.V(0, -(PLAYER_TANK_VELOCITY * delta)))
	return &result
}

func (tank *PlayerTank) MoveLeft(delta float64) *PlayerTank {
	result := *tank
	result.position = result.position.Add(pixel.V(-(PLAYER_TANK_VELOCITY * delta), 0))
	return &result
}

func (tank *PlayerTank) MoveRight(delta float64) *PlayerTank {
	result := *tank
	result.position = result.position.Add(pixel.V((PLAYER_TANK_VELOCITY * delta), 0))
	return &result
}
