package sort

import "testing"

func Sort(t *testing.T, input, expectedOutput []interface{}) {
	Quicksort(input, value)
	for i := range input {
		if input[i] != expectedOutput[i] {
			t.Fatalf("Expected %v, got %v", expectedOutput[i], input[i])
		}
	}
}

func value(el interface{}) float64 {
	return el.(float64)
}

func Test1(t *testing.T) {
	input := []interface{}{5.0, 2.0, 1.0, 5.0, 8.0, 0.0}
	expectedOutput := []interface{}{0.0, 1.0, 2.0, 5.0, 5.0, 8.0}
	Sort(t, input, expectedOutput)
}

func Test2(t *testing.T) {
	input := []interface{}{5.0}
	expectedOutput := []interface{}{5.0}
	Sort(t, input, expectedOutput)
}

func Test3(t *testing.T) {
	input := []interface{}{5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0}
	expectedOutput := []interface{}{5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0}
	Sort(t, input, expectedOutput)
}

func Test4(t *testing.T) {
	input := []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	expectedOutput := []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	Sort(t, input, expectedOutput)
}

func Test5(t *testing.T) {
	input := []interface{}{9.0, 8.0, 7.0, 6.0, 5.0, 4.0}
	expectedOutput := []interface{}{4.0, 5.0, 6.0, 7.0, 8.0, 9.0}
	Sort(t, input, expectedOutput)
}

func Test6(t *testing.T) {
	input := []interface{}{}
	expectedOutput := []interface{}{}
	Sort(t, input, expectedOutput)
}

func Test7(t *testing.T) {
	input := []interface{}{9.0, 8.0, 7.0, 10.0, 5.0, 4.0}
	expectedOutput := []interface{}{4.0, 5.0, 7.0, 8.0, 9.0, 10.0}
	Sort(t, input, expectedOutput)
}

func Test8(t *testing.T) {
	input := []interface{}{9.0, 8.0, 7.0, 0.0, 5.0, 4.0}
	expectedOutput := []interface{}{0.0, 4.0, 5.0, 7.0, 8.0, 9.0}
	Sort(t, input, expectedOutput)
}
