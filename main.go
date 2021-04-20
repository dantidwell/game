package main

import (
	"image/color"
	"log"
	"os"

	"github.com/dantidwell/game/assets"
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

	assets *assets.Pack

	input struct {
		leftButton  button
		rightButton button
		upButton    button
		downButton  button
	}

	player struct {
		dirX float64
		dirY float64

		posX float64
		posY float64
	}

	// assets go here for now
	bgColor color.Color
}

var palette = []color.Color{
	color.RGBA{255, 0, 0, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{0, 0, 255, 255},
	color.RGBA{255, 255, 0, 255},
}

func (g *game) Draw(screen *ebiten.Image) {
	// draw floor ...
	var opts ebiten.DrawImageOptions

	for y := 0; y < 240; y += 16 {
		for x := 0; x < 256; x += 16 {
			opts.GeoM.Reset()
			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(g.assets.GetImage("floor0"), &opts)
		}
	}

	// draw north wall ...
	for x := 0; x < 256; x += 16 {
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(x), 0)
		screen.DrawImage(g.assets.GetImage("wall"), &opts)
	}

	// draw border ...
	opts.GeoM.Reset()
	opts.GeoM.Translate(0, 0)
	screen.DrawImage(g.assets.GetImage("border_nw"), &opts)

	opts.GeoM.Reset()
	opts.GeoM.Translate(256-16, 0)
	screen.DrawImage(g.assets.GetImage("border_ne"), &opts)

	opts.GeoM.Reset()
	opts.GeoM.Translate(0, 240-16)
	screen.DrawImage(g.assets.GetImage("border_sw"), &opts)

	opts.GeoM.Reset()
	opts.GeoM.Translate(256-16, 240-16)
	screen.DrawImage(g.assets.GetImage("border_se"), &opts)

	for y := 16; y < 240-16; y += 16 {
		opts.GeoM.Reset()
		opts.GeoM.Translate(0, float64(y))
		screen.DrawImage(g.assets.GetImage("border_left"), &opts)

		opts.GeoM.Reset()
		opts.GeoM.Translate(256-16, float64(y))
		screen.DrawImage(g.assets.GetImage("border_right"), &opts)
	}
	for x := 16; x < 256-16; x += 16 {
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(x), 0)
		screen.DrawImage(g.assets.GetImage("border_top"), &opts)

		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(x), 240-16)
		screen.DrawImage(g.assets.GetImage("border_bottom"), &opts)
	}

	// draw the player
	opts.GeoM.Reset()
	opts.GeoM.Translate(g.player.posX, g.player.posY)
	if g.player.dirX == 0 && g.player.dirY == 1 {
		screen.DrawImage(g.assets.GetImage("dan_down"), &opts)
	} else if g.player.dirX == 0 && g.player.dirY == -1 {
		screen.DrawImage(g.assets.GetImage("dan_up"), &opts)
	} else if g.player.dirX == 1 && g.player.dirY == 0 {
		screen.DrawImage(g.assets.GetImage("dan_right"), &opts)
	} else if g.player.dirX == -1 && g.player.dirY == 0 {
		screen.DrawImage(g.assets.GetImage("dan_left"), &opts)
	}
}

func (g *game) Update(screen *ebiten.Image) error {
	g.ticks++

	g.input.upButton.update(ebiten.IsKeyPressed(ebiten.KeyUp))
	g.input.downButton.update(ebiten.IsKeyPressed(ebiten.KeyDown))
	g.input.leftButton.update(ebiten.IsKeyPressed(ebiten.KeyLeft))
	g.input.rightButton.update(ebiten.IsKeyPressed(ebiten.KeyRight))

	if g.input.downButton.pressed {
		g.player.posY += 16
		g.player.dirX, g.player.dirY = 0, 1
	} else if g.input.upButton.pressed {
		g.player.posY -= 16
		g.player.dirX, g.player.dirY = 0, -1
	} else if g.input.rightButton.pressed {
		g.player.posX += 16
		g.player.dirX, g.player.dirY = 1, 0
	} else if g.input.leftButton.pressed {
		g.player.posX -= 16
		g.player.dirX, g.player.dirY = -1, 0
	}
	return nil
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 240
}

func main() {
	f, err := os.Open("game.pak")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	myGame := game{
		assets:  assets.Load(f),
		bgColor: color.Black,
	}
	myGame.player.dirX = 1
	myGame.player.dirY = 0
	myGame.player.posX = 16
	myGame.player.posY = 16

	if err := ebiten.RunGame(&myGame); err != nil {
		panic(err)
	}
}
