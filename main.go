// test.go
package main

import (
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	TITLE                    string  = "Tank game"
	SCREEN_WIDTH             float64 = 500
	SCREEN_HEIGHT            float64 = 500
	PLAYER_TANK_TEXTURE_PATH string  = "resources/player-tank.png" // Make the textures, to be in same rotation angle
	ENEMY_TANK_TEXTURE_PATH  string  = "resources/enemy-tank.png"  // Make the textures, to be in same rotation angle
	BULLET_TEXTURE_PATH      string  = "resources/bullet_6.png"    // Make the textures, to be in same rotation angle
	BULLET_VELOCITY          float64 = 500
	TANK_ROTATION_ANGLE      float64 = 5
	PLAYER_TANK_VELOCITY     float64 = 300
	// TANK_SCALE               float64 = 1 // Make the textures, to be in same scale
	// BULLET_SCALE             float64 = 1 // Make the textures, to be in same scale
	// BULLET_INITIAL_ROTATION float64 = math.Pi / 2.0
	// TANK_INITIAL_ROTATION   float64 = 0
	VSYNC           bool = true
	SMOOTH_TEXTURES bool = false

	LEVEL_0_MAX_NUM_OF_ENEMY_TANKS         int     = 50
	LEVEL_0_ENEMY_SPAWN_OFF_TIME           float64 = 3.0 // seconds
	LEVEL_0_ENEMY_TANK_VELOCITY            float64 = PLAYER_TANK_VELOCITY + 10
	LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME float64 = 3.0 // seconds
)

/*
*
*

TODO : Use batch rendering, it is necessary for specifically a large number of bullets
*
*
*
*

*/

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  TITLE,
		Bounds: pixel.R(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT),
		VSync:  VSYNC,
	}
	win, err := pixelgl.NewWindow(cfg)
	HandleFatalError(err)

	win.SetSmooth(SMOOTH_TEXTURES)

	playerTankPicture, err := LoadPicture(PLAYER_TANK_TEXTURE_PATH)
	HandleFatalError(err)

	playerTankSprite := pixel.NewSprite(playerTankPicture, playerTankPicture.Bounds())

	playerTank := NewPlayerTank(playerTankSprite, win.Bounds().Center())
	var playerTankBullets []Bullet
	var enemyTankBullets []Bullet
	callbacks := map[pixelgl.Button]func(delta float64) *PlayerTank{
		pixelgl.KeyLeft:  playerTank.RotateAntiClockWise,
		pixelgl.KeyRight: playerTank.RotateClockWise,

		pixelgl.KeyW: playerTank.MoveUp,
		pixelgl.KeyA: playerTank.MoveLeft,
		pixelgl.KeyS: playerTank.MoveDown,
		pixelgl.KeyD: playerTank.MoveRight,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // TODO : What to do after 2262(UnixNano)???

	enemyTankPicture, err := LoadPicture(ENEMY_TANK_TEXTURE_PATH)
	HandleFatalError(err)
	enemyTankSprite := pixel.NewSprite(enemyTankPicture, enemyTankPicture.Bounds())
	enemyTanks := make([]EnemyTank, LEVEL_0_MAX_NUM_OF_ENEMY_TANKS)
	x := r.Intn(LEVEL_0_MAX_NUM_OF_ENEMY_TANKS / 2) // Initially random num.of tanks will be alive(halving it so that the generated random no. is not too much)
	i := 0
	for i = 0; i < x; i++ { // first x no.of tanks will be alive
		enemyTanks[i] = NewEnemyTank(enemyTankSprite, r.Float64()*(2.0*math.Pi), true, r.Float64()*LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME)
	}
	lastEnemyTankAlive := i - 1
	for i = (x - 1); i < len(enemyTanks); i++ { // rest of the tanks will be dead for now...
		enemyTanks[i] = NewEnemyTank(enemyTankSprite, r.Float64()*(2.0*math.Pi), false, r.Float64()*LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME)
	}

	SetPositions(enemyTanks, GetBoundingBox(playerTank.position, playerTank.tankSprite), r)

	last := time.Now()
	enemyTankSpawnTimer := 0.0

	for !win.Closed() {
		win.Clear(colornames.Bisque)

		dt := (time.Since(last).Seconds())
		last = time.Now()

		enemyTankSpawnTimer += dt
		if (enemyTankSpawnTimer >= LEVEL_0_ENEMY_SPAWN_OFF_TIME) && (lastEnemyTankAlive < len(enemyTanks)) {
			enemyTanks[lastEnemyTankAlive].alive = true // TODO : No bounds check(may panic)
			enemyTanks[lastEnemyTankAlive].position = GetPositionOfOneEnemyTank(enemyTanks[lastEnemyTankAlive], enemyTanks[:lastEnemyTankAlive], GetBoundingBox(playerTank.position, playerTank.tankSprite), r)
			lastEnemyTankAlive += 1
			enemyTankSpawnTimer = 0.0
		}

		// optimization(removing the bullets, which are out of the window)
		// range over slice will not work, as:
		// for i, _ := range ...{...}, here the maximum value of i is the length of the slice
		// i is initialized with length of the slice, but it doesn't assert new value of that langth, when the length of that slice changes
		// for i := 0; i < len(...); i++ {...} in this kind of loop the ;len(...); condition is always checked
		for i := 0; i < len(playerTankBullets); i++ {
			if playerTankBullets[i].OutOfWindow(win) {
				playerTankBullets = RemoveElementFromBulletSlice(playerTankBullets, i)
			}
		}
		for i := 0; i < len(enemyTankBullets); i++ {
			if enemyTankBullets[i].OutOfWindow(win) {
				enemyTankBullets = RemoveElementFromBulletSlice(enemyTankBullets, i)
			}
		}

		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}
		if win.JustPressed(pixelgl.KeySpace) {
			playerTankBullets = append(playerTankBullets, playerTank.Shoot())
		}

		for index, _ := range playerTankBullets {
			for idx, _ := range enemyTanks {
				enemyTankBoundingBox := GetBoundingBox(enemyTanks[idx].position, enemyTanks[idx].tankSprite)
				bulletNosePosition := playerTankBullets[index].bulletSprite.Frame().Max.Add(playerTankBullets[index].position)
				if enemyTankBoundingBox.Contains(bulletNosePosition) {
					enemyTanks[idx] = enemyTanks[idx].Die(dt)
				}
			}
		}

		for key, callbackFunc := range callbacks {
			if win.Pressed(key) {
				intersect := false
				experimentalPlayerTank := callbackFunc(dt)
				// collision detection with enemy tanks
				for index, _ := range enemyTanks {
					enemyTankBoundingBox := GetBoundingBox(enemyTanks[index].position, enemyTanks[index].tankSprite)
					playerTankBoundingBox := GetBoundingBox(experimentalPlayerTank.position, experimentalPlayerTank.tankSprite)
					if (enemyTanks[index].alive == true) && enemyTankBoundingBox.Intersects(playerTankBoundingBox) {
						intersect = true
						break
					}
				}
				if !intersect {
					// In the callbacks map, if they were declared like: playerTank.moveDown, where playerTank is an actual value, not a pointer to playerTank, then
					// the callback functions are 'bound' to that playerTank, with which they were initialized, changing the playerTank will not change the playerTank with
					// which they were initialized, so I used playerTank as a pointer.

					// TODO : playerTank = experimentalPlayerTank, will not work, why?
					playerTank.rotationAngle = experimentalPlayerTank.rotationAngle
					playerTank.position = experimentalPlayerTank.position
				}
			}
		}

		// TODO : Update first, or draw first?(most probably update first) + update order
		//		playerTank.Update()
		for index, _ := range playerTankBullets {
			playerTankBullets[index].Update(dt)
		}
		for index, _ := range enemyTankBullets {
			enemyTankBullets[index].Update(dt)
		}
		var bullet Bullet
		var experimentalEnemyTank EnemyTank
		//willUpdate := false
		for index, _ := range enemyTanks {
			if enemyTanks[index].alive && enemyTanks[index].WillUpdate() {
				enemyTanks[index].timer = time.Now()
				switch r.Intn(2) {
				case 0:
					experimentalEnemyTank = enemyTanks[index].MoveInRandomDir(dt, r)
				case 1:
					experimentalEnemyTank, bullet = enemyTanks[index].SpinAndShoot(dt, r, playerTank.position)
					enemyTankBullets = append(enemyTankBullets, bullet)
				}
				/*
					enemyTanks[index], bullet = enemyTanks[index].Update(dt, r, playerTank.position)
					if bullet != nil {
						enemyTankBullets = append(enemyTankBullets, bullet)
					}*/
				enemyTanks[index] = experimentalEnemyTank
			}
		}

		playerTank.Draw(win)
		for index, _ := range playerTankBullets {
			playerTankBullets[index].Draw(win)
		}
		for index, _ := range enemyTankBullets {
			enemyTankBullets[index].Draw(win)
		}
		for index, _ := range enemyTanks {
			enemyTanks[index].Draw(win)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
