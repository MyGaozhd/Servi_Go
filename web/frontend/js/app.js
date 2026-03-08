// app.js —— 应用入口，负责初始化和事件绑定
// 通过事件委托减少监听器数量，所有业务操作委托给 api / state / render

import { api }    from './api.js'
import { state }  from './state.js'
import { renderAll, renderList } from './render.js'

// ── 初始化 ────────────────────────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', async () => {
  await loadTodos()
  bindAddForm()
  bindFilterBar()
  bindListDelegate()
  bindModal()
})

// ── 数据加载 ──────────────────────────────────────────────────────────────────
async function loadTodos() {
  try {
    state.setTodos(await api.list() ?? [])
    refresh()
  } catch (e) {
    toast('❌ 加载失败：' + e.message)
  }
}

// ── 新增事项 ──────────────────────────────────────────────────────────────────
function bindAddForm() {
  document.getElementById('add-form').addEventListener('submit', async e => {
    e.preventDefault()
    const title    = document.getElementById('inp-title').value.trim()
    const note     = document.getElementById('inp-note').value.trim()
    const priority = document.getElementById('inp-priority').value
    if (!title) return toast('请输入事项标题')
    try {
      const todo = await api.create({ title, note, priority })
      state.prepend(todo)
      refresh()
      e.target.reset()
      toast('✅ 已添加')
    } catch (e) {
      toast('❌ ' + e.message)
    }
  })
}

// ── 筛选 Tab ──────────────────────────────────────────────────────────────────
function bindFilterBar() {
  document.querySelector('.filter-bar').addEventListener('click', e => {
    const btn = e.target.closest('.filter-btn')
    if (!btn) return
    document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'))
    btn.classList.add('active')
    state.setFilter(btn.dataset.filter)
    renderList(state.visible())
  })
}

// ── 列表事件委托（统一处理 checkbox / 编辑 / 删除）───────────────────────────
function bindListDelegate() {
  document.getElementById('todo-list').addEventListener('change', async e => {
    const el = e.target.closest('[data-check]')
    if (!el) return
    await toggleDone(+el.dataset.check, el.checked)
  })

  document.getElementById('todo-list').addEventListener('click', async e => {
    const editBtn = e.target.closest('[data-edit]')
    const delBtn  = e.target.closest('[data-del]')
    if (editBtn) openEdit(+editBtn.dataset.edit)
    if (delBtn)  await deleteTodo(+delBtn.dataset.del)
  })
}

// ── 操作：切换完成 ────────────────────────────────────────────────────────────
async function toggleDone(id, done) {
  try {
    const updated = await api.update(id, { done })
    state.replace(updated)
    refresh()
  } catch (e) {
    toast('❌ ' + e.message)
  }
}

// ── 操作：删除 ────────────────────────────────────────────────────────────────
async function deleteTodo(id) {
  if (!confirm('确认删除？')) return
  try {
    await api.remove(id)
    state.remove(id)
    refresh()
    toast('🗑️ 已删除')
  } catch (e) {
    toast('❌ ' + e.message)
  }
}

// ── 编辑 Modal ────────────────────────────────────────────────────────────────
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
      const updated = await api.update(state.editingId, { title, note, priority })
      state.replace(updated)
      refresh()
      closeModal()
      toast('✅ 已保存')
    } catch (e) {
      toast('❌ ' + e.message)
    }
  })
}

function openEdit(id) {
  const todo = state.todos.find(t => t.id === id)
  if (!todo) return
  state.startEdit(id)
  document.getElementById('modal-title').value    = todo.title
  document.getElementById('modal-note').value     = todo.note ?? ''
  document.getElementById('modal-priority').value = todo.priority
  document.getElementById('modal-backdrop').classList.add('open')
}

function closeModal() {
  state.endEdit()
  document.getElementById('modal-backdrop').classList.remove('open')
}

// ── 渲染快捷方法 ──────────────────────────────────────────────────────────────
function refresh() {
  renderAll(state.todos)
  renderList(state.visible())
}

// ── Toast ─────────────────────────────────────────────────────────────────────
let toastTimer
function toast(msg) {
  const el = document.getElementById('toast')
  el.textContent = msg
  el.classList.add('show')
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => el.classList.remove('show'), 2500)
}
