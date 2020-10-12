package t8

import "testing"

/**
  实现set
1、添加元素
2、删除元素
3、判断元素是否存在
4、元素个数
*/

func Test8_0(t *testing.T) {

	set := map[int]bool{}
	set[1] = true

	if set[1] {
		t.Log("已经存在")
	} else {
		t.Log("不存在")
	}

	t.Log(len(set))

	delete(set, 1)

	if set[1] {
		t.Log("已经存在")
	} else {
		t.Log("不存在")
	}
}
