// Package t8 演示用 map 实现 Set 数据结构。
// 原始散写的 map[int]bool 示例已封装为 IntSet 结构体，并补充交集运算测试。
package t8

import (
	"sort"
	"testing"
)

// ─────────────────────────────────────────────
// IntSet：基于 map 实现的整数集合
// ─────────────────────────────────────────────

// IntSet 是用 map[int]struct{} 实现的整数集合。
// 使用空结构体 struct{} 作为 value，相比 bool 节省内存（struct{} 占 0 字节）。
type IntSet struct {
	data map[int]struct{}
}

// NewIntSet 创建并返回一个空的 IntSet。
func NewIntSet() *IntSet {
	return &IntSet{data: make(map[int]struct{})}
}

// Add 向集合中添加元素；若元素已存在，则为幂等操作（不报错）。
func (s *IntSet) Add(val int) {
	s.data[val] = struct{}{}
}

// Remove 从集合中删除元素；若元素不存在，则为幂等操作（不报错）。
func (s *IntSet) Remove(val int) {
	delete(s.data, val)
}

// Contains 判断集合中是否存在指定元素。
func (s *IntSet) Contains(val int) bool {
	_, ok := s.data[val]
	return ok
}

// Len 返回集合中元素的个数。
func (s *IntSet) Len() int {
	return len(s.data)
}

// Values 返回集合所有元素的排序切片（保证顺序确定性，方便测试断言）。
func (s *IntSet) Values() []int {
	result := make([]int, 0, len(s.data))
	for v := range s.data {
		result = append(result, v)
	}
	sort.Ints(result)
	return result
}

// Intersection 返回两个集合的交集（新集合，不修改原集合）。
// 遍历较小集合的每个元素，判断是否存在于较大集合，时间复杂度 O(min(m,n))。
func Intersection(a, b *IntSet) *IntSet {
	// 保证遍历较小的集合
	small, large := a, b
	if a.Len() > b.Len() {
		small, large = b, a
	}

	result := NewIntSet()
	for v := range small.data {
		if large.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// ─────────────────────────────────────────────

// Test8_0 保留原有演示：用 map[int]bool 实现 set 的散写方式。
func Test8_0(t *testing.T) {

	set := map[int]bool{}
	set[1] = true

	if set[1] {
		t.Log("已经存在")
	} else {
		t.Log("不存在")
	}

	t.Log(len(set))

	delete(set, 1)

	if set[1] {
		t.Log("已经存在")
	} else {
		t.Log("不存在")
	}
}

// TestIntSet 完整测试 IntSet 的所有操作：
// Add、Contains、Remove、Len，以及边界情况（重复添加/删除不存在的元素）。
func TestIntSet(t *testing.T) {
	s := NewIntSet()

	// 初始状态：空集合
	if s.Len() != 0 {
		t.Errorf("新建集合 Len 应为 0，实际 %d", s.Len())
	}
	t.Log("新建集合为空，Len =", s.Len())

	// —— Add：添加元素 ——
	s.Add(1)
	s.Add(2)
	s.Add(3)
	t.Logf("Add(1,2,3) 后，Values = %v，Len = %d", s.Values(), s.Len())
	if s.Len() != 3 {
		t.Errorf("期望 Len=3，实际 %d", s.Len())
	}

	// 边界 case：重复添加，Len 不变
	s.Add(2)
	s.Add(2)
	t.Logf("重复 Add(2) 两次后，Len = %d（应仍为 3）", s.Len())
	if s.Len() != 3 {
		t.Errorf("重复添加后期望 Len=3，实际 %d", s.Len())
	}

	// —— Contains：判断存在 ——
	for _, v := range []int{1, 2, 3} {
		if !s.Contains(v) {
			t.Errorf("Contains(%d) 应为 true", v)
		}
	}
	if s.Contains(99) {
		t.Error("Contains(99) 应为 false")
	}
	t.Log("Contains(1/2/3)=true，Contains(99)=false ✓")

	// —— Remove：删除元素 ——
	s.Remove(2)
	t.Logf("Remove(2) 后，Values = %v，Len = %d", s.Values(), s.Len())
	if s.Contains(2) {
		t.Error("Remove(2) 后 Contains(2) 应为 false")
	}
	if s.Len() != 2 {
		t.Errorf("Remove 后期望 Len=2，实际 %d", s.Len())
	}

	// 边界 case：删除不存在的元素，不报错（幂等）
	s.Remove(99)
	t.Logf("Remove(99)（不存在）后，Len = %d（应仍为 2）", s.Len())
	if s.Len() != 2 {
		t.Errorf("删除不存在元素后期望 Len=2，实际 %d", s.Len())
	}

	// —— Values：返回排序切片 ——
	vals := s.Values()
	expected := []int{1, 3}
	for i, v := range expected {
		if vals[i] != v {
			t.Errorf("Values[%d]: 期望 %d，实际 %d", i, v, vals[i])
		}
	}
	t.Logf("Values() = %v（已排序）✓", vals)
}

// TestSetIntersection 演示求两个集合的交集。
//
//	{1,2,3} ∩ {2,3,4} = {2,3}
//
// 使用 Intersection 函数，遍历较小集合以提高效率。
func TestSetIntersection(t *testing.T) {
	a := NewIntSet()
	for _, v := range []int{1, 2, 3} {
		a.Add(v)
	}

	b := NewIntSet()
	for _, v := range []int{2, 3, 4} {
		b.Add(v)
	}

	inter := Intersection(a, b)
	got := inter.Values()
	t.Logf("{1,2,3} ∩ {2,3,4} = %v", got)

	expected := []int{2, 3}
	if len(got) != len(expected) {
		t.Errorf("交集长度：期望 %d，实际 %d", len(expected), len(got))
	}
	for i, v := range expected {
		if got[i] != v {
			t.Errorf("交集[%d]: 期望 %d，实际 %d", i, v, got[i])
		}
	}

	// 边界 case：两个不相交集合，交集为空
	c := NewIntSet()
	for _, v := range []int{10, 20} {
		c.Add(v)
	}
	empty := Intersection(a, c)
	t.Logf("{1,2,3} ∩ {10,20} = %v（空集）", empty.Values())
	if empty.Len() != 0 {
		t.Errorf("不相交集合的交集应为空，实际 Len=%d", empty.Len())
	}

	// 边界 case：与自身求交集 = 自身
	self := Intersection(a, a)
	t.Logf("{1,2,3} ∩ {1,2,3} = %v（自身）", self.Values())
	if self.Len() != a.Len() {
		t.Errorf("与自身交集长度应为 %d，实际 %d", a.Len(), self.Len())
	}

	// 边界 case：与空集合求交集 = 空集
	emptySet := NewIntSet()
	result := Intersection(a, emptySet)
	t.Logf("{1,2,3} ∩ {} = %v（空集）", result.Values())
	if result.Len() != 0 {
		t.Errorf("与空集合交集应为空，实际 Len=%d", result.Len())
	}
}
