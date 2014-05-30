package main

import (
	"fmt"
	"runtime"
)

var workers = runtime.NumCPU()

type Result struct {
	jobname    string
	resultcode int
	resultinfo string
}

type Job struct {
	jobname string
	results chan<- Result
}

func main() {

	// go语言里大多数并发程序的开始处都有这一行代码, 但这行代码最终将会是多余的,  
	// 因为go语言的运行时系统会变得足够聪明以自动适配它所运行的机器  
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 返回当前处理器的数量  
	//fmt.Println(runtime.GOMAXPROCS(0))  
	// 返回当前机器的逻辑处理器或者核心的数量  
	//fmt.Println(runtime.NumCPU())  

	// 模拟8个工作任务  
	jobnames := []string{"gerry", "wcdj", "golang", "C++", "Lua", "perl", "python", "C"}
	doRequest(jobnames)
}

func doRequest(jobnames []string) {

	// 定义需要的channels切片  
	jobs := make(chan Job, workers)
	results := make(chan Result, len(jobnames))
	done := make(chan struct{}, workers)

	// ---------------------------------------------  
	/* 
	* 下面是go协程并发处理的一个经典框架 
	*/

	// 将需要并发处理的任务添加到jobs的channel中  
	go addJobs(jobs, jobnames, results) // Executes in its own goroutine  

	// 根据cpu的数量启动对应个数的goroutines从jobs争夺任务进行处理  
	for i := 0; i < workers; i++ {
		go doJobs(done, jobs) // Each executes in its own goroutine  
	}

	// 新创建一个接受结果的routine, 等待所有worker routiines的完成结果, 并将结果通知主routine  
	go awaitCompletion(done, results)

	// 在主routine输出结果  
	processResults(results)
	// ---------------------------------------------  

}

func addJobs(jobs chan<- Job, jobnames []string, results chan<- Result) {
	for _, jobname := range jobnames {

		// 在channel中添加任务  
		jobs <- Job{jobname, results}
	}
	close(jobs)
}

func doJobs(done chan<- struct{}, jobs <-chan Job) {

	// 在channel中取出任务并计算  
	for job := range jobs {

		/* 
		* 定义类型自己的方法来处理业务逻辑, 实现框架和业务分离 
		*/
		job.Do()
	}

	// 所有任务完成后的结束标志, 一个空结构体切片  
	done <- struct{}{}
}

// 方法是作用在自定义类型的值上的一类特殊函数  
func (job Job) Do() {

	// 打印当前处理的任务名称  
	fmt.Printf("... doing work in [%s]\n", job.jobname)

	// 模拟处理结果  
	if job.jobname == "golang" {
		job.results <- Result{job.jobname, 0, "OK"}
	} else {
		job.results <- Result{job.jobname, -1, "Error"}
	}
}

func awaitCompletion(done <-chan struct{}, results chan Result) {
	for i := 0; i < workers; i++ {
		<-done
	}
	close(results)
}

func processResults(results <-chan Result) {
	for result := range results {
		fmt.Printf("done: %s,%d,%s\n", result.jobname, result.resultcode, result.resultinfo)
	}
}

