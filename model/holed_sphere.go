package model

import (
    "math"
    "gonum.org/v1/gonum/spatial/r3"
    "fishies/render"
)

type HoledSphere struct {
    Center      r3.Vec
    InitialPos  r3.Vec    
    Velocity    r3.Vec    
    Radius      float64
    EyeRadius   float64
    EyeDistance float64
    Rotation    float64
    Color       render.Color
}

func (s *HoledSphere) rotatePoint(p r3.Vec) r3.Vec {
    sin, cos := math.Sin(s.Rotation), math.Cos(s.Rotation)
    return r3.Vec{
        p.X*cos - p.Z*sin,
        p.Y,
        p.X*sin + p.Z*cos,
    }
}

func (s *HoledSphere) Intersect(origin, direction r3.Vec) (float64, render.Color, r3.Vec) {
    const (
        MAX_STEPS float64 = 100
        EPSILON   float64 = 0.001
        MAX_DIST  float64 = 20.0
    )
    
    t := 0.0
    
    for i := 0.0; i < MAX_STEPS; i++ {
        p := r3.Add(origin, r3.Scale(t, direction))
        rotated := s.rotatePoint(r3.Sub(p, s.Center))
        
        sphereDist := r3.Norm(rotated)/s.Radius - 1.0
        
        eyePos1 := r3.Vec{
            s.EyeDistance,
            0,
            -math.Sqrt(s.Radius*s.Radius - s.EyeDistance*s.EyeDistance),
        }
        eyePos2 := r3.Vec{
            -s.EyeDistance,
            0,
            -math.Sqrt(s.Radius*s.Radius - s.EyeDistance*s.EyeDistance),
        }
        
        eye1Dist := r3.Norm(r3.Sub(rotated, eyePos1)) - s.EyeRadius
        eye2Dist := r3.Norm(r3.Sub(rotated, eyePos2)) - s.EyeRadius
        
        dist := sphereDist
        
        if dist < EPSILON {
            hitPoint := r3.Add(origin, r3.Scale(t, direction))
            normal := r3.Unit(r3.Sub(hitPoint, s.Center))
            
            if eye1Dist < 0 || eye2Dist < 0 {
                return t, render.Color{2, 2, 2}, normal
            }
            return t, s.Color, normal
        }
        
        t += dist
        if t > MAX_DIST {
            break
        }
    }
    
    return -1, render.Color{0, 0, 0}, r3.Vec{}
}