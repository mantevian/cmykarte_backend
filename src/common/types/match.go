package types

import (
	"time"

	"mantevian.xyz/cmykarte_backend/common/enums"
)

type Match struct {
	PlayerRGB Player
	PlayerCMY Player
	Round     int
	Turn      enums.CardType
	Field     Field
	StartedAt time.Time
}

type MatchStartedMessage struct {
	MatchId   string         `json:"match_id"`
	PlayerRGB PublicPlayer   `json:"player_rgb"`
	PlayerCMY PublicPlayer   `json:"player_cmy"`
	Round     int            `json:"round"`
	Turn      enums.CardType `json:"turn"`
	Field     Field          `json:"field"`
}

type EndTurnRequest struct {
	MatchId string `json:"match_id"`
}

type PlayCardRequest struct {
	MatchId   string `json:"match_id"`
	CardIndex int    `json:"card_index"`
	Tile      *Vec2  `json:"tile"`
}

type MatchInfoForGroupMessage struct {
	MatchId   string         `json:"match_id"`
	PlayerRGB PublicPlayer   `json:"player_rgb"`
	PlayerCMY PublicPlayer   `json:"player_cmy"`
	Round     int            `json:"round"`
	Turn      enums.CardType `json:"turn"`
	Field     Field          `json:"field"`
}

type MatchInfoForPlayerMessage struct {
	MatchId string `json:"match_id"`
	Hand    []Card `json:"hand"`
}
