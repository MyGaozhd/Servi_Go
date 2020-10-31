package t32

import (
	"fmt"
	"sync"
	"testing"
	"unsafe"
)

func Test32_0(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			obj := getSingletonInstance()
			t.Log(unsafe.Pointer(obj))
			wg.Done()
		}()
	}
	wg.Wait()
}

type Singleton struct {
}

var singletonInstance *Singleton
var once sync.Once

//once.do 确保在多线程环境下只执行一次
func getSingletonInstance() *Singleton {
	once.Do(func() {
		fmt.Println("create onject")
		singletonInstance = new(Singleton)
	})

	return singletonInstance
}
