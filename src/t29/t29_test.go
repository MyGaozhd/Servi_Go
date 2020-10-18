package t29

import (
	"fmt"
	"sync"
	"testing"
)

func Test29_0(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)
	dataProduce(ch, &wg)
	wg.Add(2)
	dataReseiver(ch, &wg, 1)
	dataReseiver(ch, &wg, 2)
	wg.Wait()
}

func dataProduce(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		//关闭channle
		close(ch)
		wg.Done()
	}()
}

/**
  data, ok := <-ch;
  ok =true,表示channle没有关闭。false表示已经关闭
*/

func dataReseiver(ch chan int, wg *sync.WaitGroup, count int) {
	go func() {
		for {
			if data, ok := <-ch; ok {
				fmt.Println(count, data)
			} else {
				break
			}
		}
		wg.Done()
	}()
}
