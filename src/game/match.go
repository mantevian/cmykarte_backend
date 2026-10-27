package game

import (
	"fmt"
	"math/rand/v2"
	"time"
	"uuid"

	"mantevian.xyz/cmykarte_backend/common/enums"
	"mantevian.xyz/cmykarte_backend/common/types"
	"mantevian.xyz/cmykarte_backend/ws"
)

var Matches map[string]types.Match

func StartMatch(user1 types.Uuid, user2 types.Uuid) (matchId string) {
	matchId = fmt.Sprintf("match_%s", uuid.NewV4())

	var playerRGB types.Player
	var playerCMY types.Player

	player1 := types.Player{
		UserId: user1,
		Deck:   Shuffle(CreateRGBDeck()),
	}

	player2 := types.Player{
		UserId: user2,
		Deck:   Shuffle(CreateRGBDeck()),
	}

	if rand.Int()%2 == 0 {
		playerRGB = player1
		playerCMY = player2
	} else {
		playerRGB = player2
		playerCMY = player1
	}

	match := types.Match{
		PlayerRGB: playerRGB,
		PlayerCMY: playerCMY,
		Turn:      enums.CardTypeRGB,
		StartedAt: time.Now(),
	}

	Matches[matchId] = match

	ws.MoveUser(user1, matchId)
	ws.MoveUser(user2, matchId)

	ws.Broadcast(matchId, "match_started", types.MatchStartedMessage{
		MatchId: matchId,
		PlayerRGB: types.PublicPlayer{
			UserId: playerRGB.UserId,
			Type:   playerRGB.Type,
			Money:  playerRGB.Money,
		},
		PlayerCMY: types.PublicPlayer{
			UserId: playerCMY.UserId,
			Type:   playerCMY.Type,
			Money:  playerCMY.Money,
		},
	})

	return
}
