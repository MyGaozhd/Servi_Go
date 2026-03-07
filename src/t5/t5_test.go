package t5

import (
	"sort"
	"testing"
)

// Test5_0 演示 map 的基本用法
func Test5_0(t *testing.T) {
	m1 := map[string]int{"1": 1, "2": 2, "3": 3}
	t.Logf("len m1=%d", len(m1))
	t.Log(m1["1"])

	m2 := map[string]int{}
	m2["kkk"] = 1
	t.Logf("len m2=%d", len(m2))
	t.Log(m2["kkk"])

	m3 := make(map[string]int, 10)
	m3["hhh"] = 1
	t.Logf("len m3=%d", len(m3))
	t.Log(m3["hhh"])
}

// Test5_1 演示 map 判断 key 是否存在
func Test5_1(t *testing.T) {
	m2 := map[string]int{}
	t.Log(m2["key"])
	m2["key"] = 0
	t.Log(m2["key"])

	m3 := make(map[string]int, 10)
	if val, ok := m3["key"]; ok {
		t.Log(val)
	} else {
		t.Log("不存在")
	}

	m3["key"] = 1
	if val, ok := m3["key"]; ok {
		t.Log(val)
	} else {
		t.Log("不存在")
	}
}

// Test5_2 演示 map 遍历
func Test5_2(t *testing.T) {
	m1 := map[string]int{"1": 1, "2": 2, "3": 3}

	for k, v := range m1 {
		t.Log(k, v)
	}
}

// TestMapDelete 演示 map 的键删除操作：
//   - delete(m, key) 删除指定键
//   - 删除不存在的键不会 panic（静默忽略）
//   - 删除后 len 减少
func TestMapDelete(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	t.Logf("初始: len=%d, map=%v", len(m), m)

	// ── 删除存在的键 ───────────────────────────────────────────────────
	delete(m, "b")
	t.Logf("删除 'b' 后: len=%d, map=%v", len(m), m)
	if len(m) != 2 {
		t.Errorf("删除后期望 len=2，得到 %d", len(m))
	}
	if _, ok := m["b"]; ok {
		t.Error("'b' 应已被删除，但仍存在于 map 中")
	}

	// ── 删除不存在的键不会 panic ──────────────────────────────────────
	// delete 对 nil map 也不会 panic（仅对 nil map 写操作会 panic）
	delete(m, "不存在的key") // 静默忽略，不报错
	t.Logf("删除不存在键后: len=%d（len 不变）", len(m))
	if len(m) != 2 {
		t.Errorf("删除不存在键后 len 应不变，期望 2，得到 %d", len(m))
	}

	// ── 连续删除全部键 ────────────────────────────────────────────────
	delete(m, "a")
	delete(m, "c")
	t.Logf("清空后: len=%d", len(m))
	if len(m) != 0 {
		t.Errorf("全部删除后 len 应为 0，得到 %d", len(m))
	}
}

// TestMapNilPanic 演示 nil map 的行为：
//   - var m map[string]int 声明后 m 是 nil，未分配底层哈希表
//   - 读 nil map 不会 panic，返回值类型的零值
//   - 写 nil map 会 panic：assignment to entry in nil map
//   - 使用 defer + recover 捕获 panic 进行演示
func TestMapNilPanic(t *testing.T) {
	// ── 读 nil map 不 panic，返回零值 ─────────────────────────────────
	var m map[string]int // m == nil
	val := m["anyKey"]   // 不 panic，返回 int 零值 0
	t.Logf("读 nil map: m[\"anyKey\"] = %d（零值，未 panic）", val)
	if val != 0 {
		t.Errorf("读 nil map 期望返回零值 0，得到 %d", val)
	}

	// len(nil map) == 0，也不会 panic
	if len(m) != 0 {
		t.Errorf("nil map 的 len 期望 0，得到 %d", len(m))
	}

	// ── 写 nil map 会 panic ────────────────────────────────────────────
	// 通过 defer+recover 捕获 panic，验证确实会触发
	didPanic := func() (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("捕获到 panic: %v（写 nil map 导致）", r)
				panicked = true
			}
		}()
		var nilMap map[string]int
		nilMap["key"] = 1 // 这里会 panic
		return false
	}()

	if !didPanic {
		t.Error("写 nil map 应触发 panic，但未发生")
	}
	t.Log("写 nil map 会 panic，已通过 recover 捕获——正确！")

	// ── 正确做法：先用 make 初始化 ────────────────────────────────────
	m2 := make(map[string]int) // 分配了底层哈希表，不再是 nil
	m2["key"] = 42             // 安全写入
	t.Logf("make 初始化后写入: m2[\"key\"] = %d", m2["key"])
	if m2["key"] != 42 {
		t.Errorf("make 后写入期望 42，得到 %d", m2["key"])
	}
}

// TestMapUnordered 演示 map 遍历顺序不确定：
//   - 多次遍历同一 map，打印每次的 key 序列（可能不同）
//   - 如需有序输出，应将 key 提取到 slice 后排序，再按序访问
func TestMapUnordered(t *testing.T) {
	m := map[string]int{
		"banana":     3,
		"apple":      5,
		"cherry":     1,
		"date":       2,
		"elderberry": 4,
	}

	// ── 多次遍历，顺序可能不同 ────────────────────────────────────────
	// Go 运行时故意随机化 map 的遍历起点，以防止开发者依赖遍历顺序
	t.Log("map 多次遍历顺序（可能不同）：")
	for round := 1; round <= 3; round++ {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		t.Logf("  第 %d 次: %v", round, keys)
	}

	// ── 需要有序遍历时：提取 key → 排序 → 按序访问 ───────────────────
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 字典序排序

	t.Log("有序遍历（key 排序后）：")
	prevKey := ""
	for _, k := range keys {
		t.Logf("  m[%q] = %d", k, m[k])
		// 验证当前 key 确实大于等于前一个 key（即有序）
		if k < prevKey {
			t.Errorf("遍历顺序不正确：%q 应在 %q 之后", k, prevKey)
		}
		prevKey = k
	}

	// 验证有序遍历覆盖了所有键
	if len(keys) != len(m) {
		t.Errorf("有序遍历 key 数量期望 %d，得到 %d", len(m), len(keys))
	}
}
