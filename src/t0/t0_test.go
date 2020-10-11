package t0

import (
	"fmt"
	"testing"
)

func Test0_0(t *testing.T) {
	a := 1
	b := 1

	fmt.Print(a)
	for i := 0; i < 5; i++ {
		fmt.Print(" ", b)
		tem := a
		a = b
		b = tem + a
	}
	fmt.Println()
}
