package scheduler

import (
	"errors"
	"testing"

	werrors "gopcp.v2/chapter6/webcrawler/errors"
	"gopcp.v2/chapter6/webcrawler/toolkit/buffer"
	"gopcp.v2/chapter6/webcrawler/module"
)

func TestErrorGen(t *testing.T) {
	simpleErrMsg := "testing error"
	expectedErrType := werrors.ERROR_TYPE_SCHEDULER
	err := genError(simpleErrMsg)
	ce, ok := err.(werrors.CrawlerError)
	if !ok {
		t.Fatalf("Inconsistent error type: expected: %T, actual: %T",
			werrors.NewCrawlerError("", ""), err)
	}
	if ce.Type() != expectedErrType {
		t.Fatalf("Inconsistent error type string: expected: %q, actual: %q",
			expectedErrType, ce.Type())
	}
	expectedErrMsg := "crawler error: scheduler error: " + simpleErrMsg
	if ce.Error() != expectedErrMsg {
		t.Fatalf("Inconsistent error message: expected: %q, actual: %q",
			expectedErrMsg, ce.Error())
	}
}

func TestErrorGenByError(t *testing.T) {
	simpleErrMsg := "testing error"
	simpleErr := errors.New(simpleErrMsg)
	expectedErrType := werrors.ERROR_TYPE_SCHEDULER
	err := genErrorByError(simpleErr)
	ce, ok := err.(werrors.CrawlerError)
	if !ok {
		t.Fatalf("Inconsistent error type: expected: %T, actual: %T",
			werrors.NewCrawlerError("", ""), err)
	}
	if ce.Type() != expectedErrType {
		t.Fatalf("Inconsistent error type string: expected: %q, actual: %q",
			expectedErrType, ce.Type())
	}
	expectedErrMsg := "crawler error: scheduler error: " + simpleErrMsg
	if ce.Error() != expectedErrMsg {
		t.Fatalf("Inconsistent error message: expected: %q, actual: %q",
			expectedErrMsg, ce.Error())
	}
}

func TestParameterErrorGen(t *testing.T) {
	simpleErrMsg := "testing error"
	expectedErrType := werrors.ERROR_TYPE_SCHEDULER
	err := genParameterError(simpleErrMsg)
	ce, ok := err.(werrors.CrawlerError)
	if !ok {
		t.Fatalf("Inconsistent error type: expected: %T, actual: %T",
			werrors.NewCrawlerError("", ""), err)
	}
	if ce.Type() != expectedErrType {
		t.Fatalf("Inconsistent error type string: expected: %q, actual: %q",
			expectedErrType, ce.Type())
	}
	expectedErrMsg := "crawler error: scheduler error: illegal parameter: " + simpleErrMsg
	if ce.Error() != expectedErrMsg {
		t.Fatalf("Inconsistent error message: expected: %q, actual: %q",
			expectedErrMsg, ce.Error())
	}
}

func TestErrorSend(t *testing.T) {
	cerr := werrors.NewCrawlerError(
		werrors.ERROR_TYPE_SCHEDULER, "testing error")
	mid := module.MID("")
	buffer, _ := buffer.NewPool(10, 2)
	if !sendError(cerr, mid, buffer) {
		t.Fatalf("Couldn't send error! (error: %s, MID: %s, buffer: %#v)",
			cerr, mid, buffer)
	}
	err := errors.New("testing error")
	if !sendError(err, mid, buffer) {
		t.Fatalf("Couldn't send error! (error: %s, MID: %s, buffer: %#v)",
			err, mid, buffer)
	}
	mids := []module.MID{
		module.MID("D0"),
		module.MID("A0"),
		module.MID("P0"),
	}
	for _, mid := range mids {
		if !sendError(err, mid, buffer) {
			t.Fatalf("Couldn't send error! (error: %s, MID: %s, buffer: %#v)",
				err, mid, buffer)
		}
	}
	if sendError(nil, mid, buffer) {
		t.Fatalf("It still can send error with nil error!")
	}
	if sendError(err, mid, nil) {
		t.Fatalf("It still can send error with nil buffer!")
	}
	buffer.Close()
	if sendError(err, mid, buffer) {
		t.Fatalf("It still can send error with closed buffer!")
	}
}
