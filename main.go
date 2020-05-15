// main.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

//==============SETTINGS==============
const (
	//==============WINDOW SETTINGS==============
	TITLE           string = "Tank game"
	SCREEN_WIDTH    int32  = 500
	SCREEN_HEIGHT   int32  = 500
	VSYNC           bool   = true
	SMOOTH_TEXTURES bool   = true
)

const (
	//==============ANTI-ALIASING TYPES==============
	NEAREST     string = "0"
	LINEAR      string = "1"
	ANISOTROPIC string = "2"
)

const (
	//==============TEXTURE PATHS==============
	PLAYER_TANK_TEXTURE_PATH string = "resources/player-tank.png" // Make the textures, to be in same rotation angle
	ENEMY_TANK_TEXTURE_PATH  string = "resources/enemy-tank.png"  // Make the textures, to be in same rotation angle
	BULLET_TEXTURE_PATH      string = "resources/bullet_6.png"    // Make the textures, to be in same rotation angle

	//==============SOUND EFFECTS PATHS==============
	SHOOT_SOUND_PATH     string = "resources/flak_gun_sound.ogg"
	EXPLOSION_SOUND_PATH string = "resources/bombexplosion.ogg"
)

//==============EXPLOSION ANIMATION==============
const (
	EXPLOSION_ANIMATION_TEXTURE_PATH string = "resources/explosion_animation.png"

	// for now, these are hard coded(change these, if you change the texture image)
	CELL_WIDTH                   int32 = 128
	CELL_HEIGHT                  int32 = 128
	NUMBER_OF_CELLS_LENGTH_WISE  int   = 8
	NUMBER_OF_CELLS_BREADTH_WISE int   = 8
)

var EXPLOSION_ANIMATION_COORDS [NUMBER_OF_CELLS_LENGTH_WISE * NUMBER_OF_CELLS_BREADTH_WISE]sdl.Point = [NUMBER_OF_CELLS_LENGTH_WISE * NUMBER_OF_CELLS_BREADTH_WISE]sdl.Point{
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 0}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 0}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 0}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 0}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 0}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 0}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 0}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 0},
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 1}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 1}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 1}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 1}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 1}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 1}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 1}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 1},
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 2}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 2}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 2}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 2}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 2}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 2}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 2}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 2},
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 3}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 3}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 3}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 3}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 3}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 3}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 3}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 3},
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 4}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 4}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 4}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 4}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 4}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 4}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 4}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 4},
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 5}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 5}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 5}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 5}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 5}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 5}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 5}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 5},
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 6}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 6}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 6}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 6}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 6}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 6}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 6}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 6},
	sdl.Point{CELL_WIDTH * 0, CELL_HEIGHT * 7}, sdl.Point{CELL_WIDTH * 1, CELL_HEIGHT * 7}, sdl.Point{CELL_WIDTH * 2, CELL_HEIGHT * 7}, sdl.Point{CELL_WIDTH * 3, CELL_HEIGHT * 7}, sdl.Point{CELL_WIDTH * 4, CELL_HEIGHT * 7}, sdl.Point{CELL_WIDTH * 5, CELL_HEIGHT * 7}, sdl.Point{CELL_WIDTH * 6, CELL_HEIGHT * 7}, sdl.Point{CELL_WIDTH * 7, CELL_HEIGHT * 7},
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	//==============GAMEPLAY==============
	BULLET_VELOCITY               float32 = 500
	TANK_ROTATION_ANGLE           float32 = 500 // TODO : Why so low rotation on setting this to 5?(maybe due to delta calculation)
	PLAYER_TANK_VELOCITY          float32 = 300
	EXPLOSION_ANIMATION_LIFE_SPAN float32 = 0.5 // seconds

	//==============SPECIAL FLAGS==============
	CRAZY_TANKS bool = false // A special flag, toggle it, and enjoy!
)

//==============LEVEL SETTINGS==============
const (
	LEVEL_0_MAX_NUM_OF_ENEMY_TANKS         int     = 10
	LEVEL_0_ENEMY_SPAWN_OFF_TIME           float32 = 3.0 // seconds
	LEVEL_0_ENEMY_TANK_VELOCITY            float32 = 310
	LEVEL_0_ENEMY_TANK_MIN_NO_UPDATES_TIME float32 = 1.0 // seconds
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
	var renderer *sdl.Renderer
	if VSYNC {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	} else {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	}
	if err != nil {
		HandleError("Failed to create renderer: ", err)
		return ERROR_FAILED_TO_CREATE_RENDERER
	}
	defer renderer.Destroy()

	// sdl.SetHint(sdl.HINT_DEFAULT)

	if SMOOTH_TEXTURES {
		sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, LINEAR)
	}

	//==============SOUND EFFECTS==============
	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		HandleError("Cannot play audio, you may play without it: ", err)
	}
	defer mix.CloseAudio()

	shootSoundEffect := GetSoundEffect(SHOOT_SOUND_PATH)
	defer shootSoundEffect.Free()

	explosionSoundEffect := GetSoundEffect(EXPLOSION_SOUND_PATH)
	defer explosionSoundEffect.Free()

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
			X: (float32(SCREEN_WIDTH) / 2.0) - (float32(playerTankImage.W) / 2.0),  // positioning exactly at the centre of the screen
			Y: (float32(SCREEN_HEIGHT) / 2.0) - (float32(playerTankImage.H) / 2.0), // positioning exactly at the centre of the screen
			W: float32(playerTankImage.W),
			H: float32(playerTankImage.H),
		},
	}
	var playerTankBullets []Bullet

	//==============CALLBACKS==============
	callbacks := map[sdl.Scancode]func(delta float32) *PlayerTank{
		sdl.SCANCODE_LEFT:  playerTank.RotateAntiClockWise,
		sdl.SCANCODE_RIGHT: playerTank.RotateClockWise,

		sdl.SCANCODE_W: playerTank.MoveUp,
		sdl.SCANCODE_A: playerTank.MoveLeft,
		sdl.SCANCODE_S: playerTank.MoveDown,
		sdl.SCANCODE_D: playerTank.MoveRight,
	}

	//==============ENEMY TANKS==============
	enemyTankImage, enemyTankTexture, errorCode := GetTexture(ENEMY_TANK_TEXTURE_PATH, renderer)
	if errorCode != 0 {
		return errorCode
	}
	defer enemyTankImage.Free()
	defer enemyTankTexture.Destroy()

	x := 2 + r.Intn(LEVEL_0_MAX_NUM_OF_ENEMY_TANKS/2) // Initially random num.of tanks will be alive(halving it so that the generated random no. is not too much), let us make at least 2 tanks alive at first
	enemyTanks := make([]EnemyTank, x)
	i := 0
	for i = 0; i < x; i++ { // first x no.of tanks will be alive
		if CRAZY_TANKS {
			enemyTanks[i] = NewEnemyTank(enemyTankTexture, enemyTankImage.W, enemyTankImage.H, r.Float32()*360.0, GetRandomFloat32(0.0, 0.5, r))
		} else {
			enemyTanks[i] = NewEnemyTank(enemyTankTexture, enemyTankImage.W, enemyTankImage.H, r.Float32()*360.0, GetRandomFloat32(LEVEL_0_ENEMY_TANK_MIN_NO_UPDATES_TIME, LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME, r))
		}
	}

	numOfEnemyTanksSpawned := x

	SetPositionOfEnemyTanks(enemyTanks, playerTank.boundingBox, r)
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

	//==============EXPLOSION ANIMATION==============
	var explosions []Explosion
	explosionImage, explosionTexture, errorCode := GetTexture(EXPLOSION_ANIMATION_TEXTURE_PATH, renderer)
	if errorCode != 0 {
		return errorCode
	}
	defer explosionImage.Free()
	defer explosionTexture.Destroy()

	//==============MAIN LOOP==============
	keyboardState := sdl.GetKeyboardState() // for handling keyboard events
	running := true                         // Main loop flag
	playerShootedInLastFrame := false       // just a flag, to manage player tank shoot events(it ensures that the player tank will not shoot continuously, on pressing down 'space')
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

		//==============CHECKING WHETHER PLAYER HAS WON==============
		if (len(enemyTanks) == 0) /*if all the tanks has been destroyed by the player, and*/ &&
			(numOfEnemyTanksSpawned == LEVEL_0_MAX_NUM_OF_ENEMY_TANKS) /*if all the tanks has been spawned*/ {
			// TODO: Render text, instead of printing inside the console
			fmt.Println("==============PLAYER WON==============")
			sdl.Delay(1000)
			break
		}

		//==============SPAWNING NEW ENEMY TANKS==============
		enemyTankSpawnTimer += dt
		if (enemyTankSpawnTimer >= LEVEL_0_ENEMY_SPAWN_OFF_TIME) && (numOfEnemyTanksSpawned < LEVEL_0_MAX_NUM_OF_ENEMY_TANKS) {
			if CRAZY_TANKS {
				enemyTanks = append(enemyTanks, NewEnemyTank(enemyTankTexture, enemyTankImage.W, enemyTankImage.H, r.Float32()*360.0, GetRandomFloat32(0.0, 0.5, r)))
			} else {
				enemyTanks = append(enemyTanks, NewEnemyTank(enemyTankTexture, enemyTankImage.W, enemyTankImage.H, r.Float32()*360.0, GetRandomFloat32(LEVEL_0_ENEMY_TANK_MIN_NO_UPDATES_TIME, LEVEL_0_ENEMY_TANK_MAX_NO_UPDATES_TIME, r)))
			}
			IndexOfLastEnemyTank := len(enemyTanks) - 1
			enemyTanks[IndexOfLastEnemyTank].boundingBox = GetPositionOfOneEnemyTank(enemyTanks[IndexOfLastEnemyTank].boundingBox, enemyTanks[:IndexOfLastEnemyTank], playerTank.boundingBox, r)
			enemyTankSpawnTimer = 0.0
			numOfEnemyTanksSpawned += 1
		}

		//==============UPDATING ENEMY TANKS==============
		var bullet Bullet
		var experimentalEnemyTank EnemyTank
		for index := range enemyTanks {

			//==============UPDATING ANIMATION(ON EVERY FRAME)==============
			enemyTanks[index].UpdateAnimation(dt)

			//==============UPDATING POSITION AND ROTATION, SHOOTING BULLETS==============
			if enemyTanks[index].WillUpdate() { // updating based on an update timer
				switch r.Intn(3) {
				case 0:
					experimentalEnemyTank = enemyTanks[index].MoveInRandomDir(dt, r)
					if ValidPosition(experimentalEnemyTank.boundingBox, enemyTanks, playerTank.boundingBox) {
						enemyTanks[index].boundingBox = experimentalEnemyTank.boundingBox
					}
				case 1:
					enemyTanks[index].Rotate(r, sdl.FPoint{
						X: playerTank.boundingBox.X,
						Y: playerTank.boundingBox.Y,
					})
				case 2:
					bullet = enemyTanks[index].Shoot(bulletTexture, bulletImage.W, bulletImage.H)
					enemyTankBullets = append(enemyTankBullets, bullet)
					PlaySoundEffect(shootSoundEffect)
				}
			}
		}

		//==============UPDATING ENEMY TANK BULLETS==============
		for index := range enemyTankBullets {
			enemyTankBullets[index].Update(dt)
		}

		//==============OPTIMIZATON(removing the bullets, which are out of the window)==============
		// range over slice will not work, as:
		// for i, _ := range ...{...}, here the maximum value of i is the length of the slice
		// i is initialized with length of the slice, but it doesn't assert new value of that length, when the length of that slice changes
		// for i := 0; i < len(...); i++ {...} in this kind of loop the ;len(...); condition is always checked
		for i := 0; i < len(playerTankBullets); i++ {
			if !IsInsideWindow(playerTankBullets[i].boundingBox) {
				playerTankBullets = RemoveElementFromBulletSlice(playerTankBullets, i)
			}
		}
		for i := 0; i < len(enemyTankBullets); i++ {
			if !IsInsideWindow(enemyTankBullets[i].boundingBox) {
				enemyTankBullets = RemoveElementFromBulletSlice(enemyTankBullets, i)
			}
		}

		//==============EVENT HANDLING==============
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if event.GetType() == sdl.QUIT {
				running = false
			}
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == sdl.K_SPACE {
					if event.GetType() == sdl.KEYDOWN {
						if !playerShootedInLastFrame {
							playerTankBullets = append(playerTankBullets, playerTank.Shoot(bulletTexture, bulletImage.W, bulletImage.H))
							PlaySoundEffect(shootSoundEffect)
							playerShootedInLastFrame = true
						}
					}
					if event.GetType() == sdl.KEYUP {
						playerShootedInLastFrame = false
					}
				}
			}
		}

		// sdl.PumpEvents() // not required

		if keyboardState[sdl.SCANCODE_ESCAPE] == 1 {
			running = false
		}

		for key, callbackFunc := range callbacks {
			if keyboardState[key] == 1 {
				intersect := false
				experimentalPlayerTank := callbackFunc(dt)
				// collision detection with enemy tanks and window
				for index := range enemyTanks {
					if enemyTanks[index].boundingBox.HasIntersection(&experimentalPlayerTank.boundingBox) {
						intersect = true
						break
					}
				}
				// if no collision with enemy tanks and window
				if !intersect && IsInsideWindow(experimentalPlayerTank.boundingBox) {
					// In the callbacks map, if they were declared like: playerTank.moveDown, where playerTank is an actual value, not a pointer to playerTank, then
					// the callback functions are 'bound' to that playerTank, with which they were initialized, changing the playerTank will not change the playerTank with
					// which they were initialized, so I used playerTank as a pointer.

					// TODO : playerTank = experimentalPlayerTank, will not work, why?
					playerTank.boundingBox = experimentalPlayerTank.boundingBox
					playerTank.rotationAngle = experimentalPlayerTank.rotationAngle
				}
			}
		}

		//==============UPDATING PLAYER TANK BULLETS==============
		for index := range playerTankBullets {
			playerTankBullets[index].Update(dt)
		}

		//==============DESTROYING ENEMY TANKS(by player tank bullets)==============
		for index := range playerTankBullets {
			for i := 0; i < len(enemyTanks); i++ {
				bulletNosePosition := sdl.FPoint{
					playerTankBullets[index].boundingBox.X + playerTankBullets[index].boundingBox.W,
					playerTankBullets[index].boundingBox.Y + (playerTankBullets[index].boundingBox.H / 2.0),
				}
				if bulletNosePosition.InRect(&enemyTanks[i].boundingBox) {
					explosions = append(explosions, NewExplosion(enemyTanks[i].boundingBox, explosionTexture))
					enemyTanks = RemoveElementFromEnemyTankSlice(enemyTanks, i)
					PlaySoundEffect(explosionSoundEffect)
				}
			}
		}

		//==============REMOVING DIED EXPLOSION ANIMATIONS==============
		for i := 0; i < len(explosions); i++ {
			if explosions[i].died {
				explosions = RemoveElementFromExplosionSlice(explosions, i)
			}
		}

		//==============UPDATING EXPLOSION ANIMATIONS==============
		for index := range explosions {
			explosions[index].Update()
		}

		// TODO : Update first, or draw first?(most probably update first) + update order
		//		playerTank.Update()

		//==============CLEARING THE SCREEN==============
		renderer.SetDrawColor(colornames.Bisque.R, colornames.Bisque.G, colornames.Bisque.B, colornames.Bisque.A)
		renderer.Clear()
		renderer.SetDrawColor(colornames.Red.R, colornames.Red.G, colornames.Red.B, colornames.Red.A)

		//==============DRAWING==============
		for index := range explosions {
			explosions[index].Draw(renderer)
		}
		DrawTexture(renderer, playerTank.tankTexture, &playerTank.boundingBox, playerTank.rotationAngle)
		for index := range playerTankBullets {
			DrawTexture(renderer, playerTankBullets[index].bulletTexture, &playerTankBullets[index].boundingBox, playerTankBullets[index].rotationAngle)
		}
		for index := range enemyTanks {
			DrawTexture(renderer, enemyTanks[index].tankTexture, &enemyTanks[index].boundingBox, enemyTanks[index].rotationAngle)
		}
		for index := range enemyTankBullets {
			DrawTexture(renderer, enemyTankBullets[index].bulletTexture, &enemyTankBullets[index].boundingBox, enemyTankBullets[index].rotationAngle)
		}
		renderer.Present()

		//==============UPDATING FPS COUNTER==============
		fpsCounter += 1
	}

	//sdl.Quit()
	return 0
}

func main() {
	os.Exit(run())
}
