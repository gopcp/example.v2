package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

func main() {
	// 禁用GC，并保证在main函数执行结束前恢复GC。
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var count int32
	newFunc := func() interface{} {
		return atomic.AddInt32(&count, 1)
	}
	pool := sync.Pool{New: newFunc}

	// New 字段值的作用。
	v1 := pool.Get()
	fmt.Printf("Value 1: %v\n", v1)

	// 临时对象池的存取。
	pool.Put(10)
	pool.Put(11)
	pool.Put(12)
	v2 := pool.Get()
	fmt.Printf("Value 2: %v\n", v2)

	// 垃圾回收对临时对象池的影响。
	debug.SetGCPercent(100)
	runtime.GC()
	v3 := pool.Get()
	fmt.Printf("Value 3: %v\n", v3)
	pool.New = nil
	v4 := pool.Get()
	fmt.Printf("Value 4: %v\n", v4)
}
