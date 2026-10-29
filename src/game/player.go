package game

import "mantevian.xyz/cmykarte_backend/common/types"

func DrawCards(player *types.Player, count int) {
	if count <= 0 || len(player.Deck) == 0 {
		return
	}

	if count > len(player.Deck) {
		count = len(player.Deck)
	}

	player.Hand = append(player.Hand, player.Deck[:count]...)
	player.Deck = player.Deck[count:]
}
