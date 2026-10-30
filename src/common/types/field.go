package types

import "mantevian.xyz/cmykarte_backend/common/enums"

type Field struct {
	Tiles     map[Vec2]Color          `json:"tiles"`
	ClaimedBy map[Vec2]enums.CardType `json:"claimed_by"`
}
