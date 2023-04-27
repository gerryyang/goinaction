package demo

import (
	"fmt"
	"reflect"
	"testing"
)

type Student2 struct {
	Name string `json:"name1" db:"name2"`
	Age  int    `json:"age1" db:"age2"`
}

func Test2(t *testing.T) {
	var s Student2
	v := reflect.ValueOf(&s)

	// 类型
	ty := v.Type()

	// 获取字段
	for i := 0; i < ty.Elem().NumField(); i++ {
		f := ty.Elem().Field(i)
		fmt.Println(f.Tag.Get("json"))
		fmt.Println(f.Tag.Get("db"))
	}
}
