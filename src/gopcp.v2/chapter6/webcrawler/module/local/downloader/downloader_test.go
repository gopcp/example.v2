package downloader

import (
	"bufio"
	"net/http"
	"testing"

	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/module/stub"
)

func TestNew(t *testing.T) {
	mid := module.MID("D1|127.0.0.1:8080")
	httpClient := &http.Client{}
	d, err := New(mid, httpClient, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a downloader: %s (mid: %s, httpClient: %#v)",
			err, mid, httpClient)
	}
	if d == nil {
		t.Fatal("Couldn't create downloader!")
	}
	if d.ID() != mid {
		t.Fatalf("Inconsistent MID for downloader: expected: %s, actual: %s",
			mid, d.ID())
	}
	mid = module.MID("D127.0.0.1")
	d, err = New(mid, httpClient, nil)
	if err == nil {
		t.Fatalf("No error when create a downloader with illegal MID %q!", mid)
	}
	mid = module.MID("D1|127.0.0.1:8888")
	httpClient = nil
	d, err = New(mid, httpClient, nil)
	if err == nil {
		t.Fatal("No error when create a downloader with nil http client!")
	}
}

func TestDownload(t *testing.T) {
	mid := module.MID("D1|127.0.0.1:8080")
	httpClient := &http.Client{}
	d, _ := New(mid, httpClient, nil)
	url := "http://www.baidu.com/robots.txt"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	depth := uint32(0)
	req := module.NewRequest(httpReq, depth)
	resp, err := d.Download(req)
	if err != nil {
		t.Fatalf("An error occurs when downloading content: %s (req: %#v)",
			err, req)
	}
	if resp == nil {
		t.Fatalf("Couldn't create download for request %#v!",
			req)
	}
	if resp.Depth() != depth {
		t.Fatalf("Inconsistent depth: expected: %d, actual: %d",
			depth, resp.Depth())
	}
	httpResp := resp.HTTPResp()
	if httpResp == nil {
		t.Fatalf("Invalid HTTP response! (url: %s)",
			url)
	}
	body := httpResp.Body
	if body == nil {
		t.Fatalf("Invalid HTTP response body! (url: %s)",
			url)
	}
	r := bufio.NewReader(body)
	line, _, err := r.ReadLine()
	if err != nil {
		t.Fatalf("An error occurs when reading HTTP response body: %s (url: %s)",
			err, url)
	}
	lineStr := string(line)
	expectedFirstLine := "User-agent: Baiduspider"
	if lineStr != expectedFirstLine {
		t.Fatalf("Inconsistent first line of the HTTP response body: expected: %s, actual: %s (url: %s)",
			expectedFirstLine, lineStr, url)
	}
	// 测试参数有误的情况。
	_, err = d.Download(nil)
	if err == nil {
		t.Fatal("No error when download with nil request!")
	}
	url = "http:///www.baidu.com/robots.txt"
	httpReq, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	req = module.NewRequest(httpReq, 0)
	resp, err = d.Download(req)
	if err == nil {
		t.Fatalf("No error when download with invalid url %q!", url)
	}
	req = module.NewRequest(nil, 0)
	resp, err = d.Download(req)
	if err == nil {
		t.Fatal("No error when download with nil HTTP request!")
	}

}

func TestCount(t *testing.T) {
	mid := module.MID("D1|127.0.0.1:8080")
	httpClient := &http.Client{}
	// 测试初始化后的计数。
	d, _ := New(mid, httpClient, nil)
	di := d.(stub.ModuleInternal)
	if di.CalledCount() != 0 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			0, di.CalledCount())
	}
	if di.AcceptedCount() != 0 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			0, di.AcceptedCount())
	}
	if di.CompletedCount() != 0 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			0, di.CompletedCount())
	}
	if di.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, di.HandlingNumber())
	}
	// 测试处理失败时的计数。
	d, _ = New(mid, httpClient, nil)
	di = d.(stub.ModuleInternal)
	url := "http:///www.baidu.com/robots.txt"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	req := module.NewRequest(httpReq, 0)
	_, err = d.Download(req)
	if di.CalledCount() != 1 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			1, di.CalledCount())
	}
	if di.AcceptedCount() != 1 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			1, di.AcceptedCount())
	}
	if di.CompletedCount() != 0 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			0, di.CompletedCount())
	}
	if di.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, di.HandlingNumber())
	}
	// 测试参数有误时的计数。
	d, _ = New(mid, httpClient, nil)
	di = d.(stub.ModuleInternal)
	_, err = d.Download(nil)
	if di.CalledCount() != 1 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			1, di.CalledCount())
	}
	if di.AcceptedCount() != 0 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			0, di.AcceptedCount())
	}
	if di.CompletedCount() != 0 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			0, di.CompletedCount())
	}
	if di.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, di.HandlingNumber())
	}
	// 测试处理成功完成时的计数。
	d, _ = New(mid, httpClient, nil)
	di = d.(stub.ModuleInternal)
	url = "http://www.baidu.com/robots.txt"
	httpReq, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	req = module.NewRequest(httpReq, 0)
	_, err = d.Download(req)
	if di.CalledCount() != 1 {
		t.Fatalf("Inconsistent called count for internal module: expected: %d, actual: %d",
			1, di.CalledCount())
	}
	if di.AcceptedCount() != 1 {
		t.Fatalf("Inconsistent accepted count for internal module: expected: %d, actual: %d",
			1, di.AcceptedCount())
	}
	if di.CompletedCount() != 1 {
		t.Fatalf("Inconsistent completed count for internal module: expected: %d, actual: %d",
			1, di.CompletedCount())
	}
	if di.HandlingNumber() != 0 {
		t.Fatalf("Inconsistent handling number for internal module: expected: %d, actual: %d",
			0, di.HandlingNumber())
	}
}
