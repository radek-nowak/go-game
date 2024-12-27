package main

type BulletManager struct {
	bullets []*Bullet
}

func NewBulletManager() *BulletManager {
	return &BulletManager{
		bullets: []*Bullet{},
	}
}

func (bm *BulletManager) RegisterNewBullet(bullet *Bullet) {
	bm.bullets = append(bm.bullets, bullet)
}

func (bm *BulletManager) CheckCollisionsWithMeteor(meteor *Meteor) bool {
	hit := false
	for i, bullet := range bm.bullets {
		if bullet.CollisionRect().Intersects(meteor.CollisionRect()) {
			hit = true
			bm.bullets = append(bm.bullets[:i], bm.bullets[i+1:]...)
		}
	}

	return hit
}
