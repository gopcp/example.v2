package lib

import (
	"time"
)

// Caller 表示调用器的接口。
type Caller interface {
	// 构建请求。
	BuildReq() RawReq
	// 调用。
	Call(req []byte, timeoutNS time.Duration) ([]byte, error)
	// 检查响应。
	CheckResp(rawReq RawReq, rawResp RawResp) *CallResult
}
