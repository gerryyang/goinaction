## 压测以及Go代码性能分析

### 压测

原则上只压测C端使用的接口，B端后台管理直接操作数据库的接口一般不进行压测。

此处以Apache ab压测工具为例：

>ab -c 10 -n 100000 -p ./post.json -T "application/json" $url

- -c concurrency数量，也就是并发数，并不是每秒发多少个请求

- -n total request number，一共进行多少次请求

- -T Content-type header for POSTing, 指定request header的Content-Type

- -p 指定post
  body的数据文件，`./post.json`文件里面是需要post到接口的JSON数据，如： 
  
    ```json
    {
	"VenderId":1,
	"ConfigKey":"PayedOrderCancelTime",
	"RegionId": "230508"
    }
    ```
压测时可以适当调整-c参数查看接口性能

### Go代码性能分析

此处使用Go语言自带的pprof工具。官网参考：https://golang.org/pkg/runtime/pprof/

pprof有两种使用方式：

1. 生成pprof文件，根据pprof文件分析代码性能

2. 在HTTP服务中多监听一个pprof端口，使用HTTP服务实时查看性能分析

此处方式1进行性能分析为例：

#### 代码准备

首先需要添加性能分析代码，将pprof相关数据写入磁盘文件

贴一段示例代码：
```go
func programStat() {
	// 负载监控
	go func() {
		fmt.Printf(fmt.Sprintf("-----Start Monitor-----"))
		for {
			cpuFileName := fmt.Sprintf("./pprof/pprof_cpu_%v", time.Now().Minute())
			fcpu, err := os.Create(cpuFileName)
			if err != nil {
				fmt.Printf(fmt.Sprintf("Open pprof_cpu Err: %+v", err))
				break
			}
			if err = pprof.StartCPUProfile(fcpu); err != nil {
				fmt.Printf("Main", fmt.Sprintf("StartCPUProfile Err: %+v", err))
				fcpu.Close()
				break
			}
			select {
			case <-time.After(1 * time.Minute):
				pprof.StopCPUProfile()
				fcpu.Close()
				fmt.Printf("Main", fmt.Sprintf("-----Update Cpu Info-----"))
				fmem, err := os.Create(fmt.Sprintf("./pprof/pprof_mem_%v", time.Now().Minute()))
				if err != nil {
					fmt.Printf("Main", fmt.Sprintf("Open pprof_mem Err: %+v", err))
					break
				}
				if err = pprof.WriteHeapProfile(fmem); err != nil {
					fmt.Printf("Main", fmt.Sprintf("WriteHeapProfile Err: %+v", err))
					fmem.Close()
					break
				}
				fmem.Close()
				fmt.Printf("Main", fmt.Sprintf("-----Update Mem Info-----"))
				fblock, err := os.Create(fmt.Sprintf("./pprof/pprof_block_%v", time.Now().Minute()))
				if err != nil {
					fmt.Printf("Main", fmt.Sprintf("Open pprof_block Err: %+v", err))
					break
				}
				if err = pprof.Lookup("block").WriteTo(fblock, 0); err != nil {
					fmt.Printf("Main", fmt.Sprintf("Block WriteTo Err: %+v", err))
					fblock.Close()
					break
				}
				fblock.Close()
				fmt.Printf("Main", fmt.Sprintf("-----Update Block Info-----"))
			}
		}
	}()
}
```

然后在main.go的入口函数里调用就行了: 
```go
func main() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	confFile := flagSet.String("c", "app.properties", "config file name")
	flagSet.Parse(os.Args[1:])
	//性能分析
	programStat()
	framework.Run(*confFile)
}
```
注意在正式环境不建议长期加上programStat()的代码。可以放置一两天收集性能数据文件，但是不推荐一直跑。

#### 查看pprof文件

加入性能分析代码跑起来后，可以适当压测接口，以便让代码在压力环境下产出性能数据。

上面的示例代码是每分钟会生成一个pprof文件，结果如下： 

```
-rw-r--r-- 1 root root  20105 Apr  3 17:58 pprof_mem_58
-rw-r--r-- 1 root root    783 Apr  3 17:58 pprof_block_58
-rw-r--r-- 1 root root  20390 Apr  3 17:59 pprof_mem_59
-rw-r--r-- 1 root root    784 Apr  3 17:59 pprof_block_59
-rw-r--r-- 1 root root 192996 Apr  3 18:00 pprof_cpu_17
```

然后将这些pprof文件下载到本地，同时还需要将Go编译之后跑起来的二进制文件一起下载到本地。

#### 进行性能分析

使用pprof文件的方式有两种，一种是以命令行的方式，另一种是图像界面。

很明显选中后一种。

使用方法： 

> go tool pprof -http=":8081" [binary] [profile]

对应我们的就是：

> go tool pprof -http=":8081" vender_service ./pprof_mem_58

然后就会自动弹出浏览器窗口。在浏览器窗口里查看就行了。

注意：Go版本1.11及以后的版本才能运行这些命令。

如果提示缺少组件，如dot not found，请参考：https://www.graphviz.org/download/



