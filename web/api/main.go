package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	store, err := NewStore()
	if err != nil {
		log.Fatalf("初始化数据存储失败: %v", err)
	}

	h := &Handler{store: store}

	addr := ":8080"
	fmt.Printf("✅ Todo API 已启动 → http://localhost%s\n", addr)
	fmt.Println("   GET    /api/todos          查询全部事项")
	fmt.Println("   POST   /api/todos          新增事项")
	fmt.Println("   GET    /api/todos/{id}     查询单条")
	fmt.Println("   PATCH  /api/todos/{id}     更新事项")
	fmt.Println("   DELETE /api/todos/{id}     删除事项")
	fmt.Println("   GET    /api/health         健康检查")

	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatal(err)
	}
}
