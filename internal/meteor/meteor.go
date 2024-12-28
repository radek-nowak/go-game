package meteor

import (
	"game/config"
	"game/internal/collision"
	"game/internal/vector"
	"game/utils"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var BigMeteorSprite = utils.MustLoadImage("assets/meteor_detailedLarge.png")
var SmallMeteorSprite = utils.MustLoadImage("assets/meteor_detailedSmall.png")

type Type int

const (
	big Type = iota
	small
)

type Meteor struct {
	position   vector.Vector
	velocity   vector.Vector
	sprite     *ebiten.Image
	angle      float64 // TODO: is it needed?
	meteorType Type
}

func NewMeteor() *Meteor {
	sprite := BigMeteorSprite

	// Center of the screen
	target := vector.NewVector(config.ScreenWidth/2, config.ScreenHeight/2)

	// Spawn point on the edge of the screen
	r := config.ScreenWidth / 2.0
	angle := rand.Float64() * 2 * math.Pi
	spawnPoint := vector.NewVector(
		target.X()+r*math.Cos(angle),
		target.Y()+r*math.Sin(angle),
	)

	// Velocity direction toward the center
	direction := vector.NewVector(target.X()-spawnPoint.X(), target.Y()-spawnPoint.Y()).Normalize()

	// Set velocity with some random speed factor
	speed := 0.9 + rand.Float64() // TODO: increase speed and spawn freq by time //+ float64(time.Now().Second())/
	velocity := vector.NewVector(direction.X()*speed, direction.Y()*speed)

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

	vel1 := vector.NewVector(originalMeteor.velocity.X(), originalMeteor.velocity.Y()).Rotate(newAngle)
	vel2 := vector.NewVector(originalMeteor.velocity.X(), originalMeteor.velocity.Y()).Rotate(-newAngle)

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

func (meteor *Meteor) IsBig() bool {
	return meteor.meteorType == big
}

func (me *Meteor) CollisionRect() collision.Rect {
	bounds := me.sprite.Bounds()

	// Apply a margin to shrink the collision rectangle
	margin := 10.0 // Adjust as needed
	return collision.NewRect(
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
