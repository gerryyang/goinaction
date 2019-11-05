
# Golang in Action
---
[TOC]


## API Doc

http://devdocs.io/go/

https://golang.org/pkg/

## IDE

Windows: LiteIDE, [zeal](https://zealdocs.org/)

Mac: Dash (register codes can be from taobao)

## Performance

1. pprof

https://golang.org/pkg/net/http/pprof/

```
go tool pprof http://localhost:6060/debug/pprof/profile
```

2. go-torch

https://github.com/uber/go-torch

```
go-torch -u http://localhost:32775 -t 30
```

3. graphviz

http://www.graphviz.org/


## Some Projects

1. [TarsGo](https://github.com/TarsCloud/TarsGo/tree/36cf7d196afc781ba8d70f908decbdd7cceccfd3/tars)

2. 一个开源的高性能golang日志库 

[Logrus](https://github.com/Sirupsen/logrus)


## Refer

[CoolShell - GO 语言简介（上）— 语法](http://coolshell.cn/articles/8460.html)

[CoolShell - GO 语言简介（下）— 特性](http://coolshell.cn/articles/8489.html)

[10 things you (probably) don't know about Go](https://talks.golang.org/2012/10things.slide#1)

[50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)

[Golang 新手可能会踩的 50 个坑](https://segmentfault.com/a/1190000013739000)

[Some go examples](https://yourbasic.org/golang/)

[dave-high-performance-go-workshop-gopherchina-2019](https://dave.cheney.net/high-performance-go-workshop/gopherchina-2019.html)

[使用二进制形式发布go package](https://colobu.com/2018/01/10/use-binary-package-in-go/)

## Stuff

[Fast HTTP (L7) and TCP (L4) tunnels written in Go](https://github.com/inlets/inletsctl)

[Go: Memory Management and Allocation](https://medium.com/a-journey-with-go/go-memory-management-and-allocation-a7396d430f44)

[Go: How Does the Garbage Collector Mark the Memory?](https://medium.com/a-journey-with-go/go-how-does-the-garbage-collector-mark-the-memory-72cfc12c6976)

[TCMalloc : Thread-Caching Malloc](http://goog-perftools.sourceforge.net/doc/tcmalloc.html)