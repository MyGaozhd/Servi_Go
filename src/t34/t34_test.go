package t34

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Test34_0(t *testing.T) {
	t.Log("start->", runtime.NumGoroutine())
	t.Log(allResponse1())
	t.Log("end->", runtime.NumGoroutine())

	//待完善，不满足原子性、可见性、一致性
	t.Log("start->", runtime.NumGoroutine())
	t.Log(allResponse2())
	t.Log("end->", runtime.NumGoroutine())
}

func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("the result is from %d", id)
}

/**
  第一种方式
  需要全部任务完成才返回
*/
func allResponse1() string {
	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}

	finalRet := ""
	for i := 0; i < 10; i++ {
		finalRet += <-ch + "\n"
	}

	return finalRet
}

/**
  第二种方式 -- 还有问题，不满足原子性、可见性、有序性
  需要全部任务完成才返回
*/
func allResponse2() string {
	var wg sync.WaitGroup
	wg.Add(10)
	finalRet := ""
	for i := 0; i < 10; i++ {
		go func(i int) {
			finalRet += runTask(i) + "\n"
			wg.Done()
		}(i)
	}
	wg.Wait()
	return finalRet
}
