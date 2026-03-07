// Package t25 演示 sync.WaitGroup 的用法及注意事项。
//
// # WaitGroup 是什么
//
// sync.WaitGroup 是 Go 标准库提供的"等待一组 goroutine 完成"的同步原语：
//   - 内部维护一个整数计数器，初始为 0
//   - Add(n)：计数器 +n（必须在启动 goroutine 之前调用）
//   - Done()：计数器 -1（等价于 Add(-1)，通常配合 defer 使用）
//   - Wait()：阻塞直到计数器归零
//
// # 经典用法
//
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go func() {
//	    defer wg.Done()
//	    // ... 业务逻辑
//	}()
//	wg.Wait()
package t25

import (
	"sync"
	"testing"
)

// Test25_0 演示 WaitGroup + Mutex 组合：并发安全地对计数器执行 5000 次自增。
//
// 两个角色各司其职：
//   - sync.Mutex：保护共享变量 count，防止数据竞争
//   - sync.WaitGroup：等待所有 goroutine 完成，替代不可靠的 time.Sleep
//
// 最终 count 必然等于 5000，体现了两者结合的正确性。
func Test25_0(t *testing.T) {
	var mu sync.Mutex
	count := 0
	var wg sync.WaitGroup

	for i := 0; i < 5000; i++ {
		// ⚠️ Add 必须在 go 语句之前调用
		// 若写在 goroutine 内部，Wait() 可能在 Add() 之前执行，导致提前返回
		wg.Add(1)
		go func() {
			// defer 保证无论函数如何退出，Done 都会被调用（计数器 -1）
			defer wg.Done()
			// 加锁保护临界区，防止并发写入 count 时数据竞争
			mu.Lock()
			defer mu.Unlock()
			count++
		}()
	}

	// 阻塞直到 5000 个 goroutine 全部完成
	wg.Wait()

	// 有锁 + WaitGroup 保证结果精确为 5000
	if count != 5000 {
		t.Errorf("期望 count = 5000，实际 = %d", count)
	} else {
		t.Logf("count = %d ✓", count)
	}
}

// Test25_1 演示 WaitGroup 的关键注意事项与常见错误。
//
// 规则总结：
//  1. Add 必须在 goroutine 启动之前调用（否则 Wait 可能提前返回）
//  2. Done 等价于 Add(-1)，通常用 defer wg.Done() 确保调用
//  3. Wait 阻塞直到计数器归零
//  4. 计数器不能变为负数（否则 panic: sync: negative WaitGroup counter）
//  5. 同一个 WaitGroup 可以被复用，但必须等上一轮 Wait 返回后才能再次 Add
func Test25_1(t *testing.T) {
	// ── 演示 1：Done 等价于 Add(-1) ──
	t.Run("Done等价于Add(-1)", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(3)

		for i := 0; i < 3; i++ {
			go func(id int) {
				// wg.Done() 与 wg.Add(-1) 完全等价，Done 是语义更清晰的封装
				defer wg.Done()
				t.Logf("goroutine %d 执行完毕", id)
			}(i)
		}

		wg.Wait()
		t.Log("3 个 goroutine 全部完成")
	})

	// ── 演示 2：Add 必须在 go 语句之前，不能在 goroutine 内部调用 ──
	t.Run("Add必须在goroutine外部调用", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make([]int, 5)

		for i := 0; i < 5; i++ {
			wg.Add(1) // ✅ 正确：在 go 语句之前 Add
			go func(i int) {
				defer wg.Done()
				results[i] = i * 2
			}(i)

			// ❌ 错误示例（注释掉的反例）：
			// go func(i int) {
			//     wg.Add(1)  // 可能在 wg.Wait() 之后才执行 Add，导致 Wait 提前返回
			//     defer wg.Done()
			//     results[i] = i * 2
			// }(i)
		}

		wg.Wait()
		t.Logf("results = %v", results)
	})

	// ── 演示 3：计数器归零后可以复用 WaitGroup ──
	t.Run("WaitGroup可以复用", func(t *testing.T) {
		var wg sync.WaitGroup

		// 第一轮
		wg.Add(2)
		go func() { defer wg.Done(); t.Log("第一轮 goroutine A") }()
		go func() { defer wg.Done(); t.Log("第一轮 goroutine B") }()
		wg.Wait()

		// 第二轮：Wait 返回后，计数器已归零，可以再次 Add
		wg.Add(2)
		go func() { defer wg.Done(); t.Log("第二轮 goroutine A") }()
		go func() { defer wg.Done(); t.Log("第二轮 goroutine B") }()
		wg.Wait()

		t.Log("两轮复用 WaitGroup 均成功")
	})

	// ── 演示 4：计数器变负会 panic（用 recover 捕获，仅演示，生产禁止）──
	t.Run("计数器为负触发panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("捕获到预期的 panic: %v", r)
				t.Log("结论：Done/Add(-1) 调用次数不能超过 Add 的总量")
			}
		}()

		var wg sync.WaitGroup
		wg.Add(1)
		wg.Done() // 计数器归零
		wg.Done() // ← 再 Done 一次，计数器变为 -1 → panic
	})
}
