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

var patternRGB2 = map[types.Vec2]types.ModifierType{
	v(0, 0):  full,
	v(-1, 0): half,
	v(1, 0):  half,
	v(0, -1): half,
	v(0, 1):  half,
}

var patternRGB3 = map[types.Vec2]types.ModifierType{
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

var patternCMY2 = map[types.Vec2]types.ModifierType{
	v(0, 0): full,
	v(1, 0): full,
}

var patternCMY3 = map[types.Vec2]types.ModifierType{
	v(0, 0): full,
	v(1, 0): full,
	v(0, 1): full,
}

var patternCMY4 = map[types.Vec2]types.ModifierType{
	v(0, 0): full,
	v(1, 0): full,
	v(0, 1): full,
	v(1, 1): full,
}

func RGBCard(cost int, color types.Color, pattern map[types.Vec2]types.ModifierType) types.Card {
	return types.Card{
		Type:    enums.CardTypeRGB,
		Color:   color,
		Cost:    cost,
		Pattern: pattern,
	}
}

func CMYCard(cost int, color types.Color, pattern map[types.Vec2]types.ModifierType) types.Card {
	return types.Card{
		Type:    enums.CardTypeCMY,
		Color:   color,
		Cost:    cost,
		Pattern: pattern,
	}
}

func CreateRGBDeck() []types.Card {
	cards := []types.Card{
		// 1
		RGBCard(1, c(6, 0, 0), pattern1), RGBCard(1, c(6, 0, 0), pattern1), // 6
		RGBCard(1, c(0, 6, 0), pattern1), RGBCard(1, c(0, 6, 0), pattern1), // 6
		RGBCard(1, c(0, 0, 6), pattern1), RGBCard(1, c(0, 0, 6), pattern1), // 6

		RGBCard(1, c(3, 3, 0), pattern1), RGBCard(1, c(3, 3, 0), pattern1), // 6
		RGBCard(1, c(0, 3, 3), pattern1), RGBCard(1, c(0, 3, 3), pattern1), // 6
		RGBCard(1, c(3, 0, 3), pattern1), RGBCard(1, c(3, 0, 3), pattern1), // 6

		RGBCard(1, c(2, 0, 0), patternRGB2), RGBCard(1, c(2, 0, 0), patternRGB2), // 6
		RGBCard(1, c(0, 2, 0), patternRGB2), RGBCard(1, c(0, 2, 0), patternRGB2), // 6
		RGBCard(1, c(0, 0, 2), patternRGB2), RGBCard(1, c(0, 0, 2), patternRGB2), // 6

		// 2
		RGBCard(2, c(10, 0, 0), pattern1), RGBCard(2, c(10, 0, 0), pattern1), // 10
		RGBCard(2, c(0, 10, 0), pattern1), RGBCard(2, c(0, 10, 0), pattern1), // 10
		RGBCard(2, c(0, 0, 10), pattern1), RGBCard(2, c(0, 0, 10), pattern1), // 10

		RGBCard(2, c(6, 6, 0), pattern1), RGBCard(2, c(6, 6, 0), pattern1), // 12
		RGBCard(2, c(0, 6, 6), pattern1), RGBCard(2, c(0, 6, 6), pattern1), // 12
		RGBCard(2, c(6, 0, 6), pattern1), RGBCard(2, c(6, 0, 6), pattern1), // 12

		RGBCard(2, c(4, 0, 0), patternRGB2), RGBCard(2, c(4, 0, 0), patternRGB2), // 12
		RGBCard(2, c(0, 4, 0), patternRGB2), RGBCard(2, c(0, 4, 0), patternRGB2), // 12
		RGBCard(2, c(0, 0, 4), patternRGB2), RGBCard(2, c(0, 0, 4), patternRGB2), // 12

		RGBCard(2, c(2, 2, 0), patternRGB2), RGBCard(2, c(2, 2, 0), patternRGB2), // 12
		RGBCard(2, c(0, 2, 2), patternRGB2), RGBCard(2, c(0, 2, 2), patternRGB2), // 12
		RGBCard(2, c(2, 0, 2), patternRGB2), RGBCard(2, c(2, 0, 2), patternRGB2), // 12

		RGBCard(2, c(4, 4, 4), pattern1), // 12

		// 3
		RGBCard(3, c(8, 8, 0), pattern1), // 16
		RGBCard(3, c(0, 8, 8), pattern1), // 16
		RGBCard(3, c(8, 0, 8), pattern1), // 16

		RGBCard(3, c(4, 6, 6), pattern1), // 16
		RGBCard(3, c(6, 4, 6), pattern1), // 16
		RGBCard(3, c(6, 6, 4), pattern1), // 16

		RGBCard(3, c(6, 0, 0), patternRGB2), // 18
		RGBCard(3, c(0, 6, 0), patternRGB2), // 18
		RGBCard(3, c(0, 0, 6), patternRGB2), // 18

		RGBCard(3, c(4, 2, 0), patternRGB2), // 18
		RGBCard(3, c(2, 4, 0), patternRGB2), // 18
		RGBCard(3, c(0, 4, 2), patternRGB2), // 18
		RGBCard(3, c(0, 2, 4), patternRGB2), // 18
		RGBCard(3, c(4, 0, 2), patternRGB2), // 18
		RGBCard(3, c(2, 0, 4), patternRGB2), // 18

		RGBCard(3, c(4, 0, 0), patternRGB3), // 20
		RGBCard(3, c(0, 4, 0), patternRGB3), // 20
		RGBCard(3, c(0, 0, 4), patternRGB3), // 20

		RGBCard(3, c(5, 5, 5), pattern1), // 15

		// 4
		RGBCard(4, c(8, 0, 0), patternRGB2), // 24
		RGBCard(4, c(0, 8, 0), patternRGB2), // 24
		RGBCard(4, c(0, 0, 8), patternRGB2), // 24

		RGBCard(4, c(6, 3, 0), patternRGB2), // 27
		RGBCard(4, c(3, 6, 0), patternRGB2), // 27
		RGBCard(4, c(0, 6, 3), patternRGB2), // 27
		RGBCard(4, c(0, 3, 6), patternRGB2), // 27
		RGBCard(4, c(6, 0, 3), patternRGB2), // 27
		RGBCard(4, c(3, 0, 6), patternRGB2), // 27

		RGBCard(4, c(6, 0, 0), patternRGB3), // 26
		RGBCard(4, c(0, 6, 0), patternRGB3), // 26
		RGBCard(4, c(0, 0, 6), patternRGB3), // 26

		// 6
		RGBCard(6, c(5, 5, 5), patternRGB3), // 53

		// 5
		RGBCard(5, c(8, 0, 0), patternRGB3), // 32
		RGBCard(5, c(0, 8, 0), patternRGB3), // 32
		RGBCard(5, c(0, 0, 8), patternRGB3), // 32
		RGBCard(5, c(4, 4, 0), patternRGB3), // 32
		RGBCard(5, c(0, 4, 4), patternRGB3), // 32
		RGBCard(5, c(4, 0, 4), patternRGB3), // 32

		// 6
		RGBCard(6, c(10, 0, 0), patternRGB3), // 38
		RGBCard(6, c(0, 10, 0), patternRGB3), // 38
		RGBCard(6, c(0, 0, 10), patternRGB3), // 38
		RGBCard(6, c(6, 4, 0), patternRGB3),  // 38
		RGBCard(6, c(0, 6, 4), patternRGB3),  // 38
		RGBCard(6, c(4, 0, 6), patternRGB3),  // 38

		// 7
		RGBCard(7, c(12, 0, 0), patternRGB3), // 44
		RGBCard(7, c(0, 12, 0), patternRGB3), // 44
		RGBCard(7, c(0, 0, 12), patternRGB3), // 44
		RGBCard(7, c(8, 4, 0), patternRGB3),  // 44
		RGBCard(7, c(0, 8, 4), patternRGB3),  // 44
		RGBCard(7, c(4, 0, 8), patternRGB3),  // 44

		// 8
		RGBCard(8, c(12, 12, 12), patternRGB3),
		RGBCard(8, c(8, 8, 0), patternRGB3),
		RGBCard(8, c(0, 8, 8), patternRGB3),
		RGBCard(8, c(8, 0, 8), patternRGB3),
	}

	return cards
}

func CreateCMYDeck() []types.Card {
	cards := []types.Card{
		// 1
		CMYCard(1, c(6, 0, 0), pattern1), CMYCard(1, c(6, 0, 0), pattern1), // 6
		CMYCard(1, c(0, 6, 0), pattern1), CMYCard(1, c(0, 6, 0), pattern1), // 6
		CMYCard(1, c(0, 0, 6), pattern1), CMYCard(1, c(0, 0, 6), pattern1), // 6

		CMYCard(1, c(3, 3, 0), pattern1), CMYCard(1, c(3, 3, 0), pattern1), // 6
		CMYCard(1, c(0, 3, 3), pattern1), CMYCard(1, c(0, 3, 3), pattern1), // 6
		CMYCard(1, c(3, 0, 3), pattern1), CMYCard(1, c(3, 0, 3), pattern1), // 6

		CMYCard(1, c(4, 0, 0), patternCMY2), CMYCard(1, c(4, 0, 0), patternCMY2), // 8
		CMYCard(1, c(0, 4, 0), patternCMY2), CMYCard(1, c(0, 4, 0), patternCMY2), // 8
		CMYCard(1, c(0, 0, 4), patternCMY2), CMYCard(1, c(0, 0, 4), patternCMY2), // 6

		// 2
		CMYCard(2, c(10, 0, 0), pattern1), CMYCard(2, c(10, 0, 0), pattern1), // 10
		CMYCard(2, c(0, 10, 0), pattern1), CMYCard(2, c(0, 10, 0), pattern1), // 10
		CMYCard(2, c(0, 0, 10), pattern1), CMYCard(2, c(0, 0, 10), pattern1), // 10

		CMYCard(2, c(6, 6, 0), pattern1), CMYCard(2, c(6, 6, 0), pattern1), // 12
		CMYCard(2, c(0, 6, 6), pattern1), CMYCard(2, c(0, 6, 6), pattern1), // 12
		CMYCard(2, c(6, 0, 6), pattern1), CMYCard(2, c(6, 0, 6), pattern1), // 12

		CMYCard(2, c(6, 0, 0), patternCMY2), CMYCard(2, c(6, 0, 0), patternCMY2), // 12
		CMYCard(2, c(0, 6, 0), patternCMY2), CMYCard(2, c(0, 6, 0), patternCMY2), // 12
		CMYCard(2, c(0, 0, 6), patternCMY2), CMYCard(2, c(0, 0, 6), patternCMY2), // 12

		CMYCard(2, c(4, 4, 0), patternCMY2), CMYCard(2, c(4, 4, 0), patternCMY2), // 16
		CMYCard(2, c(0, 4, 4), patternCMY2), CMYCard(2, c(0, 4, 4), patternCMY2), // 16
		CMYCard(2, c(4, 0, 4), patternCMY2), CMYCard(2, c(4, 0, 4), patternCMY2), // 16

		CMYCard(2, c(4, 4, 4), pattern1), // 12

		// 3
		CMYCard(3, c(8, 8, 0), pattern1), // 16
		CMYCard(3, c(0, 8, 8), pattern1), // 16
		CMYCard(3, c(8, 0, 8), pattern1), // 16

		CMYCard(3, c(4, 6, 6), pattern1), // 16
		CMYCard(3, c(6, 4, 6), pattern1), // 16
		CMYCard(3, c(6, 6, 4), pattern1), // 16

		CMYCard(3, c(8, 0, 0), patternCMY2), // 16
		CMYCard(3, c(0, 8, 0), patternCMY2), // 16
		CMYCard(3, c(0, 0, 8), patternCMY2), // 16

		CMYCard(3, c(6, 4, 0), patternCMY2), // 20
		CMYCard(3, c(4, 6, 0), patternCMY2), // 20
		CMYCard(3, c(0, 6, 4), patternCMY2), // 20
		CMYCard(3, c(0, 4, 6), patternCMY2), // 20
		CMYCard(3, c(6, 0, 4), patternCMY2), // 20
		CMYCard(3, c(4, 0, 6), patternCMY2), // 20

		CMYCard(3, c(6, 0, 0), patternCMY3), // 18
		CMYCard(3, c(0, 6, 0), patternCMY3), // 18
		CMYCard(3, c(0, 0, 6), patternCMY3), // 18

		CMYCard(3, c(2, 2, 2), patternCMY4), // 24

		// 4
		CMYCard(4, c(8, 0, 0), patternCMY3), // 24
		CMYCard(4, c(0, 8, 0), patternCMY3), // 24
		CMYCard(4, c(0, 0, 8), patternCMY3), // 24

		CMYCard(4, c(6, 3, 0), patternCMY3), // 27
		CMYCard(4, c(3, 6, 0), patternCMY3), // 27
		CMYCard(4, c(0, 6, 3), patternCMY3), // 27
		CMYCard(4, c(0, 3, 6), patternCMY3), // 27
		CMYCard(4, c(6, 0, 3), patternCMY3), // 27
		CMYCard(4, c(3, 0, 6), patternCMY3), // 27

		CMYCard(4, c(6, 0, 0), patternCMY3), // 26
		CMYCard(4, c(0, 6, 0), patternCMY3), // 26
		CMYCard(4, c(0, 0, 6), patternCMY3), // 26

		// 6
		CMYCard(6, c(4, 4, 4), patternCMY4), // 48

		// 5
		CMYCard(5, c(8, 0, 0), patternCMY4), // 32
		CMYCard(5, c(0, 8, 0), patternCMY4), // 32
		CMYCard(5, c(0, 0, 8), patternCMY4), // 32
		CMYCard(5, c(4, 4, 0), patternCMY4), // 32
		CMYCard(5, c(0, 4, 4), patternCMY4), // 32
		CMYCard(5, c(4, 0, 4), patternCMY4), // 32

		// 6
		CMYCard(6, c(10, 0, 0), patternCMY4), // 40
		CMYCard(6, c(0, 10, 0), patternCMY4), // 40
		CMYCard(6, c(0, 0, 10), patternCMY4), // 40
		CMYCard(6, c(6, 4, 0), patternCMY4),  // 40
		CMYCard(6, c(0, 6, 4), patternCMY4),  // 40
		CMYCard(6, c(4, 0, 6), patternCMY4),  // 40

		// 7
		CMYCard(7, c(12, 0, 0), patternCMY4), // 48
		CMYCard(7, c(0, 12, 0), patternCMY4), // 48
		CMYCard(7, c(0, 0, 12), patternCMY4), // 48
		CMYCard(7, c(8, 4, 0), patternCMY4),  // 48
		CMYCard(7, c(0, 8, 4), patternCMY4),  // 48
		CMYCard(7, c(4, 0, 8), patternCMY4),  // 48

		// 8
		CMYCard(8, c(8, 8, 8), patternCMY4),
		CMYCard(8, c(12, 12, 0), patternCMY4),
		CMYCard(8, c(0, 12, 12), patternCMY4),
		CMYCard(8, c(12, 0, 12), patternCMY4),
	}

	return cards
}

func Shuffle(cards []types.Card) []types.Card {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}
