package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img *ebiten.Image
	X, Y float64
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Cauldron struct {
	*Sprite
}

type Game struct{
	player *Sprite
	cauldron *Cauldron
	enemies []*Enemy
	tilemapJSON *TilemapJSON
	tilemapImg *ebiten.Image
}

// Update is fixed by default to 60 ticks per second
func (g *Game) Update() error {

	// React to key presses
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}

	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			if sprite.X < g.player.X {
				sprite.X += 1
			} else if sprite.X > g.player.X {
				sprite.X -= 1
			}
	
			if sprite.Y < g.player.Y {
				sprite.Y += 1
			} else if sprite.Y > g.player.Y {
				sprite.Y -= 1
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// fill the screen with a sky color
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	// loop over layers
	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			srcX := (id-1) % 40
			srcY := (id-1) / 40

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(
					image.Rect(srcX, srcY, srcX + 16, srcY + 16),
				).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	// set the translation of our drawImageOptions to the player's position
	opts.GeoM.Translate(g.player.X, g.player.Y)

	// draw player
	screen.DrawImage(
		// grab a subimage of the spritesheet
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 32),
		).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()

	// draw the cauldron
	opts.GeoM.Translate(g.cauldron.X, g.cauldron.Y)
	screen.DrawImage(
		g.cauldron.Img.SubImage(
			image.Rect(26*8, 0, 28*8, 16),
		).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()

	// draw all of the sprites from the sprite array
	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 32),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 176
	// return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Legend of Zelda")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// load the player image from file
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/character.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// load the enemy image from file
	enemyImg, _, err := ebitenutil.NewImageFromFile("assets/images/NPC_test.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// load the cauldron image from file
	cauldronImg, _, err := ebitenutil.NewImageFromFile("assets/images/objects.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// load the cauldron image from file
	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/Overworld.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapJSON, err := NewTileMapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &Sprite{
			Img: playerImg,
			X: 50.0,
			Y: 50.0,
		},
		cauldron: &Cauldron{
			&Sprite{
				Img: cauldronImg,
				X: 200,
				Y: 200,
			},
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Img: enemyImg,
					X: 175,
					Y: 175,
				},
				false,
			},
			{
				&Sprite{
					Img: enemyImg,
					X: 275,
					Y: 275,
				},
				true,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg: tilemapImg,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}