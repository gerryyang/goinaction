package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func f1() {
	for {
		fmt.Println("f1")
		time.Sleep(1 * time.Second)
	}
}

func f2() {
	for {
		fmt.Println("f2")
		time.Sleep(2 * time.Second)
	}
}

func main() {

	fmt.Println("Hello World")

	go f1()
	go f2()

	// s := make([]int, 10)
	// s = append(s, 1)
	// for _, v := range s {
	// 	fmt.Println(v)
	// }

	go func() {
		fmt.Println("goroutine begin")

		// s := make([]int, 20)
		// s = append(s, 1)
		// for _, v := range s {
		// 	fmt.Println(v)
		// }

		time.Sleep(1 * time.Minute)

		fmt.Println("goroutine end")
	}()

	// http://localhost:8080/debug/pprof/
	http.ListenAndServe("localhost:8080", nil)
}
