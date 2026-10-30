package types

import "mantevian.xyz/cmykarte_backend/common/enums"

type PublicPlayer struct {
	UserId   Uuid           `json:"user_id"`
	Username string         `json:"username"`
	Type     enums.CardType `json:"type"`
	Money    int            `json:"money"`
	Score    int            `json:"score"`
}

type Player struct {
	UserId   Uuid           `json:"user_id"`
	Username string         `json:"username"`
	Type     enums.CardType `json:"type"`
	Money    int            `json:"money"`
	Score    int            `json:"score"`
	Hand     []Card         `json:"hand"`
	Deck     []Card         `json:"deck"`
}
