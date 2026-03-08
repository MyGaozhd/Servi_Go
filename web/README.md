# 日常事项记录网站

前后端分离架构：Go REST API + 原生 HTML/CSS/JS，数据持久化到本地 JSON 文件。

## 目录结构

```
web/
├── api/                  # Go 后端
│   ├── main.go           # 程序入口，启动 HTTP 服务
│   ├── todo.go           # 数据模型与请求/响应结构体
│   ├── store.go          # JSON 文件读写（并发安全）
│   └── handler.go        # HTTP 路由与处理逻辑
├── frontend/             # 前端静态资源
│   ├── index.html        # 页面主体
│   ├── css/style.css     # 样式
│   └── js/app.js         # 交互逻辑（原生 JS，无框架）
├── data/
│   └── todos.json        # 数据持久化文件（自动维护）
└── README.md
```

## 启动方式

### 1. 启动后端 API

```bash
cd web/api
go run .
# ✅ Todo API 已启动 → http://localhost:8080
```

### 2. 打开前端页面

直接用浏览器打开（需允许 CORS，推荐用任意静态文件服务器）：

```bash
# 方式一：使用 Python（无需安装额外工具）
cd web/frontend
python3 -m http.server 3000
# 访问 http://localhost:3000

# 方式二：使用 VS Code Live Server 插件，右键 index.html → Open with Live Server
```

## API 文档

| 方法     | 路径                | 说明             |
|----------|---------------------|------------------|
| `GET`    | `/api/todos`        | 获取全部事项     |
| `GET`    | `/api/todos?done=true` | 获取已完成事项 |
| `POST`   | `/api/todos`        | 新增事项         |
| `GET`    | `/api/todos/{id}`   | 获取单条事项     |
| `PATCH`  | `/api/todos/{id}`   | 更新事项         |
| `DELETE` | `/api/todos/{id}`   | 删除事项         |
| `GET`    | `/api/health`       | 健康检查         |

### 新增事项请求体

```json
{
  "title":    "买菜",
  "note":     "记得买西红柿",
  "priority": "high"
}
```

> `priority` 可选值：`low` / `medium`（默认）/ `high`

### 更新事项请求体（所有字段均为可选）

```json
{
  "title":    "新标题",
  "note":     "新备注",
  "priority": "low",
  "done":     true
}
```

## 功能列表

- ✅ 新增事项（标题、备注、优先级）
- ✅ 完成 / 取消完成
- ✅ 编辑事项
- ✅ 删除事项
- ✅ 全部 / 待完成 / 已完成 筛选
- ✅ 按优先级自动排序（高 → 中 → 低）
- ✅ 统计栏（全部 / 待完成 / 已完成数量）
- ✅ 数据持久化到 `data/todos.json`
