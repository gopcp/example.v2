package analyzer

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/module/stub"
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

func TestNew(t *testing.T) {
	mid := module.MID("D1|127.0.0.1:8080")
	parsers := []module.ParseResponse{genTestingRespParser(false)}
	a, err := New(mid, parsers, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating an analyzer: %s (mid: %s)",
			err, mid)
	}
	if a == nil {
		t.Fatal("Couldn't create analyzer!")
	}
	if a.ID() != mid {
		t.Fatalf("Inconsistent MID for analyzer: expected: %s, actual: %s",
			mid, a.ID())
	}
	if len(a.RespParsers()) != len(parsers) {
		t.Fatalf("Inconsistent response parser number for pipeline: expected: %d, actual: %d",
			len(a.RespParsers()), len(parsers))
	}
	// 测试参数有误的情况。
	mid = module.MID("D127.0.0.1")
	a, err = New(mid, parsers, nil)
	if err == nil {
		t.Fatalf("No error when create an analyzer with illegal MID %q!", mid)
	}
	mid = module.MID("D1|127.0.0.1:8080")
	parsersList := [][]module.ParseResponse{
		nil,
		[]module.ParseResponse{},
		[]module.ParseResponse{genTestingRespParser(false), nil},
	}
	for _, parsers := range parsersList {
		a, err = New(mid, parsers, nil)
		if err == nil {
			t.Fatalf("No error when create an analyzer with illegal parsers %#v!",
				parsers)
		}
	}
}

func TestAnalyze(t *testing.T) {
	number := uint32(10)
	method := "GET"
	expectedURL := "https://github.com/gopcp"
	expectedDepth := uint32(1)
	resps := getTestingResps(number, method, expectedURL, expectedDepth, t)
	mid := module.MID("D1|127.0.0.1:8080")
	parsers := []module.ParseResponse{genTestingRespParser(false)}
	a, err := New(mid, parsers, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating an analyzer: %s (mid: %s)",
			err, mid)
	}
	data := []module.Data{}
	parseErrors := []error{}
	for _, resp := range resps {
		data1, parseErrors1 := a.Analyze(resp)
		data = append(data, data1...)
		parseErrors = append(parseErrors, parseErrors1...)
	}
	for i, e := range parseErrors {
		t.Errorf("An error occurs when parsing response: %s (index: %d)",
			e, i)
	}
	var count int
	for i, d := range data {
		if d == nil {
			t.Fatalf("nil datum! (index: %d)", i)
		}
		if _, ok := d.(*module.Request); ok {
			continue
		}
		item, ok := d.(module.Item)
		if !ok {
			t.Errorf("Inconsistent datum type: expected: %T, actual: %T (index: %d)",
				module.Item{}, d, i)
		}
		if item["url"] != expectedURL {
			t.Errorf("Inconsistent URL: expected: %s, actual: %s (index: %d)",
				expectedURL, item["url"], i)
		}
		index, ok := item["index"].(int)
		if !ok {
			t.Errorf("Inconsistent index type: expected: %T, actual: %T (index: %d)",
				int(0), item["index"], i)
		}
		if index != count {
			t.Errorf("Inconsistent index: expected: %d, actual: %d (index: %d)",
				count, index, i)
		}
		depth, ok := item["depth"].(uint32)
		if !ok {
			t.Errorf("Inconsistent depth type: expected: %T, actual: %T (index: %d)",
				uint32(0), item["depth"], i)
		}
		if depth != expectedDepth {
			t.Errorf("Inconsistent depth: expected: %d, actual: %d (index: %d)",
				expectedDepth, depth, i)
		}
		count++
	}
	// 测试参数有误的情况。
	// 测试响应为nil的情况。
	_, errs := a.Analyze(nil)
	if len(errs) == 0 {
		t.Fatal("No error when download with nil response!")
	}
	// 测试HTTP响应为nil的情况。
	resp := module.NewResponse(nil, 0)
	_, errs = a.Analyze(resp)
	if len(errs) == 0 {
		t.Fatalf("No error when analyze response with illegal response %#v!",
			parsers)
	}
	// 测试HTTP请求为nil的情况。
	httpResp := &http.Response{
		Request: nil,
		Body:    nil,
	}
	resp = module.NewResponse(httpResp, 0)
	_, errs = a.Analyze(resp)
	if len(errs) == 0 {
		t.Fatalf("No error when analyze response with nil request URL!")
	}
	// 测试HTTP请求URL为nil的情况。
	httpReq, _ := http.NewRequest(method, expectedURL, nil)
	httpReq.URL = nil
	httpResp = &http.Response{
		Request: httpReq,
		Body:    nil,
	}
	resp = module.NewResponse(httpResp, 0)
	_, errs = a.Analyze(resp)
	if len(errs) == 0 {
		t.Fatalf("No error when analyze response with nil request URL!")
	}
}

func TestCount(t *testing.T) {
	mid := module.MID("D1|127.0.0.1:8080")
	// 测试初始化后的计数。
	parsers := []module.ParseResponse{genTestingRespParser(false)}
	a, _ := New(mid, parsers, nil)
	ai := a.(stub.ModuleInternal)
	if ai.CalledCount() != 0 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			0, ai.CalledCount())
	}
	if ai.AcceptedCount() != 0 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			0, ai.AcceptedCount())
	}
	if ai.CompletedCount() != 0 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			0, ai.CompletedCount())
	}
	if ai.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, ai.HandlingNumber())
	}
	// 测试处理失败时的计数。
	parsers = []module.ParseResponse{genTestingRespParser(true)}
	a, _ = New(mid, parsers, nil)
	ai = a.(stub.ModuleInternal)
	resp := getTestingResps(1, "GET", "https://github.com/gopcp", 0, t)[0]
	a.Analyze(resp)
	if ai.CalledCount() != 1 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			1, ai.CalledCount())
	}
	if ai.AcceptedCount() != 1 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			1, ai.AcceptedCount())
	}
	if ai.CompletedCount() != 0 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			0, ai.CompletedCount())
	}
	if ai.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, ai.HandlingNumber())
	}
	// 测试参数有误时的计数。
	parsers = []module.ParseResponse{genTestingRespParser(false)}
	a, _ = New(mid, parsers, nil)
	ai = a.(stub.ModuleInternal)
	resp = module.NewResponse(nil, 0)
	a.Analyze(resp)
	if ai.CalledCount() != 1 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			1, ai.CalledCount())
	}
	if ai.AcceptedCount() != 0 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			0, ai.AcceptedCount())
	}
	if ai.CompletedCount() != 0 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			0, ai.CompletedCount())
	}
	if ai.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, ai.HandlingNumber())
	}
	// 测试处理成功完成时的计数。
	parsers = []module.ParseResponse{genTestingRespParser(false)}
	a, _ = New(mid, parsers, nil)
	ai = a.(stub.ModuleInternal)
	resp = getTestingResps(1, "GET", "https://github.com/gopcp", 0, t)[0]
	a.Analyze(resp)
	if ai.CalledCount() != 1 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			1, ai.CalledCount())
	}
	if ai.AcceptedCount() != 1 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			1, ai.AcceptedCount())
	}
	if ai.CompletedCount() != 1 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			1, ai.CompletedCount())
	}
	if ai.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, ai.HandlingNumber())
	}
}

// fakeHTTPRespBody 代表伪造的HTTP响应体的模板。
var fakeHTTPRespBody = "Fake HTTP Response [%d]"

// testingRespParser 为测试专用的响应解析函数。
// 生成的函数会把响应的请求URL、响应体中的索引和响应深度存在条目中。
func genTestingRespParser(fail bool) module.ParseResponse {
	if fail {
		return func(httpResp *http.Response,
			respDepth uint32) (data []module.Data, parseErrors []error) {
			errs :=
				[]error{fmt.Errorf("Fail! (httpResp: %#v, respDepth: %#v)", httpResp, respDepth)}
			return nil, errs
		}
	}
	return func(httpResp *http.Response,
		respDepth uint32) (data []module.Data, parseErrors []error) {
		data = []module.Data{}
		parseErrors = []error{}
		item := module.Item(map[string]interface{}{})
		item["url"] = httpResp.Request.URL.String()
		bufReader := bufio.NewReader(httpResp.Body)
		line, _, err := bufReader.ReadLine()
		if err != nil {
			parseErrors = append(parseErrors, err)
			return
		}
		lineStr := string(line)
		begin := strings.LastIndex(lineStr, "[")
		end := strings.LastIndex(lineStr, "]")
		if begin < 0 ||
			end < 0 ||
			begin > end {
			err := fmt.Errorf("wrong index for index: %d, %d",
				begin, end)
			parseErrors = append(parseErrors, err)
			return
		}
		index, err := strconv.Atoi(lineStr[begin+1 : end])
		if err != nil {
			parseErrors = append(parseErrors, err)
			return
		}
		item["index"] = index
		item["depth"] = respDepth
		data = append(data, item)
		req := module.NewRequest(nil, respDepth)
		data = append(data, req)
		return
	}
}

// getTestingResps 用于生成测试专用的响应实例。
// 该函数返回的响应实例都是伪造的，只提供了测试必用的内容。
func getTestingResps(
	number uint32,
	method string,
	url string,
	depth uint32,
	t *testing.T) []*module.Response {
	httpReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (method: %s, url: %s)",
			err, method, url)
	}
	resps := []*module.Response{}
	for i := uint32(0); i < number; i++ {
		httpResp := &http.Response{
			Request: httpReq,
			Body: testingReader{
				strings.NewReader(
					fmt.Sprintf(fakeHTTPRespBody, i))},
		}
		resp := module.NewResponse(httpResp, depth)
		resps = append(resps, resp)
	}
	return resps
}
