package logger

import (
	"context"
	"github.com/rs/zerolog"
	"go-kit/src/common/fault"
	"os"
)

type Logger interface {
	Info(ctx context.Context, msg string, v ...interface{})
	Debug(ctx context.Context, msg string, v ...interface{})
	Warn(ctx context.Context, err error, msg string, v ...interface{})
	Error(ctx context.Context, err error, msg string, v ...interface{})
	Fault(ctx context.Context, err error, msg string, v ...interface{})
	Fields(kv interface{}) Logger
}

type noopLogger struct{}

func (noopLogger) Info(context.Context, string, ...interface{})         {}
func (noopLogger) Debug(context.Context, string, ...interface{})        {}
func (noopLogger) Warn(context.Context, error, string, ...interface{})  {}
func (noopLogger) Fault(context.Context, error, string, ...interface{}) {}
func (noopLogger) Error(context.Context, error, string, ...interface{}) {}

func (noopLogger) Fields(interface{}) Logger { return noopLogger{} }

type StandardLogger struct {
	l *zerolog.Logger
}

func NewLogger() Logger {
	// rename error field for search log error easier
	zerolog.ErrorFieldName = "fault"
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return StandardLogger{l: &l}
}

func (s StandardLogger) getScopedLogger(ctx context.Context) *zerolog.Logger {
	newLog := *s.l

	ctxMeta := contextValue(ctx)
	if len(ctxMeta) > 0 {
		newLog = newLog.With().Fields(ctxMeta).Logger()
	}

	return &newLog
}

func (s StandardLogger) Info(ctx context.Context, msg string, v ...interface{}) {
	l := s.getScopedLogger(ctx)

	l.Info().Msgf(msg, v...)
}

func (s StandardLogger) Error(ctx context.Context, err error, msg string, v ...interface{}) {
	l := s.getScopedLogger(ctx)
	meta := fault.ExtractMeta(err)

	if len(meta) > 0 {
		l.Error().Err(err).Interface("meta_fault", meta).Msgf(msg, v...)
	} else {
		l.Error().Err(err).Msgf(msg, v...)
	}
}

func (s StandardLogger) Debug(ctx context.Context, msg string, v ...interface{}) {
	l := s.getScopedLogger(ctx)

	l.Debug().Msgf(msg, v...)
}

func (s StandardLogger) Warn(ctx context.Context, err error, msg string, v ...interface{}) {
	l := s.getScopedLogger(ctx)
	meta := fault.ExtractMeta(err)

	if len(meta) > 0 {
		l.Warn().Err(err).Interface("meta_fault", meta).Msgf(msg, v...)
	} else {
		l.Warn().Err(err).Msgf(msg, v...)
	}
}

func (s StandardLogger) Fault(ctx context.Context, err error, msg string, v ...interface{}) {
	isServerErr := fault.IsServerErr(err)

	if isServerErr {
		s.Error(ctx, err, msg, v...)
	} else {
		s.Warn(ctx, err, msg, v...)
	}
}

func (s StandardLogger) Fields(kv interface{}) Logger {
	newLog := s.l.With().Fields(kv).Logger()

	return &StandardLogger{l: &newLog}
}
