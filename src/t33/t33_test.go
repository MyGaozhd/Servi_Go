// Package t33 演示"只需任意一个任务完成即可返回"的并发模式（对冲请求）。
//
// # goroutine 泄露问题与解决方案
//
// 写法一（无缓冲 channel，会泄露）：
//
//	ch := make(chan string)
//	// 只有第一个 goroutine 能发送，剩余 9 个永远阻塞 → 泄露
//
// 写法二（有缓冲 channel，容量 >= goroutine 数，推荐）：
//
//	ch := make(chan string, 10)
//	// 所有 goroutine 完成后都能写入缓冲区，不会阻塞 → 不泄露
package t33

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// Test33_0 演示"取最快响应"模式。
func Test33_0(t *testing.T) {
	t.Log("start->", runtime.NumGoroutine())
	t.Log(firstResponse())
	time.Sleep(time.Millisecond * 50)
	t.Log("end->", runtime.NumGoroutine())
}

// runTask 模拟耗时任务，返回包含 id 的结果字符串
func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("the result is from %d", id)
}

// firstResponse 启动 10 个并发任务，返回最先完成的那个结果。
//
// 关键设计：channel 缓冲容量 = goroutine 数量（10），确保所有 goroutine
// 完成后都能写入缓冲区后退出，不会因主 goroutine 已返回而永久阻塞。
func firstResponse() string {
	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}
	return <-ch
}
