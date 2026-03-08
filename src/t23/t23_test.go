package t23

import (
	"testing"
	"time"
)

// go加上func就可以开启协程
func Test20_0(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			t.Log(i)
		}(i)
	}

	time.Sleep(time.Second * 1)
}
