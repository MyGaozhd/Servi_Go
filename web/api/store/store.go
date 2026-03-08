// Package store 负责 Todo 数据的持久化读写。
// 该层只做"存取"，不包含任何业务校验逻辑。
package store

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/MyGaozhd/Servi_Go/web/api/model"
)

// JSONStore 将 Todo 列表持久化到本地 JSON 文件（并发安全）
type JSONStore struct {
	mu     sync.RWMutex
	todos  []model.Todo
	nextID int
	path   string // JSON 文件路径
}

// New 创建并初始化 JSONStore，若文件不存在则以空列表启动
func New(filePath string) (*JSONStore, error) {
	s := &JSONStore{path: filePath}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// ── 内部持久化 ───────────────────────────────────────────────────────────────

func (s *JSONStore) load() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		s.todos = []model.Todo{}
		s.nextID = 1
		return nil
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.todos); err != nil {
		return err
	}
	for _, t := range s.todos {
		if t.ID >= s.nextID {
			s.nextID = t.ID + 1
		}
	}
	return nil
}

func (s *JSONStore) save() error {
	data, err := json.MarshalIndent(s.todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// ── 公开 CRUD 方法 ────────────────────────────────────────────────────────────

// List 返回所有 Todo 的副本（done 为 nil 时不过滤）
func (s *JSONStore) List(done *bool) []model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if done == nil {
		out := make([]model.Todo, len(s.todos))
		copy(out, s.todos)
		return out
	}
	var out []model.Todo
	for _, t := range s.todos {
		if t.Done == *done {
			out = append(out, t)
		}
	}
	return out
}

// GetByID 按 ID 查找单条，未找到返回 model.ErrNotFound
func (s *JSONStore) GetByID(id int) (model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.todos {
		if t.ID == id {
			return t, nil
		}
	}
	return model.Todo{}, model.ErrNotFound
}

// Insert 追加一条新 Todo，返回写入后的完整对象
func (s *JSONStore) Insert(t model.Todo) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t.ID = s.nextID
	s.nextID++
	s.todos = append(s.todos, t)
	return t, s.save()
}

// Update 按 ID 更新字段，未找到返回 model.ErrNotFound
func (s *JSONStore) Update(t model.Todo) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.todos {
		if existing.ID == t.ID {
			s.todos[i] = t
			return t, s.save()
		}
	}
	return model.Todo{}, model.ErrNotFound
}

// Delete 按 ID 删除，未找到返回 model.ErrNotFound
func (s *JSONStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.todos {
		if t.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return s.save()
		}
	}
	return model.ErrNotFound
}
