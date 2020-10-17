package t22

import (
	"t22/t22_1"
	"testing"
)

// t22相应引用t22_1里面的方法需要将t22_1文件复制到D:\develop\go\go1.15\src 中
func Test22_0(t *testing.T) {
	t.Log(t22_1.T22_1(2))
}
