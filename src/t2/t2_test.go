package t2

import "testing"

//for循环测试
func Test2_0(t *testing.T) {
	a := 0
	/* 相当于java 的 while(a<5) */
	for a < 5 {
		a++
		t.Log(a)
	}
}
