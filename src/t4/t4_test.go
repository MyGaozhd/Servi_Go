package t4

import (
	"fmt"
	"testing"
)

// Test4_0 演示切片的基本用法
// 注意 slice 和数组在声明时的区别：声明数组时方括号内写明长度或用 ...，
// 声明 slice 时方括号内没有任何字符。
func Test4_0(t *testing.T) {
	// 初始化
	var s0 []int
	t.Log(len(s0), cap(s0))

	// 添加元素
	s0 = append(s0, 1)
	t.Log(len(s0), cap(s0))

	s1 := []int{0, 1, 2, 3}
	t.Log(len(s1), cap(s1))

	s2 := make([]int, 3, 5)
	t.Log(len(s2), cap(s2))
	t.Log(s2[0], s2[1], s2[2])

	s2 = append(s2, 3)
	t.Log(s2[0], s2[1], s2[2], s2[3])
}

// Test4_1 演示切片容量的增长规律：
// 当切片的 len == cap 时，再向切片中添加元素，cap 会增长（小 slice 通常翻倍）
//
//	t4_test.go:31: 1 1
//	t4_test.go:31: 2 2
//	t4_test.go:31: 3 4
//	t4_test.go:31: 4 4
//	t4_test.go:31: 5 8
//	t4_test.go:31: 6 8
//	t4_test.go:31: 7 8
//	t4_test.go:31: 8 8
//	t4_test.go:31: 9 16
//	t4_test.go:31: 10 16
func Test4_1(t *testing.T) {
	s := []int{}
	for i := 0; i < 10; i++ {
		s = append(s, i)
		t.Log(len(s), cap(s))
	}
}

// Test4_2 演示切片截取
func Test4_2(t *testing.T) {
	s := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}

	q1 := s[3:6]
	t.Log(q1, len(q1), cap(q1))
}

// TestSliceCopy 演示 copy() 与共享底层数组的陷阱：
//   - 直接 s2 := s1[:] 共享底层数组，修改 s2 元素会影响 s1
//   - copy(dst, src) 创建真正独立的副本，修改互不影响
//   - 通过打印地址验证两者底层数组指针不同
func TestSliceCopy(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}

	// ── 共享底层数组的陷阱 ─────────────────────────────────────────────
	// s2 := s1[:] 只拷贝了 slice header（指针、len、cap），
	// 底层数组与 s1 共享——修改 s2[0] 会同时影响 s1[0]
	s2 := s1[:]
	s2[0] = 999
	t.Logf("修改 s2[0]=999 后: s1=%v, s2=%v（s1 也被修改了！）", s1, s2)
	if s1[0] != 999 {
		t.Error("s1 和 s2 应共享底层数组，s1[0] 期望为 999")
	}

	// 恢复 s1
	s1[0] = 1

	// ── copy() 创建真正独立副本 ───────────────────────────────────────
	// copy(dst, src) 会把 src 的内容逐元素复制到 dst，
	// dst 必须是已分配空间的 slice（用 make 提前分配）
	s3 := make([]int, len(s1))
	n := copy(s3, s1)
	t.Logf("copy 复制了 %d 个元素: s3=%v", n, s3)

	s3[0] = 888
	t.Logf("修改 s3[0]=888 后: s1=%v, s3=%v（s1 未受影响）", s1, s3)
	if s1[0] != 1 {
		t.Errorf("copy 后修改 s3 不应影响 s1，s1[0] 期望 1，得到 %d", s1[0])
	}

	// ── 验证底层数组地址不同 ──────────────────────────────────────────
	// fmt.Sprintf("%p", slice) 打印 slice 底层数组首元素地址
	addrS1 := fmt.Sprintf("%p", s1)
	addrS3 := fmt.Sprintf("%p", s3)
	t.Logf("s1 底层数组地址: %s", addrS1)
	t.Logf("s3 底层数组地址: %s", addrS3)
	if addrS1 == addrS3 {
		t.Error("copy 后 s1 和 s3 的底层数组地址应不同")
	}
}

// TestSliceDelete 演示从切片中删除元素的惯用法：
//   - 删除下标 i 的元素：append(s[:i], s[i+1:]...)
//   - 演示删除前后的 slice 内容，以及删除第一个、中间、最后一个元素的情况
func TestSliceDelete(t *testing.T) {
	// ── 删除中间元素（下标 2）──────────────────────────────────────────
	// append(s[:i], s[i+1:]...) 会将 i+1 之后的元素"左移"覆盖到 i 的位置
	// 注意：这会修改原 slice 的底层数组内容，且 len 减 1
	s := []int{0, 1, 2, 3, 4}
	delIdx := 2
	s = append(s[:delIdx], s[delIdx+1:]...)
	t.Logf("删除下标 %d 后: %v（期望 [0 1 3 4]）", delIdx, s)
	if len(s) != 4 || s[2] != 3 {
		t.Errorf("删除后 slice 不符合预期，得到 %v", s)
	}

	// ── 删除第一个元素 ────────────────────────────────────────────────
	s2 := []string{"a", "b", "c", "d"}
	s2 = append(s2[:0], s2[1:]...)
	t.Logf("删除第一个元素后: %v（期望 [b c d]）", s2)
	if len(s2) != 3 || s2[0] != "b" {
		t.Errorf("删除第一个元素后不符合预期，得到 %v", s2)
	}

	// ── 删除最后一个元素 ──────────────────────────────────────────────
	// 直接截断 s[:len(s)-1] 即可，更高效
	s3 := []int{10, 20, 30, 40}
	s3 = s3[:len(s3)-1]
	t.Logf("删除最后一个元素后: %v（期望 [10 20 30]）", s3)
	if len(s3) != 3 || s3[len(s3)-1] != 30 {
		t.Errorf("删除最后一个元素后不符合预期，得到 %v", s3)
	}
}

// TestSliceGrowth 演示切片的扩容倍数规律：
//   - 当 len == cap 时继续 append，运行时会分配更大的底层数组（小 slice 约翻倍）
//   - 打印每次 append 后的 len 和 cap，观察扩容时机
//   - 扩容策略在 Go 1.18 后有所调整（大 slice 增长率会逐渐降低）
func TestSliceGrowth(t *testing.T) {
	s := make([]int, 0)
	prevCap := 0

	t.Log("len\tcap\t是否扩容")
	for i := 0; i < 20; i++ {
		s = append(s, i)
		grew := ""
		if cap(s) != prevCap {
			grew = "<-- 扩容"
			prevCap = cap(s)
		}
		t.Logf("%d\t%d\t%s", len(s), cap(s), grew)
	}

	// 验证最终 len 正确
	if len(s) != 20 {
		t.Errorf("期望 len=20，得到 %d", len(s))
	}

	// cap 一定 >= len
	if cap(s) < len(s) {
		t.Errorf("cap(%d) 不应小于 len(%d)", cap(s), len(s))
	}

	// ── 预分配容量避免频繁扩容 ────────────────────────────────────────
	// 如果已知元素个数，用 make([]T, 0, n) 预分配容量，
	// 可以避免多次扩容带来的内存分配和数据拷贝
	preallocated := make([]int, 0, 20)
	initialCap := cap(preallocated)
	for i := 0; i < 20; i++ {
		preallocated = append(preallocated, i)
	}
	t.Logf("预分配容量 %d，填满 20 个元素后 cap=%d（无扩容）", initialCap, cap(preallocated))
	if cap(preallocated) != 20 {
		t.Errorf("预分配后不应扩容，期望 cap=20，得到 %d", cap(preallocated))
	}
}
