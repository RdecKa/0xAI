package sort

// Quicksort sortsbthe given slice of elements in increasing order. Value of
// elements is obtained using value function. Sorting is done in-place.
func Quicksort(slice []interface{}, value func(interface{}) float64) {
	if len(slice) <= 1 {
		return
	}
	leftPart, rightPart := partition(slice, value)
	Quicksort(leftPart, value)
	Quicksort(rightPart, value)
}

func partition(slice []interface{}, value func(interface{}) float64) ([]interface{}, []interface{}) {
	pivotValue := value(slice[len(slice)/2])
	i, j := -1, len(slice)
	for i < j {
		for i++; value(slice[i]) < pivotValue; i++ {
		}
		for j--; value(slice[j]) > pivotValue; j-- {
		}
		if i < j {
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
	return slice[:i], slice[i:]
}
