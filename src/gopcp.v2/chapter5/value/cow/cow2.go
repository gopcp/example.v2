package cow

import (
	"errors"
	"fmt"
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

// intArray2 代表ConcurrentArray接口的实现类型。
type intArray2 struct {
	length  uint32
	val     atomic.Value
	version uint64
}

// NewConcurrentArray2 会创建一个ConcurrentArray类型值。
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
	newArray := make([]int, array.length)
	var v uint64
	for { // 乐观的自旋锁。
		v = atomic.LoadUint64(&array.version)
		copy(newArray, array.val.Load().([]int))
		old = newArray[index]
		newArray[index] = elem
		if atomic.CompareAndSwapUint64(&array.version, v, v+1) {
			// 在此处（上下两行代码之间）仍然可能发生中断，但是产生并发安全问题的概率不大。
			array.val.Store(newArray)
			break
		}
	}
	return
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
