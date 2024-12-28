package collision

type Collider interface {
	CollisionRect() Rect
}

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func NewRect(x, y, w, h float64) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
	}
}

func (r Rect) MaxX() float64 {
	return r.Width + r.X
}

func (r Rect) MaxY() float64 {
	return r.Height + r.Y
}

func (r Rect) Intersects(other Rect) bool {
	return r.X <= other.MaxX() &&
		r.Y <= other.MaxY() &&
		other.X <= r.MaxX() &&
		other.Y <= r.MaxY()
}
