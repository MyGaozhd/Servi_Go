package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

const dataFile = "../data/todos.json"

// Store 负责从 JSON 文件读写 Todo 列表（并发安全）
type Store struct {
	mu    sync.RWMutex
	todos []Todo
	nextID int
}

// NewStore 加载已有数据，若文件不存在则初始化空列表
func NewStore() (*Store, error) {
	s := &Store{}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// load 从磁盘读取 JSON
func (s *Store) load() error {
	data, err := os.ReadFile(dataFile)
	if errors.Is(err, os.ErrNotExist) {
		s.todos = []Todo{}
		s.nextID = 1
		return nil
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.todos); err != nil {
		return err
	}
	// 计算下一个可用 ID
	for _, t := range s.todos {
		if t.ID >= s.nextID {
			s.nextID = t.ID + 1
		}
	}
	return nil
}

// save 将内存数据持久化到磁盘
func (s *Store) save() error {
	data, err := json.MarshalIndent(s.todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

// List 返回全部事项（可按 done 过滤：nil = 不过滤）
func (s *Store) List(done *bool) []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if done == nil {
		result := make([]Todo, len(s.todos))
		copy(result, s.todos)
		return result
	}
	var result []Todo
	for _, t := range s.todos {
		if t.Done == *done {
			result = append(result, t)
		}
	}
	return result
}

// Create 新增一条事项
func (s *Store) Create(req CreateTodoReq) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Title == "" {
		return Todo{}, errors.New("title 不能为空")
	}
	if req.Priority == "" {
		req.Priority = PriorityMedium
	}

	now := time.Now()
	t := Todo{
		ID:        s.nextID,
		Title:     req.Title,
		Note:      req.Note,
		Priority:  req.Priority,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.nextID++
	s.todos = append(s.todos, t)
	return t, s.save()
}

// Update 更新指定 ID 的事项（仅修改非 nil 字段）
func (s *Store) Update(id int, req UpdateTodoReq) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.todos {
		if t.ID != id {
			continue
		}
		if req.Title != nil {
			s.todos[i].Title = *req.Title
		}
		if req.Note != nil {
			s.todos[i].Note = *req.Note
		}
		if req.Priority != nil {
			s.todos[i].Priority = *req.Priority
		}
		if req.Done != nil {
			s.todos[i].Done = *req.Done
		}
		s.todos[i].UpdatedAt = time.Now()
		return s.todos[i], s.save()
	}
	return Todo{}, errors.New("todo not found")
}

// Delete 删除指定 ID 的事项
func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.todos {
		if t.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return s.save()
		}
	}
	return errors.New("todo not found")
}
