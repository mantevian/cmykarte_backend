package game

import "mantevian.xyz/cmykarte_backend/common/types"

func DrawCards(player *types.Player, count int) {
	if count <= 0 || len(player.Deck) == 0 || len(player.Hand) >= defaultHandLimit {
		return
	}

	if count > len(player.Deck) {
		count = len(player.Deck)
	}
	if count > defaultHandLimit-len(player.Hand) {
		count = defaultHandLimit - len(player.Hand)
	}

	player.Hand = append(player.Hand, player.Deck[:count]...)
	player.Deck = player.Deck[count:]
}

func drawStartingHand(player *types.Player, rules types.MatchRules) {
	for i, card := range player.Deck {
		if card.Cost == rules.GuaranteedStartingCost {
			player.Hand = append(player.Hand, card)
			player.Deck = append(player.Deck[:i], player.Deck[i+1:]...)
			break
		}
	}
	drawCardsForRound(player, rules.StartingHandSize-len(player.Hand), 1, rules)
}

func DrawCardsForRound(player *types.Player, count int, round int) {
	drawCardsForRound(player, count, round, defaultMatchRules())
}

func drawCardsForRound(player *types.Player, count int, round int, rules types.MatchRules) {
	for range count {
		if len(player.Hand) >= rules.HandLimit {
			return
		}

		cardIndex := -1
		for i, card := range player.Deck {
			if canDrawCardAtRound(card.Cost, round, rules) {
				cardIndex = i
				break
			}
		}
		if cardIndex == -1 {
			return
		}

		player.Hand = append(player.Hand, player.Deck[cardIndex])
		player.Deck = append(player.Deck[:cardIndex], player.Deck[cardIndex+1:]...)
	}
}

func canDrawCardAtRound(cost int, round int, rules types.MatchRules) bool {
	for _, rule := range rules.RoundDraws {
		if round < rule.FromRound || (rule.ToRound != nil && round > *rule.ToRound) {
			continue
		}
		if cost < rule.MinCost || (rule.MaxCost != nil && cost > *rule.MaxCost) {
			return false
		}
		return true
	}
	return false
}
