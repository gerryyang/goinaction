package main

import (
	"fmt"
	"log"
	"rsc.io/quote"
	"example.com/greetings"
)

func main() {

	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

    fmt.Println("Hello, World!")
	
	fmt.Println(quote.Go());

	message, err := greetings.Hello("gerry")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)

	/*message, err = greetings.Hello("")
	if err != nil {
		log.Fatal(err)
	}*/

	// A slice of names.
    names := []string{"Gladys", "Samantha", "Darrin"}
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(messages)

}

