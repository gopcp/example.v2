package module

import (
	"net/http"
	"strings"
	"testing"
)

// testingReader 代表测试专用的读取器，实现了io.ReadCloser接口类型。
type testingReader struct {
	sr *strings.Reader
}

func (r testingReader) Read(b []byte) (n int, err error) {
	return r.sr.Read(b)
}

func (r testingReader) Close() error {
	return nil
}

func TestRequest(t *testing.T) {
	method := "GET"
	expectedURLStr := "https://github.com/gopcp"
	expectedHTTPReq, _ := http.NewRequest(method, expectedURLStr, nil)
	expectedDepth := uint32(0)
	req := NewRequest(expectedHTTPReq, expectedDepth)
	if req == nil {
		t.Fatal("Couldn't create request!")
	}
	if _, ok := interface{}(req).(Data); !ok {
		t.Fatalf("Request didn't implement Data!")
	}
	expectedValidity := true
	valid := req.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for request: expected: %v, actual: %v",
			expectedValidity, valid)
	}
	if req.HTTPReq() != expectedHTTPReq {
		t.Fatalf("Inconsistent HTTP request for request: expected: %#v, actual: %#v",
			expectedHTTPReq, req.HTTPReq())
	}
	if req.Depth() != expectedDepth {
		t.Fatalf("Inconsistent depth for request: expected: %d, actual: %d",
			expectedDepth, req.Depth())
	}
	expectedHTTPReq.URL = nil
	req = NewRequest(expectedHTTPReq, expectedDepth)
	expectedValidity = false
	valid = req.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for request: expected: %v, actual: %v",
			expectedValidity, valid)
	}
	req = NewRequest(nil, expectedDepth)
	expectedValidity = false
	valid = req.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for request: expected: %v, actual: %v",
			expectedValidity, valid)
	}
}

func TestResponse(t *testing.T) {
	method := "GET"
	expectedURLStr := "https://github.com/gopcp"
	httpReq, _ := http.NewRequest(method, expectedURLStr, nil)
	expectHTTPResp := &http.Response{
		Request: httpReq,
		Body:    testingReader{strings.NewReader("Test response")},
	}
	expectedDepth := uint32(0)
	resp := NewResponse(expectHTTPResp, expectedDepth)
	if resp == nil {
		t.Fatal("Couldn't create response!")
	}
	if _, ok := interface{}(resp).(Data); !ok {
		t.Fatalf("Response didn't implement Data!")
	}
	expectedValidity := true
	valid := resp.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for response: expected: %v, actual: %v",
			expectedValidity, valid)
	}
	if resp.HTTPResp() != expectHTTPResp {
		t.Fatalf("Inconsistent HTTP response for response: expected: %#v, actual: %#v",
			expectHTTPResp, resp.HTTPResp())
	}
	if resp.Depth() != expectedDepth {
		t.Fatalf("Inconsistent depth for response: expected: %d, actual: %d",
			expectedDepth, resp.Depth())
	}
	expectHTTPResp.Body = nil
	resp = NewResponse(expectHTTPResp, expectedDepth)
	expectedValidity = false
	valid = resp.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for response: expected: %v, actual: %v",
			expectedValidity, valid)
	}
	resp = NewResponse(nil, expectedDepth)
	expectedValidity = false
	valid = resp.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for response: expected: %v, actual: %v",
			expectedValidity, valid)
	}
}

func TestItem(t *testing.T) {
	item := Item(map[string]interface{}{})
	if _, ok := interface{}(item).(Data); !ok {
		t.Fatalf("Item didn't implement Data!")
	}
	expectedValidity := true
	valid := item.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for item: expected: %v, actual: %v",
			expectedValidity, valid)
	}
	item = Item(nil)
	expectedValidity = false
	valid = item.Valid()
	if valid != expectedValidity {
		t.Fatalf("Inconsistent validity for item: expected: %v, actual: %v",
			expectedValidity, valid)
	}
}
