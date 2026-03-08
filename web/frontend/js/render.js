// render.js —— 负责所有 DOM 渲染，不包含任何业务逻辑
// 通过接收纯数据驱动视图更新

// ── 配置映射（集中管理，方便日后扩展）──────────────────────────────────────
const PRIORITY_LABEL = { low: '低', medium: '中', high: '高' }

// ── 公开渲染入口 ─────────────────────────────────────────────────────────────

export function renderAll(todos, filter) {
  renderStats(todos)
  renderList(todos, filter)
}

export function renderStats(todos) {
  const total  = todos.length
  const done   = todos.filter(t => t.done).length
  document.getElementById('stat-total').textContent  = total
  document.getElementById('stat-active').textContent = total - done
  document.getElementById('stat-done').textContent   = done
}

export function renderList(visible) {
  const container = document.getElementById('todo-list')
  if (visible.length === 0) {
    container.innerHTML = `<div class="empty">暂无事项 🎉</div>`
    return
  }
  container.innerHTML = visible.map(todoCard).join('')
}

// ── 卡片模板 ─────────────────────────────────────────────────────────────────

function todoCard(t) {
  const label = PRIORITY_LABEL[t.priority] ?? t.priority
  const date  = formatDate(t.created_at)
  return `
  <div class="todo-card priority-${t.priority} ${t.done ? 'done' : ''}">
    <input type="checkbox" data-check="${t.id}" ${t.done ? 'checked' : ''}>
    <div class="todo-body">
      <div class="todo-title">${esc(t.title)}</div>
      ${t.note ? `<div class="todo-note">${esc(t.note)}</div>` : ''}
      <div class="todo-meta">
        <span class="badge ${t.priority}">${label}优先级</span>
        <span class="todo-time">${date}</span>
      </div>
    </div>
    <div class="todo-actions">
      <button class="btn-icon edit" data-edit="${t.id}" title="编辑">✏️</button>
      <button class="btn-icon"      data-del="${t.id}"  title="删除">🗑️</button>
    </div>
  </div>`
}

// ── 工具函数 ─────────────────────────────────────────────────────────────────

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' })
}

function esc(str) {
  return str.replace(/[&<>"']/g, c =>
    ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]))
}
