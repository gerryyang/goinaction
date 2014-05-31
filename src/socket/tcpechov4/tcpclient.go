package main

import (
	//"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

const (
	VERSION = "1.0.0"
)

func printVersion() {
	fmt.Println("tcpclient v" + VERSION + " by gerryyang")
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
		os.Exit(1)
	}
}

func ping(tcpAddr *net.TCPAddr, id string, times int, lockChan chan bool) {

	fmt.Println(id + " start")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DialTCP: [%s]%s\n", id, err.Error())
		lockChan <- true
		return
	}

	for i := 0; i < int(times); i++ {
		_, err := conn.Write([]byte("Ping"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Write: %s\n", err.Error())
			lockChan <- true
			return
		}

		var buf [4]byte
		_, err = conn.Read(buf[0:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read: %s\n", err.Error())
			lockChan <- true
			return
		}
		//fmt.Fprintf(os.Stdout, "%d:%s\n", n, buf)

	}
	lockChan <- true
	conn.Close()
}

func main() {

	var ServiceInfo = flag.String("s", "127.0.0.1:9001", "host:port")
	var RoutineNum = flag.Int("t", 1, "routine counts")
	var RoutineReqNum = flag.Int("c", 1, "each routine requst counts")
	var GoMaxProcs = flag.Int("core", 0, "set GOMAXPROCS variable for simultaneously")
	var printVer bool
	flag.BoolVar(&printVer, "v", false, "print version")

	flag.Parse()
	fmt.Println("ServiceInfo: ", *ServiceInfo)
	fmt.Println("RoutineNum: ", *RoutineNum)
	fmt.Println("RoutineReqNum: ", *RoutineReqNum)
	fmt.Println("GoMaxProcs: ", *GoMaxProcs)
	if printVer {
		printVersion()
		os.Exit(0)
	}

	// The GOMAXPROCS variable limits the number of operating system threads that can execute user-level Go code simultaneously.
	if *GoMaxProcs > 0 {
		runtime.GOMAXPROCS(*GoMaxProcs)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", *ServiceInfo)
	fatal(err)

	//var totalPings int = 1000000
	//var concurrentConnections int = 100
	var totalPings int = (*RoutineNum) * (*RoutineReqNum)
	var concurrentConnections int = *RoutineNum
	var pingsPerConnection int = totalPings / concurrentConnections
	var actualTotalPings int = pingsPerConnection * concurrentConnections

	lockChan := make(chan bool, concurrentConnections)

	start := time.Now()
	var name string
	for i := 0; i < concurrentConnections; i++ {
		name = fmt.Sprintf("routine_%02d", i)
		go ping(tcpAddr, name, pingsPerConnection, lockChan)
	}
	for i := 0; i < int(concurrentConnections); i++ {
		<-lockChan
	}
	elapsed := 1000000 * time.Since(start).Seconds()

	// show result info
	fmt.Println("--------------- results ---------------")
	fmt.Println("test svr info: ", *ServiceInfo)
	fmt.Println("routine counts: ", *RoutineNum)
	fmt.Println("routine reqs: ", *RoutineReqNum)
	fmt.Println("total reqs: ", totalPings)
	fmt.Println("runtime.GoMaxProcs: ", *GoMaxProcs)
	fmt.Println("time elapsed: ", elapsed, "us")
	fmt.Println("time avg: ", elapsed/float64(actualTotalPings), "us")
	fmt.Println("reqs of per seconds: ", float64(actualTotalPings)/elapsed*1000000)
}
