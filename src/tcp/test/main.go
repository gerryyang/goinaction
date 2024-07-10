package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "golang.org:80")
	if err != nil {
		fmt.Println("net.Dial err: ", err)
	}
	conn.Close()

	// 验证 conn Close 后调用 Write 是否会 crash
	if _, err := conn.Write([]byte("hello")); err != nil {
		fmt.Println("conn.Write err: ", err)
	}
}

/*
$ ./test
conn.Write err:  write tcp 9.134.129.173:50710->142.251.211.241:80: use of closed network connection
*/
