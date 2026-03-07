package t13

import (
	"fmt"
	"testing"
)

func Test13_0(t *testing.T) {
	//用来释放资源释放锁，类似java的finally
	defer clear()
	// defer + recover 捕获 panic，防止测试崩溃（类似 java 的 try-finally）
	defer func() {
		if err := recover(); err != nil {
			t.Log("recover from panic:", err)
		}
	}()
	t.Log("start")
	//抛出异常
	panic("error")
}

func clear() {
	fmt.Print("clear resources")
}
