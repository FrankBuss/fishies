package model

import (
	"fishies/render"
	"math"

	"gonum.org/v1/gonum/spatial/r3"
)

// Ground represents an infinite ground plane with a checkerboard pattern
type Ground struct {
	Y float64 // height of the ground plane
}

// NewGround creates a new ground plane at the specified Y height
func NewGround(height float64) *Ground {
	return &Ground{Y: height}
}

// Intersect calculates the intersection of a ray with the ground plane
func (g *Ground) Intersect(origin, direction r3.Vec) (float64, render.Color, r3.Vec) {
	if direction.Y == 0 {
		return -1, render.Color{0, 0, 0}, r3.Vec{}
	}

	t := (g.Y - origin.Y) / direction.Y
	if t < 0 {
		return -1, render.Color{0, 0, 0}, r3.Vec{}
	}

	hitPoint := r3.Add(origin, r3.Scale(t, direction))

	gridSize := 4.0
	x := math.Floor(hitPoint.X / gridSize)
	z := math.Floor(hitPoint.Z / gridSize)
	isEven := math.Mod(math.Abs(x+z), 2) < 1

	// Increased white brightness to compensate for lighting
	darkSquare := render.Color{0.0, 0.0, 0.0} // Pure black
	lightSquare := render.Color{.5, .5, .5}

	var color render.Color
	if isEven {
		color = darkSquare
	} else {
		color = lightSquare
	}

	return t, color, r3.Vec{0, 1, 0}
}

func (g *Ground) Update(deltaTime float64) {
	// Ground doesn't animate
}
