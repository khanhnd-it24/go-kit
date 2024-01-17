package logger

import (
	"context"
	"sync/atomic"
)

type ctxLogKey struct{}

type (
	logHolder struct {
		l Logger
	}
)

var (
	globalLogger = defaultLoggerValue()
)

func defaultLoggerValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(logHolder{l: defaultLogger()})
	return v
}

func defaultLogger() Logger {
	return noopLogger{}
}

func getLogger() Logger {
	return globalLogger.Load().(logHolder).l
}

func SetLogger(logger Logger) {
	globalLogger.Store(logHolder{l: logger})
}

func Info(c context.Context, s string, v ...interface{}) {
	getLogger().Info(c, s, v...)
}

func Error(c context.Context, e error, s string, v ...interface{}) {
	getLogger().Error(c, e, s, v...)
}

func Debug(c context.Context, s string, v ...interface{}) {
	getLogger().Debug(c, s, v...)
}

func Warn(c context.Context, err error, s string, v ...interface{}) {
	getLogger().Warn(c, err, s, v...)
}

func Fault(c context.Context, err error, s string, v ...interface{}) {
	getLogger().Fault(c, err, s, v...)
}

func Fields(kv interface{}) Logger {
	return getLogger().Fields(kv)
}

func WithContextValue(c context.Context, meta map[string]interface{}) context.Context {
	metaInCtx, ok := c.Value(ctxLogKey{}).(map[string]interface{})
	if !ok {
		return context.WithValue(c, ctxLogKey{}, meta)
	}
	mergeMap := make(map[string]interface{})
	for k, v := range metaInCtx {
		mergeMap[k] = v
	}

	for k, v := range meta {
		mergeMap[k] = v
	}

	return context.WithValue(c, ctxLogKey{}, mergeMap)
}

func contextValue(c context.Context) map[string]interface{} {
	v, ok := c.Value(ctxLogKey{}).(map[string]interface{})

	if !ok {
		return nil
	}
	return v
}
