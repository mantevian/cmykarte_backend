package game

import (
	"testing"

	"mantevian.xyz/cmykarte_backend/common/enums"
	"mantevian.xyz/cmykarte_backend/common/types"
)

func TestEndTurnDrawsForBothPlayersAtRoundEnd(t *testing.T) {
	const matchId = "test_round_draws"
	Matches[matchId] = types.Match{
		PlayerRGB: types.Player{
			UserId: "rgb-user",
			Hand:   []types.Card{{Cost: 10}},
			Deck:   []types.Card{{Cost: 1}, {Cost: 2}, {Cost: 3}, {Cost: 4}},
		},
		PlayerCMY: types.Player{
			UserId: "cmy-user",
			Hand:   []types.Card{{Cost: 10}},
			Deck:   []types.Card{{Cost: 1}, {Cost: 2}, {Cost: 3}, {Cost: 4}},
		},
		Round: 1,
		Turn:  enums.CardTypeRGB,
	}
	defer delete(Matches, matchId)

	if !EndTurn(matchId, "rgb-user") {
		t.Fatal("RGB player's end turn was rejected")
	}

	match := Matches[matchId]
	if match.Round != 1 || match.Turn != enums.CardTypeCMY {
		t.Fatalf("after RGB turn, got round %d and turn %q", match.Round, match.Turn)
	}
	if len(match.PlayerRGB.Hand) != 1 || match.PlayerRGB.Hand[0].Cost != 10 {
		t.Fatalf("RGB hand changed before the round ended: %+v", match.PlayerRGB.Hand)
	}
	if len(match.PlayerCMY.Hand) != 1 || match.PlayerCMY.Hand[0].Cost != 10 {
		t.Fatalf("CMY hand changed before the round ended: %+v", match.PlayerCMY.Hand)
	}

	if !EndTurn(matchId, "cmy-user") {
		t.Fatal("CMY player's end turn was rejected")
	}

	match = Matches[matchId]
	if match.Round != 2 || match.Turn != enums.CardTypeRGB {
		t.Fatalf("after CMY turn, got round %d and turn %q", match.Round, match.Turn)
	}
	if len(match.PlayerRGB.Hand) != 2 {
		t.Fatalf("RGB hand has %d cards after round end, want 2", len(match.PlayerRGB.Hand))
	}
	if len(match.PlayerCMY.Hand) != 2 {
		t.Fatalf("CMY hand has %d cards after round end, want 2", len(match.PlayerCMY.Hand))
	}
	if got := []int{match.PlayerRGB.Hand[0].Cost, match.PlayerRGB.Hand[1].Cost}; got[0] != 10 || got[1] != 1 {
		t.Fatalf("RGB hand did not preserve existing cards and append one draw: %v", got)
	}
	if got := []int{match.PlayerCMY.Hand[0].Cost, match.PlayerCMY.Hand[1].Cost}; got[0] != 10 || got[1] != 1 {
		t.Fatalf("CMY hand did not preserve existing cards and append one draw: %v", got)
	}
}

func TestDrawStartingHandIncludesCostOneCard(t *testing.T) {
	for name, deck := range map[string][]types.Card{
		"RGB": CreateRGBDeck(),
		"CMY": CreateCMYDeck(),
	} {
		t.Run(name, func(t *testing.T) {
			player := types.Player{Deck: Shuffle(deck)}
			drawStartingHand(&player, defaultMatchRules())

			if len(player.Hand) != 4 {
				t.Fatalf("starting hand has %d cards, want 4", len(player.Hand))
			}
			hasCostOne := false
			for _, card := range player.Hand {
				if card.Cost == 1 {
					hasCostOne = true
				}
				if card.Cost > 3 {
					t.Errorf("starting hand contains cost %d card", card.Cost)
				}
			}
			if !hasCostOne {
				t.Error("starting hand has no cost 1 card")
			}
		})
	}
}

func TestEndTurnDrawsCardsForUpcomingRound(t *testing.T) {
	tests := []struct {
		name          string
		previousRound int
		wantCount     int
		maxCost       int
	}{
		{name: "round 5", previousRound: 4, wantCount: 1, maxCost: 3},
		{name: "round 6", previousRound: 5, wantCount: 2, maxCost: 5},
		{name: "round 10", previousRound: 9, wantCount: 2, maxCost: 5},
		{name: "round 11", previousRound: 10, wantCount: 2, maxCost: 7},
		{name: "round 20", previousRound: 19, wantCount: 2, maxCost: 7},
		{name: "round 21", previousRound: 20, wantCount: 3, maxCost: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchId := "test_round_filter_" + tt.name
			deck := []types.Card{
				{Cost: 1}, {Cost: 2}, {Cost: 3}, {Cost: 4},
				{Cost: 5}, {Cost: 6}, {Cost: 7}, {Cost: 8},
			}
			Matches[matchId] = types.Match{
				PlayerRGB: types.Player{UserId: "rgb-user", Deck: append([]types.Card(nil), deck...)},
				PlayerCMY: types.Player{UserId: "cmy-user", Deck: append([]types.Card(nil), deck...)},
				Round:     tt.previousRound,
				Turn:      enums.CardTypeCMY,
			}
			defer delete(Matches, matchId)

			if !EndTurn(matchId, "cmy-user") {
				t.Fatal("CMY player's end turn was rejected")
			}
			match := Matches[matchId]
			if match.Round != tt.previousRound+1 {
				t.Fatalf("round advanced to %d, want %d", match.Round, tt.previousRound+1)
			}
			for name, player := range map[string]types.Player{
				"RGB": match.PlayerRGB,
				"CMY": match.PlayerCMY,
			} {
				if len(player.Hand) != tt.wantCount {
					t.Errorf("%s drew %d cards, want %d", name, len(player.Hand), tt.wantCount)
				}
				for _, card := range player.Hand {
					if card.Cost > tt.maxCost {
						t.Errorf("%s drew cost %d during round %d", name, card.Cost, match.Round)
					}
				}
				if len(player.Hand) > 0 && player.Hand[0].Cost != 1 {
					t.Errorf("%s first drew cost %d during round %d, want cost 1 now that minimum costs are removed", name, player.Hand[0].Cost, match.Round)
				}
			}
		})
	}
}

func TestRoundDrawRecreatesAnEmptyDeck(t *testing.T) {
	player := types.Player{
		Type: enums.CardTypeRGB,
		Deck: []types.Card{{Type: enums.CardTypeRGB, Cost: 1}},
	}

	if err := DrawCardsForRound(&player, 2, 21); err != nil {
		t.Fatalf("draw cards from recreated deck: %v", err)
	}

	if len(player.Hand) != 2 {
		t.Fatalf("hand has %d cards, want 2", len(player.Hand))
	}
	if want := len(CreateRGBDeck()) - 1; len(player.Deck) != want {
		t.Fatalf("recreated RGB deck has %d cards remaining, want %d", len(player.Deck), want)
	}
	for _, card := range player.Hand {
		if card.Type != enums.CardTypeRGB {
			t.Errorf("recreated RGB deck produced card type %q", card.Type)
		}
	}
}

func TestDrawsRespectHandLimit(t *testing.T) {
	player := types.Player{
		Hand: make([]types.Card, 9),
		Deck: []types.Card{{Cost: 4}, {Cost: 5}, {Cost: 6}},
	}

	DrawCardsForRound(&player, 3, 21)
	if len(player.Hand) != 10 {
		t.Fatalf("hand has %d cards after draw, want limit of 10", len(player.Hand))
	}
	if len(player.Deck) != 2 {
		t.Fatalf("deck has %d cards after draw, want 2", len(player.Deck))
	}
}

func TestDecksContainCostSixThroughEightCards(t *testing.T) {
	for name, deck := range map[string][]types.Card{
		"RGB": CreateRGBDeck(),
		"CMY": CreateCMYDeck(),
	} {
		t.Run(name, func(t *testing.T) {
			if len(deck) <= 90 {
				t.Fatalf("deck has %d cards, want more than 90", len(deck))
			}
			eightCostCount := 0
			hasExtraPowerfulEightCostCard := false
			for _, card := range deck {
				if card.Cost == 8 {
					eightCostCount++
					if card.Color.R == 20 || card.Color.G == 20 || card.Color.B == 20 {
						hasExtraPowerfulEightCostCard = true
					}
				}
				if card.Cost == 1 && len(card.Pattern) == 1 &&
					(card.Color == types.NewColor(4, 0, 0) || card.Color == types.NewColor(2, 2, 2)) {
					t.Errorf("deck still contains an original 1-cost card: %+v", card)
				}
			}
			if eightCostCount < 6 {
				t.Errorf("deck has %d new 8-cost cards, want at least 6", eightCostCount)
			}
			if !hasExtraPowerfulEightCostCard {
				t.Error("deck has no extra-powerful 8-cost card")
			}
			for _, cost := range []int{6, 7, 8} {
				count := 0
				for _, card := range deck {
					if card.Cost == cost {
						count++
					}
				}
				if count < 6 {
					t.Errorf("deck has %d cards of cost %d, want at least 6", count, cost)
				}
			}
		})
	}
}
