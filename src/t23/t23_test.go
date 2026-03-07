// Package t23 演示 Go 语言 goroutine（协程）的核心用法。
//
// # Goroutine 基础
//
// goroutine 是 Go 运行时管理的轻量级线程：
//   - 用 go 关键字启动：go func() { ... }()
//   - 创建成本极低（初始栈仅 ~2KB，可动态扩缩）
//   - 由 Go 调度器（M:N 模型）映射到操作系统线程上执行
//   - 主 goroutine（main/测试函数）退出时，所有子 goroutine 强制终止
//
// # 常见陷阱
//
//   - 循环变量捕获：闭包捕获的是变量引用，需通过参数传值或用 := 重声明
//   - goroutine 泄露：没有退出路径的 goroutine 会一直占用内存
//   - 依赖 time.Sleep 等待不可靠，应使用 sync.WaitGroup 或 channel 同步
package t23

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// Test23_0 演示最基础的 goroutine 启动：开 10 个协程并发打印编号。
//
// 注意事项：
//   - 循环变量 i 通过函数参数传入，避免"闭包捕获同一变量"的经典陷阱
//   - 使用 time.Sleep 等待协程执行，属于"不可靠同步"，仅用于演示
//   - 输出顺序不确定——goroutine 调度是非确定性的
func Test23_0(t *testing.T) {
	// 启动 10 个 goroutine，每个打印自己的编号
	for i := 0; i < 10; i++ {
		// 将 i 作为参数传入，而不是在闭包中直接引用循环变量
		// 若写成 go func() { t.Log(i) }() 则所有协程打印的 i 可能都是 10
		go func(i int) {
			t.Log("goroutine 编号:", i)
		}(i)
	}

	// time.Sleep 是不可靠的等待方式，仅用于演示
	// 生产代码应使用 sync.WaitGroup（见 Test23_1）
	time.Sleep(time.Second * 1)
}

// Test23_1 演示使用 sync.WaitGroup 等待所有 goroutine 完成。
//
// WaitGroup 三步法：
//  1. wg.Add(n)  — 在启动 goroutine 之前，计数器 +n
//  2. wg.Done()  — 在 goroutine 内部完成时调用，计数器 -1
//  3. wg.Wait()  — 阻塞当前 goroutine，直到计数器归零
//
// 对比 time.Sleep：WaitGroup 是精确同步，不会提前退出也不会多等。
func Test23_1(t *testing.T) {
	var wg sync.WaitGroup
	results := make([]int, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1) // ⚠️ 必须在 go 语句之前调用 Add，否则存在竞态
		go func(i int) {
			defer wg.Done() // 函数退出时自动 -1，即使 panic 也会执行
			results[i] = i * i
			t.Logf("goroutine %d 完成，结果 = %d", i, results[i])
		}(i)
	}

	// 阻塞，直到 10 个 goroutine 全部 Done
	wg.Wait()
	t.Log("所有 goroutine 已完成，results =", results)
}

// Test23_2 演示 goroutine 泄露：无缓冲 channel 没有接收方导致 goroutine 永久阻塞。
//
// 泄露原因：
//   - ch := make(chan int) 无缓冲
//   - goroutine 执行 ch <- 1 时，发送会一直阻塞，直到有人接收
//   - 若没有接收方，goroutine 永远不会退出，其占用的内存和 goroutine 栈无法释放
//
// 生产环境防止泄露的方法：
//  1. 使用带缓冲的 channel，确保发送不阻塞
//  2. 通过 context 或 done channel 给 goroutine 提供退出信号
//  3. 使用 goleak（第三方库）在测试中检测泄露
func Test23_2(t *testing.T) {
	before := runtime.NumGoroutine()
	t.Logf("启动前 goroutine 数量: %d", before)

	// 无缓冲 channel，没有接收方
	ch := make(chan int)

	// 启动 3 个 goroutine，每个都会在发送时永久阻塞 → 泄露
	for i := 0; i < 3; i++ {
		go func(i int) {
			// 发送会一直阻塞，因为没有接收方
			ch <- i // ← 永远卡在这里
		}(i)
	}

	// 稍等让 goroutine 进入阻塞状态
	time.Sleep(10 * time.Millisecond)

	after := runtime.NumGoroutine()
	t.Logf("泄露后 goroutine 数量: %d（增加了 %d 个泄露的协程）", after, after-before)

	// 验证确实发生了泄露（goroutine 数量增加）
	if after <= before {
		t.Error("预期 goroutine 数量增加，但未观察到泄露")
	}

	// 【修复方案演示】关闭 channel 或使用带缓冲的 channel 可避免上述泄露
	// 此处通过消费 channel 的方式解除阻塞，让泄露的 goroutine 得以退出
	for i := 0; i < 3; i++ {
		<-ch
	}
	time.Sleep(10 * time.Millisecond)
	t.Logf("消费后 goroutine 数量: %d（泄露已修复）", runtime.NumGoroutine())
}

// Test23_3 演示 runtime.GOMAXPROCS 控制并发使用的 CPU 核数。
//
// GOMAXPROCS 说明：
//   - Go 1.5+ 默认值 = 逻辑 CPU 数（runtime.NumCPU()）
//   - GOMAXPROCS(1)：单线程，goroutine 协作调度（无真正并行）
//   - GOMAXPROCS(n)：最多 n 个 OS 线程并行执行 goroutine
//   - 返回值是修改前的旧值
//
// 注意：GOMAXPROCS 影响的是"并行"（parallelism），
// goroutine 本身是"并发"（concurrency）的，与 GOMAXPROCS 无关。
func Test23_3(t *testing.T) {
	// 查询逻辑 CPU 数量（超线程也算）
	numCPU := runtime.NumCPU()
	t.Logf("逻辑 CPU 核数: %d", numCPU)

	// 查询当前 GOMAXPROCS（Go 1.5+ 默认等于 NumCPU）
	current := runtime.GOMAXPROCS(0) // 传 0 表示仅查询，不修改
	t.Logf("当前 GOMAXPROCS: %d", current)

	// 临时改为 2 核，演示修改效果
	old := runtime.GOMAXPROCS(2)
	t.Logf("设置 GOMAXPROCS=2，旧值为 %d", old)
	fmt.Printf("GOMAXPROCS 当前值: %d\n", runtime.GOMAXPROCS(0))

	// 恢复原值，避免影响其他测试
	runtime.GOMAXPROCS(old)
	t.Logf("已恢复 GOMAXPROCS = %d", runtime.GOMAXPROCS(0))
}
