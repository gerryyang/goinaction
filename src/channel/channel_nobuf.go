package main

import (
	"fmt"
	"time"
)

func main() {
	// 创建一个channel用以同步goroutine
	done := make(chan bool)
	count := 3

	for i := 0; i < count; i++ {

		// 在goroutine中执行输出操作
		go func(i int) {
			time.Sleep(time.Duration(i) * time.Second)
			fmt.Println("goroutine message:", i)

			// 告诉main函数执行完毕.
			// 这个channel在goroutine中是可见的
			// 因为它是在相同的地址空间执行的.
			done <- true
		}(i)
	}

	<-done // 等待goroutine结束
	fmt.Println("main function message")
}

/*
output:
gerryyang@mba:channel$./channel_nobuf
goroutine message: 0
main function message

说明：没有buffer的channel, 写管道会被阻塞, 写管道和读管道的数量必须匹配, 否则会deadlock
*/
