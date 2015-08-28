package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	VERSION    = "1.0.0"
	REQ_HEADER = "10001" // 0X2711
)

func printVersion() {
	fmt.Println("file snd" + VERSION + " by gerryyang")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func handleCommandLine() (file string, blk int64, job int64, service string) {
	var f = flag.String("f", "data.txt", "file")
	var b = flag.Int64("b", 1024, "blk")
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
	fmt.Println("blk:", *b)
	fmt.Println("job:", *j)
	fmt.Println("host:ip:", *s)
	return *f, *b, *j, *s
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8) // int64 is 8 byte
	binary.LittleEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}

func proc(req *string, reqlen int, cid int, offset int64, service string, lock_chan chan bool) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	var cid_name string
	cid_name = fmt.Sprintf("%d", cid)

	send(tcpAddr, req, reqlen, cid_name, offset, lock_chan)

}

func send(tcpAddr *net.TCPAddr, req *string, reqlen int, cid_name string, offset int64, lock_chan chan bool) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DialTCP: cid_name[%s] err[%s]\n", cid_name, err.Error())
		lock_chan <- true
		return
	}

	var cid int
	fmt.Sscanf(cid_name, "%d", &cid)

	req_header, _ := strconv.ParseInt(REQ_HEADER, 10, 64)

	var reqlen_total int64 = int64(3*8 + reqlen)

	var bytes_buf bytes.Buffer
	bytes_buf.Write(Int64ToBytes(reqlen_total))
	bytes_buf.Write(Int64ToBytes(req_header))
	bytes_buf.Write(Int64ToBytes(offset))
	bytes_buf.Write([]byte(*req))

	//fmt.Printf("send: cid_name[%s]\n", cid_name)
	//fmt.Printf("send: cid_name[%s] bytes_buf[%#v]\n", cid_name, bytes_buf.Bytes())
	//fmt.Printf("%q\n", bytes_buf.Bytes())

	// TODO here may be fail, and need to retry
	wcnt, werr := conn.Write(bytes_buf.Bytes())
	if werr != nil {
		fmt.Fprintf(os.Stderr, "Write: cid_name[%s] err[%s]\n", cid_name, err.Error())
		lock_chan <- true
		return
	}
	fmt.Printf("Write: ok cid_name[%s] wcnt[%d]\n", cid_name, wcnt)

	var ans_buf [4]byte
	_, err = conn.Read(ans_buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read: cid_name[%s] err[%s]\n", cid_name, err.Error())
		lock_chan <- true
		return
	}
	fmt.Fprintf(os.Stdout, "Read: cid_name[%s] ans_buf[%s]\n", cid_name, ans_buf)

	lock_chan <- true
	conn.Close()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, blk, job, service := handleCommandLine()
	fmt.Println("job:", job)

	fin, err := os.Open(file)
	checkError(err)
	defer fin.Close()

	var offset int64 = 0
	buf := make([]byte, blk)

	fin_info, _ := fin.Stat()
	fmt.Printf("file[%s] size[%d]\n", fin_info.Name(), fin_info.Size())

	lock_chan_len := fin_info.Size()/blk + 1
	lock_chan := make(chan bool, lock_chan_len)

	start := time.Now()
	var cid int = 0

	for {

		// limit speed
		//time.Sleep(time.Duration(1) * time.Second)

		cnt, err := fin.ReadAt(buf, offset)
		if err == io.EOF {
			if cnt == 0 {
				break
			}
			//fmt.Printf("cid[%d] offset[%d] read bytes[%d] buf[%q]\n", cid, offset, cnt, buf[:cnt])
			fmt.Printf("cid[%d] offset[%d] read bytes[%d] left\n", cid, offset, cnt)
			var req string = string(buf[:cnt])
			go proc(&req, cnt, cid, offset, service, lock_chan)
			offset += int64(cnt)
			cid++
			break

		} else if cnt != len(buf) {
			continue

		} else {
			//fmt.Printf("cid[%d] offset[%d] read bytes[%d] buf[%q]\n", cid, offset, cnt, buf[:cnt])
			fmt.Printf("cid[%d] offset[%d] read bytes[%d] all\n", cid, offset, cnt)
			var req string = string(buf[:cnt])
			go proc(&req, cnt, cid, offset, service, lock_chan)
			offset += int64(len(buf))
			cid++
		}
		/*
			if (cid % int(job)) == 0 {
				fmt.Println("unlock")
				for i := 0; i < int(job-1); i++ {
					<-lock_chan
				}
			}
		*/
	}

	// wait for all goroutine completion
	for i := 0; i < int(lock_chan_len-1); i++ {
		<-lock_chan
	}

	elapsed := 1000000 * time.Since(start).Seconds()
	fmt.Println("time elapsed: ", elapsed, "us")
}
