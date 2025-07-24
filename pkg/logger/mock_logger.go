package logger

import (
	"fmt"
	"testing"
)

type MockLogger struct {
	t *testing.T
}

func NewMockLogger(t *testing.T) *MockLogger {
	return &MockLogger{t: t}
}

type MockLogContext struct {
	t     *testing.T
	level string
	err   error
}

func (l *MockLogger) Debug() LogContext {
	return &MockLogContext{t: l.t, level: "debug"}
}

func (l *MockLogger) Info() LogContext {
	return &MockLogContext{t: l.t, level: "info"}
}

func (l *MockLogger) Warn() LogContext {
	return &MockLogContext{t: l.t, level: "warn"}
}

func (l *MockLogger) Error() LogContext {
	return &MockLogContext{t: l.t, level: "error"}
}

func (l *MockLogger) Fatal() LogContext {
	return &MockLogContext{t: l.t, level: "fatal"}
}

func (l *MockLogger) Err(err error) LogContext {
	return &MockLogContext{t: l.t, level: "error", err: err}
}

func (c *MockLogContext) Msg(message string) {
	c.t.Logf("[%s] %s (err: %v)", c.level, message, c.err)
}

func (c *MockLogContext) Msgf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	c.t.Logf("[%s] %s (err: %v)", c.level, message, c.err)
}

func (c *MockLogContext) Err(err error) LogContext {
	c.err = err
	return c
}

func (c *MockLogContext) Debug() LogContext {
	c.level = "debug"
	return c
}

func (c *MockLogContext) Info() LogContext {
	c.level = "info"
	return c
}

func (c *MockLogContext) Warn() LogContext {
	c.level = "warn"
	return c
}

func (c *MockLogContext) Error() LogContext {
	c.level = "error"
	return c
}

func (c *MockLogContext) Fatal() LogContext {
	c.level = "fatal"
	return c
}
