package t2

import "testing"

//for循环测试
func Test2_0(t *testing.T) {
	a := 0
	/* 相当于java 的 while(a<5) */
	for a < 5 {
		a++
		t.Log(a)
	}
}

// swith 语句测试
func Test2_1(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch i {
		case 1, 2:
			t.Log("1,2->", i)
		case 3, 4:
			t.Log("3,4->", i)
		default:
			t.Log("0->", i)
		}
	}
}

// swith 语句测试
func Test2_2(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch {
		case i == 1 || i == 2:
			t.Log("1,2->", i)
		case i == 3 || i == 4:
			t.Log("3,4->", i)
		default:
			t.Log("0->", i)
		}
	}
}
