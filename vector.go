package main

type Vector struct {
	x float64
	y float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{
		x: x,
		y: y,
	}
}

func (v *Vector) Add(dx float64, dy float64) {
	v.x += dx
	v.y += dy
}

func (v Vector) X() float64 {
	return v.x
}

func (v Vector) Y() float64 {
	return v.y
}
