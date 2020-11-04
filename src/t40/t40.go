package main

import (
	"fmt"
	"t39"
)

// 需要将t39文件夹下的类拷贝到go/src目录下
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
		fmt.Println("The expected is 6, but the actual is %d", ret)
	}
}
