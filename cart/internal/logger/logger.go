package logger

import (
	"context"
	"go.uber.org/zap"
	"sync"
)

var (
	globalLogger *Logger
	once         = sync.Once{}
)

type Logger struct {
	l *zap.SugaredLogger
}

func NewLogger(config zap.Config) *Logger {
	l, err := config.Build(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}

	once.Do(func() {
		globalLogger = &Logger{l: l.Sugar()}
	})

	return globalLogger
}

func NewNopLogger() *Logger {
	l := zap.NewNop() // Пустой логгер
	once.Do(func() {
		globalLogger = &Logger{l: l.Sugar()}
	})
	return globalLogger
}

type CtxKey string

const LoggerCtxKey CtxKey = "logger"

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerCtxKey, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		// Если логгер не найден, возвращаем стандартный логгер (например, без настроек)
		return nil
	}
	return logger
}

func Infow(ctx context.Context, msg string, args ...interface{}) {
	if ctx != nil {
		if l, ok := ctx.Value(LoggerCtxKey).(*Logger); ok && l != nil {
			l.l.Infow(msg, args)
			return
		}
	}

	if globalLogger == nil {
		panic("global logger is nil")
	}
	globalLogger.l.Infow(msg, args...)
}

func Debugw(ctx context.Context, msg string, args ...interface{}) {
	if ctx != nil {
		if l, ok := ctx.Value(LoggerCtxKey).(*Logger); ok && l != nil {
			l.l.Debugw(msg, args)
			return
		}
	}

	if globalLogger == nil {
		panic("global logger is nil")
	}

	globalLogger.l.Debugw(msg, args...)
}

func Warnw(ctx context.Context, msg string, args ...interface{}) {
	if ctx != nil {
		if l, ok := ctx.Value(LoggerCtxKey).(*Logger); ok && l != nil {
			l.l.Warnw(msg, args)
			return
		}
	}

	if globalLogger == nil {
		panic("global logger is nil")
	}

	globalLogger.l.Warnw(msg, args...)
}

func Errorw(ctx context.Context, msg string, args ...interface{}) {
	if l, ok := ctx.Value(LoggerCtxKey).(*Logger); ok && l != nil {
		l.l.Errorf(msg, args)
		return
	}

	if globalLogger == nil {
		panic("global logger is nil")
	}

	globalLogger.l.Errorw(msg, args...)
}

func PanicF(ctx context.Context, templateMsg string, args ...interface{}) {
	if l, ok := ctx.Value(LoggerCtxKey).(*Logger); ok && l != nil {
		l.l.Panicf(templateMsg, args...)
		return
	}

	if globalLogger == nil {
		panic("global logger is nil")
	}

	globalLogger.l.Panicf(templateMsg, args...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}
