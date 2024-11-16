package render

import (
    "gonum.org/v1/gonum/spatial/r3"
)

type Intersectable interface {
    Intersect(origin, direction r3.Vec) (float64, Color, r3.Vec)
}

type Scene struct {
    Camera  r3.Vec
    Objects []Intersectable
}

type Light struct {
    Direction r3.Vec
    Intensity float64
    Ambient   float64
}

func (s *Scene) RayTrace(origin, direction r3.Vec, light *Light) Color {
    minDist := -1.0
    var hitColor Color
    var hitNormal r3.Vec
    
    for _, obj := range s.Objects {
        dist, color, normal := obj.Intersect(origin, direction)
        if dist > 0 && (minDist < 0 || dist < minDist) {
            minDist = dist
            hitColor = color
            hitNormal = normal
        }
    }
    
    if minDist < 0 {
        return Color{0, 0, 0}
    }
    
    diffuse := r3.Dot(hitNormal, light.Direction)
    if diffuse < 0 {
        diffuse = 0
    }
    
    lighting := light.Ambient + diffuse*light.Intensity
    return Color{
        R: hitColor.R * lighting,
        G: hitColor.G * lighting,
        B: hitColor.B * lighting,
    }
}

// RenderScene renders the scene to a pixel buffer using the provided callback
func RenderScene(scene *Scene, light *Light, width, height int, aspectRatio float64, pixelCallback func(x, y int, color Color)) {
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            // Convert pixel coordinates to normalized device coordinates (-1 to 1)
            screenX := (2.0*float64(x)/float64(width) - 1.0) * aspectRatio
            screenY := 1.0 - 2.0*float64(y)/float64(height)
            
            // Create ray direction
            rayDir := r3.Unit(r3.Vec{screenX, screenY, 1.0})
            
            // Trace ray and get color
            color := scene.RayTrace(scene.Camera, rayDir, light)
            
            // Send pixel to callback
            pixelCallback(x, y, color)
        }
    }
}
