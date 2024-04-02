package main

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	// 创建一个上下文对象，设置超时时间为 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建一个命令对象
	cmd := exec.CommandContext(ctx, "ls", "-l")

	// 执行命令，并获取输出
	output, err := cmd.Output()

	// 判断命令是否执行成功
	if err != nil {
		fmt.Println("Command failed:", err)

		// 如果命令执行失败，则杀死命令的进程
		if cmd.Process != nil {
			cmd.Process.Kill()
		}

		return
	}

	// 在 defer 语句中释放进程状态相关的资源
	defer func() {
		if cmd.ProcessState != nil {
			// 获取命令的系统进程状态
			sys := cmd.ProcessState.Sys()

			// 根据不同的操作系统类型，调用不同的方法释放进程状态相关的资源
			switch sys := sys.(type) {
			case *syscall.WaitStatus:
				// Linux 系统上使用 syscall.Wait4 函数释放进程状态相关的资源
				syscall.Wait4(cmd.Process.Pid, sys, 0, nil)
			}
		}
	}()

	// 打印命令的输出
	fmt.Println(string(output))
}
