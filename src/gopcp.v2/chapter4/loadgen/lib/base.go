package lib

import (
	"time"
)

// RawReq 表示原生请求的结构。
type RawReq struct {
	ID  int64
	Req []byte
}

// RawResp 表示原生响应的结构。
type RawResp struct {
	ID     int64
	Resp   []byte
	Err    error
	Elapse time.Duration
}

// RetCode 表示结果代码的类型。
type RetCode int

// 保留 1 ~ 1000 给载荷承受方使用。
const (
	RET_CODE_SUCCESS              RetCode = 0    // 成功。
	RET_CODE_WARNING_CALL_TIMEOUT         = 1001 // 调用超时警告。
	RET_CODE_ERROR_CALL                   = 2001 // 调用错误。
	RET_CODE_ERROR_RESPONSE               = 2002 // 响应内容错误。
	RET_CODE_ERROR_CALEE                  = 2003 // 被调用方（被测软件）的内部错误。
	RET_CODE_FATAL_CALL                   = 3001 // 调用过程中发生了致命错误！
)

// GetRetCodePlain 会依据结果代码返回相应的文字解释。
func GetRetCodePlain(code RetCode) string {
	var codePlain string
	switch code {
	case RET_CODE_SUCCESS:
		codePlain = "Success"
	case RET_CODE_WARNING_CALL_TIMEOUT:
		codePlain = "Call Timeout Warning"
	case RET_CODE_ERROR_CALL:
		codePlain = "Call Error"
	case RET_CODE_ERROR_RESPONSE:
		codePlain = "Response Error"
	case RET_CODE_ERROR_CALEE:
		codePlain = "Callee Error"
	case RET_CODE_FATAL_CALL:
		codePlain = "Call Fatal Error"
	default:
		codePlain = "Unknown result code"
	}
	return codePlain
}

// CallResult 表示调用结果的结构。
type CallResult struct {
	ID     int64         // ID。
	Req    RawReq        // 原生请求。
	Resp   RawResp       // 原生响应。
	Code   RetCode       // 响应代码。
	Msg    string        // 结果成因的简述。
	Elapse time.Duration // 耗时。
}

// 声明代表载荷发生器状态的常量。
const (
	// STATUS_ORIGINAL 代表原始。
	STATUS_ORIGINAL uint32 = 0
	// STATUS_STARTING 代表正在启动。
	STATUS_STARTING uint32 = 1
	// STATUS_STARTED 代表已启动。
	STATUS_STARTED uint32 = 2
	// STATUS_STOPPING 代表正在停止。
	STATUS_STOPPING uint32 = 3
	// STATUS_STOPPED 代表已停止。
	STATUS_STOPPED uint32 = 4
)

// Generator 表示载荷发生器的接口。
type Generator interface {
	// 启动载荷发生器。
	// 结果值代表是否已成功启动。
	Start() bool
	// 停止载荷发生器。
	// 结果值代表是否已成功停止。
	Stop() bool
	// 获取状态。
	Status() uint32
	// 获取调用计数。每次启动会重置该计数。
	CallCount() int64
}
