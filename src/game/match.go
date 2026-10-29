package game

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
	"uuid"

	"mantevian.xyz/cmykarte_backend/common/enums"
	"mantevian.xyz/cmykarte_backend/common/types"
	"mantevian.xyz/cmykarte_backend/ws"
)

var (
	matchesMu sync.RWMutex
	Matches   = make(map[string]types.Match)
)

func StartMatch(queuedPlayer1 QueuedPlayer, queuedPlayer2 QueuedPlayer) (matchId string) {
	matchId = fmt.Sprintf("match_%s", uuid.NewV4())

	playerRGB := types.Player{
		UserId:   queuedPlayer1.UserId,
		Username: queuedPlayer1.Username,
		Type:     enums.CardTypeRGB,
		Money:    1,
		Deck:     Shuffle(CreateRGBDeck()),
	}
	DrawCards(&playerRGB, 5)

	playerCMY := types.Player{
		UserId:   queuedPlayer2.UserId,
		Username: queuedPlayer2.Username,
		Type:     enums.CardTypeCMY,
		Money:    1,
		Deck:     Shuffle(CreateCMYDeck()),
	}
	DrawCards(&playerCMY, 5)

	if rand.Int()%2 == 0 {
		playerRGB.UserId, playerCMY.UserId = playerCMY.UserId, playerRGB.UserId
		playerRGB.Username, playerCMY.Username = playerCMY.Username, playerRGB.Username
	}

	match := types.Match{
		PlayerRGB: playerRGB,
		PlayerCMY: playerCMY,
		Round:     1,
		Turn:      enums.CardTypeRGB,
		Field:     NewField(),
		StartedAt: time.Now(),
	}

	matchesMu.Lock()
	Matches[matchId] = match
	matchesMu.Unlock()

	ws.MoveUser(queuedPlayer1.UserId, matchId)
	ws.MoveUser(queuedPlayer2.UserId, matchId)

	ws.Broadcast(matchId, "match_started", types.MatchStartedMessage{
		MatchId: matchId,
		PlayerRGB: types.PublicPlayer{
			UserId:   playerRGB.UserId,
			Username: playerRGB.Username,
			Type:     playerRGB.Type,
			Money:    playerRGB.Money,
		},
		PlayerCMY: types.PublicPlayer{
			UserId:   playerCMY.UserId,
			Username: playerCMY.Username,
			Type:     playerCMY.Type,
			Money:    playerCMY.Money,
		},
		Round: match.Round,
		Turn:  match.Turn,
		Field: match.Field,
	})

	ws.SendToUser(playerRGB.UserId, "match_info_for_player", types.MatchInfoForPlayerMessage{
		MatchId: matchId,
		Hand:    playerRGB.Hand,
	})
	ws.SendToUser(playerCMY.UserId, "match_info_for_player", types.MatchInfoForPlayerMessage{
		MatchId: matchId,
		Hand:    playerCMY.Hand,
	})

	return
}

func EndTurn(matchId string, userId types.Uuid) bool {
	matchesMu.Lock()
	match, ok := Matches[matchId]
	if !ok {
		matchesMu.Unlock()
		return false
	}

	switch match.Turn {
	case enums.CardTypeRGB:
		if match.PlayerRGB.UserId != userId {
			matchesMu.Unlock()
			return false
		}
		match.PlayerRGB.Money = 0
		match.Turn = enums.CardTypeCMY
		match.PlayerCMY.Money = match.Round
		DrawCards(&match.PlayerCMY, 2)
	case enums.CardTypeCMY:
		if match.PlayerCMY.UserId != userId {
			matchesMu.Unlock()
			return false
		}
		match.PlayerCMY.Money = 0
		match.Round++
		match.Turn = enums.CardTypeRGB
		match.PlayerRGB.Money = match.Round
		DrawCards(&match.PlayerRGB, 2)
	default:
		matchesMu.Unlock()
		return false
	}

	Matches[matchId] = match
	matchesMu.Unlock()

	publishMatchState(matchId, match)
	return true
}

func PlayCard(matchId string, userId types.Uuid, cardIndex int) error {
	matchesMu.Lock()
	match, ok := Matches[matchId]
	if !ok {
		matchesMu.Unlock()
		return errors.New("match not found")
	}

	var activePlayer *types.Player
	switch match.Turn {
	case enums.CardTypeRGB:
		if match.PlayerRGB.UserId != userId {
			matchesMu.Unlock()
			return errors.New("it is not this player's turn")
		}
		activePlayer = &match.PlayerRGB
	case enums.CardTypeCMY:
		if match.PlayerCMY.UserId != userId {
			matchesMu.Unlock()
			return errors.New("it is not this player's turn")
		}
		activePlayer = &match.PlayerCMY
	default:
		matchesMu.Unlock()
		return errors.New("match has an invalid turn")
	}

	if cardIndex < 0 || cardIndex >= len(activePlayer.Hand) {
		matchesMu.Unlock()
		return errors.New("card index is out of range")
	}

	card := activePlayer.Hand[cardIndex]
	if card.Cost > activePlayer.Money {
		matchesMu.Unlock()
		return errors.New("not enough money to play this card")
	}

	activePlayer.Money -= card.Cost
	activePlayer.Hand = append(activePlayer.Hand[:cardIndex], activePlayer.Hand[cardIndex+1:]...)
	Matches[matchId] = match
	matchesMu.Unlock()

	publishMatchState(matchId, match)
	return nil
}

func publishMatchState(matchId string, match types.Match) {
	ws.Broadcast(matchId, "match_info_for_group", types.MatchInfoForGroupMessage{
		MatchId:   matchId,
		PlayerRGB: publicPlayer(match.PlayerRGB),
		PlayerCMY: publicPlayer(match.PlayerCMY),
		Round:     match.Round,
		Turn:      match.Turn,
		Field:     match.Field,
	})

	activePlayer := match.PlayerCMY
	if match.Turn == enums.CardTypeRGB {
		activePlayer = match.PlayerRGB
	}
	ws.SendToUser(activePlayer.UserId, "match_info_for_player", types.MatchInfoForPlayerMessage{
		MatchId: matchId,
		Hand:    activePlayer.Hand,
	})
}

func publicPlayer(player types.Player) types.PublicPlayer {
	return types.PublicPlayer{
		UserId:   player.UserId,
		Username: player.Username,
		Type:     player.Type,
		Money:    player.Money,
	}
}
