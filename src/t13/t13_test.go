package t13

import (
	"fmt"
	"testing"
)

func Test13_0(t *testing.T) {
	//用来释放资源释放锁，类似java的finally
	defer clear()
	t.Log("start")
	//抛出异常
	panic("error")
}

func clear() {
	fmt.Print("clear resources")
}
