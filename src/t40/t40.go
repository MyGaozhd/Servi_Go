package main

import (
	"fmt"
	"github.com/MyGaozhd/Servi_Go/src/t39"
)

// t40 演示 pipe-filter 架构模式（使用 t39 包）
func main() {
	spliter := t39.NewSplitFilter(",")
	converter := t39.NewTointFilter()
	sum := t39.NewSumFilter()
	sp := t39.NewStraightPipeline("p1", spliter, converter, sum)
	ret, err := sp.Process("1,2,3")
	if err != nil {
		fmt.Println(err)
	}
	if ret == 6 {
		fmt.Printf("The expected is 6, and the actual is %d\n", ret)
	}
}
