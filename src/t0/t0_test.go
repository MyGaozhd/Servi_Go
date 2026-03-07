package t0

import (
	"testing"
)

// globalVar 演示包级全局变量，所有测试函数均可直接访问
var globalVar = 100

// TestFibonacci 测试斐波那契数列的生成
func TestFibonacci(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected []int
	}{
		{
			name:     "空数列",
			count:    0,
			expected: []int{},
		},
		{
			name:     "单个数字",
			count:    1,
			expected: []int{1},
		},
		{
			name:     "前两个数字",
			count:    2,
			expected: []int{1, 1},
		},
		{
			name:     "前六个数字",
			count:    6,
			expected: []int{1, 1, 2, 3, 5, 8},
		},
		{
			name:     "前十个数字",
			count:    10,
			expected: []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55},
		},
		{
			name:     "负数输入",
			count:    -1,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateFibonacci(tt.count)

			if len(result) != len(tt.expected) {
				t.Errorf("期望长度 %d, 实际长度 %d", len(tt.expected), len(result))
				return
			}

			for i := 0; i < len(tt.expected); i++ {
				if result[i] != tt.expected[i] {
					t.Errorf("索引 %d: 期望 %d, 实际 %d", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

// generateFibonacci 生成指定数量的斐波那契数列
func generateFibonacci(count int) []int {
	if count <= 0 {
		return []int{}
	}

	result := make([]int, count)
	current, next := 1, 1

	for i := 0; i < count; i++ {
		result[i] = current
		current, next = next, current+next
	}

	return result
}

// TestVarDecl 演示变量声明的各种方式：
//   - var 声明不赋值时各类型的零值（int=0, string="", bool=false）
//   - := 简短声明，类型由右侧表达式自动推导
//   - 多变量同时赋值，以及无临时变量的 swap 技巧
//   - 全局（包级）变量
func TestVarDecl(t *testing.T) {
	// ── 零值 ──────────────────────────────────────────────────────────
	// var 声明变量不赋初值，Go 会自动用该类型的"零值"填充：
	//   int/float → 0，string → ""，bool → false，pointer/slice/map → nil
	var zeroInt int
	var zeroStr string
	var zeroBool bool

	if zeroInt != 0 {
		t.Errorf("int 零值期望 0，得到 %d", zeroInt)
	}
	if zeroStr != "" {
		t.Errorf("string 零值期望 \"\"，得到 %q", zeroStr)
	}
	if zeroBool != false {
		t.Errorf("bool 零值期望 false，得到 %v", zeroBool)
	}
	t.Logf("var 零值演示: int=%d, string=%q, bool=%v", zeroInt, zeroStr, zeroBool)

	// ── := 简短声明 ────────────────────────────────────────────────────
	// 只能在函数体内使用；同一作用域内不能对已存在的变量单独使用 :=
	// （除非同一行同时引入至少一个新变量）
	age := 18
	name := "Gopher"
	t.Logf(":= 简短声明: age=%d, name=%s", age, name)

	// ── 多变量同时赋值 ─────────────────────────────────────────────────
	// 右侧所有表达式先求值，再统一赋值给左侧——这是 swap 不需要临时变量的原因
	a, b, c := 1, 2, 3
	t.Logf("初始: a=%d, b=%d, c=%d", a, b, c)
	if a != 1 || b != 2 || c != 3 {
		t.Error("多变量赋值结果不符合预期")
	}

	// 无需临时变量即可完成 swap
	a, b = b, a
	t.Logf("swap 后: a=%d, b=%d（原 b 和 a 互换）", a, b)
	if a != 2 || b != 1 {
		t.Error("swap 结果不符合预期")
	}
	_ = c // 避免"declared and not used"编译错误

	// ── 全局变量 ───────────────────────────────────────────────────────
	// globalVar 定义在包级别，包内所有函数可直接读写
	t.Logf("全局变量: globalVar=%d", globalVar)
	if globalVar != 100 {
		t.Errorf("全局变量期望 100，得到 %d", globalVar)
	}
}

// TestConst 演示常量声明：
//   - 单个 const
//   - const 块（枚举性常量，以性别为例）
//   - 类型化常量（显式指定类型，限制赋值兼容范围）
func TestConst(t *testing.T) {
	// ── 单个 const ────────────────────────────────────────────────────
	// 常量在编译期确定值，不可取地址，不可在运行期修改
	const Pi = 3.14159
	t.Logf("单个常量 Pi = %v", Pi)
	if Pi < 3.14 || Pi > 3.15 {
		t.Error("Pi 超出合理范围")
	}

	// ── const 块（枚举性常量） ─────────────────────────────────────────
	// 将一组相关常量放在同一 const 块中，语义清晰，便于维护
	const (
		GenderUnknown = 0
		GenderMale    = 1
		GenderFemale  = 2
	)
	t.Logf("性别常量: 未知=%d, 男=%d, 女=%d", GenderUnknown, GenderMale, GenderFemale)
	if GenderMale == GenderFemale {
		t.Error("男女性别常量值不应相同")
	}

	// ── 类型化常量 ────────────────────────────────────────────────────
	// 显式声明类型后，编译器做严格的类型检查；
	// 例如 MaxSize 不能直接赋值给 int32 变量，必须显式转换
	const MaxSize int = 1024
	t.Logf("类型化常量 MaxSize = %d（类型 int）", MaxSize)
	if MaxSize != 1024 {
		t.Errorf("MaxSize 期望 1024，得到 %d", MaxSize)
	}
}

// TestIota 演示 iota：
//   - 在 const 块内从 0 开始，每定义一个常量自动加 1
//   - 位移写法（1 << (iota * 10)）优雅地实现 KB / MB / GB
func TestIota(t *testing.T) {
	// ── iota 从 0 开始自增 ─────────────────────────────────────────────
	// iota 只能在 const 块内使用；每进入一个新的 const 块，iota 重置为 0
	const (
		Sunday    = iota // 0
		Monday           // 1（表达式复用上一行，iota 自动加 1）
		Tuesday          // 2
		Wednesday        // 3
		Thursday         // 4
		Friday           // 5
		Saturday         // 6
	)
	t.Logf("iota 自增: Sunday=%d Monday=%d Saturday=%d", Sunday, Monday, Saturday)
	if Sunday != 0 || Monday != 1 || Saturday != 6 {
		t.Error("iota 自增值不符合预期")
	}

	// ── 位移写法实现存储单位 ──────────────────────────────────────────
	// 用 _ 跳过 iota=0，让 KB 从 iota=1 开始，对齐二进制量级：
	//   _  : iota=0，忽略（若不跳过 KB=1<<0=1，语义不对）
	//   KB : iota=1 → 1 << (1*10) = 1 << 10 = 1 024
	//   MB : iota=2 → 1 << (2*10) = 1 << 20 = 1 048 576
	//   GB : iota=3 → 1 << (3*10) = 1 << 30 = 1 073 741 824
	const (
		_  = iota             // 跳过 0
		KB = 1 << (iota * 10) // 1 << 10 = 1 024
		MB                    // 1 << 20 = 1 048 576（表达式自动复用，iota=2）
		GB                    // 1 << 30 = 1 073 741 824（iota=3）
	)
	t.Logf("存储单位: KB=%d, MB=%d, GB=%d", KB, MB, GB)
	if KB != 1024 {
		t.Errorf("KB 期望 1024，得到 %d", KB)
	}
	if MB != 1024*1024 {
		t.Errorf("MB 期望 %d，得到 %d", 1024*1024, MB)
	}
	if GB != 1024*1024*1024 {
		t.Errorf("GB 期望 %d，得到 %d", 1024*1024*1024, GB)
	}
}
