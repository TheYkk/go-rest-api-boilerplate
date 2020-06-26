package log

import (
	"github.com/labstack/gommon/log"
	"io"
)

//noLog
type nopLogger struct{}

func (nopLogger) Output() io.Writer                         { return nil }
func (nopLogger) SetOutput(w io.Writer)                     {}
func (nopLogger) Prefix() string                            { return "" }
func (nopLogger) SetPrefix(p string)                        {}
func (nopLogger) Level() log.Lvl                            { return log.DEBUG }
func (nopLogger) SetLevel(v log.Lvl)                        {}
func (nopLogger) SetHeader(h string)                        {}
func (nopLogger) Print(i ...interface{})                    {}
func (nopLogger) Printf(format string, args ...interface{}) {}
func (nopLogger) Printj(j log.JSON)                         {}
func (nopLogger) Debug(i ...interface{})                    {}
func (nopLogger) Debugf(format string, args ...interface{}) {}
func (nopLogger) Debugj(j log.JSON)                         {}
func (nopLogger) Info(i ...interface{})                     {}
func (nopLogger) Infof(format string, args ...interface{})  {}
func (nopLogger) Infoj(j log.JSON)                          {}
func (nopLogger) Warn(i ...interface{})                     {}
func (nopLogger) Warnf(format string, args ...interface{})  {}
func (nopLogger) Warnj(j log.JSON)                          {}
func (nopLogger) Error(i ...interface{})                    {}
func (nopLogger) Errorf(format string, args ...interface{}) {}
func (nopLogger) Errorj(j log.JSON)                         {}
func (nopLogger) Fatal(i ...interface{})                    {}
func (nopLogger) Fatalj(j log.JSON)                         {}
func (nopLogger) Fatalf(format string, args ...interface{}) {}
func (nopLogger) Panic(i ...interface{})                    {}
func (nopLogger) Panicj(j log.JSON)                         {}
func (nopLogger) Panicf(format string, args ...interface{}) {}
