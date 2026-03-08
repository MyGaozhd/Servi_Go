package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handler 持有 Store，实现所有 HTTP 路由处理
type Handler struct {
	store *Store
}

// Routes 注册路由，返回 http.Handler（含 CORS 中间件）
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	// GET    /api/todos        查询全部（?done=true/false 过滤）
	// POST   /api/todos        新增
	mux.HandleFunc("/api/todos", h.handleTodos)

	// GET    /api/todos/{id}   查询单条
	// PATCH  /api/todos/{id}   更新
	// DELETE /api/todos/{id}   删除
	mux.HandleFunc("/api/todos/", h.handleTodo)

	// 健康检查
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return corsMiddleware(mux)
}

// ---------- /api/todos ----------

func (h *Handler) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var done *bool
		if q := r.URL.Query().Get("done"); q != "" {
			v := q == "true"
			done = &v
		}
		writeJSON(w, http.StatusOK, h.store.List(done))

	case http.MethodPost:
		var req CreateTodoReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
			return
		}
		todo, err := h.store.Create(req)
		if err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, todo)

	default:
		writeErr(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// ---------- /api/todos/{id} ----------

func (h *Handler) handleTodo(w http.ResponseWriter, r *http.Request) {
	// 解析 URL 中的 id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeErr(w, http.StatusBadRequest, "无效的 id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		todos := h.store.List(nil)
		for _, t := range todos {
			if t.ID == id {
				writeJSON(w, http.StatusOK, t)
				return
			}
		}
		writeErr(w, http.StatusNotFound, "todo not found")

	case http.MethodPatch:
		var req UpdateTodoReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
			return
		}
		todo, err := h.store.Update(id, req)
		if err != nil {
			writeErr(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, todo)

	case http.MethodDelete:
		if err := h.store.Delete(id); err != nil {
			writeErr(w, http.StatusNotFound, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		writeErr(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// ---------- 辅助函数 ----------

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

// corsMiddleware 允许前端跨域访问（开发阶段放开全部来源）
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
