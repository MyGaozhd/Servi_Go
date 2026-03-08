package t11

import (
	"math/rand"
	"testing"
)

/**
1、可以有多个返回值
2、所有参数都是值传递，slice/map/channel 会有传引用的错觉
3、函数可以作为变量
4、函数可以作为参数和返回值
*/

//1、多返回值
func Test11_0(t *testing.T) {
	a, b := fun1()
	t.Log(a, b)
}

func fun1() (int, int) {
	return rand.Intn(10), rand.Intn(20)
}

//2、可变参数
func Test11_1(t *testing.T) {
	t.Log(sum(1, 2, 3, 4))
}

func sum(op ...int) int {
	s := 0
	for _, item := range op {
		s += item
	}
	return s
}
