// Package t31 演示 context 包的三种核心用法。
//
// # context 是什么
//
// context.Context 是 Go 1.7 引入的标准接口，用于：
//  1. 取消信号传播（WithCancel）
//  2. 超时/截止时间控制（WithTimeout / WithDeadline）
//  3. 请求级数据传递（WithValue）
//
// # context 树结构
//
// context 形成树形结构：父 context 取消时，所有子 context 自动级联取消。
//
//	Background ─── WithCancel ─── WithTimeout ─── WithValue
//	                  │
//	               cancel()  ← 取消此节点，其所有子孙也被取消
//
// # 使用原则
//
//  1. context 作为函数第一个参数传递（惯例命名 ctx）
//  2. 不要将 context 存入结构体字段
//  3. 不传 nil context，不确定时用 context.TODO()
//  4. WithValue 只传请求级元数据（如 traceID），不传业务参数
package t31

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Test31_0 演示 context.WithCancel 手动取消。
func Test31_0(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// 启动 5 个 goroutine，各自监听 ctx.Done()
	for i := 0; i < 5; i++ {
		go func(i int, ctx context.Context) {
			for {
				if isCancelled(ctx) {
					break
				}
				time.Sleep(time.Millisecond * 5)
			}
			fmt.Println(i, "Cancelled")
		}(i, ctx)
	}

	// 手动触发取消：ctx.Done() channel 关闭，所有 goroutine 收到取消信号
	cancel()

	// 等待所有 goroutine 退出
	time.Sleep(time.Second * 1)
}

// isCancelled 非阻塞检测 context 是否已被取消
func isCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// Test31_1 演示 context.WithTimeout 超时自动取消。
func Test31_1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-time.After(200 * time.Millisecond):
			fmt.Println("工作完成（正常结束）")
		case <-ctx.Done():
			fmt.Println("收到取消信号，提前退出:", ctx.Err())
		}
	}()

	<-done

	if ctx.Err() != nil {
		t.Logf("context 已取消，原因: %v", ctx.Err())
		if ctx.Err() == context.DeadlineExceeded {
			t.Log("✓ 确认是超时导致的取消（DeadlineExceeded）")
		}
	}
}

// Test31_2 演示 context.WithValue 在调用链中传递请求级数据。
func Test31_2(t *testing.T) {
	type contextKey string
	const (
		traceIDKey contextKey = "traceID"
		userIDKey  contextKey = "userID"
	)

	ctx := context.Background()
	ctx = context.WithValue(ctx, traceIDKey, "trace-abc-123")
	ctx = context.WithValue(ctx, userIDKey, 42)

	handler := func(ctx context.Context) {
		traceID, ok := ctx.Value(traceIDKey).(string)
		if !ok {
			t.Error("traceID 类型断言失败")
			return
		}
		userID, ok := ctx.Value(userIDKey).(int)
		if !ok {
			t.Error("userID 类型断言失败")
			return
		}
		t.Logf("处理请求 | traceID=%s | userID=%d", traceID, userID)

		notExist := ctx.Value(contextKey("notExist"))
		t.Logf("不存在的 key 返回: %v（nil）", notExist)
	}

	handler(ctx)

	ctxBase := context.Background()
	ctxWithVal := context.WithValue(ctxBase, traceIDKey, "trace-xyz")
	t.Logf("原 ctx 的 traceID: %v（nil，不受影响）", ctxBase.Value(traceIDKey))
	t.Logf("新 ctx 的 traceID: %v", ctxWithVal.Value(traceIDKey))
}
