package t28

import (
	"fmt"
	"testing"
	"time"
)

/**
  select 多路选择
*/
func Test28_0(t *testing.T) {
	select {
	case retch := <-asyncService():
		t.Log(retch)
	case <-time.After(time.Millisecond * 500):
		t.Log("time out")
	}
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
