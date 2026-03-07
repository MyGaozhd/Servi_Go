// Package t30 演示通过关闭 channel 向多个 goroutine 广播取消信号。
//
// # 取消模式的演变
//
// 早期写法（发送 n 个取消信号，一个 goroutine 消费一个）：
//
//	func cancel_1(cancelChan chan struct{}) {
//	    cancelChan <- struct{}{} // 只能通知一个 goroutine
//	}
//
// 改进写法（close channel 广播，所有 goroutine 同时收到）：
//
//	func cancel_2(cancelChan chan struct{}) {
//	    close(cancelChan) // 所有正在 select <-cancelChan 的 goroutine 都会被唤醒
//	}
//
// # close 广播原理
//
// channel 关闭后，所有对该 channel 的接收操作立即返回零值（不再阻塞），
// select 中的 case <-cancelChan 因此可以持续命中，实现广播效果。
//
// # 与 context 的关系
//
// 本方案是 context 包的前身思路，实际生产推荐用 context.WithCancel（见 t31）：
//   - context 支持层级取消（父取消自动传播到子）
//   - context 携带截止时间、超时
//   - context 是标准库，API 统一，与 http/grpc 等无缝集成
package t30

import (
	"fmt"
	"testing"
	"time"
)

// Test30_0 演示用 close(channel) 向 5 个 goroutine 广播取消信号。
//
// 工作机制：
//  1. 5 个 goroutine 各自循环，每 5ms 调用 isCancelled 检查是否被取消
//  2. isCancelled 用 select+default 实现非阻塞检测：channel 关闭则返回 true
//  3. cancel_2 调用 close(cancelChan)，channel 关闭后所有 goroutine 的检测立即返回 true
//  4. 各 goroutine 收到取消信号后退出循环，打印 "Cancelled"
//
// 关键点：close 是广播操作，无论有多少个 goroutine 在监听，都能同时收到信号。
func Test30_0(t *testing.T) {
	// make(chan struct{}, 0) 等价于 make(chan struct{})，无缓冲
	// struct{} 是零大小类型，作为纯信号 channel，不传递数据，内存开销为零
	cancelChan := make(chan struct{}, 0)

	// 启动 5 个 goroutine，各自轮询取消信号
	for i := 0; i < 5; i++ {
		go func(i int, cancelCh chan struct{}) {
			for {
				// 非阻塞检测：channel 关闭则立即返回 true
				if isCancelled(cancelCh) {
					break
				}
				// 未取消，继续工作（模拟每 5ms 做一次工作）
				time.Sleep(time.Millisecond * 5)
			}
			fmt.Println(i, "Cancelled")
		}(i, cancelChan)
	}

	// 广播取消：close channel 通知所有 goroutine
	// close 只能调用一次，重复 close 会 panic
	cancel_2(cancelChan)

	// 等待所有 goroutine 打印完毕
	time.Sleep(time.Second * 1)
}

// isCancelled 非阻塞检测 channel 是否已关闭。
//
// select+default 是非阻塞接收的惯用写法：
//   - case <-cancelChan：channel 关闭后，接收立即成功（返回零值），返回 true
//   - default：channel 未关闭且无数据，立即走 default，返回 false
func isCancelled(cancelChan chan struct{}) bool {
	select {
	case <-cancelChan:
		// channel 已关闭（或有数据），取消信号到达
		return true
	default:
		// channel 未关闭，继续工作
		return false
	}
}

// cancel_2 通过关闭 channel 向所有监听者广播取消信号。
//
// 对比 cancel_1（逐一发送）：
//
//	cancel_1：每次只通知一个 goroutine（channel 每次只有一个接收者消费）
//	cancel_2：close 后所有监听者同时感知，真正的广播
//
// close 之后不能再向 channel 发送数据，否则 panic。
func cancel_2(cancelChan chan struct{}) {
	close(cancelChan)
}
