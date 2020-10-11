package t3

import "testing"

//数组定义
func Test3_0(t *testing.T) {

	var a = [3]int{}
	t.Log(a[0], a[1])

	var b = [3]int{1, 2, 3}
	t.Log(b[1])

	c := [3]int{1, 2, 3}
	t.Log(c[2])

	d := [...]int{0, 1, 2, 3, 4}
	t.Log(len(d), d[3])

	e := [2][2]int{{1, 1}, {2, 2}}
	t.Log(e[1])
}

//数组的遍历
func Test3_1(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4}

	/* 类似java 的 for i */
	for i := 0; i < len(a); i++ {
		t.Log(a[i])
	}
	t.Log("===================")
	/* 类似java 的 for each */
	for idx, e := range a {
		t.Log(idx, e)
	}
	t.Log("===================")
	/* 类似java 的 for each，病区不关心第一个返回值*/
	for _, e := range a {
		t.Log(e)
	}
}

// 数组截取
func Test3_2(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5}
	//取前三个元素[0 1 2]
	a0 := a[:3]
	t.Log(a0)
	a1 := a[0:3]
	t.Log(a1)
	//取第三个之后的元素[3,4,5]
	a2 := a[3:]
	t.Log(a2)
}
