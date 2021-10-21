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
	// 在新版本（起码 1.15 及以后）的 Go 当中，sync.Pool 里的临时对象需要两次 GC 才会被真正清除掉。
	// 只一次 GC 的话只会让其中的临时对象被“打上记号”。
	// 更具体的说，只做一次 GC 只会使得 sync.Pool 里的临时对象被移动到池中的“备用区域”（详见 victim 字段）。
	// 在我们调用 sync.Pool 的 Get 方法时，如果 sync.Pool 的“本地区域”（详见 local 字段）当中没有可用的临时对象，
	// 那么 sync.Pool 会试图从这个“备用区域”中获取临时对象。
	// 如果“备用区域”也没有可用的临时对象，sync.Pool 才会去调用 New 函数。
	// 所以，这里的例子需要再添加一行对 runtime.GC() 函数的调用，才能使它的结果与最初的相同，并起到准确示范的作用。
	runtime.GC()

	v3 := pool.Get()
	fmt.Printf("Value 3: %v\n", v3)
	pool.New = nil
	v4 := pool.Get()
	fmt.Printf("Value 4: %v\n", v4)
}
