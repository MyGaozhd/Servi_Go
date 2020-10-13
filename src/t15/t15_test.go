package t15

import "testing"

func Test15_0(t *testing.T) {
	var p Programmer
	p = new(GoProgrammer)
	t.Log(p.BuildCode())
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
	return "hello word"
}
