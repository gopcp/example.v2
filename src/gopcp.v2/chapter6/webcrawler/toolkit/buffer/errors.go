package buffer

import "errors"

// ErrClosedBufferPool 是表示缓冲池已关闭的错误的变量。
var ErrClosedBufferPool = errors.New("closed buffer pool")

// ErrClosedBuffer 是表示缓冲器已关闭的错误的变量。
var ErrClosedBuffer = errors.New("closed buffer")
