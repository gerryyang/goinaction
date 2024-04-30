
测试输出：

```
./main
==================
WARNING: DATA RACE
Read at 0x000001216848 by goroutine 9:
  main.main.func1()
      /Users/gerry/Proj/github/goinaction/src/race_check/main.go:16 +0x32

Previous write at 0x000001216848 by goroutine 6:
  main.main.func1()
      /Users/gerry/Proj/github/goinaction/src/race_check/main.go:16 +0x4a

Goroutine 9 (running) created at:
  main.main()
      /Users/gerry/Proj/github/goinaction/src/race_check/main.go:15 +0x64

Goroutine 6 (running) created at:
  main.main()
      /Users/gerry/Proj/github/goinaction/src/race_check/main.go:15 +0x64
==================
Counter: 990
Found 1 data race(s)
```