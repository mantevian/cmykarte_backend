package types

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

func NewColor(r int, g int, b int) Color {
	return Color{
		R: r,
		G: g,
		B: b,
	}
}
