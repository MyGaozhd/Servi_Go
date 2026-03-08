package t10

import (
	"strconv"
	"strings"
	"testing"
)

//字符串分割
func Test10_0(t *testing.T) {
	s := "a,b,c"
	parts := strings.Split(s, ",")
	for _, v := range parts {
		t.Log(v)
	}
}

//字符串链接
func Test10_1(t *testing.T) {
	s := "a,b,c"
	parts := strings.Split(s, ",")
	t.Log(strings.Join(parts, "-"))
}

//字符串整数转换
func Test10_2(t *testing.T) {
	s := strconv.Itoa(10)
	t.Log("str->", s)

	if i, err := strconv.Atoi("10"); err == nil {
		t.Log(10 + i)
	}
}

//字符串替换
func Test10_3(t *testing.T) {
	s := "a,b,c"
	s = strings.ReplaceAll(s, ",", "-")
	t.Log(s)
}
