// model/scene.go
package model

import (
	"fishies/render"
	"math/rand"

	"gonum.org/v1/gonum/spatial/r3"
)

type SceneObject interface {
	render.Intersectable
	Update(deltaTime float64)
}

// Random color generator
func randomColor() render.Color {
	return render.Color{
		R: 0.3 + rand.Float64()*0.7, // Avoiding too dark colors
		G: 0.3 + rand.Float64()*0.7,
		B: 0.3 + rand.Float64()*0.7,
	}
}

func CreateScene(numFish int, ground bool) (*render.Scene, *render.Light) {
	size := numFish
	if ground {
		size++
	}
	objects := make([]render.Intersectable, size)

	// Create fish at different heights with variation
	baseHeight := 2.0
	heightStep := 2.5

	for i := 0; i < numFish; i++ {
		pos := r3.Vec{
			X: rand.Float64()*4.0 - 2.0,                                      // Random initial X position between -2 and 2
			Y: baseHeight + float64(i)*heightStep + rand.Float64()*1.0 - 0.5, // Add random variation ±0.5
			Z: rand.Float64()*4.0 - 2.0,                                      // Random initial Z position between -2 and 2
		}

		// Random size between 3.0 and 5.0 (doubled base size and increased variation)
		size := 3.0 + rand.Float64()*2.0

		objects[i] = NewFish(pos, randomColor(), size)
	}

	// Add ground if enabled
	if ground {
		objects[numFish] = NewGround(0.0)
	}

	light := &render.Light{
		Direction: r3.Unit(r3.Vec{-0.5, 2, -0.5}),
		Intensity: 0.9,
		Ambient:   0.3,
	}

	cameraY := 7.0
	cameraZ := -16.0

	return &render.Scene{
		Camera:  r3.Vec{0, cameraY, cameraZ},
		Objects: objects,
	}, light
}

func UpdateScene(scene *render.Scene, deltaTime float64) {
	for _, obj := range scene.Objects {
		if sceneObj, ok := obj.(SceneObject); ok {
			sceneObj.Update(deltaTime)
		}
	}
}
