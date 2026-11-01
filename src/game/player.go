package game

import (
	"errors"

	"mantevian.xyz/cmykarte_backend/common/enums"
	"mantevian.xyz/cmykarte_backend/common/types"
)

func DrawCards(player *types.Player, count int) error {
	if count <= 0 || len(player.Hand) >= defaultHandLimit {
		return nil
	}

	if count > defaultHandLimit-len(player.Hand) {
		count = defaultHandLimit - len(player.Hand)
	}

	for range count {
		if len(player.Deck) == 0 {
			if err := recreateDeck(player); err != nil {
				return err
			}
		}
		player.Hand = append(player.Hand, player.Deck[0])
		player.Deck = player.Deck[1:]
	}
	return nil
}

func drawStartingHand(player *types.Player, rules types.MatchRules) error {
	for i, card := range player.Deck {
		if card.Cost == rules.GuaranteedStartingCost {
			player.Hand = append(player.Hand, card)
			player.Deck = append(player.Deck[:i], player.Deck[i+1:]...)
			break
		}
	}
	return drawCardsForRound(player, rules.StartingHandSize-len(player.Hand), 1, rules)
}

func DrawCardsForRound(player *types.Player, count int, round int) error {
	return drawCardsForRound(player, count, round, defaultMatchRules())
}

func drawCardsForRound(player *types.Player, count int, round int, rules types.MatchRules) error {
	for range count {
		if len(player.Hand) >= rules.HandLimit {
			return nil
		}
		if len(player.Deck) == 0 {
			if err := recreateDeck(player); err != nil {
				return err
			}
		}

		cardIndex := -1
		for i, card := range player.Deck {
			if canDrawCardAtRound(card.Cost, round, rules) {
				cardIndex = i
				break
			}
		}
		if cardIndex == -1 {
			return nil
		}

		player.Hand = append(player.Hand, player.Deck[cardIndex])
		player.Deck = append(player.Deck[:cardIndex], player.Deck[cardIndex+1:]...)
	}
	return nil
}

func canDrawCardAtRound(cost int, round int, rules types.MatchRules) bool {
	for _, rule := range rules.RoundDraws {
		if round < rule.FromRound || (rule.ToRound != nil && round > *rule.ToRound) {
			continue
		}
		if rule.MaxCost != nil && cost > *rule.MaxCost {
			return false
		}
		return true
	}
	return false
}

func recreateDeck(player *types.Player) error {
	switch player.Type {
	case enums.CardTypeRGB:
		player.Deck = Shuffle(CreateRGBDeck())
	case enums.CardTypeCMY:
		player.Deck = Shuffle(CreateCMYDeck())
	default:
		return errors.New("player has an invalid card type")
	}
	return nil
}
