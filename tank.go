// tank.go
package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

/*

It is not easy, to create this file to follow only one function set.
Let's say: tank.moveForward(angle float64) -> this seems to be taking the tank pointer, and change it's angle.
But in my case, it will take a value receiver, change it's angle and return that. ((a tank) func moveForwardangle float64) tank)
this is because I need to check whether the tank is overlapping with other game objects or not.
This checking(or rather collision detection) has been done in the main.go file, because that is the central place, where I
have acces to all game objects, so I have to essentially return a new tank, and not mutate the receiver.

In some cases the function signature makes it very clear that it will mutate, for those cases I have used pointer receivers
like bullet.Update()

func (tank *EnemyTank) WillUpdate() bool -> this is exceptional, it looks like it's not mutating but it needs to do that
see it's usage in main.go, there is no better or elegant way of doing it in any other way(maybe)

callbacks need to be pointers, value receivers will 'bind' the function with the receiver. Changing the receiver is not possible.(see player update in main.go)


*/

type Bullet struct {
	bulletTexture *sdl.Texture
	velocity      float32
	boundingBox   sdl.FRect
	rotationAngle float32
}

func (bullet *Bullet) Update(delta float32) {
	bullet.boundingBox.X += bullet.velocity * delta * float32(math.Cos(DegreeToRadian(float64(bullet.rotationAngle))))
	bullet.boundingBox.Y += bullet.velocity * delta * float32(math.Sin(DegreeToRadian(float64(bullet.rotationAngle))))
}

type Explosion struct {
	position            sdl.Point
	explosionTexture    *sdl.Texture
	timer               time.Time
	noUpdateTime        float32
	animationCoordIndex int
	died                bool
}

func NewExplosion(tankBoundingBox sdl.FRect, explosionTexture *sdl.Texture) Explosion {
	return Explosion{
		position: sdl.Point{
			int32(tankBoundingBox.X) - ((CELL_WIDTH - int32(tankBoundingBox.W)) / 2),
			int32(tankBoundingBox.Y) - ((CELL_HEIGHT - int32(tankBoundingBox.H)) / 2),
		}, // positioning exactly at the centre of the tank
		explosionTexture:    explosionTexture,
		timer:               time.Now(),
		noUpdateTime:        EXPLOSION_ANIMATION_LIFE_SPAN / float32(len(EXPLOSION_ANIMATION_COORDS)),
		animationCoordIndex: 0,
		died:                false,
	}
}

func (explosion *Explosion) Update() {
	if explosion.animationCoordIndex == (len(EXPLOSION_ANIMATION_COORDS) - 1) {
		explosion.died = true
	} else if time.Since(explosion.timer).Seconds() >= float64(explosion.noUpdateTime) {
		explosion.timer = time.Now()
		explosion.animationCoordIndex += 1
	}
}

func (explosion Explosion) Draw(renderer *sdl.Renderer) {
	renderer.Copy(explosion.explosionTexture,
		&sdl.Rect{
			EXPLOSION_ANIMATION_COORDS[explosion.animationCoordIndex].X,
			EXPLOSION_ANIMATION_COORDS[explosion.animationCoordIndex].Y,
			CELL_WIDTH,
			CELL_HEIGHT},
		&sdl.Rect{
			explosion.position.X,
			explosion.position.Y,
			CELL_WIDTH,
			CELL_WIDTH})
}

type EnemyTank struct {
	tankTexture                  *sdl.Texture
	rotationAngle                float32
	boundingBox                  sdl.FRect
	noUpdateTime                 float32
	timer                        time.Time
	rotationAnimationTargetAngle float32
}

func NewEnemyTank(tankTexture *sdl.Texture, width int32, height int32, initialRotationAngle float32, noUpdateTime float32) EnemyTank {
	return EnemyTank{
		tankTexture:   tankTexture,
		rotationAngle: initialRotationAngle,
		boundingBox: sdl.FRect{
			X: 0.0,
			Y: 0.0,
			W: float32(width),
			H: float32(height),
		},
		noUpdateTime: noUpdateTime,
		timer:        time.Now(),
	}
}

/*func (tank EnemyTank) Update(delta float64, r *rand.Rand, playerTankPosition pixel.Vec) (EnemyTank, Bullet) {
	var bullet Bullet

	return tank, bullet
}*/

func (tank EnemyTank) MoveInRandomDir(delta float32, r *rand.Rand) EnemyTank {
	switch r.Intn(4) {
	case 0:
		tank.boundingBox.Y += LEVEL_0_ENEMY_TANK_VELOCITY * delta // DOWN
	case 1:
		tank.boundingBox.Y -= LEVEL_0_ENEMY_TANK_VELOCITY * delta // UP
	case 2:
		tank.boundingBox.X += LEVEL_0_ENEMY_TANK_VELOCITY * delta // RIGHT
	case 3:
		tank.boundingBox.X -= LEVEL_0_ENEMY_TANK_VELOCITY * delta // LEFT
	}
	return tank
}

func (tank *EnemyTank) Rotate(r *rand.Rand, playerTankPosition sdl.FPoint) {
	/*
		TODO :
		switch r.Intn(2) {
		case 0:
			displacementVector := playerTankPosition.Sub(tank.position) // SHOOT THE PLAYER
			tank.rotationAnimationTargetAngle = displacementVector.Angle()
		case 1:*/
	tank.rotationAnimationTargetAngle = r.Float32() * 360.0 // SHOOT ANYWHERE RANDOMLY
	/*}*/
	if tank.rotationAngle >= 360.0 { // reset angle to 0, or else glitches/bugs are welcome...
		tank.rotationAngle = 0.0
	}
}

func (tank EnemyTank) Shoot(bulletTexture *sdl.Texture, bulletWidth int32, bulletHeight int32) Bullet {
	return Bullet{
		bulletTexture: bulletTexture,
		velocity:      BULLET_VELOCITY,
		boundingBox: sdl.FRect{
			/*X: tank.boundingBox.X,
			Y: tank.boundingBox.Y,*/
			X: tank.boundingBox.X + (tank.boundingBox.W / 2.0) - (float32(bulletWidth) / 2.0),  // shooting from the centre of the tank, and putting the bullet'scentre at the centre of the tank
			Y: tank.boundingBox.Y + (tank.boundingBox.H / 2.0) - (float32(bulletHeight) / 2.0), // shooting from the centre of the tank, and putting the bullet's centre at the centre of the tank
			W: float32(bulletWidth),
			H: float32(bulletHeight),
		},
		rotationAngle: tank.rotationAngle,
	}
}

func (tank *EnemyTank) UpdateAnimation(delta float32) {
	if tank.rotationAngle < tank.rotationAnimationTargetAngle {
		tank.rotationAngle += TANK_ROTATION_ANGLE * delta
	}
}

/*This function seems to be very innocent, not mutating the receiver.
Actually, it changes the timer of the receiver. Be careful...*/
func (tank *EnemyTank) WillUpdate() bool {
	if time.Since(tank.timer).Seconds() >= float64(tank.noUpdateTime) {
		tank.timer = time.Now()
		return true
	}
	return false
}

/*func (tank EnemyTank) Draw(window *sdl.Window) {
	if tank.alive {
		matrix := pixel.IM
		matrix = matrix.Moved(tank.position)
		// matrix = matrix.Scaled(tank.position, 1.0) // no need to scale, when scale is 1
		matrix = matrix.Rotated(tank.position, tank.rotationAngle)
		tank.tankSprite.Draw(window, matrix)
	}
}*/

type PlayerTank struct {
	tankTexture   *sdl.Texture
	rotationAngle float32
	boundingBox   sdl.FRect
}

/*func (tank PlayerTank) Update() PlayerTank {
	return nil
}*/

func (tank *PlayerTank) Shoot(bulletTexture *sdl.Texture, bulletWidth int32, bulletHeight int32) Bullet {
	return Bullet{
		bulletTexture: bulletTexture,
		velocity:      BULLET_VELOCITY,
		boundingBox: sdl.FRect{
			X: tank.boundingBox.X + (tank.boundingBox.W / 2.0) - (float32(bulletWidth) / 2.0),  // shooting from the centre of the tank, and putting the bullet'scentre at the centre of the tank
			Y: tank.boundingBox.Y + (tank.boundingBox.H / 2.0) - (float32(bulletHeight) / 2.0), // shooting from the centre of the tank, and putting the bullet's centre at the centre of the tank
			W: float32(bulletWidth),
			H: float32(bulletHeight),
		},
		rotationAngle: tank.rotationAngle,
	}
}

func (tank *PlayerTank) RotateClockWise(delta float32) *PlayerTank {
	result := *tank                                     // making a copy
	result.rotationAngle += TANK_ROTATION_ANGLE * delta // mutating that copy
	return &result                                      // returning pointer to that copy
}

func (tank *PlayerTank) RotateAntiClockWise(delta float32) *PlayerTank {
	result := *tank                                     // making a copy
	result.rotationAngle -= TANK_ROTATION_ANGLE * delta // mutating that copy
	return &result                                      // returning pointer to that copy
}

func (tank *PlayerTank) MoveUp(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.Y -= PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copy
}

func (tank *PlayerTank) MoveDown(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.Y += PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copy
}

func (tank *PlayerTank) MoveLeft(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.X -= PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copy
}

func (tank *PlayerTank) MoveRight(delta float32) *PlayerTank {
	result := *tank                                      // making a copy
	result.boundingBox.X += PLAYER_TANK_VELOCITY * delta // mutating that copy
	return &result                                       // returning pointer to that copy
}
