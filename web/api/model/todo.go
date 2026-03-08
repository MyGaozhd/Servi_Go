// Package model 定义业务实体、值对象和请求/响应 DTO。
// 该包不依赖任何其他内部包，是整个分层架构的最底层。
package model

import (
	"errors"
	"time"
)

// ── Sentinel errors ──────────────────────────────────────────────────────────
// 使用具名 error 变量，调用方可通过 errors.Is 精确判断错误类型。

var (
	ErrNotFound     = errors.New("todo not found")
	ErrTitleEmpty   = errors.New("title 不能为空")
	ErrInvalidPriority = errors.New("priority 必须是 low / medium / high")
)

// ── Priority ─────────────────────────────────────────────────────────────────

// Priority 表示事项的优先级
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Valid 校验优先级值是否合法
func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}

// ── Todo ─────────────────────────────────────────────────────────────────────

// Todo 是核心业务实体
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Note      string    `json:"note"`
	Priority  Priority  `json:"priority"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ── DTOs（数据传输对象）──────────────────────────────────────────────────────

// CreateReq 创建事项的请求体
type CreateReq struct {
	Title    string   `json:"title"`
	Note     string   `json:"note"`
	Priority Priority `json:"priority"`
}

// UpdateReq 更新事项的请求体（全部字段为指针，nil 表示不修改）
type UpdateReq struct {
	Title    *string   `json:"title"`
	Note     *string   `json:"note"`
	Priority *Priority `json:"priority"`
	Done     *bool     `json:"done"`
}
