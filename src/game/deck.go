package game

import (
	"math/rand"

	"mantevian.xyz/cmykarte_backend/common/enums"
	"mantevian.xyz/cmykarte_backend/common/types"
)

var v = types.NewVec2
var c = types.NewColor

var full = types.ModifierTypeFull
var half = types.ModifierTypeHalf
var one = types.ModifierTypeOne

var pattern1 = map[types.Vec2]types.ModifierType{
	v(0, 0): full,
}

var pattern2 = map[types.Vec2]types.ModifierType{
	v(0, 0):  full,
	v(-1, 0): half,
	v(1, 0):  half,
	v(0, -1): half,
	v(0, 1):  half,
}

var pattern3 = map[types.Vec2]types.ModifierType{
	v(0, 0):   full,
	v(-1, 0):  half,
	v(1, 0):   half,
	v(0, -1):  half,
	v(0, 1):   half,
	v(2, 0):   one,
	v(1, 1):   one,
	v(0, 2):   one,
	v(-1, 1):  one,
	v(-2, 0):  one,
	v(-1, -1): one,
	v(0, -2):  one,
	v(1, -1):  one,
}

func RGBCard(cost int, color types.Color, pattern map[types.Vec2]types.ModifierType) types.Card {
	return types.Card{
		Type:    enums.CardTypeRGB,
		Color:   color,
		Cost:    cost,
		Pattern: pattern,
	}
}

func CreateRGBDeck() []types.Card {
	cards := [90]types.Card{
		// 1
		RGBCard(1, c(4, 0, 0), pattern1), RGBCard(1, c(4, 0, 0), pattern1), // 4
		RGBCard(1, c(0, 4, 0), pattern1), RGBCard(1, c(0, 4, 0), pattern1), // 4
		RGBCard(1, c(0, 0, 4), pattern1), RGBCard(1, c(0, 0, 4), pattern1), // 4

		RGBCard(1, c(2, 2, 0), pattern1), RGBCard(1, c(2, 2, 0), pattern1), // 4
		RGBCard(1, c(0, 2, 2), pattern1), RGBCard(1, c(0, 2, 2), pattern1), // 4
		RGBCard(1, c(2, 0, 2), pattern1), RGBCard(1, c(2, 0, 2), pattern1), // 4

		RGBCard(1, c(1, 1, 1), pattern1), RGBCard(1, c(1, 1, 1), pattern1), // 3

		// 2
		RGBCard(2, c(6, 0, 0), pattern1), RGBCard(2, c(6, 0, 0), pattern1), // 6
		RGBCard(2, c(0, 6, 0), pattern1), RGBCard(2, c(0, 6, 0), pattern1), // 6
		RGBCard(2, c(0, 0, 6), pattern1), RGBCard(2, c(0, 0, 6), pattern1), // 6

		RGBCard(2, c(3, 3, 0), pattern1), RGBCard(2, c(3, 3, 0), pattern1), // 6
		RGBCard(2, c(0, 3, 3), pattern1), RGBCard(2, c(0, 3, 3), pattern1), // 6
		RGBCard(2, c(3, 0, 3), pattern1), RGBCard(2, c(3, 0, 3), pattern1), // 6

		RGBCard(2, c(2, 0, 0), pattern2), RGBCard(2, c(2, 0, 0), pattern2), // 6
		RGBCard(2, c(0, 2, 0), pattern2), RGBCard(2, c(0, 2, 0), pattern2), // 6
		RGBCard(2, c(0, 0, 2), pattern2), RGBCard(2, c(0, 0, 2), pattern2), // 6

		RGBCard(1, c(2, 2, 2), pattern1), // 6

		// 3
		RGBCard(3, c(10, 0, 0), pattern1), RGBCard(3, c(10, 0, 0), pattern1), // 10
		RGBCard(3, c(0, 10, 0), pattern1), RGBCard(3, c(0, 10, 0), pattern1), // 10
		RGBCard(3, c(0, 0, 10), pattern1), RGBCard(3, c(0, 0, 10), pattern1), // 10

		RGBCard(3, c(6, 6, 0), pattern1), RGBCard(3, c(6, 6, 0), pattern1), // 12
		RGBCard(3, c(0, 6, 6), pattern1), RGBCard(3, c(0, 6, 6), pattern1), // 12
		RGBCard(3, c(6, 0, 6), pattern1), RGBCard(3, c(6, 0, 6), pattern1), // 12

		RGBCard(3, c(4, 0, 0), pattern2), RGBCard(3, c(4, 0, 0), pattern2), // 12
		RGBCard(3, c(0, 4, 0), pattern2), RGBCard(3, c(0, 4, 0), pattern2), // 12
		RGBCard(3, c(0, 0, 4), pattern2), RGBCard(3, c(0, 0, 4), pattern2), // 12

		RGBCard(3, c(2, 2, 0), pattern2), RGBCard(3, c(2, 2, 0), pattern2), // 12
		RGBCard(3, c(0, 2, 2), pattern2), RGBCard(3, c(0, 2, 2), pattern2), // 12
		RGBCard(3, c(2, 0, 2), pattern2), RGBCard(3, c(2, 0, 2), pattern2), // 12

		RGBCard(3, c(4, 4, 4), pattern1), // 12

		// 4
		RGBCard(4, c(8, 8, 0), pattern1), // 16
		RGBCard(4, c(0, 8, 8), pattern1), // 16
		RGBCard(4, c(8, 0, 8), pattern1), // 16

		RGBCard(4, c(4, 6, 6), pattern1), // 16
		RGBCard(4, c(6, 4, 6), pattern1), // 16
		RGBCard(4, c(6, 6, 4), pattern1), // 16

		RGBCard(4, c(6, 0, 0), pattern2), // 18
		RGBCard(4, c(0, 6, 0), pattern2), // 18
		RGBCard(4, c(0, 0, 6), pattern2), // 18

		RGBCard(4, c(4, 2, 0), pattern2), // 18
		RGBCard(4, c(2, 4, 0), pattern2), // 18
		RGBCard(4, c(0, 4, 2), pattern2), // 18
		RGBCard(4, c(0, 2, 4), pattern2), // 18
		RGBCard(4, c(4, 0, 2), pattern2), // 18
		RGBCard(4, c(2, 0, 4), pattern2), // 18

		RGBCard(4, c(4, 0, 0), pattern3), // 20
		RGBCard(4, c(0, 4, 0), pattern3), // 20
		RGBCard(4, c(0, 0, 4), pattern3), // 20

		RGBCard(4, c(5, 5, 5), pattern1), // 15

		// 5
		RGBCard(5, c(8, 0, 0), pattern2), // 24
		RGBCard(5, c(0, 8, 0), pattern2), // 24
		RGBCard(5, c(0, 0, 8), pattern2), // 24

		RGBCard(5, c(6, 3, 0), pattern2), // 27
		RGBCard(5, c(3, 6, 0), pattern2), // 27
		RGBCard(5, c(0, 6, 3), pattern2), // 27
		RGBCard(5, c(0, 3, 6), pattern2), // 27
		RGBCard(5, c(6, 0, 3), pattern2), // 27
		RGBCard(5, c(3, 0, 6), pattern2), // 27

		RGBCard(5, c(6, 0, 0), pattern3), // 26
		RGBCard(5, c(0, 6, 0), pattern3), // 26
		RGBCard(5, c(0, 0, 6), pattern3), // 26

		RGBCard(5, c(8, 8, 8), pattern1), // 24
	}

	return cards[:]
}

func Shuffle(cards []types.Card) []types.Card {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}
