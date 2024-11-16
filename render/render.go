package render

import (
    "math"
    "gonum.org/v1/gonum/spatial/r3"
)

type Scene struct {
    Camera  r3.Vec
    Objects []Intersectable
}

type Light struct {
    Direction r3.Vec  // Light direction (should be normalized)
    Intensity float64 // Light intensity
    Ambient   float64 // Ambient light level
}

type Intersectable interface {
    Intersect(origin, direction r3.Vec) (float64, Color, r3.Vec)
}

type RenderCallback func(x, y int, color Color)

func RenderScene(scene *Scene, light *Light, width, height int, aspectRatio float64, callback RenderCallback) {
    for j := 0; j < height; j++ {
        for i := 0; i < width; i++ {
            x := (2.0*(float64(i)+0.5)/float64(width) - 1.0) * aspectRatio
            y := -(2.0*(float64(j)+0.5)/float64(height) - 1.0)
            direction := r3.Unit(r3.Vec{x, y, 1})
            
            var closestT float64 = -1.0
            var closestColor Color = Black
            var closestNormal r3.Vec
            
            for _, obj := range scene.Objects {
                t, color, normal := obj.Intersect(scene.Camera, direction)
                if t > 0 && (closestT < 0 || t < closestT) {
                    closestT = t
                    closestColor = color
                    closestNormal = normal
                }
            }
            
            if closestT > 0 && closestColor != Black {
                // Calculate lighting
                diffuse := math.Max(0, r3.Dot(closestNormal, light.Direction))
                lighting := light.Ambient + (1.0-light.Ambient)*diffuse*light.Intensity
                
                // Apply lighting to color
                closestColor = Color{
                    R: closestColor.R * lighting,
                    G: closestColor.G * lighting,
                    B: closestColor.B * lighting,
                }
            }
            
            callback(i, j, closestColor)
        }
    }
}