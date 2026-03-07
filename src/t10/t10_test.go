// Package t10 演示 Go 常用字符串工具函数及高效拼接方式。
package t10

import (
	"strconv"
	"strings"
	"testing"
)

// Test10_0 保留原有演示：strings.Split 字符串分割。
func Test10_0(t *testing.T) {
	s := "a,b,c"
	parts := strings.Split(s, ",")
	for _, v := range parts {
		t.Log(v)
	}
}

// Test10_1 保留原有演示：strings.Join 字符串连接。
func Test10_1(t *testing.T) {
	s := "a,b,c"
	parts := strings.Split(s, ",")
	t.Log(strings.Join(parts, "-"))
}

// Test10_2 保留原有演示：strconv.Itoa / strconv.Atoi 整数与字符串互转。
func Test10_2(t *testing.T) {
	s := strconv.Itoa(10)
	t.Log("str->", s)

	if i, err := strconv.Atoi("10"); err == nil {
		t.Log(10 + i)
	}
}

// Test10_3 保留原有演示：strings.ReplaceAll 全局替换。
func Test10_3(t *testing.T) {
	s := "a,b,c"
	s = strings.ReplaceAll(s, ",", "-")
	t.Log(s)
}

// TestStringBuilder 演示用 strings.Builder 高效拼接字符串。
//
// 为什么要用 Builder？
//   - 直接用 + 拼接字符串，每次都会分配新内存，拼接 n 次 → O(n²) 总内存分配
//   - strings.Builder 内部维护 []byte 缓冲区，Write 操作均摊 O(1)
//   - 拼接完毕后，调用 String() 只做一次内存分配
func TestStringBuilder(t *testing.T) {
	// —— WriteString：拼接字符串片段 ——
	var b strings.Builder
	words := []string{"Go", " ", "is", " ", "awesome", "!"}
	for _, w := range words {
		b.WriteString(w)
	}
	result := b.String()
	t.Logf("WriteString 拼接结果: %q", result)
	if result != "Go is awesome!" {
		t.Errorf("期望 \"Go is awesome!\"，实际 %q", result)
	}

	// —— Reset 后复用 Builder ——
	b.Reset()
	if b.Len() != 0 {
		t.Errorf("Reset 后 Len 应为 0，实际 %d", b.Len())
	}

	// —— WriteByte：拼接单个字节 ——
	for i := byte('A'); i <= 'E'; i++ {
		b.WriteByte(i)
	}
	byteResult := b.String()
	t.Logf("WriteByte 拼接结果: %q", byteResult)
	if byteResult != "ABCDE" {
		t.Errorf("期望 \"ABCDE\"，实际 %q", byteResult)
	}

	// —— WriteRune：拼接 Unicode 字符（含中文）——
	b.Reset()
	runes := []rune{'你', '好', ',', 'G', 'o'}
	for _, r := range runes {
		b.WriteRune(r)
	}
	runeResult := b.String()
	t.Logf("WriteRune 拼接结果: %q", runeResult)
	if runeResult != "你好,Go" {
		t.Errorf("期望 \"你好,Go\"，实际 %q", runeResult)
	}

	// —— 大量拼接场景：Builder vs + ——
	// 此处仅做演示说明，不做性能基准测试
	b.Reset()
	for i := 0; i < 5; i++ {
		b.WriteString(strconv.Itoa(i))
		if i < 4 {
			b.WriteByte(',')
		}
	}
	t.Logf("拼接 0~4: %q（用 Builder 避免重复分配）", b.String())
	if b.String() != "0,1,2,3,4" {
		t.Errorf("期望 \"0,1,2,3,4\"，实际 %q", b.String())
	}
}

// TestStringContains 演示常用字符串查找与变换函数。
//
// 涵盖：Contains / HasPrefix / HasSuffix / Index / TrimSpace / ToUpper / ToLower
func TestStringContains(t *testing.T) {
	s := "  Hello, Go World!  "

	// —— TrimSpace：去掉首尾空白 ——
	trimmed := strings.TrimSpace(s)
	t.Logf("TrimSpace: %q → %q", s, trimmed)
	if trimmed != "Hello, Go World!" {
		t.Errorf("TrimSpace 结果不符，实际 %q", trimmed)
	}

	base := trimmed // 后续操作基于去空格后的字符串

	// —— Contains：包含子串 ——
	t.Logf("Contains(%q, \"Go\")   = %v", base, strings.Contains(base, "Go"))
	t.Logf("Contains(%q, \"Java\") = %v", base, strings.Contains(base, "Java"))
	if !strings.Contains(base, "Go") {
		t.Error("应包含 \"Go\"")
	}
	if strings.Contains(base, "Java") {
		t.Error("不应包含 \"Java\"")
	}

	// —— HasPrefix / HasSuffix ——
	t.Logf("HasPrefix(%q, \"Hello\") = %v", base, strings.HasPrefix(base, "Hello"))
	t.Logf("HasSuffix(%q, \"!\")     = %v", base, strings.HasSuffix(base, "!"))
	if !strings.HasPrefix(base, "Hello") {
		t.Error("应以 \"Hello\" 开头")
	}
	if !strings.HasSuffix(base, "!") {
		t.Error("应以 \"!\" 结尾")
	}

	// 边界 case：空前缀/后缀始终返回 true
	if !strings.HasPrefix(base, "") {
		t.Error("HasPrefix with empty prefix 应为 true")
	}

	// —— Index：查找子串首次出现的字节偏移 ——
	idx := strings.Index(base, "Go")
	t.Logf("Index(%q, \"Go\") = %d", base, idx)
	if idx < 0 {
		t.Error("Index 应 >= 0")
	}

	// 不存在时返回 -1
	notFound := strings.Index(base, "Rust")
	t.Logf("Index(%q, \"Rust\") = %d（-1 表示不存在）", base, notFound)
	if notFound != -1 {
		t.Errorf("期望 -1，实际 %d", notFound)
	}

	// —— ToUpper / ToLower ——
	upper := strings.ToUpper("go")
	lower := strings.ToLower("GO")
	t.Logf("ToUpper(\"go\") = %q", upper)
	t.Logf("ToLower(\"GO\") = %q", lower)
	if upper != "GO" {
		t.Errorf("ToUpper 期望 \"GO\"，实际 %q", upper)
	}
	if lower != "go" {
		t.Errorf("ToLower 期望 \"go\"，实际 %q", lower)
	}

	// —— 边界 case：中文字符串 Contains ——
	cn := "你好，世界"
	t.Logf("Contains(%q, \"世界\") = %v", cn, strings.Contains(cn, "世界"))
	if !strings.Contains(cn, "世界") {
		t.Error("应包含 \"世界\"")
	}
}
