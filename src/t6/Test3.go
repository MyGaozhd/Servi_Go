package main

import "fmt"

func main() {
	var son = &son{}
	son.Name = "123"
	fmt.Println(son.pr)
}

type father struct {
	Name string
	Age  int
}

type son struct {
	father
	hobby string
}

func (this son) pr() string {
	return this.Name
}

type son2 struct {
	father father
	hobby  string
}
