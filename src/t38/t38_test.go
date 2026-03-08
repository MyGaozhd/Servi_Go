package t38

import (
	"reflect"
	"testing"
)

//map的比较，切片的比较
func Test38_0(t *testing.T) {
	a := map[int]string{1: "a", 2: "b"}
	b := map[int]string{1: "a", 2: "b"}
	t.Log(reflect.DeepEqual(a, b))

	a1 := []int{1, 2, 3}
	a2 := []int{1, 2, 3}
	a3 := []int{3, 2, 1}
	t.Log(reflect.DeepEqual(a1, a2))
	t.Log(reflect.DeepEqual(a1, a3))
}
