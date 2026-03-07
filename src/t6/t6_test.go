package t6

import (
	"testing"
)

// ──────────────────────────────────────────────
// 辅助类型
// ──────────────────────────────────────────────

type person struct {
	name string
	age  int
}

// ──────────────────────────────────────────────
// TestValuePass 演示值传递：函数内修改不影响外部变量
// ──────────────────────────────────────────────
func TestValuePass(t *testing.T) {
	// ── int 值传递 ────────────────────────────────────────────────────
	changeInt := func(n int) { n = 999 }
	x := 1
	changeInt(x)
	if x != 1 {
		t.Errorf("值传递 int：期望 x=1，得到 %d", x)
	}
	t.Logf("int 值传递后 x=%d（未变）", x)

	// ── string 值传递 ──────────────────────────────────────────────────
	changeStr := func(s string) { s = "changed" }
	s := "hello"
	changeStr(s)
	if s != "hello" {
		t.Errorf("值传递 string：期望 s=\"hello\"，得到 %q", s)
	}
	t.Logf("string 值传递后 s=%q（未变）", s)

	// ── struct 值传递 ──────────────────────────────────────────────────
	changePerson := func(p person) { p.name = "Bob"; p.age = 99 }
	p := person{name: "Alice", age: 30}
	changePerson(p)
	if p.name != "Alice" || p.age != 30 {
		t.Errorf("值传递 struct：期望 Alice/30，得到 %v", p)
	}
	t.Logf("struct 值传递后 p=%+v（未变）", p)
}

// ──────────────────────────────────────────────
// TestPointerPass 演示指针传递：函数内修改影响外部变量
// ──────────────────────────────────────────────
func TestPointerPass(t *testing.T) {
	// ── *int 指针传递 ──────────────────────────────────────────────────
	changeIntPtr := func(n *int) { *n = 999 }
	x := 1
	changeIntPtr(&x)
	if x != 999 {
		t.Errorf("指针传递 *int：期望 x=999，得到 %d", x)
	}
	t.Logf("*int 指针传递后 x=%d（已变）", x)

	// ── *struct 指针传递 ───────────────────────────────────────────────
	changePersonPtr := func(p *person) { p.name = "Bob"; p.age = 99 }
	p := person{name: "Alice", age: 30}
	changePersonPtr(&p)
	if p.name != "Bob" || p.age != 99 {
		t.Errorf("指针传递 *struct：期望 Bob/99，得到 %v", p)
	}
	t.Logf("*struct 指针传递后 p=%+v（已变）", p)

	// ── 指针实现 swap ──────────────────────────────────────────────────
	swapByPtr := func(a, b *int) { *a, *b = *b, *a }
	a, b := 10, 20
	swapByPtr(&a, &b)
	if a != 20 || b != 10 {
		t.Errorf("指针 swap：期望 a=20 b=10，得到 a=%d b=%d", a, b)
	}
	t.Logf("指针 swap 后 a=%d, b=%d", a, b)
}

// ──────────────────────────────────────────────
// TestPointerArith 演示取地址 & 解引用 * 和双重指针 **int
// ──────────────────────────────────────────────
func TestPointerArith(t *testing.T) {
	// ── & 取地址 ──────────────────────────────────────────────────────
	n := 42
	p := &n // p 是 *int，存储 n 的地址
	t.Logf("n=%d, &n=%p, p=%p, *p=%d", n, &n, p, *p)
	if *p != 42 {
		t.Errorf("*p 期望 42，得到 %d", *p)
	}

	// 通过指针修改值
	*p = 100
	if n != 100 {
		t.Errorf("通过指针修改后 n 期望 100，得到 %d", n)
	}
	t.Logf("*p=100 后 n=%d", n)

	// ── 双重指针 **int ──────────────────────────────────────────────────
	pp := &p // pp 是 **int，存储 p 的地址
	**pp = 200
	if n != 200 {
		t.Errorf("双重指针修改后 n 期望 200，得到 %d", n)
	}
	t.Logf("**pp=200 后 n=%d", n)
}

// ──────────────────────────────────────────────
// TestNewVsMake 演示 new(T) 与 make(T) 的区别
// ──────────────────────────────────────────────
func TestNewVsMake(t *testing.T) {
	// ── new(T) 返回 *T，值为类型零值 ──────────────────────────────────
	// new 可用于任意类型，分配内存并清零，返回指针
	pi := new(int)    // *int，值为 0
	ps := new(string) // *string，值为 ""
	t.Logf("new(int)=%d, new(string)=%q", *pi, *ps)
	if *pi != 0 || *ps != "" {
		t.Error("new 应返回类型零值")
	}
	*pi = 7
	if *pi != 7 {
		t.Error("new 返回的指针应可赋值")
	}

	// ── make(T) 仅用于 slice/map/chan，返回初始化好的值（非指针）──────
	// make 会初始化内部数据结构，返回的是已可使用的引用类型

	// make slice
	sl := make([]int, 3, 5) // len=3, cap=5，元素均为 0
	t.Logf("make slice: len=%d cap=%d %v", len(sl), cap(sl), sl)
	if len(sl) != 3 || cap(sl) != 5 {
		t.Errorf("make slice 期望 len=3 cap=5，得到 len=%d cap=%d", len(sl), cap(sl))
	}

	// make map
	m := make(map[string]int, 4) // 预分配容量 4
	m["k"] = 1
	if m["k"] != 1 {
		t.Error("make map 应可直接写入")
	}

	// make chan
	ch := make(chan int, 2) // 带缓冲 channel
	ch <- 10
	ch <- 20
	if len(ch) != 2 {
		t.Errorf("make chan: 期望 len=2，得到 %d", len(ch))
	}
}

// ──────────────────────────────────────────────
// TestNilPointer 演示 nil 指针解引用触发 panic，及 recover 捕获
// ──────────────────────────────────────────────
func TestNilPointer(t *testing.T) {
	// ── nil 指针解引用触发 panic ───────────────────────────────────────
	didPanic := func() (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("捕获 panic: %v（nil 指针解引用）", r)
				panicked = true
			}
		}()
		var p *int // nil 指针
		_ = *p     // runtime panic: nil pointer dereference
		return false
	}()

	if !didPanic {
		t.Error("nil 指针解引用应触发 panic")
	}

	// ── 最佳实践：使用前检查 nil ─────────────────────────────────────
	safeDeref := func(p *int) int {
		if p == nil {
			return -1 // 默认值
		}
		return *p
	}

	var nilPtr *int
	result := safeDeref(nilPtr)
	if result != -1 {
		t.Errorf("nil 指针安全解引用期望 -1，得到 %d", result)
	}

	val := 42
	result = safeDeref(&val)
	if result != 42 {
		t.Errorf("非 nil 指针解引用期望 42，得到 %d", result)
	}
	t.Log("nil 指针安全处理演示完成")
}
