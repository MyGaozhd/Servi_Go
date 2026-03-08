const API = 'http://localhost:8080/api'

// ── 状态 ──────────────────────────────────────────────
let todos  = []
let filter = 'all'   // 'all' | 'active' | 'done'
let editingId = null

// ── 初始化 ────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
  fetchTodos()
  bindForm()
  bindFilter()
  bindModal()
})

// ── API 工具 ──────────────────────────────────────────
async function request(path, options = {}) {
  const res = await fetch(API + path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (res.status === 204) return null
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || '请求失败')
  return data
}

// ── 拉取列表 ──────────────────────────────────────────
async function fetchTodos() {
  try {
    todos = await request('/todos') ?? []
    renderAll()
  } catch (e) {
    toast('❌ 加载失败：' + e.message)
  }
}

// ── 新增事项 ──────────────────────────────────────────
function bindForm() {
  document.getElementById('add-form').addEventListener('submit', async e => {
    e.preventDefault()
    const title    = document.getElementById('inp-title').value.trim()
    const note     = document.getElementById('inp-note').value.trim()
    const priority = document.getElementById('inp-priority').value
    if (!title) return toast('请输入事项标题')
    try {
      const todo = await request('/todos', {
        method: 'POST',
        body: JSON.stringify({ title, note, priority }),
      })
      todos.unshift(todo)
      renderAll()
      e.target.reset()
      toast('✅ 已添加')
    } catch (e) {
      toast('❌ ' + e.message)
    }
  })
}

// ── 筛选 Tab ──────────────────────────────────────────
function bindFilter() {
  document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      filter = btn.dataset.filter
      document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'))
      btn.classList.add('active')
      renderList()
    })
  })
}

// ── 切换完成状态 ──────────────────────────────────────
async function toggleDone(id, done) {
  try {
    const updated = await request(`/todos/${id}`, {
      method: 'PATCH',
      body: JSON.stringify({ done }),
    })
    replaceTodo(updated)
    renderAll()
  } catch (e) {
    toast('❌ ' + e.message)
  }
}

// ── 删除 ──────────────────────────────────────────────
async function deleteTodo(id) {
  if (!confirm('确认删除？')) return
  try {
    await request(`/todos/${id}`, { method: 'DELETE' })
    todos = todos.filter(t => t.id !== id)
    renderAll()
    toast('🗑️ 已删除')
  } catch (e) {
    toast('❌ ' + e.message)
  }
}

// ── 编辑 Modal ────────────────────────────────────────
function bindModal() {
  document.getElementById('modal-cancel').addEventListener('click', closeModal)
  document.getElementById('modal-backdrop').addEventListener('click', e => {
    if (e.target === e.currentTarget) closeModal()
  })
  document.getElementById('modal-form').addEventListener('submit', async e => {
    e.preventDefault()
    const title    = document.getElementById('modal-title').value.trim()
    const note     = document.getElementById('modal-note').value.trim()
    const priority = document.getElementById('modal-priority').value
    if (!title) return toast('标题不能为空')
    try {
      const updated = await request(`/todos/${editingId}`, {
        method: 'PATCH',
        body: JSON.stringify({ title, note, priority }),
      })
      replaceTodo(updated)
      renderAll()
      closeModal()
      toast('✅ 已保存')
    } catch (e) {
      toast('❌ ' + e.message)
    }
  })
}

function openEdit(id) {
  const todo = todos.find(t => t.id === id)
  if (!todo) return
  editingId = id
  document.getElementById('modal-title').value    = todo.title
  document.getElementById('modal-note').value     = todo.note ?? ''
  document.getElementById('modal-priority').value = todo.priority
  document.getElementById('modal-backdrop').classList.add('open')
}

function closeModal() {
  editingId = null
  document.getElementById('modal-backdrop').classList.remove('open')
}

// ── 渲染 ──────────────────────────────────────────────
function renderAll() {
  renderStats()
  renderList()
}

function renderStats() {
  const total  = todos.length
  const done   = todos.filter(t => t.done).length
  const active = total - done
  document.getElementById('stat-total').textContent  = total
  document.getElementById('stat-active').textContent = active
  document.getElementById('stat-done').textContent   = done
}

function renderList() {
  const container = document.getElementById('todo-list')
  const visible = todos.filter(t => {
    if (filter === 'active') return !t.done
    if (filter === 'done')   return t.done
    return true
  })
  // 高优先级排前
  const priority = { high: 0, medium: 1, low: 2 }
  visible.sort((a, b) => {
    if (a.done !== b.done) return a.done ? 1 : -1
    return (priority[a.priority] ?? 1) - (priority[b.priority] ?? 1)
  })

  if (visible.length === 0) {
    container.innerHTML = `<div class="empty">暂无事项 🎉</div>`
    return
  }
  container.innerHTML = visible.map(todoCard).join('')

  // 绑定卡片内的事件
  container.querySelectorAll('[data-check]').forEach(el =>
    el.addEventListener('change', () => toggleDone(+el.dataset.check, el.checked)))
  container.querySelectorAll('[data-edit]').forEach(el =>
    el.addEventListener('click', () => openEdit(+el.dataset.edit)))
  container.querySelectorAll('[data-del]').forEach(el =>
    el.addEventListener('click', () => deleteTodo(+el.dataset.del)))
}

function todoCard(t) {
  const priorityMap = { low: '低', medium: '中', high: '高' }
  const date = new Date(t.created_at).toLocaleDateString('zh-CN', {
    month: 'numeric', day: 'numeric'
  })
  return `
  <div class="todo-card ${t.done ? 'done' : ''}">
    <input type="checkbox" data-check="${t.id}" ${t.done ? 'checked' : ''}>
    <div class="todo-body">
      <div class="todo-title">${escHtml(t.title)}</div>
      ${t.note ? `<div class="todo-note">${escHtml(t.note)}</div>` : ''}
      <div class="todo-meta">
        <span class="badge ${t.priority}">${priorityMap[t.priority] ?? t.priority}优先级</span>
        <span class="todo-time">${date}</span>
      </div>
    </div>
    <div class="todo-actions">
      <button class="btn-icon edit" data-edit="${t.id}" title="编辑">✏️</button>
      <button class="btn-icon"      data-del="${t.id}"  title="删除">🗑️</button>
    </div>
  </div>`
}

// ── 工具 ──────────────────────────────────────────────
function replaceTodo(updated) {
  const idx = todos.findIndex(t => t.id === updated.id)
  if (idx !== -1) todos[idx] = updated
}

function escHtml(str) {
  return str.replace(/[&<>"']/g, c =>
    ({ '&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;' }[c]))
}

let toastTimer
function toast(msg) {
  const el = document.getElementById('toast')
  el.textContent = msg
  el.classList.add('show')
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => el.classList.remove('show'), 2500)
}
