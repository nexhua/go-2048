package game

import "image/color"

var (
	BEIGE       = color.NRGBA{0xfa, 0xf8, 0xef, 0xff}
	LIGHT_BROWN = color.NRGBA{0xcd, 0xc1, 0xb4, 0xff}

	// Tile colors
	LIGHT_TAN    = color.NRGBA{0xee, 0xe4, 0xda, 0xff} // 2
	TAN          = color.NRGBA{0xed, 0xe0, 0xc8, 0xff} // 4
	ORANGE_LIGHT = color.NRGBA{0xf2, 0xb1, 0x79, 0xff} // 8
	ORANGE       = color.NRGBA{0xf5, 0x95, 0x63, 0xff} // 16
	RED_ORANGE   = color.NRGBA{0xf6, 0x7c, 0x5f, 0xff} // 32
	RED          = color.NRGBA{0xf6, 0x5e, 0x3b, 0xff} // 64
	YELLOW_LIGHT = color.NRGBA{0xed, 0xcf, 0x72, 0xff} // 128
	YELLOW       = color.NRGBA{0xed, 0xcc, 0x61, 0xff} // 256
	GOLD_LIGHT   = color.NRGBA{0xed, 0xc8, 0x50, 0xff} // 512
	GOLD         = color.NRGBA{0xed, 0xc5, 0x3f, 0xff} // 1024
	GOLD_DEEP    = color.NRGBA{0xed, 0xc2, 0x2e, 0xff} // 2048

	// For larger values
	DARK_GRAY = color.NRGBA{0x3c, 0x3a, 0x32, 0xff}

	// Text colors
	TEXT_DARK  = color.NRGBA{0x77, 0x6e, 0x65, 0xff} // for 2, 4
	TEXT_LIGHT = color.NRGBA{0xf9, 0xf6, 0xf2, 0xff} // for others
)

// TileColors maps tile values to colors
var TileColors = map[int]color.NRGBA{
	2:    LIGHT_TAN,
	4:    TAN,
	8:    ORANGE_LIGHT,
	16:   ORANGE,
	32:   RED_ORANGE,
	64:   RED,
	128:  YELLOW_LIGHT,
	256:  YELLOW,
	512:  GOLD_LIGHT,
	1024: GOLD,
	2048: GOLD_DEEP,
}

func GetColor(val int) color.NRGBA {
	if col, ok := TileColors[val]; ok {
		return col
	}

	return DARK_GRAY
}
