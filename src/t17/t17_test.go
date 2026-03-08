package t17

//1、不支持子类复写子类的方法
//2、不支持里氏替换原则

type father struct {
	Name string
	Age  int
}

type son struct {
	father
	hobby string
}
