package buffer

import (
	"fmt"
	"sync"
	"sync/atomic"

	"gopcp.v2/chapter6/webcrawler/errors"
)

// Buffer 代表FIFO的缓冲器的接口类型。
type Buffer interface {
	// Cap 用于获取本缓冲器的容量。
	Cap() uint32
	// Len 用于获取本缓冲器中的数据数量。
	Len() uint32
	// Put 用于向缓冲器放入数据。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Put(datum interface{}) (bool, error)
	// Get 用于从缓冲器获取器。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Get() (interface{}, error)
	// Close 用于关闭缓冲器。
	// 若缓冲器之前已关闭则返回false，否则返回true。
	Close() bool
	// Closed 用于判断缓冲器是否已关闭。
	Closed() bool
}

// myBuffer 代表缓冲器接口的实现类型。
type myBuffer struct {
	// ch 代表存放数据的通道。
	ch chan interface{}
	// closed 代表缓冲器的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// closingLock 代表为了消除因关闭缓冲器而产生的竞态条件的读写锁。
	closingLock sync.RWMutex
}

// NewBuffer 用于创建一个缓冲器。
// 参数size代表缓冲器的容量。
func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d", size)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &myBuffer{
		ch: make(chan interface{}, size),
	}, nil
}

func (buf *myBuffer) Cap() uint32 {
	return uint32(cap(buf.ch))
}

func (buf *myBuffer) Len() uint32 {
	return uint32(len(buf.ch))
}

func (buf *myBuffer) Put(datum interface{}) (ok bool, err error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()
	if buf.Closed() {
		return false, ErrClosedBuffer
	}
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return
}

func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.closingLock.Lock()
		close(buf.ch)
		buf.closingLock.Unlock()
		return true
	}
	return false
}

func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}
