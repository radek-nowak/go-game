package main

import (
	"fmt"
	"game/config"
	"game/internal/meteor"
	"game/internal/player"
	"game/internal/timer"
	"game/utils"
	"image/color"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	SmallMeteorHitScore = 50
	BigMeteorHitScore   = 100
)

var ScoreFont = utils.MustLoadFont("assets/Kenney Mini.ttf")

type GameActor interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Game struct {
	player           *player.Player
	meteors          []*meteor.Meteor
	bulletTimer      *timer.Timer
	meteorSpawnTimer *timer.Timer
	score            int
}

func (game *Game) Update() error {
	game.player.Update()

	// the game ends if player has 0 lives
	if !game.player.HasRemainingLives() {
		game.player = player.NewPlayer()
		game.meteors = nil
		game.player.BulletManager().Reset()
		game.score = 0
	}

	// spawn meteors
	game.meteorSpawnTimer.Update()
	if game.meteorSpawnTimer.IsReady() {
		game.meteorSpawnTimer.Reset()

		meteor := meteor.NewMeteor()
		game.meteors = append(game.meteors, meteor)
	}

	// update meteors and bullets state
	for _, m := range game.meteors {
		m.Update()
	}

	game.player.BulletManager().UpdateBullets()

	// check collisions
	for i, m := range game.meteors {
		if m.CollisionRect().Intersects(game.player.CollisionRect()) {
			game.player.Reset()
			game.meteors = nil
			game.player.BulletManager().Reset()
		}
		meteorShotDown := game.player.BulletManager().CheckCollisionsWithMeteor(m)
		if meteorShotDown {
			game.meteors = append(game.meteors[:i], game.meteors[i+1:]...)
			if m.IsBig() {
				game.score += BigMeteorHitScore
				newMeteors := meteor.NewSmallMeteors(m)
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

	game.player.BulletManager().DrawBullets(screen)

	// draw game score
	text.Draw(screen, fmt.Sprintf("%06d", game.score), ScoreFont, config.ScreenWidth/2, 50, color.White)

	// draw remaining lives
	spaceShip := utils.MustLoadImage("assets/ship_B.png")

	for i := 0; i < game.player.RemainingLives(); i++ {
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
		player:           player.NewPlayer(),
		meteors:          []*meteor.Meteor{},
		bulletTimer:      timer.NewTimer(5000 * time.Millisecond),
		meteorSpawnTimer: timer.NewTimer(2 * time.Second),
		score:            0,
	}

	err := ebiten.RunGame(g)

	if err != nil {
		panic(err)
	}

}
