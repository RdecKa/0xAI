package hex

// -----------------
// |     Color     |
// -----------------

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

func getColorFromBits(bits uint64) Color {
	if bits == 1 {
		return Red
	} else if bits == 2 {
		return Blue
	}
	return None
}

func (c Color) opponent() Color {
	if c != Red && c != Blue {
		return None
	}

	if c == Red {
		return Blue
	}
	return Red
}
