package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
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

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [4]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read: %s\n", err.Error())
			return
		}
		fmt.Fprintf(os.Stdout, "%d:%s\n", n, buf)

		if n > 0 {
			_, err = conn.Write([]byte("Pong"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Write: %s\n", err.Error())
				return
			}
		}
	}
}

func main() {

	var ServiceInfo = flag.String("s", "127.0.0.1:9001", "host:port")
	var GoMaxProcs = flag.Int("core", 0, "set GOMAXPROCS variable for simultaneously")
	var printVer bool
	flag.BoolVar(&printVer, "v", false, "print version")

	flag.Parse()
	fmt.Println("ServiceInfo: ", *ServiceInfo)
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

	listener, err := net.ListenTCP("tcp", tcpAddr)
	fatal(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Accept: %s\n", err.Error())
			continue
		}
		go handleClient(conn)
	}
}
