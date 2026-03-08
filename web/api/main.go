package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MyGaozhd/Servi_Go/web/api/handler"
	"github.com/MyGaozhd/Servi_Go/web/api/service"
	"github.com/MyGaozhd/Servi_Go/web/api/store"
)

const (
	addr     = ":8080"
	dataFile = "../data/todos.json"
)

func main() {
	// 依赖组装：store → service → handler
	s, err := store.New(dataFile)
	if err != nil {
		log.Fatalf("初始化数据存储失败: %v", err)
	}

	svc := service.New(s)
	h := handler.New(svc)

	fmt.Printf("✅ Todo API 已启动 → http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatal(err)
	}
}
