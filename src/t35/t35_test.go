// Package t35 演示对象池模式：自定义 ObjectPool 与标准库 sync.Pool。
//
// # 两种实现对比
//
// 自定义 ObjectPool（基于有缓冲 channel）：
//   - 对象数量固定，GetObj 支持超时
//   - 适合：数据库连接池等需要严格控制数量的场景
//
// sync.Pool（标准库）：
//   - 对象数量不固定，GC 时可能被清空
//   - 适合：临时对象（bytes.Buffer 等），减少 GC 压力
//   - 不适合：需要状态持久化的对象（如数据库连接）
package t35

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// ─────────────────────────────────────────────
// 自定义 ObjectPool（基于有缓冲 channel）
// ─────────────────────────────────────────────

// ResumeAbleObj 可复用对象
type ResumeAbleObj struct{}

// ObjectPool 基于有缓冲 channel 实现的固定容量对象池
type ObjectPool struct {
	bufChan chan *ResumeAbleObj
}

// NewResumeAbleObject 创建并初始化一个容量为 num 的对象池
func NewResumeAbleObject(num int) *ObjectPool {
	objectPool := ObjectPool{}
	objectPool.bufChan = make(chan *ResumeAbleObj, num)
	for i := 0; i < num; i++ {
		objectPool.bufChan <- &ResumeAbleObj{}
	}
	return &objectPool
}

// GetObj 从池中获取对象，超时则返回 error
func (p *ObjectPool) GetObj(timeout time.Duration) (*ResumeAbleObj, error) {
	select {
	case ret := <-p.bufChan:
		return ret, nil
	case <-time.After(timeout):
		return nil, errors.New("time out")
	}
}

// ReleaseObj 将对象归还到池中，池满时返回 error
func (p *ObjectPool) ReleaseObj(obj *ResumeAbleObj) error {
	select {
	case p.bufChan <- obj:
		return nil
	default:
		return errors.New("over flow")
	}
}

// Test35_0 演示自定义 ObjectPool 的获取与归还。
func Test35_0(t *testing.T) {
	pool := NewResumeAbleObject(10)

	for i := 0; i < 10; i++ {
		if obj, err := pool.GetObj(time.Millisecond * 10); err != nil {
			t.Error(err)
		} else {
			t.Log(obj)
			if err := pool.ReleaseObj(obj); err != nil {
				t.Error(err)
			}
		}
	}
	t.Log("✓ 所有对象已获取并归还")
}

// ─────────────────────────────────────────────
// sync.Pool 标准库对象池
// ─────────────────────────────────────────────

// MyObject 模拟一个可复用的临时对象
type MyObject struct {
	data []byte
}

// reset 归还前重置对象状态，避免脏数据污染下一个使用者
func (o *MyObject) reset() {
	o.data = o.data[:0]
}

var syncPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("sync.Pool: 创建新对象（池空时才调用 New）")
		return &MyObject{data: make([]byte, 0, 64)}
	},
}

// Test35_1 演示 sync.Pool 标准库对象池的使用。
func Test35_1(t *testing.T) {
	// 第一次 Get：池为空，调用 New 创建新对象
	obj1 := syncPool.Get().(*MyObject)
	t.Logf("第一次 Get，对象地址: %p", obj1)

	obj1.data = append(obj1.data, []byte("hello")...)
	t.Logf("使用对象，data = %s", obj1.data)

	obj1.reset()
	syncPool.Put(obj1)
	t.Log("对象已归还到 sync.Pool")

	// 第二次 Get：池中有对象，直接复用
	obj2 := syncPool.Get().(*MyObject)
	t.Logf("第二次 Get，对象地址: %p", obj2)
	syncPool.Put(obj2)

	// 并发场景
	var wg sync.WaitGroup
	const goroutines = 5
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			o := syncPool.Get().(*MyObject)
			o.data = append(o.data, fmt.Sprintf("goroutine-%d", id)...)
			t.Logf("goroutine %d 使用对象，data=%s", id, o.data)
			o.reset()
			syncPool.Put(o)
		}(i)
	}
	wg.Wait()
	t.Log("✓ 并发使用 sync.Pool 完成")
}
