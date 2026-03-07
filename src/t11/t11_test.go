// Package t11 演示 Go 函数作为一等公民的各种用法：
// 多返回值、可变参数、闭包、函数作为返回值、命名返回值。
package t11

import (
	"errors"
	"math/rand"
	"testing"
)

// ─────────────────────────────────────────────
// 原有辅助函数（保留）
// ─────────────────────────────────────────────

func fun1() (int, int) {
	return rand.Intn(10), rand.Intn(20)
}

func sum(op ...int) int {
	s := 0
	for _, item := range op {
		s += item
	}
	return s
}

// ─────────────────────────────────────────────

// Test11_0 保留原有演示：函数多返回值。
func Test11_0(t *testing.T) {
	a, b := fun1()
	t.Log(a, b)
}

// Test11_1 保留原有演示：可变参数函数。
func Test11_1(t *testing.T) {
	t.Log(sum(1, 2, 3, 4))
}

// ─────────────────────────────────────────────
// 新增辅助函数
// ─────────────────────────────────────────────

// makeAdder 返回一个将输入值加上 n 的函数（工厂函数 + 闭包）。
func makeAdder(n int) func(int) int {
	return func(x int) int {
		return x + n
	}
}

// divide 演示命名返回值：除法，除数为 0 时返回 error。
// 命名返回值 result/err 可在函数体内直接赋值，裸 return 自动返回它们。
func divide(a, b float64) (result float64, err error) {
	if b == 0 {
		err = errors.New("除数不能为零")
		return // 裸 return：返回 result=0, err=非nil
	}
	result = a / b
	return // 裸 return：返回 result=计算结果, err=nil
}

// ─────────────────────────────────────────────

// TestClosure 演示闭包的两个核心特性：
//  1. 闭包捕获外部变量（计数器 adder）
//  2. for 循环闭包陷阱：goroutine 中捕获循环变量的错误与正确写法
func TestClosure(t *testing.T) {
	// —— 1. 闭包作为计数器 ——
	// makeCounter 返回一个闭包，每次调用将内部计数器加 1 并返回新值。
	// counter 变量生存在闭包的"捕获环境"中，不随函数返回而消失。
	makeCounter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}

	c1 := makeCounter()
	c2 := makeCounter() // c1 与 c2 拥有各自独立的 count

	t.Logf("c1: %d, %d, %d", c1(), c1(), c1()) // 1, 2, 3
	t.Logf("c2: %d（独立计数器，从 1 开始）", c2())   // 1

	if c1() != 4 {
		t.Error("c1 第 4 次调用应返回 4")
	}
	if c2() != 2 {
		t.Error("c2 第 2 次调用应返回 2")
	}

	// —— 2. for 循环闭包陷阱 ——
	// 错误写法：所有闭包共享同一个循环变量 i，
	// 等到闭包执行时 i 已经是循环结束后的值（3）。
	wrongFuncs := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		i := i // ← 正确写法：用 := 在循环体内重新声明新变量，遮蔽外层 i
		wrongFuncs[i] = func() int { return i }
	}
	t.Log("用 := 重新声明修复循环闭包陷阱:")
	for idx, f := range wrongFuncs {
		got := f()
		t.Logf("  wrongFuncs[%d]() = %d", idx, got)
		if got != idx {
			t.Errorf("期望 %d，实际 %d", idx, got)
		}
	}

	// 正确写法二：通过函数参数传入循环变量（相当于值拷贝）
	paramFuncs := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		// 用立即调用的匿名函数，把 i 作为参数传进去
		paramFuncs[i] = func(v int) func() int {
			return func() int { return v }
		}(i)
	}
	t.Log("通过参数传值修复循环闭包陷阱:")
	for idx, f := range paramFuncs {
		got := f()
		t.Logf("  paramFuncs[%d]() = %d", idx, got)
		if got != idx {
			t.Errorf("期望 %d，实际 %d", idx, got)
		}
	}
}

// TestFuncAsReturnValue 演示函数作为返回值（工厂函数）。
//
// makeAdder(n) 返回一个函数，该函数将参数加上 n。
// 返回的函数通过闭包"记住"了 n 的值。
func TestFuncAsReturnValue(t *testing.T) {
	add5 := makeAdder(5)
	add10 := makeAdder(10)

	tests := []struct {
		fn       func(int) int
		name     string
		input    int
		expected int
	}{
		{add5, "add5", 3, 8},
		{add5, "add5", 0, 5},
		{add10, "add10", 7, 17},
		{add10, "add10", -3, 7},
	}

	for _, tc := range tests {
		got := tc.fn(tc.input)
		t.Logf("%s(%d) = %d", tc.name, tc.input, got)
		if got != tc.expected {
			t.Errorf("%s(%d): 期望 %d，实际 %d", tc.name, tc.input, tc.expected, got)
		}
	}

	// add5 和 add10 是相互独立的函数，各自捕获了不同的 n
	t.Logf("add5 与 add10 是独立函数，互不影响")
}

// TestNamedReturn 演示命名返回值的用法与注意事项。
//
// 命名返回值：
//   - 在函数签名中为返回值指定名称，相当于在函数体顶部声明零值变量
//   - 裸 return（naked return）自动返回所有命名变量的当前值
//   - 适用于：短函数、defer 修改返回值；不适用于长函数（影响可读性）
func TestNamedReturn(t *testing.T) {
	// 正常除法
	r1, err1 := divide(10, 3)
	t.Logf("divide(10, 3) = %.4f, err = %v", r1, err1)
	if err1 != nil {
		t.Errorf("不期望 error: %v", err1)
	}
	if r1 < 3.333 || r1 > 3.334 {
		t.Errorf("divide(10,3) 期望约 3.3333，实际 %.4f", r1)
	}

	// 除以零：触发命名返回值的 err 分支
	r2, err2 := divide(10, 0)
	t.Logf("divide(10, 0) = %.4f, err = %v", r2, err2)
	if err2 == nil {
		t.Error("除以零应返回 error")
	}
	if r2 != 0 {
		t.Errorf("除以零时 result 应为零值 0，实际 %.4f", r2)
	}

	// 边界 case：负数除法
	r3, err3 := divide(-6, 2)
	t.Logf("divide(-6, 2) = %.1f, err = %v", r3, err3)
	if err3 != nil || r3 != -3 {
		t.Errorf("divide(-6,2): 期望 -3，实际 %.1f", r3)
	}

	// 边界 case：被除数为 0
	r4, err4 := divide(0, 5)
	t.Logf("divide(0, 5) = %.1f, err = %v", r4, err4)
	if err4 != nil || r4 != 0 {
		t.Errorf("divide(0,5): 期望 0，实际 %.1f", r4)
	}

	// 说明：命名返回值配合 defer 可修改最终返回值
	withDefer := func() (n int) {
		defer func() { n++ }() // defer 在 return 后、函数真正退出前执行，可修改命名返回值
		return 1               // n = 1，然后 defer 将 n++ → 最终返回 2
	}
	t.Logf("命名返回值 + defer 修改: %d（期望 2）", withDefer())
	if withDefer() != 2 {
		t.Error("命名返回值配合 defer 应返回 2")
	}
}
