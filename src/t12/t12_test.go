package t12

import (
	"fmt"
	"testing"
	"time"
)

func Test12_0(t *testing.T) {
	fun1 := func(op int) int {
		time.Sleep(time.Second * 2)
		return op * op
	}
	a := Decorator(fun1)(2)
	t.Log(a)
}

//传入一个函数，并返回一个装饰过的函数
// inner 传入的函数
// Decorator 返回装饰过的函数
func Decorator(inner func(op int) int) func(op int) int {

	return func(op int) int {
		start := time.Now()
		ret := inner(op)
		fmt.Print("耗时->", time.Since(start))
		return ret
	}
}
