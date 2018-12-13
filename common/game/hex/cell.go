package hex

// DirectlyConnected checks whether two cells are directly connected (they share
// one side)
func DirectlyConnected(c1, c2 [2]int) bool {
	x1, y1 := c1[0], c1[1]
	x2, y2 := c2[0], c2[1]
	if x1 == x2 && abs(y1-y2) == 1 {
		return true
	}
	if y1 == y2 && abs(x1-x2) == 1 {
		return true
	}
	if abs(x1-x2) == 1 && abs(y1-y2) == 1 && (x1-x2)+(y1-y2) == 0 {
		return true
	}
	return false
}

// GetTwoCellsBewteen returns two cells betwen (c1[0], c1[1]) and (c2[0], c2[1])
func GetTwoCellsBewteen(c1, c2 [2]int) [2][2]int {
	x1, y1, x2, y2 := c1[0], c1[1], c2[0], c2[1]
	var nx1, nx2, ny1, ny2 int
	if (x1-x2)%2 == 0 {
		nx1 = (x1 + x2) / 2
		nx2 = nx1
	} else {
		nx1 = (x1 + x2) / 2
		nx2 = nx1 + 1
	}

	if (y1-y2)%2 == 0 {
		ny1 = (y1 + y2) / 2
		ny2 = ny1
	} else {
		ny2 = (y1 + y2) / 2
		ny1 = ny2 + 1
	}

	return [2][2]int{[2]int{nx1, ny1}, [2]int{nx2, ny2}}
}
