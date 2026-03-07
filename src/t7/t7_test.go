// Package t7 演示用 map 实现函数工厂与函数注册表（策略模式）。
package t7

import "testing"

// ─────────────────────────────────────────────
// 辅助类型
// ─────────────────────────────────────────────

// CalcFunc 是通用计算函数类型
type CalcFunc func(op int) int

// ─────────────────────────────────────────────

// Test7_0 保留原有演示：map[string]func 快速原型。
// 将不同匿名函数存入 map，通过 key 调用。
func Test7_0(t *testing.T) {
	m := map[string]func(op int) int{}
	m["1"] = func(op int) int {
		return op
	}

	m["2"] = func(op int) int {
		return op * op
	}

	m["3"] = func(op int) int {
		return op * op * op
	}

	t.Log(m["1"](2), m["2"](2), m["3"](2))
}

// TestMapFactory 演示用 map 实现函数工厂（Factory Pattern）。
//
// 工厂核心思想：
//   - 将"行为"与"名称"绑定在 map 中
//   - 调用方只需知道 key，无需关心具体实现
//   - 相比 if-else/switch，map 工厂更易扩展（开闭原则）
func TestMapFactory(t *testing.T) {
	// 构造函数工厂：名称 → 计算函数
	factory := map[string]CalcFunc{
		"identity": func(op int) int { return op },           // f(x) = x
		"square":   func(op int) int { return op * op },      // f(x) = x²
		"cube":     func(op int) int { return op * op * op }, // f(x) = x³
	}

	input := 3
	testCases := []struct {
		key      string
		expected int
	}{
		{"identity", 3},  // 3
		{"square", 9},    // 3² = 9
		{"cube", 27},     // 3³ = 27
	}

	for _, tc := range testCases {
		fn, ok := factory[tc.key]
		if !ok {
			t.Errorf("工厂中不存在 key=%q", tc.key)
			continue
		}
		got := fn(input)
		t.Logf("factory[%q](%d) = %d", tc.key, input, got)
		if got != tc.expected {
			t.Errorf("key=%q: 期望 %d，实际 %d", tc.key, tc.expected, got)
		}
	}

	// 边界 case：访问不存在的 key 返回 nil，调用 nil 函数会 panic
	didPanic := func() (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("访问不存在 key 并调用结果: panic = %v（符合预期）", r)
				panicked = true
			}
		}()
		factory["nonexistent"](input) // factory["nonexistent"] == nil → panic
		return false
	}()
	if !didPanic {
		t.Error("期望因调用 nil 函数而 panic")
	}

	// 正确做法：先用 ok-idiom 检查 key 是否存在
	if fn, ok := factory["nonexistent"]; ok {
		t.Log(fn(input))
	} else {
		t.Log("key 不存在，安全跳过（ok-idiom 防御）")
	}
}

// TestFuncRegistry 演示可扩展的函数注册表，实现策略模式简化版。
//
// 要点：
//   - 注册表是运行时可动态扩展的 map
//   - 通过名称调用对应策略，解耦"注册"与"使用"
//   - 适用场景：插件系统、中间件链、命令路由等
func TestFuncRegistry(t *testing.T) {
	// registry 是全局（包级别）注册表的局部模拟
	registry := make(map[string]CalcFunc)

	// register 向注册表添加新策略（运行时动态注册）
	register := func(name string, fn CalcFunc) {
		if fn == nil {
			t.Errorf("注册 %q 失败：fn 不能为 nil", name)
			return
		}
		registry[name] = fn
		t.Logf("已注册策略: %q", name)
	}

	// call 通过名称调用已注册的策略
	call := func(name string, op int) (int, bool) {
		fn, ok := registry[name]
		if !ok {
			return 0, false
		}
		return fn(op), true
	}

	// 初始注册：基础运算
	register("double", func(op int) int { return op * 2 })
	register("negate", func(op int) int { return -op })

	// 运行时追加新策略（体现"可扩展"）
	register("square", func(op int) int { return op * op })

	// 验证所有已注册策略
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"double", 5, 10},
		{"negate", 5, -5},
		{"square", 5, 25},
	}

	for _, tc := range tests {
		result, ok := call(tc.name, tc.input)
		if !ok {
			t.Errorf("策略 %q 未注册", tc.name)
			continue
		}
		t.Logf("call(%q, %d) = %d", tc.name, tc.input, result)
		if result != tc.expected {
			t.Errorf("策略 %q: 期望 %d，实际 %d", tc.name, tc.expected, result)
		}
	}

	// 边界 case：调用未注册的策略
	if _, ok := call("unknown", 1); ok {
		t.Error("未注册的策略不应返回 ok=true")
	} else {
		t.Log("调用未注册策略，安全返回 ok=false")
	}

	// 边界 case：覆盖已注册的策略（允许热更新）
	register("double", func(op int) int { return op * 3 }) // 改为 triple
	result, _ := call("double", 5)
	t.Logf("热更新后 call(\"double\", 5) = %d（策略已改为 ×3）", result)
	if result != 15 {
		t.Errorf("热更新后期望 15，实际 %d", result)
	}
}
