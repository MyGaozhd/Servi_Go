// Package t28 演示 select 多路 channel 选择的用法。
//
// # select 语法
//
//	select {
//	case v := <-ch1:  // ch1 有数据可读时执行
//	    ...
//	case ch2 <- v:    // ch2 可写时执行
//	    ...
//	case <-time.After(d): // 超时
//	    ...
//	default:          // 所有 case 都未就绪时立即执行（非阻塞）
//	    ...
//	}
//
// # 特性
//
//   - 没有 default 时，select 阻塞直到至少一个 case 就绪
//   - 有 default 时，所有 case 未就绪则立即走 default（非阻塞）
//   - 多个 case 同时就绪时，Go 运行时伪随机选一个（无优先级）
//   - select{} 空语句永久阻塞（常用于防止 main 退出）
package t28

import (
	"fmt"
	"testing"
	"time"
)

// Test28_0 演示 select + time.After 实现超时控制。
func Test28_0(t *testing.T) {
	select {
	case retch := <-asyncService():
		t.Log("服务返回:", retch)
	case <-time.After(time.Millisecond * 500):
		t.Log("超时：服务未在 500ms 内响应")
	}
}

// syncService 模拟耗时 1 秒的同步服务
func syncService() string {
	time.Sleep(time.Second * 1)
	return "service done!"
}

// asyncService 封装异步调用，用有缓冲 channel 避免 goroutine 泄露
func asyncService() chan string {
	retch := make(chan string, 1)
	go func() {
		ret := syncService()
		fmt.Println("syncService returned")
		retch <- ret
		fmt.Println("AsyncService exit")
	}()
	return retch
}

// Test28_1 演示 select + default 实现非阻塞 channel 操作。
func Test28_1(t *testing.T) {
	t.Run("channel无数据走default", func(t *testing.T) {
		ch := make(chan int, 1)
		select {
		case v := <-ch:
			t.Errorf("不应该收到数据: %d", v)
		default:
			t.Log("channel 为空，走 default（非阻塞）")
		}
	})

	t.Run("channel有数据走case", func(t *testing.T) {
		ch := make(chan int, 1)
		ch <- 42
		select {
		case v := <-ch:
			t.Logf("收到数据: %d", v)
		default:
			t.Error("不应该走 default")
		}
	})

	t.Run("非阻塞轮询", func(t *testing.T) {
		ch := make(chan string, 1)
		received := false
		go func() {
			time.Sleep(50 * time.Millisecond)
			ch <- "result"
		}()
		for i := 0; i < 20; i++ {
			select {
			case msg := <-ch:
				t.Logf("第 %d 次轮询，收到: %s", i+1, msg)
				received = true
			default:
				time.Sleep(10 * time.Millisecond)
			}
			if received {
				break
			}
		}
		if !received {
			t.Error("轮询结束但未收到数据")
		}
	})
}

// Test28_2 演示 select 多个 case 同时就绪时的随机选择行为。
//
// 关键结论：
//   - 当多个 case 同时就绪，select 伪随机选一个（无优先级）
//   - Go 规范明确说明：多个 case 就绪时选择是随机的
//
// 实现说明：
//   - ch1/ch2 缓冲容量均为 1
//   - 每轮循环：先各填一条数据（两 case 同时就绪），select 随机取其一，再排空另一个
func Test28_2(t *testing.T) {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	count1, count2 := 0, 0
	const rounds = 1000

	for i := 0; i < rounds; i++ {
		// 每轮循环内填充，保证两个 case 同时就绪
		ch1 <- "来自 ch1"
		ch2 <- "来自 ch2"

		select {
		case v := <-ch1:
			_ = v
			count1++
			<-ch2 // 排空另一个 channel
		case v := <-ch2:
			_ = v
			count2++
			<-ch1
		}
	}

	t.Logf("执行 %d 轮，ch1 被选中 %d 次，ch2 被选中 %d 次", rounds, count1, count2)
	t.Logf("分布比例：ch1=%.1f%%，ch2=%.1f%%", float64(count1)/rounds*100, float64(count2)/rounds*100)

	if count1 == 0 || count2 == 0 {
		t.Error("预期两个 case 都会被随机选中，但有一个从未被选中")
	} else {
		t.Log("✓ 两个 case 均被随机选中，验证了 select 的随机性")
	}
}
