# tanks
A simple tanks war game, written in [GO](https://golang.org/).

> - I got the images from here:
[tank image](https://opengameart.org/content/pixel-tank)
[bullet imge](https://opengameart.org/content/bullet-symbol)
> - (I forgot from where I got the sound effects, although I remember I used [this](https://opengameart.org) site to get the sound effects)
> - The resources I used in this project are not original, I have edited some of them
> - I used the [go-sdl2](https://godoc.org/github.com/veandco/go-sdl2) library here.

Thanks in advance for contributing to my project :relaxed:

## The game:
There will be a player tank(the green tank), and lots of enemy tanks(the red tanks).
The enemy tanks will either try to shoot the player or else shoot at a random direction.
The player will win, if it kills all the enemy tanks by shooting them.
If any of the enemy tanks shoot and kills the player tank, the player loses.
At first there will be a minumum number of enemy tanks, which will increase slowly...

## How to run:
Grab the latest stable compiled binaries [here](https://github.com/dev-abir/tanks/releases/latest)(scroll down, and check the **Assets**)

## Controls:
Press `w` to move the player tank(the green tank) forward(or up).
Press `a` to move the player tank(the green tank) left.
Press `s` to move the player tank(the green tank) backward(or down).
Press `d` to move the player tank(the green tank) right.

Press `LEFT ARROW` to rotate the player tank(the green tank) anti-clockwise.
Press `LEFT ARROW` to rotate the player tank(the green tank) clockwise.

Press `SPACE` to shoot.

## How to build:
**I would highly encourage you to understand each and every steps of the build process**
---

### On GNU/Linux:
1. Get(and install) the go compiler from golang.org.(I recommend you to use the compiler from golang.org, not the package manager's one, and I have not tested the gccgo compile, there are thousands of tutorials available in the internet, for doing this)
2. Ensure that you have [these](https://github.com/veandco/go-sdl2#requirements) requirements.
3. Get the zip of my project, extract it anywhere in your pc.
4. Go to the directory where you have extracted it, and run `go build -o tanks`.
5. To run the exectutable, run `./tanks`.

### On Windows:
Almost same as of GNU/Linux, excpet, for the last step, you should use `.\tanks`, and for the first step, you have only one option, i.e., to get your compiler from golang.org.

### On macOS:
I don't use it, BTW, you may contribute the steps....

## Build release:
- If you want to build release binaries yourself, you may use the `release_script.sh` script to do that.
- I do the development in a Linux distro, so to run this script, you need to be on a Linux distro.
- I cross compile to get builds for other OS's. Make sure you have met [these](https://github.com/veandco/go-sdl2#cross-compiling) dependencies.
- Run `sh release_script.sh` in the root directory of my project, after completition of this command, you will see the release builds, inside the `release` directory.
- You may also read [this](https://github.com/veandco/go-sdl2#static-compilation) to understand my release script.
