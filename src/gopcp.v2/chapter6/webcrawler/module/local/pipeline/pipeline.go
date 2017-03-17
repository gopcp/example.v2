package pipeline

import (
	"fmt"

	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/module/stub"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

// New 用于创建一个条目处理管道实例。
func New(
	mid module.MID,
	itemProcessors []module.ProcessItem,
	scoreCalculator module.CalculateScore) (module.Pipeline, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if itemProcessors == nil {
		return nil, genParameterError("nil item processor list")
	}
	if len(itemProcessors) == 0 {
		return nil, genParameterError("empty item processor list")
	}
	var innerProcessors []module.ProcessItem
	for i, pipeline := range itemProcessors {
		if pipeline == nil {
			err := genParameterError(fmt.Sprintf("nil item processor[%d]", i))
			return nil, err
		}
		innerProcessors = append(innerProcessors, pipeline)
	}
	return &myPipeline{
		ModuleInternal: moduleBase,
		itemProcessors: innerProcessors,
	}, nil
}

// myPipeline 代表条目处理管道的实现类型。
type myPipeline struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// itemProcessors 代表条目处理器的列表。
	itemProcessors []module.ProcessItem
	// failFast 代表处理是否需要快速失败。
	failFast bool
}

func (pipeline *myPipeline) ItemProcessors() []module.ProcessItem {
	processors := make([]module.ProcessItem, len(pipeline.itemProcessors))
	copy(processors, pipeline.itemProcessors)
	return processors
}

func (pipeline *myPipeline) Send(item module.Item) []error {
	pipeline.ModuleInternal.IncrHandlingNumber()
	defer pipeline.ModuleInternal.DecrHandlingNumber()
	pipeline.ModuleInternal.IncrCalledCount()
	var errs []error
	if item == nil {
		err := genParameterError("nil item")
		errs = append(errs, err)
		return errs
	}
	pipeline.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Process item %+v... \n", item)
	var currentItem = item
	for _, processor := range pipeline.itemProcessors {
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				break
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	if len(errs) == 0 {
		pipeline.ModuleInternal.IncrCompletedCount()
	}
	return errs
}

func (pipeline *myPipeline) FailFast() bool {
	return pipeline.failFast
}

func (pipeline *myPipeline) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}

// extraSummaryStruct 代表条目处理管道实额外信息的摘要类型。
type extraSummaryStruct struct {
	FailFast        bool `json:"fail_fast"`
	ProcessorNumber int  `json:"processor_number"`
}

func (pipeline *myPipeline) Summary() module.SummaryStruct {
	summary := pipeline.ModuleInternal.Summary()
	summary.Extra = extraSummaryStruct{
		FailFast:        pipeline.failFast,
		ProcessorNumber: len(pipeline.itemProcessors),
	}
	return summary
}
