package scheduler

import (
	"testing"

	"gopcp.v2/chapter6/webcrawler/module"
)

func TestSummaryNew(t *testing.T) {
	requestArgs := genRequestArgs([]string{}, 0)
	dataArgs := genDataArgs(10, 2, 0)
	moduleArgs := genSimpleModuleArgs(1, 1, 1, t)
	sched := NewScheduler()
	sched.Init(requestArgs, dataArgs, moduleArgs)
	summary := newSchedSummary(
		requestArgs,
		dataArgs,
		moduleArgs,
		sched.(*myScheduler))
	if summary == nil {
		t.Fatal("Couldn't new sched summary!")
	}
	summary = newSchedSummary(
		requestArgs,
		dataArgs,
		moduleArgs,
		nil)
	if summary != nil {
		t.Fatalf("It still can new sched summary with nil sched!")
	}
}

func TestSummaryStruct(t *testing.T) {
	requestArgs := genRequestArgs([]string{}, 0)
	dataArgs := genDataArgs(10, 2, 0)
	moduleArgs := genSimpleModuleArgs(3, 3, 3, t)
	sched := NewScheduler()
	sched.Init(requestArgs, dataArgs, moduleArgs)
	summary := newSchedSummary(
		requestArgs,
		dataArgs,
		moduleArgs,
		sched.(*myScheduler))
	if summary == nil {
		t.Fatal("Couldn't new sched summary!")
	}
	one := summary.Struct()
	another := summary.Struct()
	if !one.Same(another) {
		t.Fatalf("Different scheduler summaries: one: %#v, another: %#v",
			one, another)
	}
	// 测试摘要不同的情况。
	// 不同的请求参数。
	another.RequestArgs.MaxDepth = 11
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different request arguments!")
	}
	another.RequestArgs.MaxDepth = one.RequestArgs.MaxDepth
	// 不同的数据参数。
	another.DataArgs.ReqBufferCap = 11
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different data arguments!")
	}
	another.DataArgs = one.DataArgs
	// 不同的组件参数。
	another.ModuleArgs.DownloaderListSize = 11
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different module arguments summary!")
	}
	another.ModuleArgs = one.ModuleArgs
	// 不同的调度器状态。
	another.Status = "stopped"
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different status!")
	}
	another.Status = one.Status
	// 不同的下载器摘要。
	another.Downloaders = nil
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different downloaders summary!")
	}
	another.Downloaders = make([]module.SummaryStruct, len(one.Downloaders))
	copy(another.Downloaders, one.Downloaders)
	another.Downloaders = another.Downloaders[0:2]
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different downloaders summary number!")
	}
	another.Downloaders = make([]module.SummaryStruct, len(one.Downloaders))
	copy(another.Downloaders, one.Downloaders)
	another.Downloaders[0] = module.SummaryStruct{}
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different downloader summary!")
	}
	another.Downloaders = make([]module.SummaryStruct, len(one.Downloaders))
	copy(another.Downloaders, one.Downloaders)
	// 不同的分析器摘要。
	another.Analyzers = nil
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different analyzers summary!")
	}
	another.Analyzers = make([]module.SummaryStruct, len(one.Analyzers))
	copy(another.Analyzers, one.Analyzers)
	another.Analyzers = another.Analyzers[0:2]
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different analyzers summary number!")
	}
	another.Analyzers = make([]module.SummaryStruct, len(one.Analyzers))
	copy(another.Analyzers, one.Analyzers)
	another.Analyzers[0] = module.SummaryStruct{}
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different analyzer summary!")
	}
	another.Analyzers = make([]module.SummaryStruct, len(one.Analyzers))
	copy(another.Analyzers, one.Analyzers)
	// 不同的条目处理管道摘要。
	another.Pipelines = nil
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different item pipelines summary!")
	}
	another.Pipelines = make([]module.SummaryStruct, len(one.Pipelines))
	copy(another.Pipelines, one.Pipelines)
	another.Pipelines = another.Pipelines[0:2]
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different item pipelines summary number!")
	}
	another.Pipelines = make([]module.SummaryStruct, len(one.Pipelines))
	copy(another.Pipelines, one.Pipelines)
	another.Pipelines[0] = module.SummaryStruct{}
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different item pipelines summary!")
	}
	another.Pipelines = make([]module.SummaryStruct, len(one.Pipelines))
	copy(another.Pipelines, one.Pipelines)
	// 不同的请求缓冲池摘要。
	another.ReqBufferPool.Total = 10
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different request buffer summary!")
	}
	another.ReqBufferPool = one.ReqBufferPool
	// 不同的响应缓冲池摘要。
	another.RespBufferPool.Total = 11
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different response buffer summary!")
	}
	another.RespBufferPool = one.RespBufferPool
	// 不同的条目缓冲池摘要。
	another.ItemBufferPool.Total = 12
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different item buffer summary!")
	}
	another.ItemBufferPool = one.ItemBufferPool
	// 不同的错误缓冲池摘要。
	another.ErrorBufferPool.Total = 13
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different error buffer summary!")
	}
	another.ErrorBufferPool = one.ErrorBufferPool
	// 不同的URL数量。
	another.NumURL = 14
	if one.Same(another) {
		t.Fatalf("Same scheduler summaries with different URL number!")
	}
	another.NumURL = one.NumURL
	if !one.Same(another) {
		t.Fatalf("Different scheduler summaries: one: %#v, another: %#v",
			one, another)
	}
}

func TestSummaryString(t *testing.T) {
	requestArgs := genRequestArgs([]string{}, 0)
	dataArgs := genDataArgs(10, 2, 0)
	moduleArgs := genSimpleModuleArgs(2, 2, 1, t)
	sched := NewScheduler()
	sched.Init(requestArgs, dataArgs, moduleArgs)
	summary := newSchedSummary(
		requestArgs,
		dataArgs,
		moduleArgs,
		sched.(*myScheduler))
	if summary == nil {
		t.Fatal("Couldn't new sched summary!")
	}
	expectedSummaryStr := `{
    "request_args": {
        "accepted_primary_domains": [],
        "max_depth": 0
    },
    "data_args": {
        "req_buffer_cap": 10,
        "req_max_buffer_number": 2,
        "resp_buffer_cap": 10,
        "resp_max_buffer_number": 2,
        "item_buffer_cap": 10,
        "item_max_buffer_number": 2,
        "error_buffer_cap": 10,
        "error_max_buffer_number": 2
    },
    "module_args": {
        "downloader_list_size": 2,
        "analyzer_List_size": 2,
        "pipeline_list_size": 1
    },
    "status": "initialized",
    "downloaders": [
        {
            "id": "D1",
            "called": 0,
            "accepted": 0,
            "completed": 0,
            "handling": 0
        },
        {
            "id": "D2",
            "called": 0,
            "accepted": 0,
            "completed": 0,
            "handling": 0
        }
    ],
    "analyzers": [
        {
            "id": "A3",
            "called": 0,
            "accepted": 0,
            "completed": 0,
            "handling": 0
        },
        {
            "id": "A4",
            "called": 0,
            "accepted": 0,
            "completed": 0,
            "handling": 0
        }
    ],
    "pipelines": [
        {
            "id": "P5",
            "called": 0,
            "accepted": 0,
            "completed": 0,
            "handling": 0,
            "extra": {
                "fail_fast": false,
                "processor_number": 1
            }
        }
    ],
    "request_buffer_pool": {
        "buffer_cap": 10,
        "max_buffer_number": 2,
        "buffer_number": 1,
        "total": 0
    },
    "response_buffer_pool": {
        "buffer_cap": 10,
        "max_buffer_number": 2,
        "buffer_number": 1,
        "total": 0
    },
    "item_buffer_pool": {
        "buffer_cap": 10,
        "max_buffer_number": 2,
        "buffer_number": 1,
        "total": 0
    },
    "error_buffer_pool": {
        "buffer_cap": 10,
        "max_buffer_number": 2,
        "buffer_number": 1,
        "total": 0
    },
    "url_number": 0
}`
	summaryStr := summary.String()
	if summaryStr != expectedSummaryStr {
		t.Fatalf("Inconsistent sheduler summary: expected:\n%s\nactual:\n%s",
			expectedSummaryStr, summaryStr)
	}
}
