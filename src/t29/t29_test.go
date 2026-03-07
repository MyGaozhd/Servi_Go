// Package t29 演示 channel 的关闭机制与广播通知。
//
// # Channel 关闭规则
//
//   - close(ch)：关闭 channel，通知接收方没有更多数据
//   - 关闭后接收：立即返回零值（不阻塞），ok == false
//   - 关闭后发送：panic（send on closed channel）
//   - 只有发送方应该关闭 channel，接收方不应关闭
//   - close 是广播：所有正在等待的接收方都会被唤醒
//
// # 接收方检测关闭
//
// 方式一：ok 检测
//
//	data, ok := <-ch
//	if !ok { /* channel 已关闭 */ }
//
// 方式二：for-range（推荐）
//
//	for v := range ch { /* channel 关闭后自动退出循环 */ }
package t29

import (
	"fmt"
	"sync"
	"testing"
)

// Test29_0 演示生产者关闭 channel，多个消费者通过 ok 检测退出。
//
// 场景：
//   - 1 个生产者发送 0~9，发完后 close(ch)
//   - 2 个消费者并发从 ch 读取，用 data, ok := <-ch 判断是否结束
//   - close 是广播操作：ch 关闭后，所有阻塞的消费者都会被唤醒
//
// 注意：多个消费者竞争读取同一 channel，每条数据只被一个消费者消费（非广播数据）。
// close 广播的是"关闭信号"，而不是数据本身。
func Test29_0(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int)

	// 生产者：发送 10 个数，然后关闭 channel
	wg.Add(1)
	dataProduce(ch, &wg)

	// 两个消费者：并发读取，通过 ok 检测 channel 关闭
	wg.Add(2)
	dataReseiver(ch, &wg, 1)
	dataReseiver(ch, &wg, 2)

	// 等待所有生产者和消费者完成
	wg.Wait()
	t.Log("所有 goroutine 已完成")
}

// dataProduce 生产者：向 ch 发送 0~9，发完后关闭 channel
func dataProduce(ch chan int, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- i
		}
		// 关闭 channel：广播通知所有消费者没有更多数据
		// 只能由发送方关闭，多次 close 会 panic
		close(ch)
	}()
}

// dataReseiver 消费者：循环读取 ch，用 ok 判断 channel 是否已关闭
//
// ok 的语义：
//   - ok == true：成功接收到发送方发送的数据
//   - ok == false：channel 已关闭且缓冲区为空，data 为零值
func dataReseiver(ch chan int, wg *sync.WaitGroup, count int) {
	go func() {
		defer wg.Done()
		for {
			// data, ok := <-ch：ok 为 false 表示 channel 已关闭
			if data, ok := <-ch; ok {
				fmt.Println(count, data)
			} else {
				// channel 已关闭，退出循环
				break
			}
		}
	}()
}

// Test29_1 演示用 for-range 遍历 channel，比 ok 检测更简洁。
//
// for-range 遍历 channel 的特性：
//   - 自动接收 channel 中的数据，无需手动 ok 检测
//   - channel 关闭且缓冲区读完后，循环自动退出
//   - 如果 channel 未关闭，for-range 永远阻塞（goroutine 泄露！）
//
// 对比 ok 检测方式：
//
//	ok 检测：需要手动 break，代码略繁琐
//	for-range：自动退出，代码更简洁，推荐在"channel 会被关闭"的场景使用
func Test29_1(t *testing.T) {
	// 方式一：ok 检测（原始写法，与 Test29_0 同款）
	t.Run("ok检测方式", func(t *testing.T) {
		ch := make(chan int, 5)
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch) // 必须关闭，否则接收方不知道何时停止

		for {
			if v, ok := <-ch; ok {
				t.Logf("  ok方式收到: %d", v)
			} else {
				break // channel 关闭，手动退出
			}
		}
	})

	// 方式二：for-range（推荐写法）
	t.Run("forRange方式", func(t *testing.T) {
		ch := make(chan int, 5)
		for i := 0; i < 5; i++ {
			ch <- i
		}
		// 不 close 的话 for-range 永远阻塞 -> goroutine 泄露
		close(ch)

		// for-range 自动在 channel 关闭且读完后退出，无需手动 break
		for v := range ch {
			t.Logf("  range方式收到: %d", v)
		}
		t.Log("for-range 自动退出（channel 已关闭）")
	})

	// 多生产者场景：用 WaitGroup + 关闭 channel 配合 for-range
	t.Run("多生产者+forRange消费", func(t *testing.T) {
		ch := make(chan int, 20)
		var produceWg sync.WaitGroup

		// 3 个生产者并发生产，每人发 5 条
		for p := 0; p < 3; p++ {
			produceWg.Add(1)
			go func(p int) {
				defer produceWg.Done()
				for i := 0; i < 5; i++ {
					ch <- p*10 + i
				}
			}(p)
		}

		// 等所有生产者完成后关闭 channel（必须只关闭一次）
		go func() {
			produceWg.Wait()
			close(ch) // 所有生产者完成后由协调 goroutine 负责关闭
		}()

		// 消费者用 for-range，channel 关闭后自动退出
		total := 0
		for v := range ch {
			_ = v
			total++
		}
		// 3 个生产者各发 5 条，共 15 条
		t.Logf("消费者通过 for-range 共收到 %d 条数据（期望 15）", total)
		if total != 15 {
			t.Errorf("期望 15，实际 %d", total)
		}
	})
}
