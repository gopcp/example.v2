package cow

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

// ConcurrentIntArray 代表并发安全的整数数组接口。
type ConcurrentIntArray interface {
	// Set 用于设置指定索引上的元素值。
	Set(index int, elem int) (old int, err error)
	// Get 用于获取指定索引上的元素值。
	Get(index int) (elem int, err error)
	// Len 用于获取数组的长度。
	Len() int
}

// segment 是一个代表内部段的类型。
type segment struct {
	val    atomic.Value
	length int    // 段内的元素数量。
	status uint32 // 0：可读可写；1：只读。
}

func (seg *segment) init(length int) {
	seg.length = length
	seg.val.Store(make([]int, length))
}

func (seg *segment) checkIndex(index int) error {
	if index < 0 || index >= seg.length {
		return fmt.Errorf("index out of range [0, %d) in segment", seg.length)
	}
	return nil
}

func (seg *segment) set(index int, elem int) (old int, err error) {
	if err = seg.checkIndex(index); err != nil {
		return
	}
	point := 10 //TODO 此处是一个优化点，可以根据实际情况调整。
	count := 0
	for { // 简易的自旋锁。
		count++
		if !atomic.CompareAndSwapUint32(&seg.status, 0, 1) {
			if count%point == 0 {
				runtime.Gosched()
			}
			continue
		}
		defer atomic.StoreUint32(&seg.status, 0)
		newArray := make([]int, seg.length)
		copy(newArray, seg.val.Load().([]int))
		old = newArray[index]
		newArray[index] = elem
		seg.val.Store(newArray)
		return
	}
}

func (seg *segment) get(index int) (elem int, err error) {
	if err = seg.checkIndex(index); err != nil {
		return
	}
	elem = seg.val.Load().([]int)[index]
	return
}

// myIntArray 代表 ConcurrentIntArray 接口的实现类型。
type myIntArray struct {
	length    int        // 元素总数量。
	segLenStd int        // 单个内部段的标准长度。
	segments  []*segment // 内部段列表。
}

// NewConcurrentIntArray 会创建一个 ConcurrentIntArray 类型值。
func NewConcurrentIntArray(length int) ConcurrentIntArray {
	if length < 0 {
		length = 0
	}
	array := new(myIntArray)
	array.init(length)
	return array
}

func (array *myIntArray) init(length int) {
	array.length = length
	array.segLenStd = 10 //TODO 此处是一个优化点，可以根据参数值调整。
	segNum := length / array.segLenStd
	segLenTail := length % array.segLenStd
	if segLenTail > 0 {
		segNum = segNum + 1
	}
	array.segments = make([]*segment, segNum)
	for i := 0; i < segNum; i++ {
		seg := segment{}
		if i == segNum-1 && segLenTail > 0 {
			seg.init(segLenTail)
		} else {
			seg.init(array.segLenStd)
		}
		array.segments[i] = &seg
	}
}

func (array *myIntArray) Set(index int, elem int) (old int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	seg := array.segments[index/array.segLenStd]
	return seg.set(index%array.segLenStd, elem)
}

func (array *myIntArray) Get(index int) (elem int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	seg := array.segments[index/array.segLenStd]
	return seg.get(index % array.segLenStd)
}

func (array *myIntArray) Len() int {
	return array.length
}

// checkIndex 用于检查索引的有效性。
func (array *myIntArray) checkIndex(index int) error {
	if index < 0 || index >= array.length {
		return fmt.Errorf("index out of range [0, %d)", array.length)
	}
	return nil
}
