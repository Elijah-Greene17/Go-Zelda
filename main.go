package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{
	PlayerImage *ebiten.Image
	X			 float64
	Y 			 float64
}

// Update is fixed by default to 60 ticks per second
func (g *Game) Update() error {

	// React to key presses
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Y += 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.X, g.Y)

	// draw player
	screen.DrawImage(
		g.PlayerImage.SubImage(
			image.Rect(0, 0, 16, 32),
		).(*ebiten.Image),
		&opts,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Legend of Zelda")
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/character.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{PlayerImage: playerImg, X: 100, Y: 100}); err != nil {
		log.Fatal(err)
	}
}