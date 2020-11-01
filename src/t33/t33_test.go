package t33

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func Test33_0(t *testing.T) {
	t.Log("start->", runtime.NumGoroutine())
	t.Log(firstResponse())
	t.Log("end->", runtime.NumGoroutine())
}

func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("the result is from %d", id)
}

/**
  只需要任意任务完成就返回
*/
func firstResponse() string {
	/**
	写法一：
	  ch := make(chan string)此写法会出现协程泄露
	运行结果一：
	  t33_test.go:11: start-> 2
	  t33_test.go:12: the result is from 2
	  t33_test.go:13: end-> 11
	写法二：
	  ch := make(chan string, 10)防止协程泄露
	运行结果二：
	  t33_test.go:11: start-> 2
	  t33_test.go:12: the result is from 3
	  t33_test.go:13: end-> 2 **重点在此处**
	*/
	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}

	return <-ch
}
