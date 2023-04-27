package demo

import (
	"fmt"
	"reflect"
	"testing"
)

type Student1 struct {
	name  string
	age   uint8
	infos interface{}
}

func Test1(t *testing.T) {
	s := &Student1{
		name: "zhangSan",
		age:  18,
		infos: map[string]interface{}{
			"class": "class1",
			"grade": uint8(1),
			"read": func(str string) {
				fmt.Println(str)
			},
		},
	}
	options := s.infos
	fmt.Println("infos type:", reflect.TypeOf(options))
	fmt.Println("infos value:", reflect.ValueOf(options))

	fmt.Println("infos.class type:", reflect.TypeOf(options.(map[string]interface{})["class"]))
	fmt.Println("infos.class value:", reflect.ValueOf(options.(map[string]interface{})["class"]))

	fmt.Println("infos.grade type:", reflect.TypeOf(options.(map[string]interface{})["grade"]))
	fmt.Println("infos.grade value:", reflect.ValueOf(options.(map[string]interface{})["grade"]))

	fmt.Println("infos.read type:", reflect.TypeOf(options.(map[string]interface{})["read"]))
	fmt.Println("infos.read value:", reflect.ValueOf(options.(map[string]interface{})["read"]))

	read := options.(map[string]interface{})["read"]
	if reflect.TypeOf(read).Kind() == reflect.Func {
		read.(func(str string))("I am reading!")
	}
}
