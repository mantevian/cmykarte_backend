package types

type Vec2 struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewVec2(x int, y int) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}
