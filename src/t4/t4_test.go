package t4

import "testing"

// 注意slice和数组在声明时的区别：声明数组时，方括号内写明了数组的长度或使用...自动计算长度，而声明slice时，方括号内没有任何字符。
//切片的用法
func Test4_0(t *testing.T) {
	//初始化
	var s0 []int
	t.Log(len(s0), cap(s0))

	//添加元素
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

//切片容量的增长规律，当切片的len=cap的时候，再向切片中添加元素，此时cap会增长为原始cap的两倍
/**
  t4_test.go:31: 1 1
  t4_test.go:31: 2 2
  t4_test.go:31: 3 4
  t4_test.go:31: 4 4
  t4_test.go:31: 5 8
  t4_test.go:31: 6 8
  t4_test.go:31: 7 8
  t4_test.go:31: 8 8
  t4_test.go:31: 9 16
  t4_test.go:31: 10 16
*/
func Test4_1(t *testing.T) {
	s := []int{}
	for i := 0; i < 10; i++ {
		s = append(s, i)
		t.Log(len(s), cap(s))
	}
}

func Test4_2(t *testing.T) {
	s := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}

	q1 := s[3:6]
	t.Log(q1, len(q1), cap(q1))
}
