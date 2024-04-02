package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func perfStat() {

	if err := os.MkdirAll("./pprof", 0755); err != nil {
		fmt.Printf("MkdirAll err: %+v\n", err)
		return
	}

	go func() {
		fmt.Printf("Start Monitor\n")

		for {
			fmt.Printf("Update CPU Info\n")
			cpuFile := fmt.Sprintf("./pprof/pprof_cpu_%v", time.Now().Minute())
			fcpu, err := os.Create(cpuFile)
			if err != nil {
				fmt.Printf("Create err: %+v\n", err)
				break
			}

			if err = pprof.StartCPUProfile(fcpu); err != nil {
				fmt.Printf("StartCPUProfile err: %+v\n", err)
				fcpu.Close()
				break
			}

			// sleep 1 min
			time.Sleep(1 * time.Minute)

			pprof.StopCPUProfile()
			fcpu.Close()

			fmt.Printf("Update MEM Info\n")
			memFile := fmt.Sprintf("./pprof/pprof_mem_%v", time.Now().Minute())
			fmem, err := os.Create(memFile)
			if err != nil {
				fmt.Printf("Create err: %+v\n", err)
				break
			}
			if err = pprof.WriteHeapProfile(fmem); err != nil {
				fmt.Printf("WriteHeapProfile err: %+v\n", err)
				fmem.Close()
				break
			}
			fmem.Close()

			fmt.Printf("Update BLOCK Info\n")
			blockFile := fmt.Sprintf("./pprof/pprof_block_%v", time.Now().Minute())
			fblock, err := os.Create(blockFile)
			if err != nil {
				fmt.Printf("Create err: %+v\n", err)
				break
			}
			if err = pprof.Lookup("block").WriteTo(fblock, 0); err != nil {
				fmt.Printf("Block WriteTo err: %+v\n", err)
				fblock.Close()
				break
			}
			fblock.Close()
		}
	}()
}

func main() {

	// 性能监控
	perfStat()

	time.Sleep(10 * time.Minute)

}
