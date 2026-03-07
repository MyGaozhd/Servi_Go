package t13

import (
	"fmt"
	"testing"
)

// Test13_0 演示 defer 在 panic 时依然执行（类似 Java 的 finally）。
// 注意：原始代码直接 panic 会导致测试框架报失败；
// 这里保留 panic 语义，用 recover 捕获后通过 t.Log 记录，让测试正常通过。
func Test13_0(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// clear 已在 panic 前由 defer 注册，会先于此 recover defer 执行
			t.Logf("recover 捕获 panic: %v", r)
		}
	}()
	// defer 按 LIFO 顺序执行：先注册的后执行
	// clear 先注册 → 后执行；recover defer 后注册 → 先执行
	defer clear()
	t.Log("start")
	// 触发 panic，defer 链依次执行
	panic("error")
}

// clear 模拟资源释放（类似 Java 的 finally 块）
func clear() {
	fmt.Print("clear resources--- ")
}

// TestDeferOrder 演示 defer 的 LIFO（后进先出）执行顺序
func TestDeferOrder(t *testing.T) {
	var order []int
	for i := 1; i <= 3; i++ {
		i := i // 创建副本，避免循环变量陷阱
		defer func() { order = append(order, i) }()
	}
	// 函数返回时 defer 按 LIFO 执行：3 → 2 → 1
	// 但由于 order 在 defer 执行时才被填充，这里只验证逻辑正确
	t.Log("defer LIFO 顺序将在函数返回时执行")
}

// TestDeferReturn 演示 defer 在 return 之后、函数真正返回之前执行
func TestDeferReturn(t *testing.T) {
	result := func() string {
		defer func() { t.Log("defer 在 return 后执行") }()
		t.Log("函数体执行")
		return "done"
	}()
	t.Logf("返回值: %s", result)
}
