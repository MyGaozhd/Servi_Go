package t20

import (
	"errors"
	"testing"
)

func Test20_0(t *testing.T) {
	_, err := errorMethod()
	if err != nil {
		t.Log(err)
	}
}

func errorMethod() (int, error) {
	return 0, errors.New("方法错误")
}
