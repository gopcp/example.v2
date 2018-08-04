package module

import (
	"sync/atomic"
)

// defaultFakeDownloader 代表默认的仿造下载器。
var defaultFakeDownloader = NewFakeDownloader(MID("D0"), CalculateScoreSimple)

// defaultFakeAnalyzer 代表默认的仿造分析器。
var defaultFakeAnalyzer = NewFakeAnalyzer(MID("A1"), CalculateScoreSimple)

// defaultFakePipeline 代表默认的仿造条目处理管道。
var defaultFakePipeline = NewFakePipeline(MID("P2"), CalculateScoreSimple)

// fakeModules 代表默认仿造组件的切片。
var fakeModules = []Module{
	defaultFakeDownloader,
	defaultFakeAnalyzer,
	defaultFakePipeline,
}

// defaultFakeModuleMap 代表组件类型与默认仿造实例的映射。
var defaultFakeModuleMap = map[Type]Module{
	TYPE_DOWNLOADER: defaultFakeDownloader,
	TYPE_ANALYZER:   defaultFakeAnalyzer,
	TYPE_PIPELINE:   defaultFakePipeline,
}

// fakeModuleFuncMap 代表组件类型与仿造实例生成函数的映射。
var fakeModuleFuncMap = map[Type]func(mid MID) Module{
	TYPE_DOWNLOADER: func(mid MID) Module {
		return NewFakeDownloader(mid, CalculateScoreSimple)
	},
	TYPE_ANALYZER: func(mid MID) Module {
		return NewFakeAnalyzer(mid, CalculateScoreSimple)
	},
	TYPE_PIPELINE: func(mid MID) Module {
		return NewFakePipeline(mid, CalculateScoreSimple)
	},
}

// fakeModule 代表仿造的组件。
type fakeModule struct {
	// mid 代表组件ID。
	mid MID
	// score 代表组件评分。
	score uint64
	// count 代表组件基础计数。
	count uint64
	// scoreCalculator 代表评分计算器。
	scoreCalculator CalculateScore
}

func (fm *fakeModule) ID() MID {
	return fm.mid
}

func (fm *fakeModule) Addr() string {
	parts, err := SplitMID(fm.mid)
	if err == nil {
		return parts[2]
	}
	return ""
}

func (fm *fakeModule) Score() uint64 {
	return atomic.LoadUint64(&fm.score)
}

func (fm *fakeModule) SetScore(score uint64) {
	atomic.StoreUint64(&fm.score, score)
}

func (fm *fakeModule) ScoreCalculator() CalculateScore {
	return fm.scoreCalculator
}

func (fm *fakeModule) CalledCount() uint64 {
	return fm.count + 10
}

func (fm *fakeModule) AcceptedCount() uint64 {
	return fm.count + 8
}

func (fm *fakeModule) CompletedCount() uint64 {
	return fm.count + 6
}

func (fm *fakeModule) HandlingNumber() uint64 {
	return fm.count + 2
}

func (fm *fakeModule) Counts() Counts {
	return Counts{
		fm.CalledCount(),
		fm.AcceptedCount(),
		fm.CompletedCount(),
		fm.HandlingNumber(),
	}
}

func (fm *fakeModule) Summary() SummaryStruct {
	return SummaryStruct{}
}

// NewFakeAnalyzer 用于创建一个仿造的分析器实例。
func NewFakeAnalyzer(mid MID, scoreCalculator CalculateScore) Analyzer {
	return &fakeAnalyzer{
		fakeModule: fakeModule{
			mid:             mid,
			scoreCalculator: scoreCalculator,
		},
	}
}

// fakeAnalyzer 代表分析器的仿造类型。
type fakeAnalyzer struct {
	// fakeModule 代表仿造的组件实例。
	fakeModule
}

func (analyzer *fakeAnalyzer) RespParsers() []ParseResponse {
	return nil
}

func (analyzer *fakeAnalyzer) Analyze(resp *Response) (dataList []Data, errorList []error) {
	return
}

// NewFakeDownloader 用于创建一个仿造的下载器实例。
func NewFakeDownloader(mid MID, scoreCalculator CalculateScore) Downloader {
	return &fakeDownloader{
		fakeModule: fakeModule{
			mid:             mid,
			scoreCalculator: scoreCalculator,
		},
	}
}

// fakeDownloader 代表下载器的实现类型。
type fakeDownloader struct {
	// fakeModule 代表仿造的组件实例。
	fakeModule
}

func (downloader *fakeDownloader) Download(req *Request) (*Response, error) {
	return nil, nil
}

// NewFakePipeline 用于创建一个仿造的条目处理管道实例。
func NewFakePipeline(mid MID, scoreCalculator CalculateScore) Pipeline {
	return &fakePipeline{
		fakeModule: fakeModule{
			mid:             mid,
			scoreCalculator: scoreCalculator,
		},
	}
}

// fakePipeline 代表条目处理管道的实现类型。
type fakePipeline struct {
	// fakeModule 代表仿造的组件实例。
	fakeModule
	// failFast 代表处理是否需要快速失败。
	failFast bool
}

func (pipeline *fakePipeline) ItemProcessors() []ProcessItem {
	return nil
}

func (pipeline *fakePipeline) Send(item Item) []error {
	return nil
}

func (pipeline *fakePipeline) FailFast() bool {
	return pipeline.failFast
}

func (pipeline *fakePipeline) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}
