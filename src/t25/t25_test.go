package t25

import (
	"sync"
	"testing"
)

/**
  var wg sync.WaitGroup
  wg.Add(1)
  wg.Done()
  wg.Wait()
*/
func Test25_0(t *testing.T) {
	var mu sync.Mutex
	count := 0
	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				mu.Unlock()
			}()
			mu.Lock()
			count++
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log(count)
}
