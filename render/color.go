package render

type Color struct {
    R, G, B float64
}

var (
    Black = Color{0, 0, 0}
    White = Color{1, 1, 1}
    Red = Color{1, 0, 0}
    Green = Color{0, 1, 0}
    Blue = Color{0, 0, 1}
    Yellow = Color{1, 1, 0}
    Cyan = Color{0, 1, 1}
    Magenta = Color{1, 0, 1}
)