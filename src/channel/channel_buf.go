package main

import (
	"fmt"
)

func main() {
	// 创建一个channel用以同步goroutine
	done := make(chan bool, 2)
	count := 3

	for i := 0; i < count; i++ {

		// 在goroutine中执行输出操作
		go func(i int) {
			done <- true
			fmt.Println("goroutine message:", i)

		}(i)
	}

	<-done // 等待goroutine结束
	fmt.Println("main function message")
}

/*
output:
gerryyang@mba:channel$./channel_buf
goroutine message: 0
goroutine message: 1
main function message

说明：有buffer的channel, 假设管道buffer的大小为n, 则可以连续写管道n次（非阻塞）, 否则会deadlock, 读管道同理
*/
