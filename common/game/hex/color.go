package hex

// -----------------
// |     Color     |
// -----------------

// Color represents a player in a game
type Color byte

// enum for players (colors)
const (
	None Color = 0
	Red  Color = 1
	Blue Color = 2
)

func (c Color) String() string {
	switch c {
	case Red:
		return "r"
	case Blue:
		return "b"
	default:
		return "."
	}
}

// GetColorFromBits reads the color from two bits
func GetColorFromBits(bits uint32) Color {
	if bits == 1 {
		return Red
	} else if bits == 2 {
		return Blue
	}
	return None
}

// Opponent returns the opponent of the given color
func (c Color) Opponent() Color {
	if c != Red && c != Blue {
		return None
	}

	if c == Red {
		return Blue
	}
	return Red
}
