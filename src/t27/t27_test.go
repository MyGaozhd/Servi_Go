// Package t27 演示有缓冲 channel 的用法及 channel 方向类型。
//
// # 有缓冲 Channel
//
// make(chan T, n) 创建容量为 n 的有缓冲 channel：
//   - 缓冲未满时，发送方不阻塞，直接写入缓冲区返回
//   - 缓冲已满时，发送方阻塞（等待接收方消费）
//   - 缓冲为空时，接收方阻塞（等待发送方写入）
//
// # 与无缓冲 Channel 对比
//
//	make(chan T)    → 无缓冲，发送方必须等接收方（同步）
//	make(chan T, n) → 有缓冲，发送方在缓冲未满时不阻塞（解耦）
//
// # Channel 方向类型
//
//	chan T     双向 channel（可读可写）
//	chan<- T   只写 channel（只能发送）
//	<-chan T   只读 channel（只能接收）
//
// 方向类型常用于函数签名，编译器静态检查用途，提升代码安全性。
package t27

import (
	"fmt"
	"testing"
	"time"
)

// Test27_0 演示有缓冲 channel 的异步服务模式。
//
// 关键差异（对比 t26 无缓冲版本）：
//   - retch := make(chan string, 1) 缓冲容量为 1
//   - goroutine 执行 retch <- ret 时，缓冲未满，立即写入并继续执行
//   - "AsyncService exit" 的打印不再等待主 goroutine 接收，可以提前出现
//
// 这使得发送方和接收方在时间上完全解耦——发送方不需要等接收方就绪。
func Test27_0(t *testing.T) {
	// 启动异步服务（内部使用有缓冲 channel，发送不阻塞）
	retch := asyncService()

	// 与 asyncService 并发执行
	otherService()

	// 接收结果（此时缓冲中已有数据，立即返回）
	result := <-retch
	t.Log("asyncService 返回:", result)
}

// otherService 模拟另一个耗时 1 秒的服务
func otherService() {
	time.Sleep(time.Second * 1)
	fmt.Println("otherService done!")
}

// syncService 模拟耗时 1 秒的同步服务
func syncService() string {
	time.Sleep(time.Second * 1)
	return "service done!"
}

// asyncService 使用有缓冲 channel 封装异步调用。
// 与无缓冲版本区别：发送方写入缓冲后无需等待接收方，立即继续执行。
func asyncService() chan string {
	// 有缓冲 channel，容量为 1：发送不阻塞（只要缓冲未满）
	retch := make(chan string, 1)
	go func() {
		ret := syncService()
		fmt.Println("syncService returned")
		retch <- ret // 缓冲未满，立即写入，不阻塞
		// 有缓冲：这行会在发送后立即执行，不等主 goroutine 接收
		fmt.Println("AsyncService exit")
	}()
	return retch
}

// Test27_1 演示 channel 的方向类型，提升代码安全性和可读性。
//
// 方向类型的价值：
//   - chan<- T（只写）：传入函数后只能发送，编译器禁止接收操作
//   - <-chan T（只读）：传入函数后只能接收，编译器禁止发送操作
//   - 在函数签名中明确意图，让调用者一眼看懂数据流向
//   - 编译时静态检查，避免误操作
func Test27_1(t *testing.T) {
	// 创建双向 channel（可读可写），然后分别传给 producer 和 consumer
	ch := make(chan int, 5) // 有缓冲，避免 producer 阻塞

	done := make(chan struct{})

	// producer 只能向 ch 发送（chan<- int）
	go producer(ch, 5)

	// consumer 只能从 ch 接收（<-chan int）
	go func() {
		consumer(ch, t)
		close(done)
	}()

	<-done
	t.Log("生产者消费者演示完成")
}

// producer 只写 channel：向 ch 发送 n 个整数后关闭 channel。
// 函数签名中 chan<- int 表明此函数只能发送，不能接收。
// 若在函数体内写 <-ch（接收），编译器会报错：invalid operation。
func producer(ch chan<- int, n int) {
	for i := 0; i < n; i++ {
		ch <- i
		fmt.Printf("生产: %d\n", i)
	}
	close(ch) // 关闭 channel，通知消费者没有更多数据
}

// consumer 只读 channel：从 ch 接收数据直到 channel 关闭。
// 函数签名中 <-chan int 表明此函数只能接收，不能发送。
// 若在函数体内写 ch <- v（发送），编译器会报错：invalid operation。
func consumer(ch <-chan int, t *testing.T) {
	for v := range ch { // for-range 在 channel 关闭后自动退出
		t.Logf("消费: %d", v)
	}
}
