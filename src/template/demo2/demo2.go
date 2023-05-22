package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"fmt"

	"gopkg.in/yaml.v2"
)

func loadYaml(filename string) (interface{}, error) {
	yamlMap := make(map[interface{}]interface{})
	if dat, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else if err = yaml.Unmarshal(dat, &yamlMap); err != nil {
		return nil, err
	}
	return yamlMap, nil
}

func main() {

	// 加载模板文件
	tmpStr, err := ioutil.ReadFile("test.cfg")
	if err != nil {
	    fmt.Printf("ioutil.ReadFile err")
		return
	}
	fmt.Printf("%v\n", tmpStr)

	// 加载数据
	meta, err := loadYaml("meta.yaml")
	if err != nil {
	    fmt.Printf("loadYaml err")
		return
	}
	fmt.Printf("%v\n", meta)

	// 模板替换
	t, err := template.New("test").Parse(string(tmpStr))
	err = t.Execute(os.Stdout, meta)
	if err != nil {
		return
	}
}
