package arg

import (
	"reflect"
	"testing"
)

func TestRemoveIndex(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wanted := []int{0, 1, 2, 3, 4, 8, 9, 10}
	result := RemoveIndex(input, 5)
	result = RemoveIndex(result, 5)
	result = RemoveIndex(result, 5)

	if !reflect.DeepEqual(result, wanted) {
		t.Errorf("wanted '%+v' got '%+v'", wanted, result)
	}
}
