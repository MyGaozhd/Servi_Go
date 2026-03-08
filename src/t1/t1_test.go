package t1

import "testing"

// Test1_0 演示清零运算符 &^
// 如果右边是 1，不管左边是多少，最终结果都是 0
// 如果右边是 0，左边是多少结果就是多少
func Test1_0(t *testing.T) {
	a := 1 &^ 1
	b := 1 &^ 0
	c := 0 &^ 1
	d := 0 &^ 0
	t.Log(a, b, c, d)
}

// Test1_1 与上面的写法效果相同。所以 &^ 运算符相当于一种语法糖
func Test1_1(t *testing.T) {
	a := 1 & (^1)
	b := 1 & (^0)
	c := 0 & (^1)
	d := 0 & (^0)
	t.Log(a, b, c, d)
}

// TestBitOps 演示所有基础位运算：
//   - & (AND)：两个位都为 1 才为 1
//   - | (OR)：两个位中有 1 就为 1
//   - ^ (XOR)：两个位相同为 0，不同为 1
//   - << 左移：每左移 1 位相当于 *2
//   - >> 右移：每右移 1 位相当于 /2（整除）
//   - ^ 单目取反：对所有位取反（在 Go 中用单目 ^，而非 ~）
func TestBitOps(t *testing.T) {
	// ── & AND ─────────────────────────────────────────────────────────
	// 0101 (5)
	// 0011 (3)
	// ---- AND
	// 0001 (1)
	andResult := 5 & 3
	t.Logf("5 & 3 = %d（二进制: %04b & %04b = %04b）", andResult, 5, 3, andResult)
	if andResult != 1 {
		t.Errorf("5 & 3 期望 1，得到 %d", andResult)
	}

	// ── | OR ──────────────────────────────────────────────────────────
	// 0101 (5)
	// 0011 (3)
	// ---- OR
	// 0111 (7)
	orResult := 5 | 3
	t.Logf("5 | 3 = %d（二进制: %04b | %04b = %04b）", orResult, 5, 3, orResult)
	if orResult != 7 {
		t.Errorf("5 | 3 期望 7，得到 %d", orResult)
	}

	// ── ^ XOR ─────────────────────────────────────────────────────────
	// 0101 (5)
	// 0011 (3)
	// ---- XOR
	// 0110 (6)
	xorResult := 5 ^ 3
	t.Logf("5 ^ 3 = %d（二进制: %04b ^ %04b = %04b）", xorResult, 5, 3, xorResult)
	if xorResult != 6 {
		t.Errorf("5 ^ 3 期望 6，得到 %d", xorResult)
	}

	// ── << 左移 ────────────────────────────────────────────────────────
	// 左移 1 位 = ×2；左移 n 位 = ×(2^n)
	// 1 << 3 = 1 × 8 = 8
	leftShift := 1 << 3
	t.Logf("1 << 3 = %d（相当于 1×2×2×2）", leftShift)
	if leftShift != 8 {
		t.Errorf("1 << 3 期望 8，得到 %d", leftShift)
	}

	// ── >> 右移 ────────────────────────────────────────────────────────
	// 右移 1 位 = ÷2（整除）；16 >> 2 = 16÷4 = 4
	rightShift := 16 >> 2
	t.Logf("16 >> 2 = %d（相当于 16÷4）", rightShift)
	if rightShift != 4 {
		t.Errorf("16 >> 2 期望 4，得到 %d", rightShift)
	}

	// ── ^ 单目取反 ─────────────────────────────────────────────────────
	// Go 用 ^ 作为单目取反运算符（其他语言用 ~）
	// ^0 在 int 下为 -1（所有位全 1 的补码表示）
	// ^1 = -2（最低位从 1 变 0，其余全 1）
	notZero := ^0
	notOne := ^1
	t.Logf("^0 = %d, ^1 = %d", notZero, notOne)
	if notZero != -1 {
		t.Errorf("^0 期望 -1，得到 %d", notZero)
	}
	if notOne != -2 {
		t.Errorf("^1 期望 -2，得到 %d", notOne)
	}
}

// TestBitClear 演示 &^ 清零运算符的实际应用：
//   - 用 &^ 清除指定的 bit 位（如清除某组权限位）
//   - 验证：清除后相应位变为 0，其余位不变
func TestBitClear(t *testing.T) {
	// 定义权限 bit 常量（使用 iota 位移，每个常量占 1 个 bit）
	const (
		PermRead    = 1 << iota // 0001 = 1
		PermWrite               // 0010 = 2
		PermExecute             // 0100 = 4
	)

	// 初始拥有全部权限：Read | Write | Execute = 0111 = 7
	perm := PermRead | PermWrite | PermExecute
	t.Logf("初始权限: %03b（Read|Write|Execute = %d）", perm, perm)
	if perm != 7 {
		t.Errorf("初始权限期望 7，得到 %d", perm)
	}

	// 用 &^ 清除写权限（PermWrite 对应的 bit 清零，其余 bit 不变）
	// perm  = 0111
	// Write = 0010
	// 结果  = 0101 = 5（保留 Read 和 Execute）
	perm = perm &^ PermWrite
	t.Logf("清除写权限后: %03b（= %d，期望保留 Read+Execute）", perm, perm)
	if perm != PermRead|PermExecute {
		t.Errorf("清除写权限后期望 %d，得到 %d", PermRead|PermExecute, perm)
	}

	// 验证：写权限已清除，读/执行权限仍存在
	if perm&PermWrite != 0 {
		t.Error("写权限应已被清除")
	}
	if perm&PermRead == 0 {
		t.Error("读权限不应被清除")
	}
	if perm&PermExecute == 0 {
		t.Error("执行权限不应被清除")
	}

	// 连续清除多个 bit：同时清除 Read 和 Execute
	perm = perm &^ (PermRead | PermExecute)
	t.Logf("清除全部权限后: %03b（= %d）", perm, perm)
	if perm != 0 {
		t.Errorf("清除全部权限后期望 0，得到 %d", perm)
	}
}
