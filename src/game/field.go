package game

import (
	"mantevian.xyz/cmykarte_backend/common/enums"
	"mantevian.xyz/cmykarte_backend/common/types"
)

const fieldSize = 5

func NewField() types.Field {
	tiles := make(map[types.Vec2]types.Color, fieldSize*fieldSize)
	claimedBy := make(map[types.Vec2]enums.CardType, fieldSize*fieldSize)
	for x := 0; x < fieldSize; x++ {
		for y := 0; y < fieldSize; y++ {
			tile := types.NewVec2(x, y)
			tiles[tile] = types.NewColor(10, 10, 10)
			claimedBy[tile] = ""
		}
	}

	return types.Field{Tiles: tiles, ClaimedBy: claimedBy}
}
