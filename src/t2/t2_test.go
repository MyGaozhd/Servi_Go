package t2

import "testing"

// Test2_0 演示 for 循环的 while 写法
func Test2_0(t *testing.T) {
	a := 0
	/* 相当于 java 的 while(a<5) */
	for a < 5 {
		a++
		t.Log(a)
	}
}

// Test2_1 演示 switch 语句（带具体值的 case）
func Test2_1(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch i {
		case 1, 2:
			t.Log("1,2->", i)
		case 3, 4:
			t.Log("3,4->", i)
		default:
			t.Log("0->", i)
		}
	}
}

// Test2_2 演示 switch 无条件写法（相当于 switch true）
func Test2_2(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch {
		case i == 1 || i == 2:
			t.Log("1,2->", i)
		case i == 3 || i == 4:
			t.Log("3,4->", i)
		default:
			t.Log("0->", i)
		}
	}
}

// TestForRange 演示 for-range 的各种遍历场景：
//   - 遍历 array（下标 + 值）
//   - 遍历 slice（下标 + 值）
//   - 遍历 map（key + value，顺序不固定）
//   - 遍历 string（按 rune 迭代，正确处理多字节 UTF-8 字符，而非逐 byte）
//   - 只要下标不要值：for i := range arr
func TestForRange(t *testing.T) {
	// ── 遍历 array ─────────────────────────────────────────────────────
	arr := [3]string{"Go", "Python", "Java"}
	t.Log("遍历 array:")
	for i, v := range arr {
		t.Logf("  arr[%d] = %s", i, v)
	}

	// ── 遍历 slice ─────────────────────────────────────────────────────
	sl := []int{10, 20, 30, 40}
	sum := 0
	for i, v := range sl {
		t.Logf("  sl[%d] = %d", i, v)
		sum += v
	}
	if sum != 100 {
		t.Errorf("slice 求和期望 100，得到 %d", sum)
	}

	// ── 遍历 map ───────────────────────────────────────────────────────
	// 注意：map 的遍历顺序每次运行可能不同，不能依赖顺序做断言
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	total := 0
	t.Log("遍历 map（顺序不固定）:")
	for k, v := range m {
		t.Logf("  m[%q] = %d", k, v)
		total += v
	}
	if total != 6 {
		t.Errorf("map value 求和期望 6，得到 %d", total)
	}

	// ── 遍历 string ────────────────────────────────────────────────────
	// for-range 遍历字符串时 v 是 rune（Unicode 码点），而非 byte；
	// 对多字节 UTF-8 字符（如中文），下标会跳跃而非连续
	s := "Go语言"
	t.Log("遍历 string（按 rune，不是 byte）:")
	var runes []rune
	for i, r := range s {
		t.Logf("  s[%d] = %c（rune=%d）", i, r, r)
		runes = append(runes, r)
	}
	// "Go语言" 含 4 个 rune，但 byte 长度为 2 + 3 + 3 = 8
	if len(runes) != 4 {
		t.Errorf("期望 4 个 rune，得到 %d 个", len(runes))
	}
	if len(s) != 8 { // byte 长度
		t.Errorf("byte 长度期望 8，得到 %d", len(s))
	}

	// ── 只要下标不要值：for i := range ────────────────────────────────
	t.Log("只取下标:")
	for i := range arr {
		t.Logf("  下标 %d", i)
	}

	// ── 只要值不要下标：for _, v := range ─────────────────────────────
	t.Log("只取值:")
	for _, v := range sl {
		t.Logf("  值 %d", v)
	}

	// ── 坑：循环变量取地址 ─────────────────────────────────────────────
	// Go 1.22 之前，range 的循环变量在整个循环中复用同一块内存；
	// 直接对 v 取地址会导致所有指针指向相同位置（最后一个值）。
	// 修复方式：在循环体内用 v := v 创建局部变量副本再取地址。
	ptrs := make([]*int, len(sl))
	for i, v := range sl {
		v := v // 遮蔽循环变量，创建独立副本
		ptrs[i] = &v
	}
	for i, p := range ptrs {
		if *p != sl[i] {
			t.Errorf("ptrs[%d] 期望 %d，得到 %d", i, sl[i], *p)
		}
	}
}

// TestSwitchFallthrough 演示 fallthrough 关键字：
//   - Go switch 默认不穿透（与 C / Java 不同，无需手动 break）
//   - 显式写 fallthrough 才继续执行下一个 case 的代码体
//   - 关键坑：fallthrough 不会重新检查下一 case 的条件，直接执行其代码体
func TestSwitchFallthrough(t *testing.T) {
	// ── 默认不穿透 ─────────────────────────────────────────────────────
	x := 1
	var result []string
	switch x {
	case 1:
		result = append(result, "case 1")
	case 2:
		result = append(result, "case 2") // 不会被执行
	default:
		result = append(result, "default")
	}
	t.Logf("默认不穿透: result=%v", result)
	if len(result) != 1 || result[0] != "case 1" {
		t.Errorf("期望只执行 case 1，得到 %v", result)
	}

	// ── 显式 fallthrough ───────────────────────────────────────────────
	// fallthrough 必须是 case 代码体的最后一条语句；
	// 执行后无条件进入下一个 case 的代码体
	var log []string
	switch x { // x == 1
	case 1:
		log = append(log, "case 1")
		fallthrough
	case 2:
		// x != 2，但因 fallthrough 进入，条件不再检查
		log = append(log, "case 2 body（fallthrough 进入）")
		fallthrough
	case 3:
		log = append(log, "case 3 body（fallthrough 进入）")
		// 无 fallthrough，到此停止
	case 4:
		log = append(log, "case 4（不应被执行）")
	}
	t.Logf("显式 fallthrough: log=%v", log)
	if len(log) != 3 {
		t.Errorf("期望 case 1/2/3 共 3 项，得到 %d 项: %v", len(log), log)
	}

	// ── 坑：fallthrough 不检查下一 case 的条件 ────────────────────────
	y := 10
	var trap []string
	switch y {
	case 10:
		trap = append(trap, "matched case 10")
		fallthrough
	case 99: // y != 99，条件不满足，但 fallthrough 照样执行其代码体
		trap = append(trap, "case 99 代码体被执行（条件从未检查）")
	}
	t.Logf("fallthrough 不检查条件: trap=%v", trap)
	if len(trap) != 2 {
		t.Errorf("期望 trap 长度 2，得到 %d", len(trap))
	}
}
