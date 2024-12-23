package main

import (
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

const (
	ScreenWidth  = 200
	ScreenHeight = 200
)

func mustLoadImage(name string) *ebiten.Image {
	file, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

type Rotation struct {
	R float64
}

type Game struct {
	player *Player
}

func (game *Game) Update() error {
	game.player.Update()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.player.Draw(screen)
}

func (ga *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth - 100, outsideHeight - 100
}

func main() {

	g := &Game{
		player: NewPlayer(),
	}

	err := ebiten.RunGame(g)

	if err != nil {
		panic(err)
	}

}
