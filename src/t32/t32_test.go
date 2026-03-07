// Package t32 演示 sync.Once 实现"只执行一次"的线程安全保证。
package t32

import (
	"fmt"
	"sync"
	"testing"
	"unsafe"
)

// Singleton 演示单例模式的目标结构体
type Singleton struct{}

var singletonInstance *Singleton
var once sync.Once

// getSingletonInstance 线程安全的单例获取函数。
func getSingletonInstance() *Singleton {
	once.Do(func() {
		fmt.Println("create onject")
		singletonInstance = new(Singleton)
	})
	return singletonInstance
}

// Test32_0 演示单例模式：10 个并发 goroutine 调用 getSingletonInstance，
// 验证所有 goroutine 拿到的是同一个指针（地址相同）。
func Test32_0(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obj := getSingletonInstance()
			t.Log(unsafe.Pointer(obj))
		}()
	}
	wg.Wait()
	t.Log("✓ 所有 goroutine 拿到的是同一个单例对象")
}

// Test32_1 演示 once.Do 只执行一次的核心语义。
func Test32_1(t *testing.T) {
	t.Run("串行多次调用", func(t *testing.T) {
		var o sync.Once
		count := 0
		for i := 0; i < 5; i++ {
			o.Do(func() {
				count++
				t.Log("初始化函数执行了（只应执行一次）")
			})
		}
		if count != 1 {
			t.Errorf("期望 count=1，实际 count=%d", count)
		} else {
			t.Logf("count = %d ✓", count)
		}
	})

	t.Run("并发调用只执行一次", func(t *testing.T) {
		var o sync.Once
		var mu sync.Mutex
		count := 0
		var wg sync.WaitGroup

		const goroutines = 100
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				o.Do(func() {
					mu.Lock()
					count++
					mu.Unlock()
					t.Log("并发场景：初始化函数执行了（只应执行一次）")
				})
			}()
		}
		wg.Wait()

		if count != 1 {
			t.Errorf("期望 count=1，实际 count=%d", count)
		} else {
			t.Logf("count = %d ✓（%d 个并发 goroutine，只执行了一次）", count, goroutines)
		}
	})

	t.Run("不同函数传入只有第一次执行", func(t *testing.T) {
		var o sync.Once
		result := ""
		o.Do(func() { result = "第一次" })
		o.Do(func() { result = "第二次" })
		o.Do(func() { result = "第三次" })

		t.Logf("result = %q（只有第一次的函数被执行）", result)
		if result != "第一次" {
			t.Errorf("期望 '第一次'，实际 '%s'", result)
		}
	})
}
