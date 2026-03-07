package t22

import (
	"testing"

	"github.com/MyGaozhd/Servi_Go/src/t22/t22_1"
)

// Test22_0 演示跨子包调用：t22 引用 t22_1 子包中的函数。
// t22_1.T22_1(n) 返回 n² 和 nil error。
func Test22_0(t *testing.T) {
	result, err := t22_1.T22_1(2)
	if err != nil {
		t.Fatalf("T22_1 返回意外 error: %v", err)
	}
	if result != 4 {
		t.Errorf("T22_1(2) 期望 4，得到 %d", result)
	}
	t.Logf("T22_1(2) = %d", result)
}
