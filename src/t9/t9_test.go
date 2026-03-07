// Package t9 演示 Go 字符串的基础特性：
// UTF-8 编码、byte/rune 区别、两种遍历方式、不可变性。
package t9

import (
	"testing"
)

// Test9_0 保留原有演示：字符串基本使用与 UTF-8 字节长度。
func Test9_0(t *testing.T) {
	var s string
	t.Log(s)
	t.Log(len(s))

	s = "hello"
	t.Log("==============================")
	t.Log(s)
	t.Log(len(s))

	s = "\xE4\xB8\xA5"
	t.Log("==============================")
	t.Log(s)
	t.Log(len(s))

	s = "严"
	t.Log("==============================")
	t.Log(s)
	t.Log(len(s))
}

// Test9_1 保留原有演示：用 range 按 rune 遍历中文字符串。
func Test9_1(t *testing.T) {
	s := "中华人民共和国"
	for _, v := range s {
		t.Log(string(v))
	}
}

// TestByteVsRune 演示 byte 与 rune 的本质区别。
//
// 核心结论：
//   - Go string 底层是 []byte（UTF-8 编码的字节序列）
//   - len(s) 返回字节数，不是字符数
//   - ASCII 字符：1字符 = 1字节
//   - 中文字符：1字符 = 3字节（UTF-8）
//   - rune（int32 别名）代表一个 Unicode 码点（字符）
func TestByteVsRune(t *testing.T) {
	// —— ASCII 字符串 ——
	ascii := "hello"
	t.Logf("[ASCII] 字符串: %q", ascii)
	t.Logf("[ASCII] len（字节数）= %d", len(ascii))       // 5
	t.Logf("[ASCII] rune 数量   = %d", len([]rune(ascii))) // 5
	if len(ascii) != len([]rune(ascii)) {
		t.Error("ASCII 字符串字节数与 rune 数应相等")
	}

	// —— 中文字符串 ——
	chinese := "严"
	t.Logf("[中文] 字符串: %q", chinese)
	t.Logf("[中文] len（字节数）= %d（UTF-8 中文字符占 3 字节）", len(chinese)) // 3
	t.Logf("[中文] rune 数量   = %d（只有 1 个字符）", len([]rune(chinese)))  // 1
	if len(chinese) != 3 {
		t.Errorf("期望 \"严\" 字节数为 3，实际 %d", len(chinese))
	}
	if len([]rune(chinese)) != 1 {
		t.Errorf("期望 \"严\" rune 数为 1，实际 %d", len([]rune(chinese)))
	}

	// —— 混合字符串 ——
	mixed := "Go语言"
	t.Logf("[混合] 字符串: %q", mixed)
	t.Logf("[混合] len（字节数）= %d（2×ASCII + 2×3字节中文 = 8）", len(mixed)) // 2+6=8
	t.Logf("[混合] rune 数量   = %d（4 个字符）", len([]rune(mixed)))            // 4
	if len(mixed) != 8 {
		t.Errorf("期望字节数 8，实际 %d", len(mixed))
	}
	if len([]rune(mixed)) != 4 {
		t.Errorf("期望 rune 数 4，实际 %d", len([]rune(mixed)))
	}

	// —— []byte vs []rune 转换 ——
	t.Log("---- []byte 转换（按字节） ----")
	bs := []byte(chinese) // "严" → [0xE4, 0xB8, 0xA5]
	for i, b := range bs {
		t.Logf("  bs[%d] = 0x%X", i, b)
	}

	t.Log("---- []rune 转换（按字符） ----")
	rs := []rune(chinese) // "严" → [0x4E25]
	for i, r := range rs {
		t.Logf("  rs[%d] = U+%04X（%s）", i, r, string(r))
	}
}

// TestStringIterate 演示两种字符串遍历方式的行为差异。
//
// 方式一：for i := 0; i < len(s); i++ → 按字节遍历，中文会产生乱码字节
// 方式二：for _, r := range s          → 按 rune 遍历，正确识别多字节字符
func TestStringIterate(t *testing.T) {
	s := "Hi,中国"

	// —— 方式一：按字节遍历 ——
	// len(s) = 2(Hi) + 1(,) + 3(中) + 3(国) = 9
	t.Log("=== 方式一：按字节遍历（for i < len(s)）===")
	t.Logf("字节总数: %d", len(s))
	for i := 0; i < len(s); i++ {
		t.Logf("  s[%d] = 0x%02X（%q）", i, s[i], rune(s[i]))
	}
	// 注意：中文字节范围 0x80~0xFF，直接 string(s[i]) 会是乱码

	// —— 方式二：按 rune 遍历（推荐处理含中文的字符串）——
	t.Log("=== 方式二：按 rune 遍历（for _, r := range s）===")
	runeCount := 0
	for i, r := range s {
		t.Logf("  字节偏移[%d]，rune = U+%04X（%s）", i, r, string(r))
		runeCount++
	}
	t.Logf("rune 总数: %d（正确字符数）", runeCount)

	expectedRuneCount := len([]rune(s)) // 5
	if runeCount != expectedRuneCount {
		t.Errorf("期望 rune 数 %d，实际 %d", expectedRuneCount, runeCount)
	}

	// —— 坑点演示：用字节索引取字符 ——
	// s[3] 是"中"的第一个字节，string(s[3]) 不是完整的"中"
	firstByteOfZhong := s[3]
	t.Logf("s[3] = 0x%02X，string(s[3]) = %q（不是完整的\"中\"）", firstByteOfZhong, string(firstByteOfZhong))
}

// TestStringImmutable 演示字符串的不可变性。
//
// 关键点：
//   - Go 字符串是只读的字节序列，不能通过下标修改
//   - 若要修改，必须先转为 []byte 或 []rune，修改后再转回 string
func TestStringImmutable(t *testing.T) {
	s := "hello"
	t.Logf("原始字符串: %q", s)

	// s[0] = 'H'  ← 编译错误：cannot assign to s[0]（此处仅注释说明，不实际编写）
	t.Log("s[0] = 'H' 会产生编译错误：cannot assign to s[0]（字符串不可变）")

	// —— 方式一：通过 []byte 修改（适合 ASCII 操作）——
	bs := []byte(s)    // 拷贝一份字节切片
	bs[0] = 'H'        // 修改字节
	modified := string(bs)
	t.Logf("[]byte 修改后: %q（原始 s 仍为 %q）", modified, s)
	if modified != "Hello" {
		t.Errorf("期望 \"Hello\"，实际 %q", modified)
	}
	if s != "hello" {
		t.Errorf("原始 s 不应被修改，实际 %q", s)
	}

	// —— 方式二：通过 []rune 修改（适合含中文的操作）——
	s2 := "你好世界"
	rs := []rune(s2)   // 拷贝一份 rune 切片
	rs[2] = '中'        // 将第 3 个字符替换为"中"
	rs[3] = '文'        // 将第 4 个字符替换为"文"
	s2Modified := string(rs)
	t.Logf("[]rune 修改后: %q（原始 s2 仍为 %q）", s2Modified, s2)
	if s2Modified != "你好中文" {
		t.Errorf("期望 \"你好中文\"，实际 %q", s2Modified)
	}
	if s2 != "你好世界" {
		t.Errorf("原始 s2 不应被修改，实际 %q", s2)
	}

	// —— 边界 case：string 赋值是安全的，两个变量相互独立 ——
	a := "foo"
	b := a   // b 是 a 的一个独立副本（值语义）
	b = "bar"
	t.Logf("a=%q，b=%q（互不影响）", a, b)
	if a != "foo" {
		t.Errorf("a 不应被修改，实际 %q", a)
	}
}
