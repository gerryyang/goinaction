package main

import (
	"net"
	"runtime"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [4]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		if n > 0 {
			_, err = conn.Write([]byte("Pong"))
			if err != nil {
				return
			}
		}
	}
}

func main() {
	runtime.GOMAXPROCS(4)

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ":1201")
	listener, _ := net.ListenTCP("tcp", tcpAddr)

	for {
		conn, _ := listener.Accept()
		go handleClient(conn)
	}
}
