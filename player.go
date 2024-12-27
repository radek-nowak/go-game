package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var SpaceShip = mustLoadImage("assets/ship_B.png")
var BulletSprite = mustLoadImage("assets/star_tiny.png")

const topVelocity = 1.0

type Player struct {
	position      Vector
	velocity      Vector
	rotation      Rotation
	sprite        *ebiten.Image
	bulletTimer   *Timer
	bulletManager *BulletManager
}

func NewPlayer() *Player {
	sprite := SpaceShip

	bounds := sprite.Bounds()

	hw := float64(bounds.Dx() / 2)
	hh := float64(bounds.Dy() / 2)

	return &Player{
		position:      *NewVector(ScreenWidth/2-hw, ScreenHeight/2-hh),
		velocity:      *NewVector(0, 0),
		rotation:      Rotation{},
		sprite:        sprite,
		bulletTimer:   NewTimer(200 * time.Millisecond),
		bulletManager: &BulletManager{},
	}
}

func (player *Player) Update() error {

	rotationSpeed := 2.0

	xk := math.Cos((270.0 + player.rotation.R) * math.Pi / 180.0)
	yk := math.Sin((270.0 + player.rotation.R) * math.Pi / 180.0)

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.rotation.R += rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.rotation.R -= rotationSpeed

	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		acceleration := 0.03
		player.velocity.Add(xk*acceleration, yk*acceleration)
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		player.shoot()
	}

	if player.position.X() > ScreenWidth {
		player.position = *NewVector(0, player.position.Y())
	}
	if player.position.X() < 0 {
		player.position = *NewVector(ScreenWidth, player.position.Y())
	}
	if player.position.Y() > ScreenHeight {
		player.position = *NewVector(player.position.X(), 0)
	}
	if player.position.Y() < 0 {
		player.position = *NewVector(player.position.X(), ScreenHeight)
	}

	// TODO: top velocity
	// TODO: increase acceleration in the early phase of the movement, and then decrease
	player.position.Add(player.velocity.x, player.velocity.y)

	return nil
}

func (player *Player) shoot() {
	player.bulletTimer.Update()
	if player.bulletTimer.IsReady() {
		player.bulletTimer.Reset()
		bullet := player.NewBullet()
		player.bulletManager.bullets = append(player.bulletManager.bullets, bullet)
	}
}

func (player *Player) Draw(screen *ebiten.Image) {
	bounds := player.sprite.Bounds()
	op := &ebiten.DrawImageOptions{}

	hw := float64(bounds.Dx() / 2)
	hh := float64(bounds.Dy() / 2)

	op.GeoM.Translate(-hw, -hh)
	op.GeoM.Rotate(player.rotation.R * math.Pi / 180.0)
	op.GeoM.Translate(hw, hh)
	op.GeoM.Translate(player.position.X(), player.position.Y())

	screen.DrawImage(player.sprite, op)
}

func (pl *Player) CollisionRect() Rect {
	bounds := pl.sprite.Bounds()

	// Apply a margin to shrink the collision rectangle
	margin := 10.0 // Adjust based on how much smaller you want the rectangle
	return NewRect(
		pl.position.X()+margin,
		pl.position.Y()+margin,
		float64(bounds.Dx())-2*margin,
		float64(bounds.Dy())-2*margin,
	)
}

type Bullet struct {
	position Vector
	velocity Vector
	// rotation Rotation// TODO: for the trajectory
	sprite *ebiten.Image
}

func (bu *Bullet) CollisionRect() Rect {
	bounds := bu.sprite.Bounds()

	// Apply a margin to shrink the collision rectangle
	margin := 10.0 // Adjust based on how much smaller you want the rectangle
	return NewRect(
		bu.position.X()+margin,
		bu.position.Y()+margin,
		float64(bounds.Dx())-2*margin,
		float64(bounds.Dy())-2*margin,
	)
}

func (bubullet *Bullet) Update() error {
	bubullet.position.Add(bubullet.velocity.x, bubullet.velocity.y)
	return nil
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// options.GeoM.Scale(0.5, 0.5)// TODO: scaling shifts bullets away from ship
	options.GeoM.Translate(bullet.position.X(), bullet.position.Y())
	screen.DrawImage(bullet.sprite, options)
}

func (p *Player) NewBullet() *Bullet {
	xk := math.Cos((270.0 + p.rotation.R) * math.Pi / 180.0)
	yk := math.Sin((270.0 + p.rotation.R) * math.Pi / 180.0)
	k := NewVector(xk*5, yk*5)
	k.Add(p.velocity.x, p.velocity.y)
	// v := p.velocity.Normalize()
	// v.Add(3, 3)
	return &Bullet{
		position: p.position,
		velocity: *k,
		sprite:   BulletSprite,
	}
}
