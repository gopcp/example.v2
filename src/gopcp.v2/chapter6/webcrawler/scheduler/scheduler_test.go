package scheduler

import (
	"net/http"
	"runtime"
	"testing"
	"time"

	"gopcp.v2/chapter5/cmap"
	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/toolkit/buffer"
)

// snGen 代表序列号生成器。
var snGen = module.NewSNGenertor(1, 0)

func TestSchedNew(t *testing.T) {
	sched := NewScheduler()
	if sched == nil {
		t.Fatal("Couldn't create scheduler!")
	}
}

func TestSchedInit(t *testing.T) {
	requestArgs := genRequestArgs([]string{"bing.com"}, 0)
	dataArgs := genDataArgs(10, 2, 1)
	moduleArgs := genSimpleModuleArgs(3, 2, 1, t)
	sched := NewScheduler()
	err := sched.Init(
		requestArgs,
		dataArgs,
		moduleArgs)
	if err != nil {
		t.Fatalf("An error occurs when initializing scheduler: %s", err)
	}
	err = sched.Init(
		requestArgs,
		dataArgs,
		moduleArgs)
	if err != nil {
		t.Fatalf("An error occurs when initializing scheduler: %s", err)
	}
	// 测试请求参数异常时的情况。
	invalidRequestArgs := genRequestArgs(nil, 0)
	err = sched.Init(
		invalidRequestArgs,
		dataArgs,
		moduleArgs)
	if err == nil {
		t.Fatalf("No error when initialize scheduler with illegal request arguments %s!",
			invalidRequestArgs)
	}
	// 测试数据参数异常时的情况。
	sched = NewScheduler()
	invalidDataArgs := genDataArgs(0, 0, 0)
	err = sched.Init(
		requestArgs,
		invalidDataArgs,
		moduleArgs)
	if err == nil {
		t.Fatalf("No error when initialize scheduler with illegal data arguments %s!",
			invalidDataArgs)
	}
	// 测试组件参数异常时的情况。
	invalidModuleArgsList := []ModuleArgs{
		ModuleArgs{},
		genSimpleModuleArgs(-1, 1, 1, t),
		genSimpleModuleArgs(1, -1, 1, t),
		genSimpleModuleArgs(1, 1, -1, t),
		ModuleArgs{
			Downloaders: genSimpleDownloaders(2, true, snGen, t),
			Analyzers:   genSimpleAnalyzers(2, false, snGen, t),
			Pipelines:   genSimplePipelines(2, false, snGen, t),
		},
		ModuleArgs{
			Downloaders: genSimpleDownloaders(2, false, snGen, t),
			Analyzers:   genSimpleAnalyzers(2, true, snGen, t),
			Pipelines:   genSimplePipelines(2, false, snGen, t),
		},
		ModuleArgs{
			Downloaders: genSimpleDownloaders(2, false, snGen, t),
			Analyzers:   genSimpleAnalyzers(2, false, snGen, t),
			Pipelines:   genSimplePipelines(2, true, snGen, t),
		},
	}
	for _, invalidModuleArgs := range invalidModuleArgsList {
		sched = NewScheduler()
		dataArgs := genDataArgs(10, 2, 1)
		err = sched.Init(
			requestArgs,
			dataArgs,
			invalidModuleArgs)
		if err == nil {
			t.Fatalf("No error when initialize scheduler with illegal module arguments %s!",
				invalidModuleArgs)
		}
	}
	invalidModuleArgsList = []ModuleArgs{
		genSimpleModuleArgs(-2, 1, 1, t),
		genSimpleModuleArgs(1, -2, 1, t),
		genSimpleModuleArgs(1, 1, -2, t),
	}
	for _, invalidModuleArgs := range invalidModuleArgsList {
		sched = NewScheduler()
		err = sched.Init(
			requestArgs,
			dataArgs,
			invalidModuleArgs)
		if err != nil {
			t.Fatalf("An error occurs when initializing scheduler: %s",
				err)
		}
	}
}

func TestSchedStart(t *testing.T) {
	sched := NewScheduler()
	requestArgs := genRequestArgs([]string{}, 0)
	dataArgs := genDataArgs(10, 2, 1)
	moduleArgs := genSimpleModuleArgs(3, 2, 1, t)
	err := sched.Init(
		requestArgs,
		dataArgs,
		moduleArgs)
	if err != nil {
		t.Fatalf("An error occurs when initializing scheduler: %s", err)
	}
	url := "http://cn.bing.com/search?q=golang"
	firstHTTPReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	err = sched.Start(firstHTTPReq)
	if err != nil {
		t.Fatalf("An error occurs when starting scheduler: %s",
			err)
	}
	// 测试首个请求异常时的情况。
	sched.Stop()
	err = sched.Start(nil)
	if err == nil {
		t.Fatalf("No error when start scheduler with nil HTTP request!")
	}
	sched.Stop()
	firstHTTPReq.Host = ""
	err = sched.Start(firstHTTPReq)
	if err == nil {
		t.Fatalf("No error when start scheduler with empty HTTP host!")
	}
	sched.Stop()
}

func TestSchedStop(t *testing.T) {
	sched := NewScheduler()
	requestArgs := genRequestArgs([]string{}, 0)
	dataArgs := genDataArgs(10, 2, 1)
	moduleArgs := genSimpleModuleArgs(3, 2, 1, t)
	err := sched.Init(
		requestArgs,
		dataArgs,
		moduleArgs)
	if err != nil {
		t.Fatalf("An error occurs when initializing scheduler: %s", err)
	}
	url := "http://cn.bing.com/search?q=golang"
	firstHTTPReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	err = sched.Start(firstHTTPReq)
	if err != nil {
		t.Fatalf("An error occurs when starting scheduler: %s", err)
	}
	if err = sched.Stop(); err != nil {
		t.Fatalf("An error occurs when stopping scheduler: %s", err)
	}
}

func TestSchedStatus(t *testing.T) {
	// 准备初始化参数。
	requestArgs := genRequestArgs([]string{"bing.com"}, 0)
	dataArgs := genDataArgs(10, 2, 1)
	moduleArgs := genSimpleModuleArgs(3, 2, 1, t)
	sched := NewScheduler()
	// 准备启动参数。
	url := "http://cn.bing.com/search?q=golang"
	firstHTTPReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	// 测试未初始化状态下的启动。
	if err = sched.Start(firstHTTPReq); err == nil {
		t.Fatal("No error when start scheduler before initialize!")
	}
	// 测试未初始化状态下的停止。
	if err = sched.Stop(); err == nil {
		t.Fatal("No error when stop scheduler before initialize!")
	}
	// 测试未初始化状态下的初始化。
	if err := sched.Init(requestArgs, dataArgs, moduleArgs); err != nil {
		t.Fatalf("An error occurs when initializing scheduler: %s", err)
	}
	// 测试重复初始化。
	if err := sched.Init(requestArgs, dataArgs, moduleArgs); err != nil {
		t.Fatalf("An error occurs when repeatedly initializing scheduler: %s", err)
	}
	// 测试已初始化状态下的停止。
	if err = sched.Stop(); err == nil {
		t.Fatal("No error when stop scheduler after initialize!")
	}
	// 测试已初始化状态下的启动。
	if err = sched.Start(firstHTTPReq); err != nil {
		t.Fatalf("An error occurs when starting scheduler after initialize: %s", err)
	}
	// 测试重复启动。
	if err = sched.Start(firstHTTPReq); err == nil {
		t.Fatal("No error when repeatedly start scheduler!")
	}
	// 测试已启动状态下的初始化。
	if err := sched.Init(requestArgs, dataArgs, moduleArgs); err == nil {
		t.Fatal("No error when initialize scheduler after start!")
	}
	// 测试已启动状态下的停止。
	if err = sched.Stop(); err != nil {
		t.Fatalf("An error occurs when stopping scheduler after start: %s", err)
	}
	// 测试重复停止。
	if err = sched.Stop(); err == nil {
		t.Fatal("No error when repeatedly stop scheduler!")
	}
	// 测试已停止状态下的初始化。
	if err = sched.Init(requestArgs, dataArgs, moduleArgs); err != nil {
		t.Fatalf("An error occurs when initializing scheduler after stop: %s", err)
	}
}

func TestSchedSimple(t *testing.T) {
	requestArgs := genRequestArgs([]string{}, 0)
	dataArgs := genDataArgs(10, 2, 1)
	moduleArgs := genSimpleModuleArgs(3, 2, 1, t)
	url := "http://cn.bing.com/search?q=golang"
	firstHTTPReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	sched := NewScheduler()
	err = sched.Init(
		requestArgs,
		dataArgs,
		moduleArgs)
	if err != nil {
		t.Fatalf("An error occurs when initializing scheduler: %s",
			err)
	}
	err = sched.Start(firstHTTPReq)
	if err != nil {
		t.Fatalf("An error occurs when starting scheduler: %s",
			err)
	}
	// 观察错误。
	go func() {
		errChan := sched.ErrorChan()
		for {
			err, ok := <-errChan
			if !ok {
				break
			}
			t.Errorf("An error occurs when running scheduler: %s", err)
		}
	}()
	// 打印有变化的摘要信息。
	// go func() {
	// 	var count int
	// 	var prevSummary string
	// 	for {
	// 		summary := sched.Summary().String()
	// 		if prevSummary == "" || summary != prevSummary {
	// 			logger.Infof("-- Summary[%d]:\n%s",
	// 				count, summary)
	// 			prevSummary = summary
	// 			count++
	// 		}
	// 		time.Sleep(time.Millisecond * 500)
	// 	}
	// }()
	// 若连续5秒空闲，则停止调度器。
	var count int
	max := 5
	tickCh := time.Tick(time.Second)
	for _ = range tickCh {
		if sched.Idle() {
			count++
			logger.Infof("Increase idle count, and value is %d.", count)
		}
		if count >= max {
			logger.Infof("The idle count is equal or greater than %d.", max)
			break
		}
	}
	if err = sched.Stop(); err != nil {
		t.Fatalf("An error occurs when stopping scheduler: %s",
			err)
	}
	_, ok := <-sched.ErrorChan()
	if ok {
		t.Fatalf("The error channel has not been closed in stopped scheduler!")
	}
	select {
	case <-sched.ErrorChan():
		t.Logf("Closed error channel.")
	}
	logger.Infof("-- Final summary:\n%s", sched.Summary())
}

func TestSchedSendReq(t *testing.T) {
	requestArgs := genRequestArgs([]string{}, 0)
	dataArgs := genDataArgs(10, 2, 1)
	moduleArgs := genSimpleModuleArgs(3, 2, 1, t)
	sched := NewScheduler()
	err := sched.Init(
		requestArgs,
		dataArgs,
		moduleArgs)
	if err != nil {
		t.Fatalf("An error occurs when initializing scheduler: %s", err)
	}
	url := "http://cn.bing.com/search?q=golang"
	firstHTTPReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	err = sched.Start(firstHTTPReq)
	if err != nil {
		t.Fatalf("An error occurs when starting scheduler: %s",
			err)
	}
	mySched := sched.(*myScheduler)
	urlMapLen := mySched.urlMap.Len()
	if urlMapLen != 1 {
		t.Fatalf("Inconsistent URL map length: expected: %d, actual: %d",
			1, urlMapLen)
	}
	// 测试参数无效的情况。
	if mySched.sendReq(nil) {
		t.Fatalf("It still can send nil request!")
	}
	url = "http://cn.bing.com/images/search?q=golang"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("An error occurs when creating a HTTP request: %s (url: %s)",
			err, url)
	}
	req := module.NewRequest(httpReq, 0)
	// 测试URL重复的情况。
	if !mySched.sendReq(req) {
		t.Fatalf("Couldn't send request! (request: %#v)",
			req)
	}
	if mySched.sendReq(req) {
		t.Fatalf("It still can send repeated request!")
	}
	mySched.urlMap, _ = cmap.NewConcurrentMap(16, nil)
	// 测试scheme不匹配的情况。
	httpReq.URL.Scheme = "tcp"
	if mySched.sendReq(req) {
		t.Fatalf("It still can send request with unsupported URL scheme!")
	}
	// 测试URL无效的情况。
	httpReq.URL = nil
	if mySched.sendReq(req) {
		t.Fatalf("It still can send request with nil URL!")
	}
	// 测试调度器已停止的情况。
	sched.Stop()
	time.Sleep(time.Millisecond * 500)
	if mySched.sendReq(nil) {
		t.Fatalf("It still can send request in stopped scheduler!")
	}
}

func TestSendResp(t *testing.T) {
	// 测试响应无效的情况。
	buffer, _ := buffer.NewPool(10, 2)
	if sendResp(nil, buffer) {
		t.Fatalf("It still can send nil response!")
	}
	// 测试响应无效的情况。
	httpReq, _ :=
		http.NewRequest("GET", "https://github.com/gopcp", nil)
	httpResp := &http.Response{
		Request: httpReq,
		Body:    nil,
	}
	resp := module.NewResponse(httpResp, 0)
	buffer.Close()
	done := sendResp(resp, buffer)
	runtime.Gosched()
	if done {
		t.Fatalf("It still can send response with closed buffer!")
	}
}

func TestSendItem(t *testing.T) {
	// 测试响应无效的情况。
	buffer, _ := buffer.NewPool(10, 2)
	if sendItem(nil, buffer) {
		t.Fatalf("It still can send nil item!")
	}
	// 测试响应无效的情况。
	item := module.Item(map[string]interface{}{})
	buffer.Close()
	done := sendItem(item, buffer)
	runtime.Gosched()
	if done {
		t.Fatalf("It still can send item with closed buffer!")
	}
}
