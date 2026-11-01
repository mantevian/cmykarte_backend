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
const defaultHandLimit = 10

var (
	matchesMu sync.RWMutex
	Matches   = make(map[string]types.Match)
)

func StartMatch(queuedPlayer1 QueuedPlayer, queuedPlayer2 QueuedPlayer) (matchId string, err error) {
	matchId = fmt.Sprintf("match_%s", uuid.NewV4())
	rules := defaultMatchRules()

	playerRGB := types.Player{
		UserId:   queuedPlayer1.UserId,
		Username: queuedPlayer1.Username,
		Type:     enums.CardTypeRGB,
		Money:    1,
		Deck:     Shuffle(CreateRGBDeck()),
	}
	if err := drawStartingHand(&playerRGB, rules); err != nil {
		return "", fmt.Errorf("draw RGB starting hand: %w", err)
	}

	playerCMY := types.Player{
		UserId:   queuedPlayer2.UserId,
		Username: queuedPlayer2.Username,
		Type:     enums.CardTypeCMY,
		Money:    1,
		Deck:     Shuffle(CreateCMYDeck()),
	}
	if err := drawStartingHand(&playerCMY, rules); err != nil {
		return "", fmt.Errorf("draw CMY starting hand: %w", err)
	}

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
		Rules:     rules,
		StartedAt: time.Now(),
	}

	if err := persistMatchStart(matchId, match); err != nil {
		return "", fmt.Errorf("persist match start: %w", err)
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
		Rules:     match.Rules,
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
	return EndTurnWithError(matchId, userId) == nil
}

func EndTurnWithError(matchId string, userId types.Uuid) error {
	matchesMu.Lock()
	match, ok := Matches[matchId]
	if !ok || match.Finished {
		matchesMu.Unlock()
		return errors.New("match not found or has ended")
	}

	candidate := cloneMatch(match)
	var activePlayer types.Player
	switch candidate.Turn {
	case enums.CardTypeRGB:
		if candidate.PlayerRGB.UserId != userId {
			matchesMu.Unlock()
			return errors.New("it is not this player's turn")
		}
		activePlayer = candidate.PlayerRGB
		candidate.PlayerRGB.Money = 0
		candidate.Turn = enums.CardTypeCMY
		candidate.PlayerCMY.Money = candidate.Round
	case enums.CardTypeCMY:
		if candidate.PlayerCMY.UserId != userId {
			matchesMu.Unlock()
			return errors.New("it is not this player's turn")
		}
		activePlayer = candidate.PlayerCMY
		candidate.PlayerCMY.Money = 0
		candidate.Round++
		candidate.Turn = enums.CardTypeRGB
		candidate.PlayerRGB.Money = candidate.Round
		if candidate.Rules.HandLimit == 0 {
			candidate.Rules = defaultMatchRules()
		}
		for _, rule := range candidate.Rules.RoundDraws {
			if candidate.Round >= rule.FromRound && (rule.ToRound == nil || candidate.Round <= *rule.ToRound) {
				if err := drawCardsForRound(&candidate.PlayerRGB, rule.CardsPerRound, candidate.Round, candidate.Rules); err != nil {
					matchesMu.Unlock()
					return fmt.Errorf("draw RGB round cards: %w", err)
				}
				if err := drawCardsForRound(&candidate.PlayerCMY, rule.CardsPerRound, candidate.Round, candidate.Rules); err != nil {
					matchesMu.Unlock()
					return fmt.Errorf("draw CMY round cards: %w", err)
				}
				break
			}
		}
	default:
		matchesMu.Unlock()
		return errors.New("match has an invalid turn")
	}

	if err := persistMatchTransition(matchId, candidate, matchEvent{
		eventType: "turn_ended",
		player:    activePlayer,
		round:     match.Round,
		outcome:   map[string]any{"round_advanced": candidate.Round != match.Round},
	}, "", nil); err != nil {
		matchesMu.Unlock()
		return fmt.Errorf("persist turn end: %w", err)
	}
	Matches[matchId] = candidate
	matchesMu.Unlock()

	publishMatchState(matchId, candidate)
	return nil
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

	candidate := cloneMatch(match)
	var activePlayer *types.Player
	switch candidate.Turn {
	case enums.CardTypeRGB:
		if candidate.PlayerRGB.UserId != userId {
			matchesMu.Unlock()
			return errors.New("it is not this player's turn")
		}
		activePlayer = &candidate.PlayerRGB
	case enums.CardTypeCMY:
		if candidate.PlayerCMY.UserId != userId {
			matchesMu.Unlock()
			return errors.New("it is not this player's turn")
		}
		activePlayer = &candidate.PlayerCMY
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

	changes, err := cardFieldChanges(candidate.Field, card, tile)
	if err != nil {
		matchesMu.Unlock()
		return err
	}

	rgbScoreBefore := candidate.PlayerRGB.Score
	cmyScoreBefore := candidate.PlayerCMY.Score
	initializeFieldClaims(&candidate.Field)
	claimedTiles := 0
	for _, change := range changes {
		if candidate.Field.ClaimedBy[change.tile] != "" {
			continue
		}

		color := clampColor(change.color)
		candidate.Field.Tiles[change.tile] = color
		switch color {
		case types.NewColor(20, 20, 20):
			candidate.Field.ClaimedBy[change.tile] = enums.CardTypeRGB
			candidate.PlayerRGB.Score++
			claimedTiles++
		case types.NewColor(0, 0, 0):
			candidate.Field.ClaimedBy[change.tile] = enums.CardTypeCMY
			candidate.PlayerCMY.Score++
			claimedTiles++
		}
	}
	activePlayer.Money -= card.Cost
	activePlayer.Hand = append(activePlayer.Hand[:cardIndex], activePlayer.Hand[cardIndex+1:]...)
	var winner *types.Player
	if candidate.PlayerRGB.Score >= winningScore {
		winner = &candidate.PlayerRGB
	} else if candidate.PlayerCMY.Score >= winningScore {
		winner = &candidate.PlayerCMY
	}
	if winner != nil {
		candidate.Finished = true
	}
	endReason := ""
	if winner != nil {
		endReason = "victory"
	}
	outcome := map[string]any{
		"tiles_changed":   len(changes),
		"tiles_claimed":   claimedTiles,
		"rgb_score_delta": candidate.PlayerRGB.Score - rgbScoreBefore,
		"cmy_score_delta": candidate.PlayerCMY.Score - cmyScoreBefore,
		"winner_side":     "",
	}
	if winner != nil {
		outcome["winner_side"] = winner.Type
	}
	if err := persistMatchTransition(matchId, candidate, matchEvent{
		eventType: "card_played",
		player:    *activePlayer,
		round:     candidate.Round,
		card:      &card,
		target:    &tile,
		outcome:   outcome,
	}, endReason, winner); err != nil {
		matchesMu.Unlock()
		return fmt.Errorf("persist card play: %w", err)
	}
	Matches[matchId] = candidate
	matchesMu.Unlock()

	publishMatchState(matchId, candidate)
	if winner != nil {
		ws.Broadcast(matchId, "match_ended", types.MatchEndedMessage{
			MatchId:   matchId,
			Winner:    publicPlayer(*winner),
			PlayerRGB: publicPlayer(candidate.PlayerRGB),
			PlayerCMY: publicPlayer(candidate.PlayerCMY),
			Field:     candidate.Field,
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
	if match.Rules.HandLimit == 0 {
		match.Rules = defaultMatchRules()
	}
	ws.Broadcast(matchId, "match_info_for_group", types.MatchInfoForGroupMessage{
		MatchId:   matchId,
		PlayerRGB: publicPlayer(match.PlayerRGB),
		PlayerCMY: publicPlayer(match.PlayerCMY),
		Round:     match.Round,
		Turn:      match.Turn,
		Field:     match.Field,
		Rules:     match.Rules,
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

func defaultMatchRules() types.MatchRules {
	round5 := 5
	round10 := 10
	round20 := 20
	cost3 := 3
	cost5 := 5
	cost7 := 7
	return types.MatchRules{
		StartingHandSize:       4,
		GuaranteedStartingCost: 1,
		HandLimit:              defaultHandLimit,
		RoundDraws: []types.RoundDrawRule{
			{FromRound: 1, ToRound: &round5, MaxCost: &cost3, CardsPerRound: 1},
			{FromRound: 6, ToRound: &round10, MaxCost: &cost5, CardsPerRound: 2},
			{FromRound: 11, ToRound: &round20, MaxCost: &cost7, CardsPerRound: 2},
			{FromRound: 21, CardsPerRound: 3},
		},
	}
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
