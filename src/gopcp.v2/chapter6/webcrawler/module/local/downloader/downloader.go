package downloader

import (
	"net/http"

	"gopcp.v2/chapter6/webcrawler/module"
	"gopcp.v2/chapter6/webcrawler/module/stub"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

// New 用于创建一个下载器实例。
func New(
	mid module.MID,
	client *http.Client,
	scoreCalculator module.CalculateScore) (module.Downloader, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, genParameterError("nil http client")
	}
	return &myDownloader{
		ModuleInternal: moduleBase,
		httpClient:     *client,
	}, nil
}

// myDownloader 代表下载器的实现类型。
type myDownloader struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// httpClient 代表下载用的HTTP客户端。
	httpClient http.Client
}

func (downloader *myDownloader) Download(req *module.Request) (*module.Response, error) {
	downloader.ModuleInternal.IncrHandlingNumber()
	defer downloader.ModuleInternal.DecrHandlingNumber()
	downloader.ModuleInternal.IncrCalledCount()
	if req == nil {
		return nil, genParameterError("nil request")
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		return nil, genParameterError("nil HTTP request")
	}
	downloader.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Do the request (URL: %s, depth: %d)... \n", httpReq.URL, req.Depth())
	httpResp, err := downloader.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	downloader.ModuleInternal.IncrCompletedCount()
	return module.NewResponse(httpResp, req.Depth()), nil
}
