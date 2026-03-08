package t18

import "testing"

func Test18_0(t *testing.T) {
	var g Programmer
	g = new(GoProgrammer)
	t.Log(g.BuildCode())

	var j Programmer
	j = new(JAVAProgrammer)
	t.Log(j.BuildCode())
}

//定义一个接口
type Programmer interface {
	BuildCode() string
}

//定义一个实现类
type GoProgrammer struct {
	name string
}

//实现固定接口
func (p *GoProgrammer) BuildCode() string {
	return "go hello word"
}

//定义一个实现类
type JAVAProgrammer struct {
	name string
}

//实现固定接口
func (p *JAVAProgrammer) BuildCode() string {
	return "java hello word"
}
