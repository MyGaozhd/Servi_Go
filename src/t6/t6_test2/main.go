package main

import (
	"fmt"
)

func main() {
	test2_1()
	fmt.Println("==============================")
	test2_2()
	fmt.Println("==============================")
	test2_3()
	fmt.Println("==============================")
	test2_4()
}

func test2_1() {
	var a [10]int

	for i := 0; i < len(a); i++ {
		a[i] = i
	}

	for i := 0; i < len(a); i++ {
		fmt.Println(a[i])
	}
}

func test2_2() {
	b := [3]int{1, 2, 3}
	modify(b)
	fmt.Println(b)
}

func modify(array [3]int) {
	fmt.Println(len(array))
	for i := 0; i < len(array); i++ {
		array[i] = len(array) - i
	}
	fmt.Println(array)
}

/**
数组切片
*/
func test2_3() {
	mySlice := make([]int, 5)
	fmt.Println(len(mySlice))
	fmt.Println(cap(mySlice))
	for i := 0; i < len(mySlice); i++ {
		mySlice[i] = i
	}
	fmt.Println(mySlice)
	mySlice = append(mySlice, 5)

	fmt.Println(mySlice)
}

/**
数组切片
*/
func test2_4() {

	var a = [5]int{1, 2, 3, 4, 5}
	fmt.Println(a)
	mySlice := a[:]
	modify2_4(mySlice)
	fmt.Println(a)
	fmt.Println(mySlice)
}

func modify2_4(array []int) {
	fmt.Println(len(array))
	for i := 0; i < len(array); i++ {
		array[i] = len(array) - i
	}
	fmt.Println(array)
}
