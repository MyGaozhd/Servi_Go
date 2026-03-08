// state.js —— 封装应用状态，提供读写方法
// 将 todos / filter / editingId 收拢到单一对象，避免全局散变量

export const state = {
  todos:     /** @type {Array} */   [],
  filter:    /** @type {string} */  'all',   // 'all' | 'active' | 'done'
  editingId: /** @type {number|null} */ null,

  // ── todos ──────────────────────────────────────────
  setTodos(list) { this.todos = list },

  /** 在列表头部插入新 todo */
  prepend(todo) { this.todos.unshift(todo) },

  /** 用 updated 替换列表中同 id 的条目 */
  replace(updated) {
    const i = this.todos.findIndex(t => t.id === updated.id)
    if (i !== -1) this.todos[i] = updated
  },

  /** 从列表中移除指定 id */
  remove(id) { this.todos = this.todos.filter(t => t.id !== id) },

  /** 按当前 filter 和优先级排序，返回可展示的列表副本 */
  visible() {
    const PRIORITY = { high: 0, medium: 1, low: 2 }
    return this.todos
      .filter(t => {
        if (this.filter === 'active') return !t.done
        if (this.filter === 'done')   return t.done
        return true
      })
      .slice()
      .sort((a, b) => {
        if (a.done !== b.done) return a.done ? 1 : -1
        return (PRIORITY[a.priority] ?? 1) - (PRIORITY[b.priority] ?? 1)
      })
  },

  // ── filter ─────────────────────────────────────────
  setFilter(f) { this.filter = f },

  // ── editingId ──────────────────────────────────────
  startEdit(id) { this.editingId = id },
  endEdit()     { this.editingId = null },
}
