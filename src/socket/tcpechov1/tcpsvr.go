package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(4)

	//service := "127.0.0.1:9001"
	service := ":9001"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
		/*
			daytime := time.Now().String()
			// do not care return value
			conn.Write([]byte(daytime))
			// close client
			conn.Close()
		*/
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	daytime := time.Now().String()
	conn.Write([]byte(daytime))
	//time.Sleep(time.Duration(100) * time.Second)
}
