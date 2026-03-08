package t27

import (
	"fmt"
	"testing"
	"time"
)

/**
  有buffer的channel
  retch := make(chan string, 1)
  不必须调用 <-retch 后，retch <- ret后面代码才会执行
*/
func Test27_0(t *testing.T) {
	retch := asyncService()
	otherService()
	t.Log(<-retch)
}

func otherService() {
	time.Sleep(time.Second * 1)
	fmt.Println("otherService done!")
}

func syncService() string {
	time.Sleep(time.Second * 1)
	return " service done!"
}

func asyncService() chan string {
	retch := make(chan string, 1)
	go func() {
		ret := syncService()
		fmt.Println("syncService returend")
		retch <- ret
		fmt.Println("AsyncService exit")
	}()
	return retch
}
