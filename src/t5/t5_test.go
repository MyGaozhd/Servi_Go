package t5

import (
	"testing"
)

//map 的用法
func Test5_0(t *testing.T) {
	m1 := map[string]int{"1": 1, "2": 2, "3": 3}
	t.Logf("len m1=%d", len(m1))
	t.Log(m1["1"])

	m2 := map[string]int{}
	m2["kkk"] = 1
	t.Logf("len m2=%d", len(m2))
	t.Log(m2["kkk"])

	m3 := make(map[string]int, 10)
	m3["hhh"] = 1
	t.Logf("len m3=%d", len(m3))
	t.Log(m3["hhh"])
}

//map 判断key
func Test5_1(t *testing.T) {
	m2 := map[string]int{}
	t.Log(m2["key"])
	m2["key"] = 0
	t.Log(m2["key"])

	m3 := make(map[string]int, 10)
	if val, ok := m3["key"]; ok {
		t.Log(val)
	} else {
		t.Log("不存在")
	}

	m3["key"] = 1
	if val, ok := m3["key"]; ok {
		t.Log(val)
	} else {
		t.Log("不存在")
	}
}

//map 遍历
func Test5_2(t *testing.T) {
	m1 := map[string]int{"1": 1, "2": 2, "3": 3}

	for k, v := range m1 {
		t.Log(k, v)
	}
}
