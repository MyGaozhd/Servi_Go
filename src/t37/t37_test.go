package t37

import (
	"reflect"
	"testing"
)

func Test37_0(t *testing.T) {
	u := &User{"servi", 31}
	//获取名称方式一
	t.Log(reflect.ValueOf(*u).FieldByName("name"))

	if name, ok := reflect.TypeOf(*u).FieldByName("name"); !ok {
		t.Log("field to get name field")
	} else {
		t.Log(name.Tag.Get("format"))
	}
	t.Log(u)
	setAge := reflect.ValueOf(u).MethodByName("SetAge")
	setAge.Call([]reflect.Value{reflect.ValueOf(1)})
	t.Log(u)
}

type User struct {
	name string `format:"normal"`
	age  int
}

func (u *User) setName(name string) *User {
	u.name = name
	return u
}

//反射调用必须大写
func (u *User) SetAge(newAge int) *User {
	u.age = newAge
	return u
}
