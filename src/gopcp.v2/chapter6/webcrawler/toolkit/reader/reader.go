package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// MultipleReader 代表多重读取器的接口。
type MultipleReader interface {
	// Reader 用于获取一个可关闭读取器的实例。
	// 后者会持有本多重读取器中的数据。
	Reader() io.ReadCloser
}

// myMultipleReader 代表多重读取器的实现类型。
type myMultipleReader struct {
	data []byte
}

// NewMultipleReader 用于新建并返回一个多重读取器的实例。
func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var data []byte
	var err error
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader: couldn't create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMultipleReader{
		data: data,
	}, nil
}

func (rr *myMultipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}
