package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	//"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	//daytime := time.Now().String()
	//fmt.Println(daytime)

	//var buff string
	var buff [1024]byte

	{
		//iret, err := conn.Read([]byte(buff))
		iret, err := conn.Read(buff[0:])
		if err != nil {
			fmt.Println("Read err")
			return
		}
		if iret > 0 {
			////req := buff
			//fmt.Println("iret: ", iret, " req: ", req)
			////fmt.Fprintf(os.Stdout, "%d:%s\n", iret, req)

			// test timeout
			s := "nice"
			//conn.Write([]byte(daytime))
			conn.Write([]byte(s))

			//time.Sleep(time.Duration(100) * time.Second)
		}
	}

}

func main() {
	runtime.GOMAXPROCS(4)

	service := "127.0.0.1:9001"
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
