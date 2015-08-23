package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"time"
)

const (
	VERSION = "1.0.0"
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

func handleCommandLine() (file string, job int64, service string) {
	var f = flag.String("f", "data.txt", "file")
	var j = flag.Int64("j", 1, "job")
	var s = flag.String("s", "127.0.0.1:9001", "host:ip")
	var printVer bool
	flag.BoolVar(&printVer, "v", false, "print version")

	flag.Parse()
	if printVer {
		printVersion()
		os.Exit(0)
	}
	fmt.Println("file:", *f)
	fmt.Println("job:", *j)
	fmt.Println("host:ip:", *s)
	return *f, *j, *s
}

func proc(req *string, cid int, job int64, service string, lockChan chan bool) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	var cid_name string
	cid_name = fmt.Sprintf("%d", cid)

	send(tcpAddr, req, cid_name, job, lockChan)

}

func send(tcpAddr *net.TCPAddr, req *string, cid_name string, job int64, lockChan chan bool) {
	fmt.Printf("[%s] start\n", cid_name)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s]DialTCP: %s\n", cid_name, err.Error())
		lockChan <- true
		return
	}

	var cid int
	fmt.Sscanf(cid_name, "%d", &cid)
	var offset int = cid * int(job)

	var bytes_buf bytes.Buffer
	bytes_buf.WriteByte(byte(offset))
	bytes_buf.Write([]byte(*req))

	fmt.Printf("%s Write[%#v]\n", cid_name, bytes_buf.Bytes())
	//fmt.Printf("%q\n", bytes_buf.Bytes())

	_, err = conn.Write(bytes_buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s]Write: %s\n", cid_name, err.Error())
		lockChan <- true
		return
	}

	var ans_buf [4]byte
	_, err = conn.Read(ans_buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s]Read: %s\n", cid_name, err.Error())
		lockChan <- true
		return
	}
	fmt.Fprintf(os.Stdout, "[%s]Read: %s\n", cid_name, ans_buf)

	lockChan <- true
	conn.Close()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, job, service := handleCommandLine()

	fin, err := os.Open(file)
	checkError(err)
	defer fin.Close()

	var offset int64 = 0
	buf := make([]byte, job)

	lockChan := make(chan bool, job)
	start := time.Now()

	var cid int = 0
	cnt, err := fin.ReadAt(buf, offset)
	for cid = 0; err != io.EOF; cid++ {
		cnt, err = fin.ReadAt(buf, offset)
		if cnt != 0 {
			offset += job
			fmt.Printf("read %d bytes: %q\n", cnt, buf[:cnt])
			var req string = string(buf[:cnt])
			go proc(&req, cid, job, service, lockChan)
		}
	}
	for i := 0; i < int(cid-1); i++ {
		<-lockChan
	}
	elapsed := 1000000 * time.Since(start).Seconds()
	fmt.Println("time elapsed: ", elapsed, "us")

}
