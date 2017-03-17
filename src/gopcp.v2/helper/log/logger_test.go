package log

import (
	"os"
	"testing"

	"gopcp.v2/helper/log/base"
)

func TestLogger(t *testing.T) {
	// Test default logger
	logger := DLogger()
	if logger == nil {
		t.Fatal("The default logger is invalid!")
	}
	if logger.Name() != "logrus" {
		t.Fatalf("Inconsistent logger type: expected: %s, actual: %s",
			"logrus", logger.Name())
	}
	t.Logf("The default logger: %#v\n", logger)

	// Test logger based on logrus
	loggerType := base.TYPE_LOGRUS
	loggerLevel := base.LEVEL_DEBUG
	logFormat := base.FORMAT_JSON
	options := []base.Option{
		base.OptWithLocation{Value: true},
	}
	logger = Logger(
		loggerType,
		loggerLevel,
		logFormat,
		os.Stderr,
		options)
	if logger == nil {
		t.Fatal("The logrus logger is invalid!")
	}
	if logger.Name() != "logrus" {
		t.Fatalf("Inconsistent logger type: expected: %s, actual: %s",
			"logrus", logger.Name())
	}
	if logger.Level() != loggerLevel {
		t.Fatalf("Inconsistent log level: expected: %d, actual: %d",
			loggerLevel, logger.Level())
	}
	if logger.Format() != logFormat {
		t.Fatalf("Inconsistent log format: expected: %s, actual: %s",
			logFormat, logger.Format())
	}
	t.Logf("The logrus logger: %#v\n", logger)
}
