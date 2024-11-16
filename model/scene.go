package model

import (
    "math"
    "gonum.org/v1/gonum/spatial/r3"
    "fishies/render"
)

func CreateScene() (*render.Scene, *render.Light) {
    spheres := []render.Intersectable{
        &HoledSphere{
            Center: r3.Vec{-8, 5, 0},
            InitialPos: r3.Vec{-8, 5, 0},
            Velocity: r3.Vec{0, 0, 0},
            Color:  render.Red,
            Radius: 3.0,
            EyeRadius: 0.4,
            EyeDistance: 0.6,
        },
        &HoledSphere{
            Center: r3.Vec{0, 3, 0},
            InitialPos: r3.Vec{0, 3, 0},
            Velocity: r3.Vec{0, 0, 0},
            Color:  render.Cyan,
            Radius: 5.0,
            EyeRadius: 1.1,
            EyeDistance: 1.3,
        },
        &HoledSphere{
            Center: r3.Vec{8, 4, 0},
            InitialPos: r3.Vec{8, 4, 0},
            Velocity: r3.Vec{0, 0, 0},
            Color:  render.Yellow,
            Radius: 2.0,
            EyeRadius: 0.3,
            EyeDistance: 0.6,
        },
    }
    
    light := &render.Light{
        Direction: r3.Unit(r3.Vec{-1, -1, -1}),
        Intensity: 1.0,
        Ambient:   0.2,
    }
    
    return &render.Scene{
        Camera: r3.Vec{0, 0, -10},
        Objects: spheres,
    }, light
}

var totalElapsedTime float64 = 0.0

func UpdateScene(scene *render.Scene, deltaTime float64) {
    const (
        GRAVITY float64 = -9.81
        FLOOR_Y float64 = -5.0
    )

    totalElapsedTime += deltaTime

    // Camera movement using accumulated time
    scene.Camera = r3.Vec{
        3 * math.Sin(totalElapsedTime)*0.5,
        1 * math.Sin(float64(totalElapsedTime)*0.25),
        -10 + 2*math.Cos(totalElapsedTime)*0.5,
    }

    // Update spheres
    for _, obj := range scene.Objects {
        if sphere, ok := obj.(*HoledSphere); ok {
            // Doubled rotation speed
            sphere.Rotation += deltaTime * 0.8

            // Apply gravity using actual deltaTime
            sphere.Velocity.Y += deltaTime * GRAVITY

            // Update position using actual deltaTime
            newPos := r3.Add(sphere.Center, r3.Scale(deltaTime, sphere.Velocity))
            sphere.Center = newPos

            // Floor collision with perfect elastic reflection
            if sphere.Center.Y-sphere.Radius < FLOOR_Y {
                // Set position exactly at floor level to prevent any loss of energy
                sphere.Center.Y = FLOOR_Y + sphere.Radius
                
                // Perfect elastic reflection
                sphere.Velocity.Y = math.Abs(sphere.Velocity.Y)  // Force positive velocity
            }
        }
    }
}