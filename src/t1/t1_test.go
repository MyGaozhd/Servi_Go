package t1

import "testing"

//清零运算符 &^
// 如果右边是1 不管左边是多少最终结果都是0
// 如果右边是0  左边是多少结果就是多少
func Test1_0(t *testing.T) {
	a := 1 &^ 1
	b := 1 &^ 0
	c := 0 &^ 1
	d := 0 &^ 0
	t.Log(a, b, c, d)
}

//与上面的写法效果相同。所以 &^ 运算符 相当于一种语法糖
func Test1_1(t *testing.T) {
	a := 1 & (^1)
	b := 1 & (^0)
	c := 0 & (^1)
	d := 0 & (^0)
	t.Log(a, b, c, d)
}
