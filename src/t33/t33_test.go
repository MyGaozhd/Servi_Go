package t33

import (
	"fmt"
	"testing"
	"time"
)

func Test33_0(t *testing.T) {
	t.Log(firstResponse())
}

func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("the result is from %d", id)
}

/**
  只需要任意任务完成就返回
*/
func firstResponse() string {
	ch := make(chan string)
	for i := 0; i < 10; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}

	return <-ch
}
