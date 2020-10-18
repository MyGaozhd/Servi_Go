package t24

import (
	"sync"
	"testing"
	"time"
)

func Test24_0(t *testing.T) {
	var mu sync.Mutex
	count := 0
	for i := 0; i < 5000; i++ {
		go func() {
			defer func() {
				mu.Unlock()
			}()
			mu.Lock()
			count++
		}()
	}
	time.Sleep(time.Second * 2)
	t.Log(count)
}
