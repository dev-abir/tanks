// test.go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/veandco/go-sdl2/sdl"
)

//==============SETTINGS==============
const (
	//==============WINDOW SETTINGS==============
	TITLE           string = "Tank game"
	SCREEN_WIDTH    int32  = 500
	SCREEN_HEIGHT   int32  = 500
	VSYNC           bool   = true
	SMOOTH_TEXTURES bool   = false

	//==============TEXTURE PATHS==============
	PLAYER_TANK_TEXTURE_PATH string = "resources/player-tank.png" // Make the textures, to be in same rotation angle
	ENEMY_TANK_TEXTURE_PATH  string = "resources/enemy-tank.png"  // Make the textures, to be in same rotation angle
	BULLET_TEXTURE_PATH      string = "resources/bullet_6.png"    // Make the textures, to be in same rotation angle

	//==============GAMEPLAY==============
	BULLET_VELOCITY      float32 = 500
	TANK_ROTATION_ANGLE  float32 = 500 // TODO : Why so low rotation on setting this to 5?(maybe due to delta calculation)
	PLAYER_TANK_VELOCITY float32 = 300
)

//==============LEVEL SETTINGS==============
const (
	LEVEL_0_MAX_NUM_OF_ENEMY_TANKS         int     = 30
	LEVEL_0_ENEMY_SPAWN_OFF_TIME           float32 = 3.0 // seconds
	LEVEL_0_ENEMY_TANK_VELOCITY            float32 = PLAYER_TANK_VELOCITY + 10
	LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME float32 = 3.0 // seconds
)

/*
TODO : Use batch rendering, it is necessary for specifically a large number of bullets
*/

func run() int {

	//==============VARS==============
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // TODO : What to do after 2262(UnixNano)???

	/* TODO : Most probably depricated
	//==============SDL INIT==============
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		HandleError("Failed to initialize sdl:", err)
		os.Exit(ERROR_FAILED_TO_INIT_SDL)
	}*/

	//==============CREATE WINDOW==============
	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		HandleError("Failed to create window: ", err)
		return ERROR_FAILED_TO_CREATE_WINDOW
	}
	defer window.Destroy()

	//==============CREATE RENDERER==============
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		HandleError("Failed to create renderer: ", err)
		return ERROR_FAILED_TO_CREATE_RENDERER
	}
	defer renderer.Destroy()

	/* TODO:(GIVEN THE SAME vsync FLAGS IN sdl.createrenderer)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")*/

	//==============PLAYER TANK==============
	playerTankImage, playerTankTexture, errorCode := GetTexture(PLAYER_TANK_TEXTURE_PATH, renderer)
	if errorCode != 0 {
		return errorCode
	}
	defer playerTankImage.Free()
	defer playerTankTexture.Destroy()

	playerTank := &PlayerTank{
		tankTexture:   playerTankTexture,
		rotationAngle: 0.0,
		boundingBox: sdl.FRect{
			X: float32(SCREEN_WIDTH) / 2.0,
			Y: float32(SCREEN_HEIGHT) / 2.0,
			W: float32(playerTankImage.W),
			H: float32(playerTankImage.H),
		},
	}
	var playerTankBullets []Bullet

	//==============CALLBACKS==============
	callbacks := map[sdl.Keycode]func(delta float32) *PlayerTank{
		sdl.K_LEFT:  playerTank.RotateAntiClockWise,
		sdl.K_RIGHT: playerTank.RotateClockWise,

		sdl.K_w: playerTank.MoveUp,
		sdl.K_a: playerTank.MoveLeft,
		sdl.K_s: playerTank.MoveDown,
		sdl.K_d: playerTank.MoveRight,
	}

	//==============ENEMY TANKS==============
	enemyTankImage, enemyTankTexture, errorCode := GetTexture(ENEMY_TANK_TEXTURE_PATH, renderer)
	if errorCode != 0 {
		return errorCode
	}
	defer enemyTankImage.Free()
	defer enemyTankTexture.Destroy()

	enemyTanks := make([]EnemyTank, LEVEL_0_MAX_NUM_OF_ENEMY_TANKS)
	x := r.Intn(LEVEL_0_MAX_NUM_OF_ENEMY_TANKS / 2) // Initially random num.of tanks will be alive(halving it so that the generated random no. is not too much)
	i := 0
	for i = 0; i < x; i++ { // first x no.of tanks will be alive
		enemyTanks[i] = NewEnemyTank(enemyTankTexture, enemyTankImage.W, enemyTankImage.H, r.Float32()*(2.0*math.Pi), true, r.Float32()*LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME)
	}
	lastEnemyTankAlive := i - 1
	for i = (x - 1); i < len(enemyTanks); i++ { // rest of the tanks will be dead for now...
		enemyTanks[i] = NewEnemyTank(enemyTankTexture, enemyTankImage.W, enemyTankImage.H, r.Float32()*(2.0*math.Pi), false, r.Float32()*LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME)
	}
	SetPositions(enemyTanks, playerTank.boundingBox, r)
	var enemyTankBullets []Bullet

	last := time.Now() // for calculating dt(delta)
	var enemyTankSpawnTimer float32 = 0.0
	fpsCounter := 0
	var fpsTimer float32 = 0.0

	//==============BULLET==============
	bulletImage, bulletTexture, errorCode := GetTexture(BULLET_TEXTURE_PATH, renderer)
	if errorCode != 0 {
		return errorCode
	}
	defer bulletImage.Free()
	defer bulletTexture.Destroy()

	//==============MAIN LOOP==============
	running := true // Main loop flag
	for running {

		//==============CALCULATING dt(DELTA)==============
		dt := float32(time.Since(last).Seconds())
		last = time.Now()

		//==============PRINTING FPS==============
		fpsTimer += dt
		if fpsTimer >= 1.0 {
			fmt.Println("Current FPS: ", fpsCounter)
			fpsTimer = 0.0
			fpsCounter = 0
		}

		//==============SPAWNING NEW ENEMY TANKS==============
		enemyTankSpawnTimer += dt
		if (enemyTankSpawnTimer >= LEVEL_0_ENEMY_SPAWN_OFF_TIME) && (lastEnemyTankAlive < len(enemyTanks)) {
			enemyTanks[lastEnemyTankAlive].alive = true
			enemyTanks[lastEnemyTankAlive].boundingBox = GetPositionOfOneEnemyTank(enemyTanks[lastEnemyTankAlive].boundingBox, enemyTanks[:lastEnemyTankAlive], playerTank.boundingBox, r)
			lastEnemyTankAlive += 1
			enemyTankSpawnTimer = 0.0
		}

		//==============UPDATING ENEMY TANKS==============
		var bullet Bullet
		var experimentalEnemyTank EnemyTank
		for index, _ := range enemyTanks {
			if enemyTanks[index].alive && enemyTanks[index].WillUpdate() {
				switch r.Intn(2) {
				case 0:
					experimentalEnemyTank = enemyTanks[index].MoveInRandomDir(dt, r)
				case 1:
					experimentalEnemyTank, bullet = enemyTanks[index].SpinAndShoot(dt, r, sdl.FPoint{
						X: playerTank.boundingBox.X,
						Y: playerTank.boundingBox.Y,
					}, bulletTexture, bulletImage.W, bulletImage.H)
					enemyTankBullets = append(enemyTankBullets, bullet)
				}
				if ValidPosition(experimentalEnemyTank.boundingBox, enemyTanks[:index], playerTank.boundingBox) {
					enemyTanks[index] = experimentalEnemyTank
				}
			}
		}

		//==============UPDATING ENEMY TANK BULLETS==============
		for index, _ := range enemyTankBullets {
			enemyTankBullets[index].Update(dt)
		}

		//==============OPTIMIZATON(removing the bullets, which are out of the window)==============
		// range over slice will not work, as:
		// for i, _ := range ...{...}, here the maximum value of i is the length of the slice
		// i is initialized with length of the slice, but it doesn't assert new value of that langth, when the length of that slice changes
		// for i := 0; i < len(...); i++ {...} in this kind of loop the ;len(...); condition is always checked
		for i := 0; i < len(playerTankBullets); i++ {
			if playerTankBullets[i].OutOfWindow() {
				playerTankBullets = RemoveElementFromBulletSlice(playerTankBullets, i)
			}
		}
		for i := 0; i < len(enemyTankBullets); i++ {
			if enemyTankBullets[i].OutOfWindow() {
				enemyTankBullets = RemoveElementFromBulletSlice(enemyTankBullets, i)
			}
		}

		//==============EVENT HANDLING==============
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
				if t.Keysym.Sym == sdl.K_SPACE {
					playerTankBullets = append(playerTankBullets, playerTank.Shoot(bulletTexture, bulletImage.W, bulletImage.H))
				}
				for key, callbackFunc := range callbacks {
					if t.Keysym.Sym == key {
						intersect := false
						experimentalPlayerTank := callbackFunc(dt)
						// collision detection with enemy tanks and window
						for index, _ := range enemyTanks {
							if (enemyTanks[index].alive == true) && enemyTanks[index].boundingBox.HasIntersection(&playerTank.boundingBox) {
								intersect = true
								break
							}
						}
						// if no collision with enemy tanks and window
						if !intersect &&
							(playerTank.boundingBox.X > 0.0) &&
							(playerTank.boundingBox.Y > 0.0) &&
							((playerTank.boundingBox.X + playerTank.boundingBox.W) < float32(SCREEN_WIDTH)) &&
							((playerTank.boundingBox.Y + playerTank.boundingBox.H) < float32(SCREEN_HEIGHT)) {
							// In the callbacks map, if they were declared like: playerTank.moveDown, where playerTank is an actual value, not a pointer to playerTank, then
							// the callback functions are 'bound' to that playerTank, with which they were initialized, changing the playerTank will not change the playerTank with
							// which they were initialized, so I used playerTank as a pointer.

							// TODO : playerTank = experimentalPlayerTank, will not work, why?
							playerTank.boundingBox = experimentalPlayerTank.boundingBox
							playerTank.rotationAngle = experimentalPlayerTank.rotationAngle
						}
					}
				}
			}
		}

		//==============UPDATING PLAYER TANK BULLETS==============
		for index, _ := range playerTankBullets {
			playerTankBullets[index].Update(dt)
		}

		//==============DESTROYING ENEMY TANKS(by player tank bullets)==============
		for index, _ := range playerTankBullets {
			for idx, _ := range enemyTanks {
				bulletNosePosition := sdl.FPoint{
					playerTankBullets[index].boundingBox.X + playerTankBullets[index].boundingBox.W,
					playerTankBullets[index].boundingBox.Y + (playerTankBullets[index].boundingBox.H / 2.0),
				}
				if bulletNosePosition.InRect(&enemyTanks[idx].boundingBox) {
					enemyTanks[idx].Die(dt)
				}
			}
		}

		// TODO : Update first, or draw first?(most probably update first) + update order
		//		playerTank.Update()

		//==============CLEARING THE SCREEN==============
		renderer.SetDrawColor(colornames.Bisque.R, colornames.Bisque.G, colornames.Bisque.B, colornames.Bisque.A)
		renderer.Clear()

		//==============DRAWING==============
		//renderer.CopyExF(playerTank.tankTexture, nil, &playerTank.boundingBox, float64(playerTank.rotationAngle), nil, sdl.FLIP_NONE)
		/*

			TODO : For some reason CopyExF is not working.......

		*/

		renderer.CopyEx(playerTank.tankTexture, nil, &sdl.Rect{
			int32(playerTank.boundingBox.X),
			int32(playerTank.boundingBox.Y),
			int32(playerTank.boundingBox.W),
			int32(playerTank.boundingBox.H)}, float64(playerTank.rotationAngle), nil, sdl.FLIP_NONE)
		for index, _ := range playerTankBullets {
			renderer.CopyEx(playerTankBullets[index].bulletTexture, nil, &sdl.Rect{
				int32(playerTankBullets[index].boundingBox.X),
				int32(playerTankBullets[index].boundingBox.Y),
				int32(playerTankBullets[index].boundingBox.W),
				int32(playerTankBullets[index].boundingBox.H)}, float64(playerTankBullets[index].rotationAngle), nil, sdl.FLIP_NONE)
		}
		for index, _ := range enemyTanks {
			if enemyTanks[index].alive {
				renderer.CopyEx(enemyTanks[index].tankTexture, nil, &sdl.Rect{
					int32(enemyTanks[index].boundingBox.X),
					int32(enemyTanks[index].boundingBox.Y),
					int32(enemyTanks[index].boundingBox.W),
					int32(enemyTanks[index].boundingBox.H)}, float64(enemyTanks[index].rotationAngle), nil, sdl.FLIP_NONE)
			}
		}
		for index, _ := range enemyTankBullets {
			renderer.CopyEx(enemyTankBullets[index].bulletTexture, nil, &sdl.Rect{
				int32(enemyTankBullets[index].boundingBox.X),
				int32(enemyTankBullets[index].boundingBox.Y),
				int32(enemyTankBullets[index].boundingBox.W),
				int32(enemyTankBullets[index].boundingBox.H)}, float64(enemyTankBullets[index].rotationAngle), nil, sdl.FLIP_NONE)
		}
		/* Just for debugging....*/
		renderer.SetDrawColor(colornames.Red.R, colornames.Red.G, colornames.Red.B, colornames.Red.A)
		for index, _ := range enemyTanks {
			if enemyTanks[index].alive {
				renderer.DrawRect(&sdl.Rect{
					int32(enemyTanks[index].boundingBox.X),
					int32(enemyTanks[index].boundingBox.Y),
					int32(enemyTanks[index].boundingBox.W),
					int32(enemyTanks[index].boundingBox.H),
				})
			}
		}
		renderer.Present()

		//==============UPDATING FPS COUNTER==============
		fpsCounter += 1
	}

	//sdl.Quit()
	return 0
}

func main() {
	os.Exit(run()) // TODO : What does it do?
}
