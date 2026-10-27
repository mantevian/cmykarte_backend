package types

import "mantevian.xyz/cmykarte_backend/common/enums"

type ModifierType = string

const (
	ModifierTypeFull ModifierType = "full"
	ModifierTypeHalf ModifierType = "half"
	ModifierTypeOne  ModifierType = "one"
	ModifierTypeNone ModifierType = "none"
)

type Card struct {
	Type    enums.CardType        `json:"type"`
	Color   Color                 `json:"color"`
	Pattern map[Vec2]ModifierType `json:"pattern"`
	Cost    int                   `json:"cost"`
}
