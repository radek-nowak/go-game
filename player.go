package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var SpaceShip = mustLoadImage("assets/ship_B.png")

type Player struct {
	position Vector
	velocity Vector
	rotation Rotation
	sprite   *ebiten.Image
}

func NewPlayer() *Player {
	sprite := SpaceShip

	bounds := sprite.Bounds()

	hw := float64(bounds.Dx() / 2)
	hh := float64(bounds.Dy() / 2)

	return &Player{
		position: *NewVector(ScreenWidth/2-hw, ScreenHeight/2-hh),
		velocity: *NewVector(0, 0),
		rotation: Rotation{},
		sprite:   sprite,
	}

}

func (player *Player) Update() error {

	// speed := 1.0
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

		// r := math.Sqrt(math.Pow(player.position.X(), 2) + math.Pow(player.position.Y(), 2))
		// x := r * math.Cos((270.0+player.rotation.R)*math.Pi/180.0)
		// y := r * math.Sin((270.0+player.rotation.R)*math.Pi/180.0)

		acceleration := 0.01
		player.velocity.Add(xk*acceleration, yk*acceleration)
	}

	// r := math.Sqrt(math.Pow(player.position.X(), 2) + math.Pow(player.position.Y(), 2))
	// x := r * xk
	// y := r * yk

	// player.position.Add(x*speed, y*speed)
	// TODO: top velocity
	// TODO: maybe increase acceleration in the early phase of the movement, and then decrease
	player.position.Add(player.velocity.x, player.velocity.y)
	// __AUTO_GENERATED_PRINT_VAR_START__
	fmt.Println(fmt.Sprintf("Update velocity: %v", player.velocity)) // __AUTO_GENERATED_PRINT_VAR_END__

	return nil
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
