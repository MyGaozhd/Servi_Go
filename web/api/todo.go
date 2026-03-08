package main

import "time"

// Priority 优先级
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Todo 单条日常事项
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`               // 事项标题
	Note      string    `json:"note"`                // 备注（可选）
	Priority  Priority  `json:"priority"`             // 优先级：low / medium / high
	Done      bool      `json:"done"`                 // 是否完成
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTodoReq 创建事项的请求体
type CreateTodoReq struct {
	Title    string   `json:"title"`
	Note     string   `json:"note"`
	Priority Priority `json:"priority"`
}

// UpdateTodoReq 更新事项的请求体
type UpdateTodoReq struct {
	Title    *string   `json:"title"`
	Note     *string   `json:"note"`
	Priority *Priority `json:"priority"`
	Done     *bool     `json:"done"`
}
