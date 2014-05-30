package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"time"
)

var ServiceInfo = flag.String("s", "127.0.0.1:9001", "input host:port")
var RoutineNum = flag.Int("t", 1, "input routine counts")
var RoutineReqNum = flag.Int("c", 1, "input each routine requst counts")

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func routine(id string, delay time.Duration, tcpAddr *net.TCPAddr, lock chan bool) {

	for i := 0; i < *RoutineReqNum; i++ {

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		checkError(err)

		req := id + ": hello"
		_, err = conn.Write([]byte(req))
		checkError(err)

		// ReadAll reads from r until an error or EOF and returns the data it read.
		// A successful call returns err == nil, not err == EOF.
		// Because ReadAll is defined to read from src until EOF, it does not treat an EOF from Read as an error to be reported.
		//result, err := ioutil.ReadAll(conn)
		_, err = ioutil.ReadAll(conn)
		checkError(err)

		//fmt.Println(string(result))

		conn.Close()
	}
	lock <- true
}

func main() {

	/*
		if len(os.Args) != 2 {
			fmt.Fprintf(os.Stderr, "Usage: %s host:port\n", os.Args[0])
			os.Exit(1)
		}
	*/
	flag.Parse()
	fmt.Println("ServiceInfo: ", *ServiceInfo)
	fmt.Println("RoutineNum: ", *RoutineNum)
	fmt.Println("RoutineReqNum: ", *RoutineReqNum)

	//service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *ServiceInfo)
	checkError(err)

	var name string
	begin := time.Now()
	lock := make(chan bool, *RoutineNum)
	for i := 0; i < *RoutineNum; i++ {
		name = fmt.Sprintf("routine_%02d", i)
		go routine(name, time.Duration(rand.Intn(3))*time.Second, tcpAddr, lock)
	}
	for i := 0; i < *RoutineNum; i++ {
		<-lock
	}
	cost := time.Since(begin).Seconds()
	end := time.Now()
	fmt.Println("last: ", end.Sub(begin))

	/*var input string
	fmt.Scanln(&input)*/
	fmt.Println("avg: ", 1000000*cost/float64((*RoutineNum)*(*RoutineReqNum)), "us")

	//time.Sleep(time.Duration(100) * time.Second)
	os.Exit(0)

}
