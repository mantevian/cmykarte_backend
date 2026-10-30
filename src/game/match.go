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

const winningScore = 13

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
	DrawCards(&playerRGB, 3)

	playerCMY := types.Player{
		UserId:   queuedPlayer2.UserId,
		Username: queuedPlayer2.Username,
		Type:     enums.CardTypeCMY,
		Money:    1,
		Deck:     Shuffle(CreateCMYDeck()),
	}
	DrawCards(&playerCMY, 3)

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
		MatchId:   matchId,
		PlayerRGB: publicPlayer(playerRGB),
		PlayerCMY: publicPlayer(playerCMY),
		Round:     match.Round,
		Turn:      match.Turn,
		Field:     match.Field,
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
	if !ok || match.Finished {
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
		DrawCards(&match.PlayerCMY, 2)
	default:
		matchesMu.Unlock()
		return false
	}

	Matches[matchId] = match
	matchesMu.Unlock()

	publishMatchState(matchId, match)
	return true
}

func PlayCard(matchId string, userId types.Uuid, cardIndex int, tile types.Vec2) error {
	matchesMu.Lock()
	match, ok := Matches[matchId]
	if !ok {
		matchesMu.Unlock()
		return errors.New("match not found")
	}
	if match.Finished {
		matchesMu.Unlock()
		return errors.New("match has ended")
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

	changes, err := cardFieldChanges(match.Field, card, tile)
	if err != nil {
		matchesMu.Unlock()
		return err
	}

	initializeFieldClaims(&match.Field)
	for _, change := range changes {
		if match.Field.ClaimedBy[change.tile] != "" {
			continue
		}

		color := clampColor(change.color)
		match.Field.Tiles[change.tile] = color
		switch color {
		case types.NewColor(20, 20, 20):
			match.Field.ClaimedBy[change.tile] = enums.CardTypeRGB
			match.PlayerRGB.Score++
		case types.NewColor(0, 0, 0):
			match.Field.ClaimedBy[change.tile] = enums.CardTypeCMY
			match.PlayerCMY.Score++
		}
	}
	activePlayer.Money -= card.Cost
	activePlayer.Hand = append(activePlayer.Hand[:cardIndex], activePlayer.Hand[cardIndex+1:]...)
	var winner *types.Player
	if match.PlayerRGB.Score >= winningScore {
		winner = &match.PlayerRGB
	} else if match.PlayerCMY.Score >= winningScore {
		winner = &match.PlayerCMY
	}
	if winner != nil {
		match.Finished = true
	}
	Matches[matchId] = match
	matchesMu.Unlock()

	publishMatchState(matchId, match)
	if winner != nil {
		ws.Broadcast(matchId, "match_ended", types.MatchEndedMessage{
			MatchId:   matchId,
			Winner:    publicPlayer(*winner),
			PlayerRGB: publicPlayer(match.PlayerRGB),
			PlayerCMY: publicPlayer(match.PlayerCMY),
			Field:     match.Field,
		})
	}
	return nil
}

type tileColorChange struct {
	tile  types.Vec2
	color types.Color
}

func cardFieldChanges(field types.Field, card types.Card, anchor types.Vec2) ([]tileColorChange, error) {
	if _, ok := field.Tiles[anchor]; !ok {
		return nil, errors.New("tile is outside the field")
	}

	var delta types.Color
	switch card.Type {
	case enums.CardTypeRGB:
		delta = card.Color
	case enums.CardTypeCMY:
		delta = types.Color{
			R: -(card.Color.G/2 + card.Color.B/2),
			G: -(card.Color.R/2 + card.Color.B/2),
			B: -(card.Color.R/2 + card.Color.G/2),
		}
	default:
		return nil, errors.New("card has an invalid type")
	}

	changes := make([]tileColorChange, 0, len(card.Pattern))
	for offset, modifier := range card.Pattern {
		var applied types.Color
		switch modifier {
		case types.ModifierTypeFull:
			applied = delta
		case types.ModifierTypeHalf:
			applied = types.Color{R: delta.R / 2, G: delta.G / 2, B: delta.B / 2}
		case types.ModifierTypeOne:
			applied = types.Color{R: 1, G: 1, B: 1}
		case types.ModifierTypeNone:
			continue
		default:
			return nil, errors.New("card has an invalid pattern modifier")
		}

		target := types.NewVec2(anchor.X+offset.X, anchor.Y+offset.Y)
		color, ok := field.Tiles[target]
		if !ok {
			continue
		}
		changes = append(changes, tileColorChange{
			tile:  target,
			color: types.Color{R: color.R + applied.R, G: color.G + applied.G, B: color.B + applied.B},
		})
	}

	return changes, nil
}

func initializeFieldClaims(field *types.Field) {
	if field.ClaimedBy == nil {
		field.ClaimedBy = make(map[types.Vec2]enums.CardType, len(field.Tiles))
	}
	for tile := range field.Tiles {
		if _, ok := field.ClaimedBy[tile]; !ok {
			field.ClaimedBy[tile] = ""
		}
	}
}

func clampColor(color types.Color) types.Color {
	return types.Color{
		R: clampChannel(color.R),
		G: clampChannel(color.G),
		B: clampChannel(color.B),
	}
}

func clampChannel(channel int) int {
	if channel < 0 {
		return 0
	}
	if channel > 20 {
		return 20
	}
	return channel
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

	ws.SendToUser(match.PlayerRGB.UserId, "match_info_for_player", types.MatchInfoForPlayerMessage{
		MatchId: matchId,
		Hand:    match.PlayerRGB.Hand,
	})
	ws.SendToUser(match.PlayerCMY.UserId, "match_info_for_player", types.MatchInfoForPlayerMessage{
		MatchId: matchId,
		Hand:    match.PlayerCMY.Hand,
	})
}

func publicPlayer(player types.Player) types.PublicPlayer {
	return types.PublicPlayer{
		UserId:   player.UserId,
		Username: player.Username,
		Type:     player.Type,
		Money:    player.Money,
		Score:    player.Score,
	}
}
