package t22

import (
	"github.com/MyGaozhd/Servi_Go/src/t22/t22_1"
	"testing"
)

// t22 演示包的引用，使用 Go Modules 标准 import 路径（原注释：旧GOPATH风格需拷贝文件到src目录，已修正）
func Test22_0(t *testing.T) {
	t.Log(t22_1.T22_1(2))
}
