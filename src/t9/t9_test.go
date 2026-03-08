package t9

import "testing"

//字符串使用
func Test9_0(t *testing.T) {
	var s string
	t.Log(s)
	t.Log(len(s))

	s = "hello"
	t.Log("==============================")
	t.Log(s)
	t.Log(len(s))

	s = "\xE4\xB8\xA5"
	t.Log("==============================")
	t.Log(s)
	t.Log(len(s))

	s = "严"
	t.Log("==============================")
	t.Log(s)
	t.Log(len(s))
}

//字符串遍历
func Test9_1(t *testing.T) {
	s := "中华人民共和国"
	for _, v := range s {
		t.Log(string(v))
	}
}
