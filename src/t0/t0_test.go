package t0

import (
	"testing"
)

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
