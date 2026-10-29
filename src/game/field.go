package game

import "mantevian.xyz/cmykarte_backend/common/types"

const fieldSize = 5

func NewField() types.Field {
	tiles := make(map[types.Vec2]types.Color, fieldSize*fieldSize)
	for x := 0; x < fieldSize; x++ {
		for y := 0; y < fieldSize; y++ {
			tiles[types.NewVec2(x, y)] = types.NewColor(10, 10, 10)
		}
	}

	return types.Field{Tiles: tiles}
}
