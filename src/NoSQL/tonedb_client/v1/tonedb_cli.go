package main

import (
	"./tonedb"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
	//"net"
	//"errors"
)

const (
	VERSION = "1.0.0"
)

func printVersion() {
	fmt.Println("tonecli v" + VERSION + " by T-ONE")
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
		os.Exit(1)
	}
}

func monitorGoroutine() {
	for {
		fmt.Println("NumGoroutine: ", runtime.NumGoroutine())

		// Gosched yields the processor, allowing other goroutines to run.
		// It does not suspend the current goroutine, so execution resumes automatically.
		runtime.Gosched()

		// sleep for a while
		time.Sleep(time.Duration(3) * time.Second)
	}
}

func routine(service string, mode string, id string, times int, lockChan chan bool) {

	fmt.Println(id + " start")

	db, err := tonedb.ConnectService(service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tonedb.ConnectService: [%s]%s\n", id, err.Error())
		lockChan <- true
		return
	}

	defer db.Close()
	var val interface{}

	for i := 0; i < int(times); i++ {

		if mode == "get" {
			val, err = db.Get("a")
			if err != nil {
				fmt.Fprintf(os.Stderr, "db.Get: [%s]%s\n", id, err.Error())
				lockChan <- true
				return
			}
			fmt.Printf("%s\n", val)

		} else if mode == "multi_get" {

			keys := []string{}
			keys = append(keys, "c")
			keys = append(keys, "d")
			val, err = db.Do("multi_get", "a", "b", keys)
			if err != nil {
				fmt.Fprintf(os.Stderr, "db.Do multi_get: [%s]%s\n", id, err.Error())
				lockChan <- true
				return
			}
			fmt.Printf("%s\n", val)

		} else if mode == "set" {

			val, err = db.Set("a", "gerry")
			if err != nil {
				fmt.Fprintf(os.Stderr, "db.Set: [%s]%s\n", id, err.Error())
				lockChan <- true
				return
			}
			fmt.Printf("%s\n", val)
			val, err = db.Set("b", "wcdj")
			if err != nil {
				fmt.Fprintf(os.Stderr, "db.Set: [%s]%s\n", id, err.Error())
				lockChan <- true
				return
			}
			fmt.Printf("%s\n", val)

		} else {

			keys := []string{}
			keys = append(keys, "c")
			keys = append(keys, "d")
			val, err = db.Do("multi_get", "a", "b", keys)
			if err != nil {
				fmt.Fprintf(os.Stderr, "db.Do multi_get: [%s]%s\n", id, err.Error())
				lockChan <- true
				return
			}
			fmt.Printf("%s\n", val)

			db.Set("a", "gerry")
			val, err = db.Get("a")
			fmt.Printf("%s\n", val)

			db.Del("a")
			val, err = db.Get("a")
			if err != nil {
				fmt.Fprintf(os.Stderr, "db.Get: [%s]%s\n", id, err.Error())
				lockChan <- true
				return
			}
			fmt.Printf("%s\n", val)

			db.Do("zset", "z", "a", 3)
			db.Do("multi_zset", "z", "b", -2, "c", 5, "d", 3)
			resp, err := db.Do("zrange", "z", 0, 10)
			if err != nil {
				fmt.Fprintf(os.Stderr, "db.Do zrange: [%s]%s\n", id, err.Error())
				lockChan <- true
				return
			}
			if len(resp)%2 != 1 {
				fmt.Fprintf(os.Stderr, "[%s]bad response\n", id)
				lockChan <- true
				return
			}

			fmt.Printf("Status: %s\n", resp[0])
			for i := 1; i < len(resp); i += 2 {
				fmt.Printf("  %s : %3s\n", resp[i], resp[i+1])
			}
		}
	}

	lockChan <- true

}

func main() {

	var ServiceInfo = flag.String("s", "127.0.0.1:9001", "host:port")
	var Mode = flag.String("m", "get", "get|multi_get|set")
	var RoutineNum = flag.Int("t", 1, "routine counts")
	var RoutineReqNum = flag.Int("c", 1, "each routine requst counts")
	var GoMaxProcs = flag.Int("core", 0, "set GOMAXPROCS variable for simultaneously")
	var printVer bool
	flag.BoolVar(&printVer, "v", false, "print version")

	flag.Parse()
	fmt.Println("ServiceInfo: ", *ServiceInfo)
	fmt.Println("Mode: ", *Mode)
	fmt.Println("RoutineNum: ", *RoutineNum)
	fmt.Println("RoutineReqNum: ", *RoutineReqNum)
	fmt.Println("GoMaxProcs: ", *GoMaxProcs)
	if printVer {
		printVersion()
		os.Exit(0)
	}

	// The GOMAXPROCS variable limits the number of operating system threads that can execute user-level Go code simultaneously.
	fmt.Println("NumCPU: ", runtime.NumCPU())
	if *GoMaxProcs > 0 {
		runtime.GOMAXPROCS(*GoMaxProcs)
	}

	// monitor routines status
	//go monitorGoroutine()

	//reqs := 1000000
	//concurrentConnections := 100
	reqs := (*RoutineNum) * (*RoutineReqNum)
	concurrentConnections := *RoutineNum

	lockChan := make(chan bool, concurrentConnections)

	start := time.Now()
	var name string
	for i := 0; i < concurrentConnections; i++ {
		name = fmt.Sprintf("routine_%02d", i)
		go routine(*ServiceInfo, *Mode, name, *RoutineReqNum, lockChan)
	}
	for i := 0; i < int(concurrentConnections); i++ {
		<-lockChan
	}
	elapsed := 1000000 * time.Since(start).Seconds()

	// show result info
	fmt.Println("\n--------------- results ---------------")
	fmt.Println("test svr info: ", *ServiceInfo)
	fmt.Println("Mode: ", *Mode)
	fmt.Println("NumCPU: ", runtime.NumCPU())
	fmt.Println("runtime.GoMaxProcs: ", *GoMaxProcs)
	fmt.Println("routine counts: ", *RoutineNum)
	fmt.Println("each routine reqs: ", *RoutineReqNum)
	fmt.Println("total reqs: ", reqs)
	fmt.Println("time elapsed: ", elapsed, "us")
	fmt.Println("avg time: ", elapsed/float64(reqs), "us")
	//fmt.Println("reqs of per seconds(QPS): ", float64(reqs)/elapsed*1000000)
}
