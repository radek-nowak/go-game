package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/*
var assets embed.FS

const (
	ScreenWidth         = 1600
	ScreenHeight        = 900
	SmallMeteorHitScore = 50
	BigMeteorHitScore   = 100
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

var ScoreFont = mustLoadFont("assets/Kenney Mini.ttf")

func mustLoadFont(name string) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
}

type Rotation struct {
	R float64
}

type GameActor interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Game struct {
	player           *Player
	meteors          []*Meteor
	bulletTimer      *Timer
	meteorSpawnTimer *Timer
	score            int
}

func (game *Game) Update() error {
	game.player.Update()

	game.meteorSpawnTimer.Update()
	if game.meteorSpawnTimer.IsReady() {
		game.meteorSpawnTimer.Reset()

		meteor := NewMeteor()
		game.meteors = append(game.meteors, meteor)
	}

	for _, m := range game.meteors {
		m.Update()
	}

	for _, b := range game.player.bulletManager.bullets {
		b.Update()
	}

	for i, m := range game.meteors {
		if m.CollisionRect().Intersects(game.player.CollisionRect()) {
			game.player = NewPlayer()
			game.meteors = nil
			game.player.bulletManager.bullets = nil
			game.score = 0
		}
		meteorShotDown := game.player.bulletManager.CheckCollisionsWithMeteor(m)
		if meteorShotDown {
			game.meteors = append(game.meteors[:i], game.meteors[i+1:]...)
			if m.meteorType == big {
				game.score += BigMeteorHitScore
				newMeteors := NewSmallMeteors(m)
				game.meteors = append(game.meteors, newMeteors...)
			} else {
				game.score += SmallMeteorHitScore
			}
		}
	}

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.player.Draw(screen)

	for _, m := range game.meteors {
		m.Draw(screen)
	}

	for _, b := range game.player.bulletManager.bullets {
		b.Draw(screen)
	}

	text.Draw(screen, fmt.Sprintf("%06d", game.score), ScoreFont, ScreenWidth/2, 50, color.White)
}

func (ga *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth - 100, outsideHeight - 100
}

func main() {

	g := &Game{
		player:           NewPlayer(),
		meteors:          []*Meteor{},
		bulletTimer:      NewTimer(5000 * time.Millisecond),
		meteorSpawnTimer: NewTimer(5 * time.Second),
		score:            0,
	}

	err := ebiten.RunGame(g)

	if err != nil {
		panic(err)
	}

}
