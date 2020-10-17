package t21

import (
	"errors"
	"fmt"
	"testing"
)

//recover 相当于java中的cache方法，能捕获到异常
func Test21_0(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover from error")
		}
	}()

	if _, err := errorMethod1(); err != nil {
		t.Log("出现异常")
	}

	panic("退出")
}

func errorMethod1() (int, error) {
	return 0, errors.New("方法错误")
}
