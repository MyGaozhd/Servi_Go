// Package t24 演示 sync.Mutex（互斥锁）和 sync.RWMutex（读写锁）的使用。
//
// # 为什么需要锁
//
// 多个 goroutine 并发读写同一变量时，会出现"数据竞争"（data race）：
//   - count++ 在 CPU 层面是"读-改-写"三步操作，不是原子的
//   - 两个 goroutine 可能同时读到相同的旧值，各自加 1 后写回，导致一次自增丢失
//   - 结果：最终 count < 预期值
//
// # 检测竞态
//
// 运行以下命令可启用竞态检测器（race detector）：
//
//	go test -race ./src/t24/
//
// race detector 会在运行时动态检测并发访问冲突，输出详细报告。
package t24

import (
	"sync"
	"testing"
	"time"
)

// Test24_0 演示 sync.Mutex 互斥锁保护共享变量，确保并发安全。
//
// Mutex 使用模式（惯用写法）：
//
//	mu.Lock()
//	defer mu.Unlock()  // 用 defer 确保即使 panic 也能解锁
//	// 临界区代码
//
// 结果：5000 个 goroutine 各自 count++，最终 count == 5000。
func Test24_0(t *testing.T) {
	var mu sync.Mutex
	count := 0

	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 加锁：进入临界区，其他 goroutine 在此阻塞等待
			mu.Lock()
			// defer 解锁：函数退出时自动释放，即使发生 panic 也能正确解锁
			defer mu.Unlock()
			count++
		}()
	}

	// 用 WaitGroup 等待所有协程完成，比 time.Sleep 更精确
	wg.Wait()

	// 有了互斥锁，结果必然是 5000
	if count != 5000 {
		t.Errorf("期望 count = 5000，实际 = %d（存在数据竞争！）", count)
	} else {
		t.Logf("count = %d ✓（互斥锁保证了结果正确）", count)
	}
}

// Test24_1 演示不加锁时的数据竞争现象。
//
// 数据竞争说明：
//   - count++ 非原子：底层是 LOAD → ADD → STORE 三条指令
//   - 两个 goroutine 可能同时 LOAD 到相同的旧值，各自 +1 后 STORE
//   - 导致某些自增操作"丢失"，最终 count 通常 < 5000
//
// 检测方式：go test -race ./src/t24/ 会报告竞态并打印调用栈。
//
// ⚠️ 注意：此测试故意制造数据竞争，仅用于演示，生产代码中禁止这样写。
func Test24_1(t *testing.T) {
	count := 0

	// 用 time.Sleep 代替 WaitGroup，避免引入额外同步（保留竞态现象）
	for i := 0; i < 5000; i++ {
		go func() {
			count++ // ← 无保护的并发写，存在数据竞争
		}()
	}

	// 等待大部分 goroutine 完成（不精确，仅用于演示）
	time.Sleep(time.Second)

	// 由于数据竞争，count 通常 < 5000
	t.Logf("不加锁的 count = %d（期望 5000，数据竞争导致部分更新丢失）", count)
	if count == 5000 {
		// 极小概率下也可能等于 5000（侥幸未发生竞争），但不可依赖
		t.Log("本次运行碰巧 = 5000，但这是不可靠的，加 -race 可检测到竞争")
	}
}

// Test24_2 演示 sync.RWMutex（读写锁）在读多写少场景下的优势。
//
// RWMutex 规则：
//   - RLock() / RUnlock()：读锁，允许多个 goroutine 同时持有（并发读）
//   - Lock() / Unlock()：写锁，独占，写时不允许任何读或写
//
// 适用场景：
//   - 读操作远多于写操作（如内存缓存、配置读取）
//   - 多个并发读不会互斥，吞吐量远高于 Mutex
//
// 对比 Mutex：
//   - Mutex：读写都互斥，同一时刻只有一个 goroutine 进入临界区
//   - RWMutex：多个读可并发，写时独占；读多写少时性能更好
func Test24_2(t *testing.T) {
	var rw sync.RWMutex
	cache := make(map[string]string) // 模拟内存缓存
	var wg sync.WaitGroup

	// ── 写操作（独占锁）──
	wg.Add(1)
	go func() {
		defer wg.Done()
		rw.Lock()         // 写锁：独占，阻塞所有读和写
		defer rw.Unlock() // 写锁释放
		cache["key"] = "value"
		t.Log("写入缓存: key=value")
	}()

	wg.Wait() // 等写完再读

	// ── 读操作（共享锁）：多个 goroutine 同时持有读锁 ──
	const readCount = 5
	var readWg sync.WaitGroup
	readWg.Add(readCount)

	for i := 0; i < readCount; i++ {
		go func(id int) {
			defer readWg.Done()
			rw.RLock()         // 读锁：允许多个 goroutine 同时持有
			defer rw.RUnlock() // 读锁释放
			val := cache["key"]
			t.Logf("goroutine %d 读取缓存: key=%s", id, val)
		}(i)
	}

	readWg.Wait()
	t.Log("RWMutex 演示完成：多读并发，写时独占")
}
