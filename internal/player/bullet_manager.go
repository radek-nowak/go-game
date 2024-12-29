package player

import (
	"game/internal/meteor"

	"github.com/hajimehoshi/ebiten/v2"
)

type BulletManager struct {
	bullets []*Bullet
}

func NewBulletManager() *BulletManager {
	return &BulletManager{
		bullets: []*Bullet{},
	}
}

func (bm *BulletManager) Reset() {
	bm.bullets = nil
}

func (bm *BulletManager) RegisterNewBullet(bullet *Bullet) {
	bm.bullets = append(bm.bullets, bullet)
}

func (bm *BulletManager) UpdateBullets() error {
	for _, b := range bm.bullets {
		b.Update()
	}

	return nil
}

func (bm *BulletManager) DrawBullets(screen *ebiten.Image) {
	for _, b := range bm.bullets {
		b.Draw(screen)
	}
}

func (bm *BulletManager) CheckCollisionsWithMeteor(meteor *meteor.Meteor) bool {
	hit := false
	for i, bullet := range bm.bullets {
		if bullet.CollisionRect().Intersects(meteor.CollisionRect()) {
			hit = true
			bm.bullets = append(bm.bullets[:i], bm.bullets[i+1:]...)
		}
	}

	return hit
}
