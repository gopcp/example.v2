package cow

import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
)

// ConcurrentArray2 代表并发安全的整数数组接口。
type ConcurrentArray2 interface {
	// Set 用于设置指定索引上的元素值。
	Set(index uint32, elem int) (old int, err error)
	// Get 用于获取指定索引上的元素值。
	Get(index uint32) (elem int, err error)
	// Len 用于获取数组的长度。
	Len() uint32
}

// intArray2 代表ConcurrentArray2接口的实现类型。
type intArray2 struct {
	length uint32
	val    atomic.Value
	status uint32 // 0：可读可写；1：只读。
}

// NewConcurrentArray2 会创建一个ConcurrentArray2类型值。
func NewConcurrentArray2(length uint32) ConcurrentArray2 {
	array := intArray2{}
	array.length = length
	array.val.Store(make([]int, array.length))
	return &array
}

func (array *intArray2) Set(index uint32, elem int) (old int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}
	for { // 一个简易的自旋锁。
		if !atomic.CompareAndSwapUint32(&array.status, 0, 1) {
			runtime.Gosched()
			continue
		}
		defer atomic.StoreUint32(&array.status, 0)
		newArray := make([]int, array.length)
		copy(newArray, array.val.Load().([]int))
		old = newArray[index]
		newArray[index] = elem
		array.val.Store(newArray)
		return
	}
}

func (array *intArray2) Get(index uint32) (elem int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}
	elem = array.val.Load().([]int)[index]
	return
}

func (array *intArray2) Len() uint32 {
	return array.length
}

// checkIndex 用于检查索引的有效性。
func (array *intArray2) checkIndex(index uint32) error {
	if index >= array.length {
		return fmt.Errorf("index out of range [0, %d)", array.length)
	}
	return nil
}

// checkValue 用于检查原子值中是否已存有值。
func (array *intArray2) checkValue() error {
	v := array.val.Load()
	if v == nil {
		return errors.New("invalid int array")
	}
	return nil
}
