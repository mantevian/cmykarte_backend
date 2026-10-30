package game

import (
	"encoding/json"
	"testing"

	"mantevian.xyz/cmykarte_backend/common/enums"
	"mantevian.xyz/cmykarte_backend/common/types"
)

func TestPlayCardUsesTileFromWebSocketRequest(t *testing.T) {
	const matchId = "test_card_request_tile"
	var request types.PlayCardRequest
	if err := json.Unmarshal([]byte(`{"match_id":"test_card_request_tile","card_index":2,"tile":{"x":1,"y":3}}`), &request); err != nil {
		t.Fatalf("failed to decode play-card request: %v", err)
	}
	if request.Tile == nil || *request.Tile != (types.Vec2{X: 1, Y: 3}) {
		t.Fatalf("decoded tile is %+v, want (1,3)", request.Tile)
	}

	card := types.Card{
		Type:  enums.CardTypeRGB,
		Color: types.NewColor(4, 2, 1),
		Cost:  1,
		Pattern: map[types.Vec2]types.ModifierType{
			types.NewVec2(0, 0): types.ModifierTypeFull,
			types.NewVec2(1, 0): types.ModifierTypeFull,
		},
	}
	match := types.Match{
		PlayerRGB: types.Player{
			UserId: "rgb-user",
			Money:  1,
			Hand:   []types.Card{{}, {}, card},
		},
		PlayerCMY: types.Player{UserId: "cmy-user"},
		Turn:      enums.CardTypeRGB,
		Field:     NewField(),
	}
	Matches[matchId] = match
	defer delete(Matches, matchId)

	if err := PlayCard(request.MatchId, "rgb-user", request.CardIndex, *request.Tile); err != nil {
		t.Fatalf("PlayCard returned error: %v", err)
	}

	updated := Matches[matchId].Field.Tiles
	if got := updated[types.NewVec2(1, 3)]; got != (types.Color{R: 14, G: 12, B: 11}) {
		t.Errorf("selected tile has color %+v, want {R:14 G:12 B:11}", got)
	}
	if got := updated[types.NewVec2(2, 3)]; got != (types.Color{R: 14, G: 12, B: 11}) {
		t.Errorf("pattern tile has color %+v, want {R:14 G:12 B:11}", got)
	}
	if got := updated[types.NewVec2(0, 0)]; got != (types.Color{R: 10, G: 10, B: 10}) {
		t.Errorf("unselected origin tile changed unexpectedly: %+v", got)
	}
}

func TestPlayCardRequestRejectsTileWithoutCoordinates(t *testing.T) {
	var request types.PlayCardRequest
	if err := json.Unmarshal([]byte(`{"match_id":"test","card_index":0,"tile":{}}`), &request); err == nil {
		t.Fatal("play-card request accepted a tile without coordinates")
	}
}

func TestPlayRGBCardAppliesPatternModifiers(t *testing.T) {
	const matchId = "test_rgb_card_effects"
	anchor := types.NewVec2(0, 0)
	right := types.NewVec2(1, 0)
	down := types.NewVec2(0, 1)
	match := types.Match{
		PlayerRGB: types.Player{
			UserId: "rgb-user",
			Money:  3,
			Hand: []types.Card{{
				Type:  enums.CardTypeRGB,
				Color: types.NewColor(8, 4, 2),
				Cost:  2,
				Pattern: map[types.Vec2]types.ModifierType{
					types.NewVec2(0, 0):  types.ModifierTypeFull,
					types.NewVec2(1, 0):  types.ModifierTypeHalf,
					types.NewVec2(0, 1):  types.ModifierTypeOne,
					types.NewVec2(-1, 0): types.ModifierTypeFull,
				},
			}},
		},
		PlayerCMY: types.Player{UserId: "cmy-user"},
		Turn:      enums.CardTypeRGB,
		Field: types.Field{Tiles: map[types.Vec2]types.Color{
			anchor: types.NewColor(10, 10, 10),
			right:  types.NewColor(10, 10, 10),
			down:   types.NewColor(10, 10, 10),
		}},
	}
	Matches[matchId] = match
	defer delete(Matches, matchId)

	if err := PlayCard(matchId, "rgb-user", 0, anchor); err != nil {
		t.Fatalf("PlayCard returned error: %v", err)
	}

	got := Matches[matchId]
	if color := got.Field.Tiles[anchor]; color != (types.Color{R: 18, G: 14, B: 12}) {
		t.Errorf("full modifier produced color %+v, want {R:18 G:14 B:12}", color)
	}
	if color := got.Field.Tiles[right]; color != (types.Color{R: 14, G: 12, B: 11}) {
		t.Errorf("half modifier produced color %+v, want {R:14 G:12 B:11}", color)
	}
	if color := got.Field.Tiles[down]; color != (types.Color{R: 11, G: 11, B: 11}) {
		t.Errorf("one modifier produced color %+v, want {R:11 G:11 B:11}", color)
	}
	if got.PlayerRGB.Money != 1 || len(got.PlayerRGB.Hand) != 0 {
		t.Errorf("played card was not charged and removed: money=%d hand=%d", got.PlayerRGB.Money, len(got.PlayerRGB.Hand))
	}
}

func TestPlayCMYCardSubtractsConvertedColor(t *testing.T) {
	const matchId = "test_cmy_card_effects"
	anchor := types.NewVec2(2, 2)
	match := types.Match{
		PlayerCMY: types.Player{
			UserId: "cmy-user",
			Money:  1,
			Hand: []types.Card{{
				Type:  enums.CardTypeCMY,
				Color: types.NewColor(4, 6, 8),
				Cost:  1,
				Pattern: map[types.Vec2]types.ModifierType{
					types.NewVec2(0, 0): types.ModifierTypeFull,
				},
			}},
		},
		PlayerRGB: types.Player{UserId: "rgb-user"},
		Turn:      enums.CardTypeCMY,
		Field: types.Field{Tiles: map[types.Vec2]types.Color{
			anchor: types.NewColor(20, 30, 40),
		}},
	}
	Matches[matchId] = match
	defer delete(Matches, matchId)

	if err := PlayCard(matchId, "cmy-user", 0, anchor); err != nil {
		t.Fatalf("PlayCard returned error: %v", err)
	}

	got := Matches[matchId].Field.Tiles[anchor]
	if got != (types.Color{R: 13, G: 20, B: 20}) {
		t.Errorf("CMY effect produced color %+v, want {R:13 G:20 B:20}", got)
	}
}

func TestPlayCMYCardClaimsTileAtMinimumColor(t *testing.T) {
	const matchId = "test_cmy_claim"
	tile := types.NewVec2(0, 0)
	match := types.Match{
		PlayerCMY: types.Player{
			UserId: "cmy-user",
			Money:  1,
			Hand: []types.Card{{
				Type:    enums.CardTypeCMY,
				Color:   types.NewColor(2, 2, 2),
				Cost:    1,
				Pattern: map[types.Vec2]types.ModifierType{types.NewVec2(0, 0): types.ModifierTypeFull},
			}},
		},
		PlayerRGB: types.Player{UserId: "rgb-user"},
		Turn:      enums.CardTypeCMY,
		Field: types.Field{Tiles: map[types.Vec2]types.Color{
			tile: types.NewColor(1, 1, 1),
		}},
	}
	Matches[matchId] = match
	defer delete(Matches, matchId)

	if err := PlayCard(matchId, "cmy-user", 0, tile); err != nil {
		t.Fatalf("CMY PlayCard returned error: %v", err)
	}

	got := Matches[matchId]
	if color := got.Field.Tiles[tile]; color != (types.Color{R: 0, G: 0, B: 0}) {
		t.Errorf("claimed tile has color %+v, want {R:0 G:0 B:0}", color)
	}
	if got.Field.ClaimedBy[tile] != enums.CardTypeCMY || got.PlayerCMY.Score != 1 {
		t.Errorf("CMY claim state is owner=%q score=%d, want CMY and 1", got.Field.ClaimedBy[tile], got.PlayerCMY.Score)
	}
}

func TestPlayCardClampsColorsClaimsTilesAndLocksThem(t *testing.T) {
	const matchId = "test_tile_claim_and_lock"
	claimed := types.NewVec2(0, 0)
	capped := types.NewVec2(1, 0)
	match := types.Match{
		PlayerRGB: types.Player{
			UserId: "rgb-user",
			Money:  2,
			Hand: []types.Card{{
				Type:  enums.CardTypeRGB,
				Color: types.NewColor(4, 4, 4),
				Cost:  1,
				Pattern: map[types.Vec2]types.ModifierType{
					types.NewVec2(0, 0): types.ModifierTypeFull,
					types.NewVec2(1, 0): types.ModifierTypeFull,
				},
			}},
		},
		PlayerCMY: types.Player{
			UserId: "cmy-user",
			Money:  1,
			Hand: []types.Card{{
				Type:  enums.CardTypeCMY,
				Color: types.NewColor(4, 4, 4),
				Cost:  1,
				Pattern: map[types.Vec2]types.ModifierType{
					types.NewVec2(0, 0): types.ModifierTypeFull,
				},
			}},
		},
		Turn: enums.CardTypeRGB,
		Field: types.Field{Tiles: map[types.Vec2]types.Color{
			claimed: types.NewColor(19, 19, 19),
			capped:  types.NewColor(19, 2, 1),
		}},
	}
	Matches[matchId] = match
	defer delete(Matches, matchId)

	if err := PlayCard(matchId, "rgb-user", 0, claimed); err != nil {
		t.Fatalf("RGB PlayCard returned error: %v", err)
	}

	got := Matches[matchId]
	if color := got.Field.Tiles[claimed]; color != (types.Color{R: 20, G: 20, B: 20}) {
		t.Errorf("claimed tile has color %+v, want {R:20 G:20 B:20}", color)
	}
	if got.Field.ClaimedBy[claimed] != enums.CardTypeRGB || got.PlayerRGB.Score != 1 {
		t.Errorf("RGB claim state is owner=%q score=%d, want RGB and 1", got.Field.ClaimedBy[claimed], got.PlayerRGB.Score)
	}
	if color := got.Field.Tiles[capped]; color != (types.Color{R: 20, G: 6, B: 5}) {
		t.Errorf("unclaimed tile has color %+v, want {R:20 G:6 B:5}", color)
	}
	if got.Field.ClaimedBy[capped] != "" {
		t.Errorf("tile with only one capped channel was claimed by %q", got.Field.ClaimedBy[capped])
	}
	if len(got.Field.ClaimedBy) != len(got.Field.Tiles) {
		t.Errorf("claim map has %d entries for %d tiles", len(got.Field.ClaimedBy), len(got.Field.Tiles))
	}

	got.Turn = enums.CardTypeCMY
	Matches[matchId] = got
	if err := PlayCard(matchId, "cmy-user", 0, claimed); err != nil {
		t.Fatalf("CMY PlayCard returned error: %v", err)
	}
	got = Matches[matchId]
	if color := got.Field.Tiles[claimed]; color != (types.Color{R: 20, G: 20, B: 20}) {
		t.Errorf("locked tile changed to %+v", color)
	}
	if got.PlayerCMY.Score != 0 {
		t.Errorf("CMY received score %d for a locked tile, want 0", got.PlayerCMY.Score)
	}
}

func TestPlayCardEndsMatchAtWinningScore(t *testing.T) {
	const matchId = "test_match_winner"
	pattern := make(map[types.Vec2]types.ModifierType, winningScore)
	tiles := make(map[types.Vec2]types.Color, winningScore)
	for i := 0; i < winningScore; i++ {
		tile := types.NewVec2(i%5, i/5)
		pattern[tile] = types.ModifierTypeFull
		tiles[tile] = types.NewColor(19, 19, 19)
	}

	match := types.Match{
		PlayerRGB: types.Player{
			UserId: "rgb-user",
			Money:  1,
			Hand: []types.Card{{
				Type:    enums.CardTypeRGB,
				Color:   types.NewColor(1, 1, 1),
				Cost:    1,
				Pattern: pattern,
			}},
		},
		PlayerCMY: types.Player{UserId: "cmy-user"},
		Turn:      enums.CardTypeRGB,
		Field:     types.Field{Tiles: tiles},
	}
	Matches[matchId] = match
	defer delete(Matches, matchId)

	if err := PlayCard(matchId, "rgb-user", 0, types.NewVec2(0, 0)); err != nil {
		t.Fatalf("winning PlayCard returned error: %v", err)
	}

	got := Matches[matchId]
	if !got.Finished || got.PlayerRGB.Score != winningScore {
		t.Errorf("match finished=%t, RGB score=%d; want finished at score %d", got.Finished, got.PlayerRGB.Score, winningScore)
	}
	if err := PlayCard(matchId, "cmy-user", 0, types.NewVec2(0, 0)); err == nil {
		t.Error("PlayCard accepted an action after the match ended")
	}
	if EndTurn(matchId, "cmy-user") {
		t.Error("EndTurn accepted an action after the match ended")
	}
}

func TestPlayCardRejectsAnchorOutsideFieldWithoutConsumingCard(t *testing.T) {
	const matchId = "test_invalid_card_tile"
	match := types.Match{
		PlayerRGB: types.Player{
			UserId: "rgb-user",
			Money:  1,
			Hand: []types.Card{{
				Type:    enums.CardTypeRGB,
				Color:   types.NewColor(1, 1, 1),
				Cost:    1,
				Pattern: map[types.Vec2]types.ModifierType{types.NewVec2(0, 0): types.ModifierTypeFull},
			}},
		},
		Turn:  enums.CardTypeRGB,
		Field: NewField(),
	}
	Matches[matchId] = match
	defer delete(Matches, matchId)

	if err := PlayCard(matchId, "rgb-user", 0, types.NewVec2(-1, 0)); err == nil {
		t.Fatal("PlayCard accepted an anchor outside the field")
	}

	got := Matches[matchId]
	if got.PlayerRGB.Money != 1 || len(got.PlayerRGB.Hand) != 1 {
		t.Errorf("invalid play consumed resources: money=%d hand=%d", got.PlayerRGB.Money, len(got.PlayerRGB.Hand))
	}
}
