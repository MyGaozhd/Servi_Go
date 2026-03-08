package t36

import (
	"fmt"
	"reflect"
	"testing"
)

func Test36_0(t *testing.T) {

	checkType("a")
}

func checkType(v interface{}) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.String:
		fmt.Println("type is string")
	default:
		fmt.Println("unknown type")
	}
}
