// Package t26 演示无缓冲 channel 的同步语义。
//
// # 无缓冲 Channel 基础
//
// make(chan T) 创建一个容量为 0 的无缓冲 channel：
//   - 发送方（ch <- v）：阻塞，直到有接收方准备好接收
//   - 接收方（v := <-ch）：阻塞，直到有发送方准备好发送
//   - 本质：一个"同步点"——发送和接收必须同时就绪，数据才能传递
//
// # 与有缓冲 Channel 对比
//
//	无缓冲：make(chan T)     → 发送立即阻塞，必须等接收方
//	有缓冲：make(chan T, n)  → 缓冲未满时发送不阻塞
//
// # 适用场景
//
//   - 等待 goroutine 完成某项操作（相当于信号量）
//   - goroutine 间的"握手"同步
//   - 实现请求-响应模式
package t26

import (
	"fmt"
	"testing"
	"time"
)

// Test26_0 演示无缓冲 channel 的异步服务模式。
//
// 流程说明：
//  1. asyncService 在新 goroutine 中执行耗时操作，通过无缓冲 channel 返回结果
//  2. 主 goroutine 调用 asyncService 后，立即执行 otherService（并发）
//  3. <-retch 阻塞等待结果——此时无缓冲 channel 确保了同步
//
// 关键点：asyncService 内的 retch <- ret 必须等到主 goroutine 执行 <-retch 才能继续，
// 因此 "AsyncService exit" 的打印一定在 <-retch 之后。
func Test26_0(t *testing.T) {
	// 启动异步服务（内部开 goroutine，返回无缓冲 channel）
	retch := asyncService()

	// 与 asyncService 并发执行（此时 asyncService 的 goroutine 也在跑）
	otherService()

	// 阻塞，直到 asyncService 发送结果
	// 无缓冲：此处接收 <-retch 与 goroutine 内的 retch <- ret 形成同步点
	result := <-retch
	t.Log("asyncService 返回:", result)
}

// otherService 模拟另一个耗时 1 秒的服务（与 asyncService 并发执行）
func otherService() {
	time.Sleep(time.Second * 1)
	fmt.Println("otherService done!")
}

// syncService 模拟耗时 1 秒的同步服务，返回字符串结果
func syncService() string {
	time.Sleep(time.Second * 1)
	return "service done!"
}

// asyncService 将 syncService 包装为异步调用，通过无缓冲 channel 传递结果
func asyncService() chan string {
	// 无缓冲 channel：发送方必须等待接收方就绪
	retch := make(chan string)
	go func() {
		ret := syncService()
		fmt.Println("syncService returned")
		// 此行阻塞，直到主 goroutine 执行 <-retch
		retch <- ret
		// 只有在主 goroutine 接收之后，才会执行到这里
		fmt.Println("AsyncService exit")
	}()
	return retch
}

// Test26_1 演示无缓冲 channel 的同步语义细节。
//
// 核心特性验证：
//  1. 发送方阻塞，直到有接收方
//  2. 接收方阻塞，直到有发送方
//  3. 无缓冲 channel 是 goroutine 间的同步点（rendezvous point）
//
// 对比实验：先收后发 vs 先发后收，体现双向阻塞特性。
func Test26_1(t *testing.T) {
	// ── 场景 1：接收方先等待，发送方随后发送 ──
	t.Run("接收方先阻塞，等发送方", func(t *testing.T) {
		ch := make(chan int) // 无缓冲
		done := make(chan struct{})

		go func() {
			// 接收方：阻塞直到有人发送
			val := <-ch
			t.Logf("接收到值: %d", val)
			close(done)
		}()

		// 稍等，让接收方先进入阻塞状态
		time.Sleep(5 * time.Millisecond)

		// 发送：此时接收方已在等待，发送立即成功
		t.Log("开始发送...")
		ch <- 42
		<-done // 等待接收方打印完毕
		t.Log("发送完成，接收方已退出阻塞")
	})

	// ── 场景 2：发送方先阻塞，接收方随后接收 ──
	t.Run("发送方先阻塞，等接收方", func(t *testing.T) {
		ch := make(chan int) // 无缓冲
		sent := make(chan struct{})

		go func() {
			// 发送方：阻塞直到有人接收
			t.Log("发送方尝试发送，进入阻塞...")
			ch <- 99 // 阻塞，直到主 goroutine 接收
			t.Log("发送方解除阻塞，发送成功")
			close(sent)
		}()

		// 稍等，让发送方先进入阻塞状态
		time.Sleep(5 * time.Millisecond)

		// 接收：此时发送方已在阻塞，接收立即成功
		val := <-ch
		t.Logf("接收到值: %d", val)
		<-sent
	})

	// ── 场景 3：用无缓冲 channel 作为 goroutine 完成信号 ──
	t.Run("用channel作为完成信号", func(t *testing.T) {
		// 惯用法：chan struct{} 作为纯信号，不传数据，内存零开销
		done := make(chan struct{})

		go func() {
			time.Sleep(10 * time.Millisecond)
			t.Log("goroutine 工作完成")
			close(done) // 关闭 channel 广播完成信号
		}()

		// 等待 goroutine 完成信号
		<-done
		t.Log("主 goroutine 收到完成信号")
	})
}
