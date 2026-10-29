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
			Deck:   []types.Card{{Cost: 11}, {Cost: 12}},
		},
		PlayerCMY: types.Player{
			UserId: "cmy-user",
			Hand:   []types.Card{{Cost: 20}},
			Deck:   []types.Card{{Cost: 21}, {Cost: 22}},
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
	if len(match.PlayerCMY.Hand) != 1 || match.PlayerCMY.Hand[0].Cost != 20 {
		t.Fatalf("CMY hand changed before the round ended: %+v", match.PlayerCMY.Hand)
	}

	if !EndTurn(matchId, "cmy-user") {
		t.Fatal("CMY player's end turn was rejected")
	}

	match = Matches[matchId]
	if match.Round != 2 || match.Turn != enums.CardTypeRGB {
		t.Fatalf("after CMY turn, got round %d and turn %q", match.Round, match.Turn)
	}
	if len(match.PlayerRGB.Hand) != 3 {
		t.Fatalf("RGB hand has %d cards after round end, want 3", len(match.PlayerRGB.Hand))
	}
	if len(match.PlayerCMY.Hand) != 3 {
		t.Fatalf("CMY hand has %d cards after round end, want 3", len(match.PlayerCMY.Hand))
	}
	if got := []int{match.PlayerRGB.Hand[0].Cost, match.PlayerRGB.Hand[1].Cost, match.PlayerRGB.Hand[2].Cost}; got[0] != 10 || got[1] != 11 || got[2] != 12 {
		t.Fatalf("RGB hand did not preserve existing cards and append two draws: %v", got)
	}
	if got := []int{match.PlayerCMY.Hand[0].Cost, match.PlayerCMY.Hand[1].Cost, match.PlayerCMY.Hand[2].Cost}; got[0] != 20 || got[1] != 21 || got[2] != 22 {
		t.Fatalf("CMY hand did not preserve existing cards and append two draws: %v", got)
	}
}
