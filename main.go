package main

import (
	"embed"
	"fmt"
	"image/color"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

//go:embed assets/*
var assets embed.FS

const (
	ScreenWidth         = 1600
	ScreenHeight        = 900
	SmallMeteorHitScore = 50
	BigMeteorHitScore   = 100
)


var ScoreFont = mustLoadFont("assets/Kenney Mini.ttf")


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

	// the game ends if player has 0 lives
	if !game.player.HasRemainingLives() {
		game.player = NewPlayer()
		game.meteors = nil
		game.player.bulletManager.bullets = nil
		game.score = 0
	}

	// spawn meteors
	game.meteorSpawnTimer.Update()
	if game.meteorSpawnTimer.IsReady() {
		game.meteorSpawnTimer.Reset()

		meteor := NewMeteor()
		game.meteors = append(game.meteors, meteor)
	}

	// update meteors and bullets state
	for _, m := range game.meteors {
		m.Update()
	}

	for _, b := range game.player.bulletManager.bullets {
		b.Update()
	}

	// check collisions
	for i, m := range game.meteors {
		if m.CollisionRect().Intersects(game.player.CollisionRect()) {
			game.player.Reset()
			game.meteors = nil
			game.player.bulletManager.bullets = nil
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

	// draw game score
	text.Draw(screen, fmt.Sprintf("%06d", game.score), ScoreFont, ScreenWidth/2, 50, color.White)

	// draw remaining lives
	spaceShip := mustLoadImage("assets/ship_B.png")

	for i := 0; i < game.player.lives; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(10*float64(i), 0)
		screen.DrawImage(spaceShip, op)
	}
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
