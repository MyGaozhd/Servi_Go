package t16

import (
	"fmt"
	"testing"
	"time"
)

func Test16_0(t *testing.T) {
	fun1 := func(op int) int {
		time.Sleep(time.Second * 2)
		return op * op
	}
	a := decorator(fun1)(2)
	t.Log(a)
}

//定义一个函数类型
type InnerFunc func(op int) int

//传入一个函数，并返回一个装饰过的函数
// inner 传入的函数
// Decorator 返回装饰过的函数
func decorator(inner InnerFunc) InnerFunc {

	return func(op int) int {
		start := time.Now()
		ret := inner(op)
		fmt.Print("耗时->", time.Since(start))
		return ret
	}
}
