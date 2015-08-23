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
	VERSION         = "1.0.0"
	MAX_BUFFER_SIZE = 1024 * 1024
)

func printVersion() {
	fmt.Println("file sender" + VERSION + " by gerryyang")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func handleCommandLine() (file string, service string) {
	var f = flag.String("f", "data2.txt", "file")
	var s = flag.String("s", "127.0.0.1:9001", "host:ip")
	var printVer bool
	flag.BoolVar(&printVer, "v", false, "print version")

	flag.Parse()
	if printVer {
		printVersion()
		os.Exit(0)
	}
	fmt.Println("file:", *f)
	fmt.Println("host:ip:", *s)
	return *f, *s
}

func writeFile(file string, offset int, s_real_req string) {

	fout, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	checkError(err)
	defer fout.Close()

	var bytes_buf bytes.Buffer
	bytes_buf.WriteString(s_real_req)

	var write_len int = len(s_real_req)
	cnt, err := fout.WriteAt(bytes_buf.Bytes(), int64(offset))
	if cnt != write_len {
		fmt.Fprintf(os.Stderr, "%d cnt[%d] WriteAt: %s\n", offset, cnt, err.Error())
		return
	}
	fmt.Printf("WriteAt offset[%d] bytes[%d] s_real_req[%#v]\n", offset, cnt, s_real_req)
}

func handleClient(conn net.Conn, file string) {
	defer conn.Close()

	var buf [MAX_BUFFER_SIZE]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read: %s\n", err.Error())
			return
		}
		fmt.Printf("len:%d req:%#v\n", n, buf[0:n])

		var bytes_buf bytes.Buffer
		bytes_buf.Write(buf[0:n])
		len, _ := bytes_buf.ReadByte()

		var offset int = int(len)
		var s_real_req string = string(bytes_buf.Bytes())
		fmt.Printf("file[%s] offset[%d] s_real_req[%#v]\n", file, offset, s_real_req)

		writeFile(file, offset, s_real_req)

		if n > 0 {
			_, err = conn.Write([]byte("ok"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Write: %s\n", err.Error())
				return
			}
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, service := handleCommandLine()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Accept: %s\n", err.Error())
			continue
		}
		go handleClient(conn, file)
	}
}
