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
	Rules     MatchRules
	Finished  bool
	StartedAt time.Time
}

type MatchRules struct {
	StartingHandSize       int             `json:"starting_hand_size"`
	GuaranteedStartingCost int             `json:"guaranteed_starting_cost"`
	HandLimit              int             `json:"hand_limit"`
	RoundDraws             []RoundDrawRule `json:"round_draws"`
}

type RoundDrawRule struct {
	FromRound     int  `json:"from_round"`
	ToRound       *int `json:"to_round,omitempty"`
	MinCost       int  `json:"min_cost"`
	MaxCost       *int `json:"max_cost,omitempty"`
	CardsPerRound int  `json:"cards_per_round"`
}

type MatchStartedMessage struct {
	MatchId   string         `json:"match_id"`
	PlayerRGB PublicPlayer   `json:"player_rgb"`
	PlayerCMY PublicPlayer   `json:"player_cmy"`
	Round     int            `json:"round"`
	Turn      enums.CardType `json:"turn"`
	Field     Field          `json:"field"`
	Rules     MatchRules     `json:"rules"`
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
	Rules     MatchRules     `json:"rules"`
}

type MatchEndedMessage struct {
	MatchId   string       `json:"match_id"`
	Winner    PublicPlayer `json:"winner"`
	PlayerRGB PublicPlayer `json:"player_rgb"`
	PlayerCMY PublicPlayer `json:"player_cmy"`
	Field     Field        `json:"field"`
}

type MatchInfoForPlayerMessage struct {
	MatchId string `json:"match_id"`
	Hand    []Card `json:"hand"`
}
