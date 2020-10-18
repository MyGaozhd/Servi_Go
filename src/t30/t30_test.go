package t30

import (
	"fmt"
	"testing"
	"time"
)

func Test30_0(t *testing.T) {
	cancelChan := make(chan struct{}, 0)
	for i := 0; i < 5; i++ {
		go func(i int, cancelCh chan struct{}) {
			for {
				if isCancelled(cancelCh) {
					break
				}
				time.Sleep(time.Millisecond * 5)
			}
			fmt.Println(i, "Cancelled")
		}(i, cancelChan)
	}
	cancel_2(cancelChan)
	time.Sleep(time.Second * 1)
}

func isCancelled(cancelChan chan struct{}) bool {
	select {
	case <-cancelChan:
		return true
	default:
		return false
	}
}

//通知五个协程channel关闭。close(channel)可以被每个channel监听到，而不是阻塞在 isCancelled的 <-cancelChan 方法
func cancel_2(cancelChan chan struct{}) {
	close(cancelChan)
}
