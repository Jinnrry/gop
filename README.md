# Go Pretty Print Value

> 本项目核心代码来自于[https://github.com/ysmood/got](https://github.com/ysmood/got)

更优雅的打印出Go结构体

# Install

`go get github.com/Jinnrry/gop`

# Demo

```go

package main

import "github.com/Jinnrry/gop"

func main() {
	v:=[]interface{}{
		map[int]int{1:1},
		"hello",
		&struct {
            k string
			v int
		}{"key",1},
}
	gop.Print(v)
}

```
输出结果
```
[]interface {}/* len=3 cap=3 */{
    map[int]int{
        1: 1,
    },
    "hello",
    &struct { k string; v int }/* len=2 */{
        k: "key",
        v: 1,
    },
}
```
