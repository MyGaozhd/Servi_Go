package t19

import (
	"fmt"
	"testing"
)

func Test19_0(t *testing.T) {
	emptyInterface("111")
	emptyInterface(1)
	emptyInterface(true)
}
func emptyInterface(p interface{}) {

	switch v := p.(type) {
	case string:
		fmt.Println("string", v)
	case int:
		fmt.Println("int", v)
	default:
		fmt.Println(" unknow type")
	}
}
