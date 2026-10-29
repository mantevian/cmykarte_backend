package types

import (
	"fmt"
	"strconv"
	"strings"
)

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

func (v Vec2) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d,%d", v.X, v.Y)), nil
}

func (v *Vec2) UnmarshalText(text []byte) error {
	xText, yText, ok := strings.Cut(string(text), ",")
	if !ok {
		return fmt.Errorf("invalid Vec2 key %q: expected x,y", text)
	}

	x, err := strconv.Atoi(xText)
	if err != nil {
		return fmt.Errorf("invalid Vec2 x coordinate %q: %w", xText, err)
	}
	y, err := strconv.Atoi(yText)
	if err != nil {
		return fmt.Errorf("invalid Vec2 y coordinate %q: %w", yText, err)
	}

	v.X = x
	v.Y = y
	return nil
}
