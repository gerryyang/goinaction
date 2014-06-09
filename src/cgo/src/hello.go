package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "../inc/hello.h"

void hello(const char *name)
{
	printf("%s\n", name);
	return;
}
*/
import "C"

import (
	"fmt"
	"time"
)

func Hello(s string) {
	cs := C.CString(s)
	C.hello(cs)
}

func main() {

	Seed(1000)

	fmt.Println(int(C.random()))
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println(int(C.random()))

	fmt.Println("getpid:", int(C.getpid()))
	C.puts(C.CString("call C puts"))

	Hello("call C's go wrapper func")
	C.hello(C.CString("call C hello func"))
}

func Seed(i int) {
	C.srandom(C.uint(i))
}
