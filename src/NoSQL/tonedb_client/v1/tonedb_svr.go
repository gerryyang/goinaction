package main

import (
	"bytes"
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
	fmt.Println("tonesvr v" + VERSION + " by T-ONE")
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
		os.Exit(1)
	}
}

func send(conn net.Conn, args ...interface{}) error {
	var buf bytes.Buffer
	for _, arg := range args {
		var s string
		switch arg := arg.(type) {
		case string:
			s = arg
		case []byte:
			s = string(arg)
		case []string:
			for _, s := range arg {
				buf.WriteString(fmt.Sprintf("%d", len(s)))
				buf.WriteByte('\n')
				buf.WriteString(s)
				buf.WriteByte('\n')
			}
			continue
		case int:
			s = fmt.Sprintf("%d", arg)
		case int64:
			s = fmt.Sprintf("%d", arg)
		case float64:
			s = fmt.Sprintf("%f", arg)
		case bool:
			if arg {
				s = "1"
			} else {
				s = "0"
			}
		case nil:
			s = ""
		default:
			return fmt.Errorf("bad arguments")
		}
		buf.WriteString(fmt.Sprintf("%d", len(s)))
		buf.WriteByte('\n')
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	_, err := conn.Write(buf.Bytes())
	return err
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// 可限制单个Value大小小于1MB
	var buf [1024*1024 + 64]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read: %s\n", err.Error())
			return
		}
		fmt.Fprintf(os.Stdout, "%d:%s\n", n, buf)

		// TODO: tonedb operation

		if n > 0 {

			//_, err = conn.Write([]byte("Pong"))
			//err := send(conn, "not_found")
			err := send(conn, "ok", "find it")
			if err != nil {
				return
			}

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
