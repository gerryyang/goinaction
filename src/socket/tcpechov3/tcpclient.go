package main

import (
	"fmt"
	"net"
	"runtime"
	"time"
)

func ping(times int, lockChan chan bool) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "localhost:1201")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)

	for i := 0; i < int(times); i++ {
		_, _ = conn.Write([]byte("Ping"))
		var buff [4]byte
		_, _ = conn.Read(buff[0:])
	}
	lockChan <- true
	conn.Close()
}

func main() {
	runtime.GOMAXPROCS(4)

	var totalPings int = 1000000
	var concurrentConnections int = 100
	var pingsPerConnection int = totalPings / concurrentConnections
	var actualTotalPings int = pingsPerConnection * concurrentConnections

	lockChan := make(chan bool, concurrentConnections)

	start := time.Now()
	for i := 0; i < concurrentConnections; i++ {
		go ping(pingsPerConnection, lockChan)
	}
	for i := 0; i < int(concurrentConnections); i++ {
		<-lockChan
	}
	elapsed := 1000000 * time.Since(start).Seconds()
	fmt.Println(elapsed/float64(actualTotalPings), "us")
}
