package demo

import (
	"fmt"
	"reflect"
	"testing"
)

type Student3 struct {
	Name string `json:"name1" db:"name2"`
	Age  int    `json:"age1" db:"age2"`
}

func Test3(t *testing.T) {
	s := &Student3{
		Name: "zhangSan",
		Age:  18,
	}
	v := reflect.ValueOf(s)

	fmt.Println("set ability of v:", v.CanSet())           // false
	fmt.Println("set ability of Elem:", v.Elem().CanSet()) // true

	if v.Elem().CanSet() {
		for i := 0; i < v.Elem().NumField(); i++ {
			switch v.Elem().Field(i).Kind() {
			case reflect.String:
				v.Elem().Field(i).Set(reflect.ValueOf("lisi"))
			case reflect.Int:
				v.Elem().Field(i).Set(reflect.ValueOf(20))
			}
		}
	}

	fmt.Println("v: ", v)
	fmt.Println("student: ", v.Interface().(*Student3))
}
