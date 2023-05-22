package main

import (
	"os"
	"text/template"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// Define a template
	tmpl := "Name: {{.Name}}, Age: {{.Age}}\n"

	// Create a new template and parse the template string
	t := template.Must(template.New("person").Parse(tmpl))

	// Create a new person object
	p := Person{Name: "Alice", Age: 30}

	// Execute the template with the person object as input
	err := t.Execute(os.Stdout, p)
	if err != nil {
		panic(err)
	}
}
