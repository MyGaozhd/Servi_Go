// Package handler 负责 HTTP 路由注册与请求处理。
// 该层只做协议转换：解析请求 → 调用 service → 序列化响应。
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/MyGaozhd/Servi_Go/web/api/model"
	"github.com/MyGaozhd/Servi_Go/web/api/service"
)

// Handler 持有业务服务层
type Handler struct {
	svc *service.TodoService
}

// New 创建 Handler
func New(svc *service.TodoService) *Handler {
	return &Handler{svc: svc}
}

// Routes 注册所有路由，返回套了 CORS 中间件的 http.Handler
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/todos", h.handleTodos)   // GET / POST
	mux.HandleFunc("/api/todos/", h.handleTodo)   // GET / PATCH / DELETE /{id}
	mux.HandleFunc("/api/health", h.handleHealth)

	return CORS(mux)
}

// ── /api/health ──────────────────────────────────────────────────────────────

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ── /api/todos ───────────────────────────────────────────────────────────────

func (h *Handler) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTodos(w, r)
	case http.MethodPost:
		h.createTodo(w, r)
	default:
		WriteErr(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

func (h *Handler) listTodos(w http.ResponseWriter, r *http.Request) {
	var done *bool
	if q := r.URL.Query().Get("done"); q != "" {
		v := q == "true"
		done = &v
	}
	WriteJSON(w, http.StatusOK, h.svc.List(done))
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req model.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErr(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
		return
	}
	todo, err := h.svc.Create(req)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, todo)
}

// ── /api/todos/{id} ──────────────────────────────────────────────────────────

func (h *Handler) handleTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getTodo(w, r, id)
	case http.MethodPatch:
		h.updateTodo(w, r, id)
	case http.MethodDelete:
		h.deleteTodo(w, r, id)
	default:
		WriteErr(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

func (h *Handler) getTodo(w http.ResponseWriter, _ *http.Request, id int) {
	todo, err := h.svc.GetByID(id)
	if err != nil {
		WriteErr(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, todo)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request, id int) {
	var req model.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErr(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
		return
	}
	todo, err := h.svc.Update(id, req)
	if err != nil {
		WriteErr(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, todo)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, _ *http.Request, id int) {
	if err := h.svc.Delete(id); err != nil {
		WriteErr(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── 工具 ─────────────────────────────────────────────────────────────────────

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		WriteErr(w, http.StatusBadRequest, "无效的 id")
		return 0, false
	}
	return id, true
}
