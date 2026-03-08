// api.js —— 负责所有与后端的 HTTP 通信
// 对外暴露语义化方法，调用方无需关心 URL 拼接和 fetch 细节

const BASE = 'http://localhost:8080/api'

async function _request(path, options = {}) {
  const res = await fetch(BASE + path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (res.status === 204) return null
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || '请求失败')
  return data
}

export const api = {
  /** 获取全部事项，传 done=true/false 可过滤 */
  list: (done) => {
    const qs = done === undefined ? '' : `?done=${done}`
    return _request(`/todos${qs}`)
  },
  /** 根据 ID 获取单条 */
  get: (id) => _request(`/todos/${id}`),
  /** 创建事项 */
  create: (body) => _request('/todos', { method: 'POST', body: JSON.stringify(body) }),
  /** 局部更新事项 */
  update: (id, body) => _request(`/todos/${id}`, { method: 'PATCH', body: JSON.stringify(body) }),
  /** 删除事项 */
  remove: (id) => _request(`/todos/${id}`, { method: 'DELETE' }),
}
