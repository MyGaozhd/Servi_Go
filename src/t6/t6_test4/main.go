package main

import "fmt"

func main() {
	var base base
	var base1 base1
	base1 = 1
	base = base1
	base.pl()

	var base2 base2
	base2 = 2
	base = base2
	base.pl()
}

type base interface {
	pl()
}

type base1 int

func (this base1) pl() {
	fmt.Println(this)
}

type base2 int

func (this base2) pl() {
	fmt.Println(this)
}
