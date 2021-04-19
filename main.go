package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten"
)

type button struct{ down, pressed, released bool }

func (b *button) update(down bool) {
	b.pressed = !b.down && down
	b.released = b.down && !down
	b.down = down
}

type game struct {
	ticks uint64

	input struct {
		leftButton  button
		rightButton button
		upButton    button
		downButton  button
	}

	player struct {
		crouching bool

		walking      bool
		walkingFrame int
	}

	// assets go here for now
	bgColor         color.Color
	marioStanding   *ebiten.Image
	marioCrouching  *ebiten.Image
	marioWalkFrames []*ebiten.Image
}

var palette = []color.Color{
	color.RGBA{255, 0, 0, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{0, 0, 255, 255},
	color.RGBA{255, 255, 0, 255},
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(&color.RGBA{
		R: 147,
		G: 187,
		B: 236,
		A: 255,
	})
	if g.player.crouching {
		screen.DrawImage(g.marioCrouching, &ebiten.DrawImageOptions{})
	} else if g.player.walking {
		screen.DrawImage(g.marioWalkFrames[g.player.walkingFrame], &ebiten.DrawImageOptions{})
	} else {
		screen.DrawImage(g.marioStanding, &ebiten.DrawImageOptions{})
	}
}

func (g *game) Update(screen *ebiten.Image) error {
	g.ticks++

	g.input.upButton.update(ebiten.IsKeyPressed(ebiten.KeyUp))
	g.input.downButton.update(ebiten.IsKeyPressed(ebiten.KeyDown))
	g.input.leftButton.update(ebiten.IsKeyPressed(ebiten.KeyLeft))
	g.input.rightButton.update(ebiten.IsKeyPressed(ebiten.KeyRight))

	if g.input.downButton.down {
		g.player.crouching = true
	} else {
		g.player.crouching = false
	}
	if g.input.rightButton.down {
		g.player.walking = true
		g.player.walkingFrame = int(g.ticks % 30 / 10)
	} else {
		g.player.walking = false
		g.player.walkingFrame = 0
	}

	return nil
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 240
}

func main() {
	assets, err := ioutil.ReadFile("game.dat")
	if err != nil {
		log.Fatal(err)
	}

	myGame := game{
		bgColor:         color.Black,
		marioWalkFrames: make([]*ebiten.Image, 3),
	}

	skip := 4 * 16 * 32
	myGame.marioStanding, _ = ebiten.NewImageFromImage(&image.RGBA{
		Pix:    assets[0*skip : 1*skip],
		Stride: 4 * 16,
		Rect:   image.Rect(0, 0, 16, 32),
	}, ebiten.FilterDefault)
	myGame.marioCrouching, _ = ebiten.NewImageFromImage(&image.RGBA{
		Pix:    assets[1*skip : 2*skip],
		Stride: 4 * 16,
		Rect:   image.Rect(0, 0, 16, 32),
	}, ebiten.FilterDefault)
	myGame.marioWalkFrames[0], _ = ebiten.NewImageFromImage(&image.RGBA{
		Pix:    assets[2*skip : 3*skip],
		Stride: 4 * 16,
		Rect:   image.Rect(0, 0, 16, 32),
	}, ebiten.FilterDefault)
	myGame.marioWalkFrames[1], _ = ebiten.NewImageFromImage(&image.RGBA{
		Pix:    assets[3*skip : 4*skip],
		Stride: 4 * 16,
		Rect:   image.Rect(0, 0, 16, 32),
	}, ebiten.FilterDefault)
	myGame.marioWalkFrames[2], _ = ebiten.NewImageFromImage(&image.RGBA{
		Pix:    assets[4*skip : 5*skip],
		Stride: 4 * 16,
		Rect:   image.Rect(0, 0, 16, 32),
	}, ebiten.FilterDefault)

	if err := ebiten.RunGame(&myGame); err != nil {
		panic(err)
	}
}
