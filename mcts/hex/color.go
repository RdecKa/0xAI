package hex

// -----------------
// |     Color     |
// -----------------

type color byte

// enum for players (colors)
const (
	None color = 0
	Red  color = 1
	Blue color = 2
)

func (c color) String() string {
	switch c {
	case Red:
		return "r"
	case Blue:
		return "b"
	default:
		return "."
	}
}

func getColorFromBits(bits uint64) color {
	if bits == 1 {
		return Red
	} else if bits == 2 {
		return Blue
	}
	return None
}
