package game

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"mantevian.xyz/cmykarte_backend/common/types"
)

const rulesetVersion = "1"

var (
	persistenceMu sync.RWMutex
	persistenceDB *sql.DB
)

type matchEvent struct {
	eventType string
	player    types.Player
	round     int
	card      *types.Card
	target    *types.Vec2
	outcome   any
}

func SetPersistenceDB(db *sql.DB) {
	persistenceMu.Lock()
	persistenceDB = db
	persistenceMu.Unlock()
}

func CloseInterruptedMatches(db *sql.DB) error {
	_, err := db.Exec(
		`update matches
		set status = 'abandoned',
			ended_at = now(),
			duration_seconds = extract(epoch from now() - started_at),
			end_reason = 'server_restart'
		where status = 'in_progress'`,
	)
	return err
}

func currentDB() *sql.DB {
	persistenceMu.RLock()
	defer persistenceMu.RUnlock()
	return persistenceDB
}

func persistMatchStart(matchID string, match types.Match) error {
	db := currentDB()
	if db == nil {
		return nil
	}

	rules, err := json.Marshal(match.Rules)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`insert into matches (
			id, player_rgb_id, player_cmy_id, first_player_id, status,
			started_at, rounds_played, build_version, ruleset_version, rules
		) values ($1, $2, $3, $4, 'in_progress', $5, $6, $7, $8, $9)`,
		matchID,
		match.PlayerRGB.UserId,
		match.PlayerCMY.UserId,
		match.PlayerRGB.UserId,
		match.StartedAt,
		match.Round,
		buildVersion(),
		rulesetVersion,
		string(rules),
	)
	return err
}

func persistMatchTransition(matchID string, match types.Match, event matchEvent, endReason string, winner *types.Player) error {
	db := currentDB()
	if db == nil {
		return nil
	}

	outcome, err := json.Marshal(event.outcome)
	if err != nil {
		return err
	}
	var cardColor, cardPattern any
	var cardCost, cardType, targetX, targetY any
	if event.card != nil {
		cardColorJSON, marshalErr := json.Marshal(event.card.Color)
		err = marshalErr
		if err != nil {
			return err
		}
		cardColor = string(cardColorJSON)
		cardPatternJSON, marshalErr := json.Marshal(event.card.Pattern)
		err = marshalErr
		if err != nil {
			return err
		}
		cardPattern = string(cardPatternJSON)
		cardCost = event.card.Cost
		cardType = event.card.Type
	}
	if event.target != nil {
		targetX = event.target.X
		targetY = event.target.Y
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var turnNumber int
	if err := tx.QueryRow(
		`select count(*) from match_events where match_id = $1 and event_type = 'turn_ended'`,
		matchID,
	).Scan(&turnNumber); err != nil {
		return err
	}
	turnNumber++
	if _, err := tx.Exec(
		`insert into match_events (
			match_id, event_type, round, turn_number, player_id, player_side,
			card_cost, card_type, card_color, card_pattern, target_x, target_y,
			score_rgb, score_cmy, outcome
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		matchID,
		event.eventType,
		event.round,
		turnNumber,
		event.player.UserId,
		event.player.Type,
		cardCost,
		cardType,
		cardColor,
		cardPattern,
		targetX,
		targetY,
		match.PlayerRGB.Score,
		match.PlayerCMY.Score,
		string(outcome),
	); err != nil {
		return err
	}

	status := "in_progress"
	var endedAt, durationSeconds, winnerID, winnerSide any
	if winner != nil {
		status = "completed"
		winnerID = winner.UserId
		winnerSide = winner.Type
	} else if endReason != "" {
		status = "abandoned"
	}
	if status != "in_progress" {
		now := time.Now()
		endedAt = now
		durationSeconds = now.Sub(match.StartedAt).Seconds()
	}
	var persistedEndReason any
	if endReason != "" {
		persistedEndReason = endReason
	}
	if _, err := tx.Exec(
		`update matches set
			status = $2, winner_id = $3, winner_side = $4, ended_at = $5,
			duration_seconds = $6, rounds_played = $7, final_score_rgb = $8,
			final_score_cmy = $9, end_reason = $10
		where id = $1`,
		matchID,
		status,
		winnerID,
		winnerSide,
		endedAt,
		durationSeconds,
		match.Round,
		match.PlayerRGB.Score,
		match.PlayerCMY.Score,
		persistedEndReason,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func buildVersion() string {
	if version := os.Getenv("BUILD_VERSION"); version != "" {
		return version
	}
	return "unknown"
}

func cloneMatch(match types.Match) types.Match {
	match.PlayerRGB.Hand = append([]types.Card(nil), match.PlayerRGB.Hand...)
	match.PlayerRGB.Deck = append([]types.Card(nil), match.PlayerRGB.Deck...)
	match.PlayerCMY.Hand = append([]types.Card(nil), match.PlayerCMY.Hand...)
	match.PlayerCMY.Deck = append([]types.Card(nil), match.PlayerCMY.Deck...)

	tiles := make(map[types.Vec2]types.Color, len(match.Field.Tiles))
	for tile, color := range match.Field.Tiles {
		tiles[tile] = color
	}
	match.Field.Tiles = tiles

	claims := make(map[types.Vec2]string, len(match.Field.ClaimedBy))
	for tile, side := range match.Field.ClaimedBy {
		claims[tile] = side
	}
	match.Field.ClaimedBy = claims
	return match
}

func AbandonMatchesForUser(userID types.Uuid) error {
	matchesMu.Lock()
	defer matchesMu.Unlock()

	var failures []error
	for matchID, match := range Matches {
		if match.Finished || (match.PlayerRGB.UserId != userID && match.PlayerCMY.UserId != userID) {
			continue
		}

		candidate := cloneMatch(match)
		candidate.Finished = true
		player := candidate.PlayerRGB
		if candidate.PlayerCMY.UserId == userID {
			player = candidate.PlayerCMY
		}
		if err := persistMatchTransition(matchID, candidate, matchEvent{
			eventType: "match_abandoned",
			player:    player,
			round:     candidate.Round,
			outcome:   map[string]string{"reason": "player_disconnected"},
		}, "player_disconnected", nil); err != nil {
			failures = append(failures, err)
			continue
		}
		Matches[matchID] = candidate
	}
	return errors.Join(failures...)
}
