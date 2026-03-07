package t0

import (
	"testing"
)

// TestFibonacci 测试斐波那契数列的生成
func TestFibonacci(t *testing.T) {
	// 期望的斐波那契数列前6个数
	expected := []int{1, 1, 2, 3, 5, 8}

	// 生成斐波那契数列
	result := generateFibonacci(6)

	// 验证结果
	if len(result) != len(expected) {
		t.Errorf("期望长度 %d, 实际长度 %d", len(expected), len(result))
	}

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Errorf("索引 %d: 期望 %d, 实际 %d", i, expected[i], result[i])
		}
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
