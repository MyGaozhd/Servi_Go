package main

import "fmt"

var f0 = "f0"

const f1, f2, f3 = "fff1", true, 0
const (
	unknown = 0
	male    = 1
	female  = 2
)

const (
	x = iota
	y = iota
	z
)

func main() {
	test1()
	fmt.Println("==============================")
	test2()
	fmt.Println("==============================")
	test3()
	fmt.Println("==============================")
	test4()
	fmt.Println("==============================")
	test5()
	fmt.Println("==============================")
	test6()
}

func test1() {
	fmt.Println("123")

	var a = "servi"
	fmt.Println(a)

	b := true
	fmt.Println(b)

	f := "123456a"
	fmt.Println(f)

	var v1, v2, v3 int
	v1, v2, v3 = 1, 2, 3
	fmt.Println(v1, v2, v3)

	fmt.Println(f0, f1, f2, f3)

	fmt.Println(unknown, male, female)

	fmt.Println(x, y, z)
}

func test2() {
	a := 10
	b := 10
	fmt.Println(a * b)

	a++
	fmt.Println(a)
	a--
	fmt.Println(a)
}

func test3() {
	a := 10
	b := 20
	if a < b {
		fmt.Println("a*b ", a*b)
	}

	/* 定义局部变量 */
	var grade = "B"
	var marks = 90

	switch marks {
	case 90:
		grade = "A"
	case 80:
		grade = "B"
	case 50, 60, 70:
		grade = "C"
	default:
		grade = "D"
	}

	switch {
	case grade == "A":
		fmt.Printf("优秀!\n")
	case grade == "B", grade == "C":
		fmt.Printf("良好\n")
	case grade == "D":
		fmt.Printf("及格\n")
	case grade == "F":
		fmt.Printf("不及格\n")
	default:
		fmt.Printf("差\n")
	}
	fmt.Printf("你的等级是 %s\n", grade)
}

func test4() {
	x, y := 10, 20
	fmt.Println(x, y)
	swap1(x, y)
	fmt.Println(x, y)
	x, y = swap1(x, y)
	fmt.Println(x, y)
}

func swap1(x int, y int) (int, int) {
	temp := x
	x = y
	y = temp

	return x, y
}

func test5() {
	x, y := 10, 20
	fmt.Println(x, y)
	swap2(&x, &y)
	fmt.Println(x, y)

	a, b := 10, 20
	a, b = swap2(&a, &b)
	fmt.Println(a, b)
}

func swap2(x *int, y *int) (int, int) {
	temp := *x
	*x = *y
	*y = temp

	return *x, *y
}

func test6() {
	a := func2(1, func1)
	fmt.Println(a)
}

type F func(x int) int

func func1(x int) int {

	return x + 1
}

func func2(x int, f F) int {

	return f(x)
}
