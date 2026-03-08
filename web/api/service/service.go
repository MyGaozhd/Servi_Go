// Package service 封装业务逻辑：输入校验、字段默认值、编排 store 操作。
// Handler 层只调用 Service，不直接操作 Store。
package service

import (
	"strings"
	"time"

	"github.com/MyGaozhd/Servi_Go/web/api/model"
	"github.com/MyGaozhd/Servi_Go/web/api/store"
)

// TodoService 提供所有 Todo 业务操作
type TodoService struct {
	store *store.JSONStore
}

// New 创建 TodoService
func New(s *store.JSONStore) *TodoService {
	return &TodoService{store: s}
}

// List 查询事项列表（done 为 nil 时返回全部）
func (svc *TodoService) List(done *bool) []model.Todo {
	return svc.store.List(done)
}

// GetByID 查询单条事项
func (svc *TodoService) GetByID(id int) (model.Todo, error) {
	return svc.store.GetByID(id)
}

// Create 校验请求并创建新事项
func (svc *TodoService) Create(req model.CreateReq) (model.Todo, error) {
	// 业务校验
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return model.Todo{}, model.ErrTitleEmpty
	}
	if req.Priority == "" {
		req.Priority = model.PriorityMedium // 默认中优先级
	}
	if !req.Priority.Valid() {
		return model.Todo{}, model.ErrInvalidPriority
	}

	now := time.Now()
	t := model.Todo{
		Title:     req.Title,
		Note:      strings.TrimSpace(req.Note),
		Priority:  req.Priority,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return svc.store.Insert(t)
}

// Update 按 ID 更新事项（仅修改请求中的非 nil 字段）
func (svc *TodoService) Update(id int, req model.UpdateReq) (model.Todo, error) {
	t, err := svc.store.GetByID(id)
	if err != nil {
		return model.Todo{}, err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return model.Todo{}, model.ErrTitleEmpty
		}
		t.Title = title
	}
	if req.Note != nil {
		t.Note = strings.TrimSpace(*req.Note)
	}
	if req.Priority != nil {
		if !req.Priority.Valid() {
			return model.Todo{}, model.ErrInvalidPriority
		}
		t.Priority = *req.Priority
	}
	if req.Done != nil {
		t.Done = *req.Done
	}
	t.UpdatedAt = time.Now()

	return svc.store.Update(t)
}

// Delete 删除指定 ID 的事项
func (svc *TodoService) Delete(id int) error {
	return svc.store.Delete(id)
}
