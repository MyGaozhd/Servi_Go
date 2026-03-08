package t14

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

func Test14_0(t *testing.T) {
	s := Student{"001", "servi", 1}
	fmt.Println("原始e->", unsafe.Pointer(&s.name))
	t.Log(s.toString())
	t.Log(s.to_String())
}

type Student struct {
	id   string
	name string
	sex  int
}

//此种方法在实例对应方法被调用时，实例成员会进行值复制
func (s Student) toString() string {
	fmt.Println("toString[e]->", unsafe.Pointer(&s.name))
	var budilder strings.Builder
	budilder.WriteString("student->")
	budilder.WriteString(",")
	budilder.WriteString(s.id)
	budilder.WriteString(",")
	budilder.WriteString(s.name)
	budilder.WriteString(",")
	budilder.WriteString(strconv.Itoa(s.sex))
	return budilder.String()
}

//通常情况下为了避免内存拷贝我们使用此种定义方式
func (s *Student) to_String() string {
	fmt.Println("to_String[e]->", unsafe.Pointer(&s.name))
	var budilder strings.Builder
	budilder.WriteString("student->")
	budilder.WriteString(",")
	budilder.WriteString(s.id)
	budilder.WriteString(",")
	budilder.WriteString(s.name)
	budilder.WriteString(",")
	budilder.WriteString(strconv.Itoa(s.sex))
	return budilder.String()
}
