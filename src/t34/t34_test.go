// Package t34 演示"需要所有任务完成才返回"的两种并发模式。
//
// # 两种实现方式对比
//
// 方式一（allResponse1，channel 聚合）：无数据竞争，推荐
// 方式二（allResponse2，WaitGroup + 共享变量）：存在竞争隐患，用于对比说明
package t34

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// Test34_0 演示并对比两种"等待所有任务"的实现方式。
func Test34_0(t *testing.T) {
	t.Log("allResponse1 start->", runtime.NumGoroutine())
	t.Log(allResponse1())
	t.Log("allResponse1 end->", runtime.NumGoroutine())

	t.Log("allResponse2 start->", runtime.NumGoroutine())
	t.Log(allResponse2())
	t.Log("allResponse2 end->", runtime.NumGoroutine())
}

// runTask 模拟耗时任务，返回包含 id 的结果字符串
func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("the result is from %d", id)
}

// allResponse1 方式一：用有缓冲 channel 汇总所有结果（推荐）。
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

// allResponse2 方式二：用 WaitGroup 等待，共享变量聚合结果（存在竞争）。
//
// 问题：goroutine 间的"读旧值→加工→写新值"不满足原子性，可能导致结果丢失。
func allResponse2() string {
	var wg sync.WaitGroup
	wg.Add(10)
	finalRet := ""
	for i := 0; i < 10; i++ {
		go func(i int) {
			finalRet = syncBuildString(finalRet, runTask(i)+"\n")
			wg.Done()
		}(i)
	}
	wg.Wait()
	return finalRet
}

var mu sync.Mutex

// syncBuildString 在锁内拼接两个字符串（但调用者读取 finalRet 在锁外，仍有竞争）。
func syncBuildString(str1 string, str2 string) string {
	mu.Lock()
	defer mu.Unlock()
	return str1 + str2
}
