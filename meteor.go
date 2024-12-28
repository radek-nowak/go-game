package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var BigMeteorSprite = mustLoadImage("assets/meteor_detailedLarge.png")
var SmallMeteorSprite = mustLoadImage("assets/meteor_detailedSmall.png")

type Type int

const (
	big Type = iota
	small
)

type Meteor struct {
	position   Vector
	velocity   Vector
	sprite     *ebiten.Image
	angle      float64 // TODO: is it needed?
	meteorType Type
}

func NewMeteor() *Meteor {
	sprite := BigMeteorSprite

	// Center of the screen
	target := NewVector(ScreenWidth/2, ScreenHeight/2)

	// Spawn point on the edge of the screen
	r := ScreenWidth / 2.0
	angle := rand.Float64() * 2 * math.Pi
	spawnPoint := NewVector(
		target.X()+r*math.Cos(angle),
		target.Y()+r*math.Sin(angle),
	)

	// Velocity direction toward the center
	direction := NewVector(target.X()-spawnPoint.X(), target.Y()-spawnPoint.Y()).Normalize()

	// Set velocity with some random speed factor
	speed := 0.9 + rand.Float64() // TODO: increase speed and spawn freq by time //+ float64(time.Now().Second())/
	velocity := NewVector(direction.X()*speed, direction.Y()*speed)

	return &Meteor{
		position:   *spawnPoint,
		velocity:   *velocity,
		sprite:     sprite,
		angle:      angle,
		meteorType: big,
	}
}

func NewSmallMeteors(originalMeteor *Meteor) []*Meteor {
	sprite := SmallMeteorSprite

	angle := originalMeteor.angle

	minAngle, maxAngle := 10.0, 45.0

	newAngle := minAngle + rand.Float64()*(maxAngle-minAngle)

	vel1 := NewVector(originalMeteor.velocity.X(), originalMeteor.velocity.Y()).Rotate(newAngle)
	vel2 := NewVector(originalMeteor.velocity.X(), originalMeteor.velocity.Y()).Rotate(-newAngle)

	return []*Meteor{
		{
			position:   originalMeteor.position,
			velocity:   *vel1,
			sprite:     sprite,
			angle:      angle,
			meteorType: small,
		},
		{
			position:   originalMeteor.position,
			velocity:   *vel2,
			sprite:     sprite,
			angle:      angle,
			meteorType: small,
		},
	}

}

func (me *Meteor) CollisionRect() Rect {
	bounds := me.sprite.Bounds()

	// Apply a margin to shrink the collision rectangle
	margin := 10.0 // Adjust as needed
	return NewRect(
		me.position.X()+margin,
		me.position.Y()+margin,
		float64(bounds.Dx())-2*margin,
		float64(bounds.Dy())-2*margin,
	)
}

func (meteor *Meteor) Update() error {
	// Update the position based on velocity
	meteor.position.Add(meteor.velocity.X(), meteor.velocity.Y())
	return nil
}

func (me *Meteor) Draw(screen *ebiten.Image) {
	bounds := me.sprite.Bounds()
	op := &ebiten.DrawImageOptions{}

	hw := float64(bounds.Dx() / 2)
	hh := float64(bounds.Dy() / 2)

	op.GeoM.Translate(-hw, -hh)
	op.GeoM.Rotate(math.Pi / 180.0)
	op.GeoM.Translate(hw, hh)
	op.GeoM.Translate(me.position.X(), me.position.Y())

	screen.DrawImage(me.sprite, op)
}
